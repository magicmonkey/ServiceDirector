package main

import (
	"ServiceRegistry"
	"Interfaces/http"
	"Interfaces/update"
	"Interfaces/persistor"
	"Interfaces/replication"
	"math/rand"
	"time"
	"flag"
	"log"
)

func main() {

	// Seed the RNG.  This isn't cryptography, it doesn't matter if the RNG is predictable.
	rand.Seed(time.Now().UnixNano())

	role := flag.String("role", "master", "master or slave")
	masterAddr := flag.String("master", "127.0.0.1:8083", "Address of master to connect to")
	httpAddr := flag.String("httpAddr", "0.0.0.0:8081", "Address to bind the HTTP interface to")
	httpUpdateAddr := flag.String("httpUpdateAddr", "0.0.0.0:8082", "Address to bind the HTTP update interface to")
	flag.Parse()

	log.Println("[Main] Role is", *role)

	switch (*role) {
	case "master":
		runMaster(*httpAddr, *httpUpdateAddr)
	case "slave":
		runSlave(*masterAddr, *httpAddr)
	}
}

func runMaster(httpAddr string, httpUpdateAddr string) {

	var sr ServiceRegistry.ServiceRegistry

	// The Persistor is the thing which saves any updates to Redis
	// Also the initial ServiceRegistry is loaded from it
	sru1 := make(chan ServiceRegistry.ServiceRegistry, 10)
	p := persistor.NewPersistor()
	go p.Listen(sru1)

	// The Master is the thing which allows slaves to connect and get updates
	sr = p.LoadServiceRegistry("FirstRegistry")

	sru2 := make(chan ServiceRegistry.ServiceRegistry, 10)
	go replication.StartListener(sru2)

	sru3 := make(chan ServiceRegistry.ServiceRegistry, 10)

	sr.RegisterUpdateChannel(sru1)
	sr.RegisterUpdateChannel(sru2)
	sr.RegisterUpdateChannel(sru3)

	requestUpdate := make(chan bool, 10)

	finished1 := make(chan bool, 10)
	finished2 := make(chan bool, 10)
	go http.RunHTTP(sru3, httpAddr, finished1, requestUpdate)
	go update.RunHTTP(&sr, httpUpdateAddr, finished2, requestUpdate)

	for {
		select {
		case <-finished1:
			return
		case <-finished2:
			return
		case <-requestUpdate:
			log.Println("[Main] Someone has requested an update")
			sr.SendRegistryUpdate()
		case <-time.After(30*time.Second):
			sru1 <- sr
		}
	}

}

func runSlave(masterAddr string, httpAddr string) {

	var sr *ServiceRegistry.ServiceRegistry
	sru1 := make(chan *ServiceRegistry.ServiceRegistry)
	go replication.StartSlave(masterAddr, sru1)

	sru2 := make(chan ServiceRegistry.ServiceRegistry)

	log.Println("[Main] Master is", masterAddr)
	finished := make(chan bool)
	requestUpdate := make(chan bool)
	go http.RunHTTP(sru2, httpAddr, finished, requestUpdate)

	for {
		select {
		case <-finished:
			log.Println("[Main] HTTP server has exited, so I might as well quit")
			return
		case sr = <-sru1:
			log.Println("[Main] Got updated service registry")
			sru2 <- *sr
		}
	}

	//	sr := ServiceRegistry.NewServiceRegistry("FirstRegistry", sru)
	//	sr.GenerateTestData()

}

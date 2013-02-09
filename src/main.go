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

	var sr *ServiceRegistry.ServiceRegistry

	// The Persistor is the thing which saves any updates to Redis
	// Also the initial ServiceRegistry is loaded from it
	sru2 := make(chan *ServiceRegistry.ServiceRegistry)
	p := persistor.NewPersistor()
	go p.Listen(sru2)

	// The Master is the thing which allows slaves to connect and get updates
	sr = p.LoadServiceRegistry("FirstRegistry")

	sru1 := make(chan *ServiceRegistry.ServiceRegistry)
	go replication.StartListener(sru1, sr)
	sr.RegisterUpdateChannel(sru2)
	sr.RegisterUpdateChannel(sru1)

	c1 := make(chan bool)
	c2 := make(chan bool)
	go http.RunHTTP(&sr, httpAddr, c1)
	go update.RunHTTP(sr, httpUpdateAddr, c2)

	for {
		select {
		case <-c1:
			return
		case <-c2:
			return
		case <-time.After(3*time.Second):
			sru1 <- sr
		}
	}
}

func runSlave(masterAddr string, httpAddr string) {

	var sr *ServiceRegistry.ServiceRegistry
	sru1 := make(chan *ServiceRegistry.ServiceRegistry)

	go replication.StartSlave(masterAddr, sru1)

	log.Println("[Main] Master is", masterAddr)
	c1 := make(chan bool)
	go http.RunHTTP(&sr, httpAddr, c1)

	for {
		select {
		case <-c1:
			log.Println("[Main] HTTP server has exited, so I might as well quit")
			return
		case sr = <-sru1:
			log.Println("[Main] Got updated service registry")
			//			log.Println(sr)
		}
	}

	//	sr := ServiceRegistry.NewServiceRegistry("FirstRegistry", sru)
	//	sr.GenerateTestData()

}

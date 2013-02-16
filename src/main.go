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

// The Master is the thing which allows slaves to connect and get updates
func runMaster(httpAddr string, httpUpdateAddr string) {

	var sr *ServiceRegistry.ServiceRegistry

	// The Persistor is the thing which saves any updates to Redis
	// Also the initial ServiceRegistry is loaded from it
	p := persistor.NewPersistor()
	sru1 := p.Listen()

	sr = p.LoadServiceRegistry("FirstRegistry")

	rep := replication.NewMaster()
	sru2 := rep.StartListener()

	h := http.NewBalancer()
	finished1, requestUpdate, sru3 := h.RunHTTP(httpAddr)

	u := update.NewUpdater()
	finished2 := u.RunHTTP(httpUpdateAddr, sr)

	sr.RegisterUpdateChannel(sru1)
	sr.RegisterUpdateChannel(sru2)
	sr.RegisterUpdateChannel(sru3)

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
			sru1 <- *sr
		}
	}

}

func runSlave(masterAddr string, httpAddr string) {

	var sr *ServiceRegistry.ServiceRegistry
	sru1 := make(chan *ServiceRegistry.ServiceRegistry)
	go replication.StartSlave(masterAddr, sru1)

	log.Println("[Main] Master is", masterAddr)

	h := http.NewBalancer()
	finished, requestUpdate, sru2 := h.RunHTTP(httpAddr)

	for {
		select {
		case <-finished:
			log.Println("[Main] HTTP server has exited, so I might as well quit")
			return
		case <-requestUpdate:
			log.Println("[Main] Ignoring request for an update")
		case sr = <-sru1:
			log.Println("[Main] Got updated service registry")
			sru2 <- *sr
		}
	}

	//	sr := ServiceRegistry.NewServiceRegistry("FirstRegistry", sru)
	//	sr.GenerateTestData()

}

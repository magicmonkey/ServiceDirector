package main

import (
	"ServiceRegistry"
	"Interfaces/http"
	"Interfaces/update"
	"Interfaces/persistor"
	"Interfaces/replication"
	"math/rand"
	"time"
)

func main() {

	// Seed the RNG.  This isn't cryptography, it doesn't matter if the RNG is predictable.
	rand.Seed(time.Now().UnixNano())

	// The Persistor is the thing which saves any updates to Redis
	sru2 := make(chan *ServiceRegistry.ServiceRegistry)
	p := persistor.NewPersistor()
	go p.Listen(sru2)

	// The Master is the thing which allows slaves to connect and get updates
	sr := p.LoadServiceRegistry("FirstRegistry")

	// TODO: Make it so that only a master does saving
	sru1 := make(chan *ServiceRegistry.ServiceRegistry)
	go replication.StartListener(sru1)

	sr.RegisterUpdateChannel(sru2)
	sr.RegisterUpdateChannel(sru1)

	//	sr := ServiceRegistry.NewServiceRegistry("FirstRegistry", sru)
	//	sr.GenerateTestData()

	c1 := make(chan bool)
	c2 := make(chan bool)
	go http.RunHTTP(sr, c1)
	go update.RunHTTP(sr, c2)

	select {
	case <-c1:
	case <-c2:
	}

}

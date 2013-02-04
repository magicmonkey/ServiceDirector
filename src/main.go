package main

import (
	"ServiceRegistry"
	"Interfaces/http"
	"Interfaces/update"
	"Interfaces/persistor"
	"math/rand"
	"time"
)

func main() {

	// Seed the RNG.  This isn't cryptography, it doesn't matter if the RNG is predictable.
	rand.Seed(time.Now().UnixNano())

	sru := make(chan *ServiceRegistry.ServiceRegistry)

	// TODO: Make it so that only a master does saving

	// The Persistor is the thing which saves any updates to Redis
	p := persistor.NewPersistor(sru)
	go p.Listen()

	// The Master is the thing which allows slaves to connect and get updates


	sr := p.LoadServiceRegistry("FirstRegistry", sru)

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

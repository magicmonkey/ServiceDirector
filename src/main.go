package main

import (
	"ServiceRegistry"
	"Interfaces/http"
	"Interfaces/update"
	"math/rand"
	"time"
)

func main() {

	// Seed the RNG.  This isn't cryptography, it doesn't matter if the RNG is predictable.
	rand.Seed(time.Now().UnixNano())

	sr := ServiceRegistry.ServiceRegistry{}

	sr.GenerateTestData()

	c1 := make(chan bool)
	c2 := make(chan bool)

	go http.RunHTTP(&sr, c1)
	go update.RunHTTP(&sr, c2)

	<-c1
	<-c2

}


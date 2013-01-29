// Package to update entries in the Service Registry
package update

import (
	"ServiceRegistry"
	"net/http"
	"fmt"
)

func getServiceHandler(sr *ServiceRegistry.ServiceRegistry) (http.HandlerFunc) {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("Got update (%v) request for %v\n", r.Method, r.URL.Path)

		switch r.Method {
		case "POST":
			if r.Header.Get("Content-Type") != "application/json" {
				http.Error(w, "Server only understands application/json as the content-type", 400)
				return
			}
		case "PUT":
			fmt.Println("You're creating a resource?")
		}

	}
}

// Runs the actual HTTP server, ie spawns the goroutine via http.ListenAndServe
func RunHTTP(sr *ServiceRegistry.ServiceRegistry, c chan bool) {
	sm := http.NewServeMux()
	sm.HandleFunc("/services/", getServiceHandler(sr))
	listenAddr := ":8082"
	fmt.Printf("Starting HTTP server for updates on %v\n", listenAddr)
	if e := http.ListenAndServe(listenAddr, sm); e != nil {
		panic(e)
	}
	c<-true
}

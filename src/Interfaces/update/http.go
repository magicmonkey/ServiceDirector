// Package to update entries in the Service Registry
package update

import (
	"ServiceRegistry"
	"net/http"
	"fmt"
)

func getServiceHandler(sr *ServiceRegistry.ServiceRegistry) (http.HandlerFunc) {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}

// Runs the actual HTTP server, ie spawns the goroutine via http.ListenAndServe
func RunHTTP(sr *ServiceRegistry.ServiceRegistry) {
	http.HandleFunc("/services/", getServiceHandler(sr))
	listenAddr := ":8082"
	fmt.Printf("Starting HTTP server for updates on %v\n", listenAddr)
	if e := http.ListenAndServe(listenAddr, nil); e != nil {
		panic(e)
	}
}

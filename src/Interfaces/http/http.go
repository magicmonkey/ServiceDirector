// Provides an HTTP interface into the Service Registry
package http

import (
	"fmt"
	"net/http"
	"ServiceRegistry"
	"strings"
)

// Returns a handler for /services for the given Service Registry, allowing the URLs beneath /services to refer to the
// services in that registry
func getServiceHandler(sr *ServiceRegistry.ServiceRegistry) (http.HandlerFunc) {
	return func(w http.ResponseWriter, r *http.Request) {
		if pathParts := strings.Split(r.URL.Path, "/"); len(pathParts) > 2 {
			fmt.Fprintf(w, "Service requested: %v\n", pathParts[2])
			svc := sr.GetServiceWithName(pathParts[2])
			svs := svc.GetLocationForVersion(sr.GetVersionFromString(pathParts[3]))
			fmt.Fprintf(w, "Versions: %v\n", svc.Versions)
			fmt.Fprintf(w, "Chosen version: %v\n", svs)
		}
	}
}

// Runs the actual HTTP server, ie spawns the goroutine via http.ListenAndServe
func RunHTTP(sr *ServiceRegistry.ServiceRegistry) {
	http.HandleFunc("/services/", getServiceHandler(sr))
	listenAddr := ":8081"
	fmt.Printf("Starting HTTP server on %v\n", listenAddr)
	if e := http.ListenAndServe(listenAddr, nil); e != nil {
		panic(e)
	}
}

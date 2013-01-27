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
			svc := sr.GetServiceWithName(pathParts[2], false)
			if svc == nil {
				http.Error(w, fmt.Sprintf("No service found with name %v", pathParts[2]), 400)
				return
			}
			svs := svc.GetLocationForVersion(sr.GetVersionFromString(pathParts[3]))
			if svs == nil {
				http.Error(w, fmt.Sprintf("Found service with name %v but could not find an instance with version %v", pathParts[2], pathParts[3]), 400)
				return
			}
			fmt.Fprintln(w, svs.Location.String())
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

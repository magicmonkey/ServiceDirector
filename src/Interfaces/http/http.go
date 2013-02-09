// Provides an HTTP interface into the Service Registry
package http

import (
	"fmt"
	"log"
	"net/http"
	"ServiceRegistry"
	"strings"
)

// Returns a handler for /services for the given Service Registry, allowing the URLs beneath /services to refer to the
// services in that registry
func getServiceHandler(sr **ServiceRegistry.ServiceRegistry) (http.HandlerFunc) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("[HTTP] Got request for %v\n", r.URL.Path)
		pathParts := strings.Split(r.URL.Path, "/")

		if len(pathParts) == 4 {

			// /services/TestService/1.2.4
			svc, _ := (*sr).GetServiceWithName(pathParts[2], false)
			if svc == nil {
				http.Error(w, fmt.Sprintf(`No service found with name %v`, pathParts[2]), 400)
				return
			}
			svs := svc.GetLocationForVersionString(pathParts[3])
			if svs == nil {
				http.Error(w, fmt.Sprintf(`Found service with name %v but could not find an instance with version %v`, pathParts[2], pathParts[3]), 400)
				return
			}
			fmt.Fprintln(w, svs.Location)
			return
		}

		http.Error(w, `Request should be in the format /services/<service>/<version>`, 400)
		return

	}
}

// Runs the actual HTTP server, ie calls http.ListenAndServe
func RunHTTP(sr **ServiceRegistry.ServiceRegistry, listenAddr string, c chan bool) {
	sm := http.NewServeMux()
	sm.HandleFunc("/services/", getServiceHandler(sr))
	log.Printf("[HTTP] Starting HTTP server on %v\n", listenAddr)
	if e := http.ListenAndServe(listenAddr, sm); e != nil {
		panic(e)
	}
	c<-true
}

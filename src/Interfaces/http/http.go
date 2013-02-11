// Provides an HTTP interface into the Service Registry
package http

import (
	"fmt"
	"log"
	"net/http"
	"ServiceRegistry"
	"strings"
)

type httpBalancer struct {
	sr     ServiceRegistry.ServiceRegistry
}

// Returns a handler for /services for the given Service Registry, allowing the URLs beneath /services to refer to the
// services in that registry
func (b *httpBalancer) serviceHandler () (http.HandlerFunc) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("[HTTP] Got request for %v\n", r.URL.Path)
		pathParts := strings.Split(r.URL.Path, "/")

		if len(pathParts) == 4 {

			// /services/TestService/1.2.4
			svc, _ := b.sr.GetServiceWithName(pathParts[2], false)
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

func (b *httpBalancer) listenForUpdates (srChan chan ServiceRegistry.ServiceRegistry) {
	for {
		b.sr = <- srChan
		log.Println("[HTTP] Got an updated service registry")
	}
}

// Runs the actual HTTP server, ie calls http.ListenAndServe
func RunHTTP(srChan chan ServiceRegistry.ServiceRegistry, listenAddr string, finished chan bool, requestUpdate chan bool) {
	h := new(httpBalancer)
	go h.listenForUpdates(srChan)
	requestUpdate<-true
	sm := http.NewServeMux()
	sm.HandleFunc("/services/", h.serviceHandler())
	log.Printf("[HTTP] Starting HTTP server on [%v]\n", listenAddr)
	if e := http.ListenAndServe(listenAddr, sm); e != nil {
		panic(e)
	}
	finished<-true
}

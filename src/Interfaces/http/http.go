// Provides an HTTP interface into the Service Registry
package http

import (
	"fmt"
	"log"
	"net/http"
	"ServiceRegistry"
	"github.com/gorilla/mux"
)

type httpBalancer struct {
	sr ServiceRegistry.ServiceRegistry
}

func NewBalancer() (h *httpBalancer) {
	h = new(httpBalancer)
	return
}

// Returns a handler for /services for the given Service Registry, allowing the URLs beneath /services to refer to the
// services in that registry
func (b *httpBalancer) serviceHandler () (http.HandlerFunc) {

	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("[HTTP] Got request for %v\n", r.URL.Path)
		vars := mux.Vars(r)

		// /services/TestService/1.2.4
		svc, _ := b.sr.GetServiceWithName(vars["service"], false)
		if svc == nil {
			http.Error(w, fmt.Sprintf(`No service found with name %v`, vars["service"]), 400)
			return
		}
		svs := svc.GetLocationForVersionString(vars["version"])
		if svs == nil {
			http.Error(w, fmt.Sprintf(`Found service with name %v but could not find an instance with version %v`, vars["service"], vars["version"]), 400)
			return
		}
		fmt.Fprintln(w, svs.Location)
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
func (h *httpBalancer) RunHTTP(listenAddr string) (finished chan bool, requestUpdate chan bool, updateChannel chan ServiceRegistry.ServiceRegistry) {
	finished = make(chan bool, 10)
	requestUpdate = make(chan bool, 10)
	updateChannel = make(chan ServiceRegistry.ServiceRegistry, 10)
	go h.doRunHTTP(listenAddr, finished, requestUpdate, updateChannel)
	return
}

func (h *httpBalancer) doRunHTTP(listenAddr string, finished chan bool, requestUpdate chan bool, updateChannel chan ServiceRegistry.ServiceRegistry) {
	go h.listenForUpdates(updateChannel)
	requestUpdate<-true

	r := mux.NewRouter()

	sm := http.NewServeMux()
	r.HandleFunc("/services/{service}/{version}", h.serviceHandler())

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, `Request should be in the format /services/<service>/<version>`, 400)
	})

	sm.Handle("/", r)

	log.Printf("[HTTP] Starting HTTP server on [%v]\n", listenAddr)
	if e := http.ListenAndServe(listenAddr, sm); e != nil {
		panic(e)
	}
	finished<-true
}

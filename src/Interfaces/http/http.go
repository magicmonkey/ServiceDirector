package http

import (
	"fmt"
	"net/http"
	"ServiceRegistry"
	"strings"
)

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

func RunHTTP(sr *ServiceRegistry.ServiceRegistry) {
	http.HandleFunc("/services/", getServiceHandler(sr))
	listenAddr := ":8081"
	fmt.Printf("Starting HTTP server on %v\n", listenAddr)
	if e := http.ListenAndServe(listenAddr, nil); e != nil {
		panic(e)
	}
}

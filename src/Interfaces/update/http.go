// Package to update entries in the Service Registry
package update

import (
	`ServiceRegistry`
	`net/http`
	`net/url`
	`fmt`
	`strings`
	`encoding/json`
)

func getServiceHandler(sr *ServiceRegistry.ServiceRegistry) (http.HandlerFunc) {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("Update interface: Got %v request for %v\n", r.Method, r.URL.Path)

		switch r.Method {

		case `POST`:

			pathParts := strings.Split(r.URL.Path, "/")
			if len(pathParts) != 3 {
				http.Error(w, `Should only POST to /services/<service-name>`, 400)
				return
			}

			svc, _ := sr.GetServiceWithName(pathParts[2], false)

			if svc == nil {
				http.Error(w, `Cannot find service, hav you created it via a PUT request?`, 404)
				return
			}

			if r.Header.Get(`Content-Type`) != `application/json` {
				http.Error(w, `Server only understands application/json as the content-type`, 400)
				return
			}

			type SubmittedInstance struct {
				Version  string
				Location string
			}

			req := new(SubmittedInstance)
			json.NewDecoder(r.Body).Decode(req)
			u, _ := url.Parse(req.Location)
			svc.AddServiceInstance(ServiceRegistry.NewVersion(req.Version), ServiceRegistry.NewLocation(u))
			w.WriteHeader(201)
			fmt.Fprintf(w, "Added location %v to service %v version %v\n", req.Location, svc.Name, req.Version)

		case `PUT`:
			pathParts := strings.Split(r.URL.Path, "/")
			if len(pathParts) != 3 {
				http.Error(w, `Should only PUT to /services/<service-name>`, 400)
				return
			}
			// Currently don't care about the body, but we may do later
			svc, created := sr.GetServiceWithName(pathParts[2], true)
			switch created {
			case true:
				w.WriteHeader(201)
				fmt.Fprintf(w, `Created resource %v`, svc.Name)
			case false:
				w.WriteHeader(200)
				fmt.Fprintf(w, `Resource %v already exists`, svc.Name)
			}
		}
	}
}

// Runs the actual HTTP server, ie spawns the goroutine via http.ListenAndServe
func RunHTTP(sr *ServiceRegistry.ServiceRegistry, c chan bool) {
	sm := http.NewServeMux()
	sm.HandleFunc(`/services/`, getServiceHandler(sr))
	listenAddr := `:8082`
	fmt.Printf("Starting HTTP server for updates on %v\n", listenAddr)
	if e := http.ListenAndServe(listenAddr, sm); e != nil {
		panic(e)
	}
	c<-true
}

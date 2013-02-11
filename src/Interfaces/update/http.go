// Package to update entries in the Service Registry
package update

import (
	`ServiceRegistry`
	`net/http`
	`fmt`
	`strings`
	`encoding/json`
	`log`
)

func getServiceHandler(sr *ServiceRegistry.ServiceRegistry) (http.HandlerFunc) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("[HTTP update] Got %v request for %v\n", r.Method, r.URL.Path)

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
			svc.AddServiceInstance(ServiceRegistry.NewVersion(req.Version), ServiceRegistry.NewLocation(req.Location))
			w.WriteHeader(201)
			fmt.Fprintf(w, "Added location %v to service %v version %v\n", req.Location, svc.Name, req.Version)

		case `PUT`:
			pathParts := strings.Split(r.URL.Path, "/")
			if len(pathParts) != 3 {
				http.Error(w, `Should only PUT to /services/<service-name>`, 400)
				return
			}
			// Currently don't care about the body, but we may do later

			if svc, created := sr.GetServiceWithName(pathParts[2], true); created {
				w.WriteHeader(201)
				fmt.Fprintf(w, "Created resource %v\n", svc.Name)
			} else {
				w.WriteHeader(200)
				fmt.Fprintf(w, `Resource %v already exists`, svc.Name)
			}

		case `GET`:
			pathParts := strings.Split(r.URL.Path, "/")

			switch len(pathParts) {

			case 4:
				// /services/TestService/1.2.4
				svc, _ := sr.GetServiceWithName(pathParts[2], false)
				if svc == nil {
					http.Error(w, fmt.Sprintf(`No service found with name %v`, pathParts[2]), 400)
					return
				}
				svs := svc.GetLocationsForVersionString(pathParts[3])
				if svs == nil {
					http.Error(w, fmt.Sprintf(`Found service with name %v but could not find an instance with version %v`, pathParts[2], pathParts[3]), 400)
					return
				}

				type EncodedServiceInstance struct {
					Name    string
					Version string
					URI     *[]string
				}

				thingToEncode := EncodedServiceInstance{svc.Name, pathParts[3], new([]string)}
				enc := json.NewEncoder(w)
				for _, value := range svs {
					*thingToEncode.URI = append(*thingToEncode.URI, value.Location)
				}
				enc.Encode(thingToEncode)

			case 3:
				// /services/TestService or /services/

				if pathParts[2] == "" {
					// /services/ : List the services available

					type EncodedService struct {
						Name       string
						URI        string
					}

					enc := json.NewEncoder(w)
					var thingToEncode []EncodedService
					for _, value := range sr.Services {
						thingToEncode = append(thingToEncode, EncodedService{value.Name, fmt.Sprintf("/services/%v", value.Name)})
					}
					enc.Encode(thingToEncode)
					//					enc.Encode(sr.Services)
				}

				if pathParts[2] != "" {
					// /services/TestService : List the versions of the service
					svc, _ := sr.GetServiceWithName(pathParts[2], false)
					if svc == nil {
						http.Error(w, fmt.Sprintf(`No service found with name %v`, pathParts[2]), 400)
						return
					}

					type EncodedService struct {
						Version       string
						URI           string
					}

					var thingToEncode []EncodedService
					enc := json.NewEncoder(w)
					for _, value := range svc.Versions {
						thingToEncode = append(thingToEncode, EncodedService{string(value.Version), fmt.Sprintf("/services/%v/%v", svc.Name, value.Version)})
					}
					enc.Encode(thingToEncode)

				}

			}

		}
	}
}

// Runs the actual HTTP server, ie spawns the goroutine via http.ListenAndServe
func RunHTTP(sr *ServiceRegistry.ServiceRegistry, listenAddr string, finished chan bool, requestUpdate chan bool) {
	sm := http.NewServeMux()
	sm.HandleFunc(`/services/`, getServiceHandler(sr))
	log.Printf("[HTTP update] Starting HTTP server for updates on %v\n", listenAddr)
	if e := http.ListenAndServe(listenAddr, sm); e != nil {
		panic(e)
	}
	log.Println("[HTTP update] *** FINISHED ***")
	finished<-true
}

// Package to update entries in the Service Registry
package update

import (
	`ServiceRegistry`
	`net/http`
	`fmt`
	`encoding/json`
	`log`
	"github.com/gorilla/mux"
)

type Updater struct {
}

func NewUpdater() (u *Updater) {
	u = new(Updater)
	return
}

func makeHandler(fn func (http.ResponseWriter, *http.Request, *ServiceRegistry.ServiceRegistry, map[string]string), sr *ServiceRegistry.ServiceRegistry) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		fn(w, r, sr, vars)
	}
}

func getHandleAddServiceInstance(w http.ResponseWriter, r *http.Request, sr *ServiceRegistry.ServiceRegistry, vars map[string]string) {
	svc, _ := sr.GetServiceWithName(vars["service"], false)
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
}

func getHandleCreateService(w http.ResponseWriter, r *http.Request, sr *ServiceRegistry.ServiceRegistry, vars map[string]string) {
	if svc, created := sr.GetServiceWithName(vars["service"], true); created {
		w.WriteHeader(201)
		fmt.Fprintf(w, "Created resource %v\n", svc.Name)
	} else {
		w.WriteHeader(200)
		fmt.Fprintf(w, `Resource %v already exists`, svc.Name)
	}
}
func getHandleGetServiceInstance(w http.ResponseWriter, r *http.Request, sr *ServiceRegistry.ServiceRegistry, vars map[string]string) {
	svc, _ := sr.GetServiceWithName(vars["service"], false)
	if svc == nil {
		http.Error(w, fmt.Sprintf(`No service found with name %v`, vars["service"]), 400)
		return
	}
	svs := svc.GetLocationsForVersionString(vars["version"])
	if svs == nil {
		http.Error(w, fmt.Sprintf(`Found service with name %v but could not find an instance with version %v`, vars["service"], vars["version"]), 400)
		return
	}
	type EncodedServiceInstance struct {
		Name    string
		Version string
		URI     *[]string
	}
	thingToEncode := EncodedServiceInstance{svc.Name, vars["version"], new([]string)}
	enc := json.NewEncoder(w)
	for _, value := range svs {
		*thingToEncode.URI = append(*thingToEncode.URI, value.Location)
	}
	enc.Encode(thingToEncode)
}

func getHandleGetServices(w http.ResponseWriter, r *http.Request, sr *ServiceRegistry.ServiceRegistry, vars map[string]string) {
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
}

func getHandleGetVersions(w http.ResponseWriter, r *http.Request, sr *ServiceRegistry.ServiceRegistry, vars map[string]string) {
	svc, _ := sr.GetServiceWithName(vars["service"], false)
	if svc == nil {
		http.Error(w, fmt.Sprintf(`No service found with name %v`, vars["service"]), 400)
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

// Runs the actual HTTP server, ie spawns the goroutine via http.ListenAndServe
func (u *Updater) RunHTTP(listenAddr string, sr *ServiceRegistry.ServiceRegistry) (finished chan bool) {
	finished = make(chan bool, 10)
	go u.doRunHTTP(sr, listenAddr, finished)
	return
}

func (u *Updater) doRunHTTP(sr *ServiceRegistry.ServiceRegistry, listenAddr string, finished chan bool) {
	sm := http.NewServeMux()
	r := mux.NewRouter()

	r.HandleFunc("/services/", makeHandler(getHandleGetServices, sr)).Methods("GET")
	r.HandleFunc("/services/{service}", makeHandler(getHandleAddServiceInstance, sr)).Methods("POST")
	r.HandleFunc("/services/{service}", makeHandler(getHandleCreateService, sr)).Methods("PUT")
	r.HandleFunc("/services/{service}", makeHandler(getHandleGetVersions, sr)).Methods("GET")
	r.HandleFunc("/services/{service}/{version}", makeHandler(getHandleGetServiceInstance, sr)).Methods("GET")

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, "You probably wanted /services/")
		})

	sm.Handle("/", r)

	log.Printf("[HTTP update] Starting HTTP server for updates on %v\n", listenAddr)
	if e := http.ListenAndServe(listenAddr, sm); e != nil {
		panic(e)
	}
	log.Println("[HTTP update] *** FINISHED ***")
	finished<-true
}

package main

import (
	"net/url"
	"ServiceRegistry"
	"Interfaces/http"
	"fmt"
	"math/rand"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	sr := ServiceRegistry.ServiceRegistry{}

	a := ServiceRegistry.NewService("TestService")
	b := ServiceRegistry.NewService("SomeOtherService")
	sr.Services = append(sr.Services, a)
	sr.Services = append(sr.Services, b)

	//	Create the data structures

	var u *url.URL

	u, _ = url.Parse("http://10.1.0.1:1237/blah")
	b.AddServiceInstance(ServiceRegistry.Version{1, 0, 0}, ServiceRegistry.NewLocation(u))

	u, _ = url.Parse("http://10.1.0.1:1238/blah")
	b.AddServiceInstance(ServiceRegistry.Version{1, 0, 0}, ServiceRegistry.NewLocation(u))

	u, _ = url.Parse("http://10.1.0.1:1239/blah")
	b.AddServiceInstance(ServiceRegistry.Version{1, 0, 0}, ServiceRegistry.NewLocation(u))

	u, _ = url.Parse("http://localhost:123/blah")
	a.AddServiceInstance(ServiceRegistry.Version{1, 24, 37}, ServiceRegistry.NewLocation(u))

	u, _ = url.Parse("https://kevin.valinor.local/blahz")
	a.AddServiceInstance(ServiceRegistry.Version{1, 24, 37}, ServiceRegistry.NewLocation(u))

	u, _ = url.Parse("https://kevin.valinor.local/blahqwertyz")
	a.AddServiceInstance(ServiceRegistry.Version{2, 24, 37}, ServiceRegistry.NewLocation(u))

	u, _ = url.Parse("https://kevin.valinor.local/sfvslfvsf")
	a.AddServiceInstance(ServiceRegistry.Version{2, 24, 37}, ServiceRegistry.NewLocation(u))

	fmt.Printf(a.Name + ": %d versions\n", len(a.Versions))
	fmt.Printf(b.Name + ": %d versions\n", len(b.Versions))

	// Interrogate the data structures
	fmt.Println(sr.GetServiceWithName("TestService").GetLocationsForVersion(ServiceRegistry.Version{1, 24, 37})[0].Location)

	// Interrogate the data structures
	fmt.Println(sr.GetServiceWithName("TestService").GetLocationForVersion(ServiceRegistry.Version{1, 24, 37}))

	http.RunHTTP(&sr)

}


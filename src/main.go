package main

import (
	"net/url"
	"ServiceRegistry"
	"Interfaces/http"
	"Interfaces/update"
	"fmt"
	"math/rand"
	"time"
)

func main() {

	// Seed the RNG.  This isn't cryptography, it doesn't matter if the RNG is predictable.
	rand.Seed(time.Now().UnixNano())

	sr := ServiceRegistry.ServiceRegistry{}

	a := ServiceRegistry.NewService("TestService")
	b := ServiceRegistry.NewService("SomeOtherService")
	sr.Services = append(sr.Services, a)
	sr.Services = append(sr.Services, b)

	//	Create the data structures

	var u *url.URL

	u, _ = url.Parse("http://10.1.0.1:1237/blah")
	b.AddServiceInstance(ServiceRegistry.NewVersion("1.0.0"), ServiceRegistry.NewLocation(u))

	u, _ = url.Parse("http://10.1.0.1:1238/blah")
	b.AddServiceInstance(ServiceRegistry.NewVersion("1.0.0"), ServiceRegistry.NewLocation(u))

	u, _ = url.Parse("http://10.1.0.1:1239/blah")
	b.AddServiceInstance(ServiceRegistry.NewVersion("1.0.0"), ServiceRegistry.NewLocation(u))

	u, _ = url.Parse("http://localhost:123/blah")
	a.AddServiceInstance(ServiceRegistry.NewVersion("1.24.37"), ServiceRegistry.NewLocation(u))

	u, _ = url.Parse("https://kevin.valinor.local/blahz")
	a.AddServiceInstance(ServiceRegistry.NewVersion("1.24.37"), ServiceRegistry.NewLocation(u))

	u, _ = url.Parse("https://kevin.valinor.local/blahqwertyz")
	a.AddServiceInstance(ServiceRegistry.NewVersion("2.24.37"), ServiceRegistry.NewLocation(u))

	u, _ = url.Parse("https://kevin.valinor.local/sfvslfvsf")
	a.AddServiceInstance(ServiceRegistry.NewVersion("2.24.37"), ServiceRegistry.NewLocation(u))

	fmt.Printf(a.Name + ": %d versions\n", len(a.Versions))
	fmt.Printf(b.Name + ": %d versions\n", len(b.Versions))

	// Interrogate the data structures
	fmt.Println(sr.GetServiceWithName("TestService", false).GetLocationsForVersion(ServiceRegistry.NewVersion("1.24.37"))[0].Location)

	// Interrogate the data structures
	fmt.Println(sr.GetServiceWithName("TestService", false).GetLocationForVersion(ServiceRegistry.Version("1.24.37")))

	http.RunHTTP(&sr)
	update.RunHTTP(&sr)

}


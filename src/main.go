package main

import "fmt"
import "ServiceRegistry"
import "net/url"

func main() {
	fmt.Println("Hello")

	a := ServiceRegistry.NewService("TestService")
	b := ServiceRegistry.NewService("SomeOtherService")

	var u *url.URL
	var v ServiceRegistry.Version
	var loc ServiceRegistry.ServiceLocation

	u, _ = url.Parse("http://localhost:123/blah")
	v = ServiceRegistry.NewVersion(1, 24, 37)
	loc = ServiceRegistry.NewLocation(u)
	a.AddServiceInstance(v, loc)

	v = ServiceRegistry.NewVersion(1, 0, 0)

	u, _ = url.Parse("http://10.1.0.1:1237/blah")
	loc = ServiceRegistry.NewLocation(u)
	b.AddServiceInstance(v, loc)

	u, _ = url.Parse("http://10.1.0.1:1238/blah")
	loc = ServiceRegistry.NewLocation(u)
	b.AddServiceInstance(v, loc)

	u, _ = url.Parse("http://10.1.0.1:1239/blah")
	loc = ServiceRegistry.NewLocation(u)
	b.AddServiceInstance(v, loc)

	u, _ = url.Parse("https://kevin.valinor.local/blahz")
	v = ServiceRegistry.NewVersion(1, 24, 37)
	loc = ServiceRegistry.NewLocation(u)
	a.AddServiceInstance(v, loc)

	fmt.Printf(a.Name + ": %d versions\n", len(a.Versions))
	fmt.Printf(b.Name + ": %d versions\n", len(b.Versions))
}

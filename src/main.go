package main

import "fmt"
import "ServiceRegistry"
import "net/url"

func main() {
	fmt.Println("Hello")

	a := ServiceRegistry.NewService("TestService")
	b := ServiceRegistry.NewService("SomeOtherService")

	u, _ := url.Parse("http://localhost:123/blah")
	version := ServiceRegistry.NewVersion(1, 24, 37)
	loc := ServiceRegistry.NewLocation(u)
	a.AddServiceInstance(version, loc)

	fmt.Printf(a.Name + ": %d versions\n", len(a.Versions))
	fmt.Printf(b.Name + "\n")
}

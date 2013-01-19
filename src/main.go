package main

import "fmt"
import "ServiceRegistry"

func main() {
	fmt.Println("Hello")
	a := new(ServiceRegistry.Service)
	fmt.Println(a.Versions)
}

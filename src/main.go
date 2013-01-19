/**
 * Created with IntelliJ IDEA.
 * User: kevin
 * Date: 16/01/13
 * Time: 00:21
 * To change this template use File | Settings | File Templates.
 */

package main

import "fmt"
import "ServiceRegistry"

func main() {
	fmt.Println("Hello")
	a := new(ServiceRegistry.Service)
	fmt.Println(a.Versions)
}

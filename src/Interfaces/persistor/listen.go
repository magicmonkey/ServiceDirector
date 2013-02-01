package persistor

import (
	"ServiceRegistry"
	"fmt"
)

type Persistor struct {

}

func (p *Persistor) Listen(sr *ServiceRegistry.ServiceRegistry, okToContinue chan bool) {
	fmt.Println("Persistor: Listening for updates...")
	c := sr.MakeUpdateChannel()
	okToContinue <- true
	for {
		update := <-c
		fmt.Println("Got an update!", update.Message)
	}
}

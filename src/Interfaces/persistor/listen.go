package persistor

import (
	"ServiceRegistry"
	"fmt"
	"github.com/vmihailenco/redis"
	"sync"
	"encoding/gob"
	"bytes"
	//	"os"
)

type Persistor struct {
	serviceRegistryUpdateChannel *chan ServiceRegistry.ServiceRegistryUpdate
	serviceUpdateChannel *chan ServiceRegistry.ServiceUpdate
	redis *redis.Client
	connected bool
}

func NewPersistor(sruc *chan ServiceRegistry.ServiceRegistryUpdate, suc *chan ServiceRegistry.ServiceUpdate) (*Persistor) {
	p := Persistor{sruc, suc, nil, false}
	return &p
}

func (p *Persistor) getRedis() (*redis.Client) {
	// Guard in case someone else is opening Redis
	m := sync.Mutex{}
	m.Lock()
	if !p.connected {
		fmt.Println("Opening Redis...")
		p.redis = redis.NewTCPClient("localhost:6379", "", 0)
		p.connected = true
	}
	m.Unlock()
	return p.redis
}

func (p *Persistor) Listen() {
	fmt.Println("Persistor: Listening for updates...")
	for {
		select {
		case msg1 := <-*p.serviceUpdateChannel:
			p.saveService(msg1.Service)
		case msg2 := <-*p.serviceRegistryUpdateChannel:
			p.saveServiceRegistry(msg2.ServiceRegistry)

		}
	}

}

func (p *Persistor) saveService(s *ServiceRegistry.Service) {

	fmt.Println("Saving a service")

	c := p.getRedis()
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	enc.Encode(s)
	err := c.Set(fmt.Sprintf("service-%v", s.Name), buf.String())
	fmt.Println(err.Val())
}

func (p *Persistor) saveServiceRegistry(sr *ServiceRegistry.ServiceRegistry) {
	fmt.Println("Saving a service registry")

	c := p.getRedis()
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	enc.Encode(sr)
	err := c.Set(fmt.Sprintf("serviceregistry-%v", sr.Name), buf.String())
	fmt.Println(err.Val())
}

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
	serviceRegistryUpdateChannel *chan *ServiceRegistry.ServiceRegistry
	redis *redis.Client
	connected bool
}

func NewPersistor(sruc *chan *ServiceRegistry.ServiceRegistry) (*Persistor) {
	p := Persistor{sruc, nil, false}
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
		msg1 := <-*p.serviceRegistryUpdateChannel
		p.saveServiceRegistry(msg1)
	}

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

func (p *Persistor) LoadServiceRegistry(name string) (*ServiceRegistry.ServiceRegistry) {
	c := p.getRedis()
	srBytes := c.Get(fmt.Sprintf("serviceregistry-%v", name))
	buf := bytes.NewBuffer([]byte(srBytes.Val()))
	dec := gob.NewDecoder(buf)
	sr := new(ServiceRegistry.ServiceRegistry)
	dec.Decode(&sr)
	return sr
}

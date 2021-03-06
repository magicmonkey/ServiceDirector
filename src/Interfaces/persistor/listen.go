package persistor

import (
	"ServiceRegistry"
	"fmt"
	"github.com/vmihailenco/redis"
	"sync"
	"encoding/gob"
	"bytes"
	"log"
	//	"os"
)

type Persistor struct {
	redis *redis.Client
	connected bool
}

func NewPersistor() (*Persistor) {
	p := Persistor{nil, false}
	return &p
}

func (p *Persistor) getRedis() (*redis.Client) {
	// Guard in case someone else is opening Redis
	redisAddr := "localhost:6379"
	m := sync.Mutex{}
	m.Lock()
	if !p.connected {
		log.Printf("[Persistor] Opening Redis at [%v]\n", redisAddr)
		p.redis = redis.NewTCPClient(redisAddr, "", 0)
		p.connected = true
	}
	m.Unlock()
	return p.redis
}

func (p *Persistor) Listen() (updateChannel chan ServiceRegistry.ServiceRegistry) {
	updateChannel = make(chan ServiceRegistry.ServiceRegistry, 10)
	go p.doListen(updateChannel)
	return
}

func (p *Persistor) doListen(updateChannel chan ServiceRegistry.ServiceRegistry) {
	log.Println("[Persistor] Listening for updates...")
	for {
		sr := <-updateChannel
		log.Printf("[Persistor] Saving updated service registry [%v]\n", sr.Name)
		p.saveServiceRegistry(&sr)
	}
}

func (p *Persistor) saveServiceRegistry(sr *ServiceRegistry.ServiceRegistry) {
	c := p.getRedis()
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	enc.Encode(sr)
	c.Set(fmt.Sprintf("serviceregistry-%v", sr.Name), buf.String())
	buf.Reset()
}

func (p *Persistor) LoadServiceRegistry(name string) (*ServiceRegistry.ServiceRegistry) {
	log.Printf("[Persistor] Loading a service registry called [%v]\n", name)
	c := p.getRedis()
	srBytes := c.Get(fmt.Sprintf("serviceregistry-%v", name))
	buf := bytes.NewBuffer([]byte(srBytes.Val()))
	dec := gob.NewDecoder(buf)
	sr := new(ServiceRegistry.ServiceRegistry)
	dec.Decode(sr)
	sr.Name = name

	// Reconnect the up-tree references
	for _, value := range sr.Services {
		value.SetServiceRegistry(sr)
	}

	return sr
}

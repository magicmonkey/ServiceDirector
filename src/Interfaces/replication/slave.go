package replication

import (
	"ServiceRegistry"
	"net"
	"fmt"
	"log"
	"encoding/json"
)

func StartSlave(masterAddr string, sru1 chan *ServiceRegistry.ServiceRegistry) {
	log.Println("[Replication slave] Starting replication slave")
	conn, err := net.Dial("tcp", masterAddr)
	if err != nil {
		// handle error
	}
	var sr *ServiceRegistry.ServiceRegistry
	fmt.Fprintf(conn, "Hello!\n")
	dec := json.NewDecoder(conn)

	for {
		dec.Decode(&sr)
		sru1 <- sr
	}
}

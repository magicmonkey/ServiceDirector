package replication

import (
	"net"
	"log"
	"bytes"
	"fmt"
	"io"
	"ServiceRegistry"
	"encoding/gob"
)

type listener struct {
	connections []net.Conn
}

// Starts the replication listener, and sends any channel updates down the network connection
func StartListener(sruc chan *ServiceRegistry.ServiceRegistry) {
	l := new(listener)
	ln, err := net.Listen("tcp", ":8083")
	if err != nil {
		panic(err)
	}
	go l.listenForUpdates(sruc)
	for {
		conn, err := ln.Accept()
		if err != nil {
			panic(err)
			continue
		}
		l.handleConnection(conn)
		go l.readFromConnection(conn)
	}

}

func (l *listener) listenForUpdates(sruc chan *ServiceRegistry.ServiceRegistry) {
	var buf bytes.Buffer
	log.Println("[Replication master] Listening for updates...")
	for {
		msg1 := <-sruc
		log.Println("[Replication master] Got an update")
		enc := gob.NewEncoder(&buf)
		enc.Encode(msg1)
		l.Write(buf.String())
	}
}

func (l *listener) handleConnection(conn net.Conn) {
	log.Println("Got a connection")
	l.connections = append(l.connections, conn)
}

// Writes a message to each connected client
func (l *listener) Write(s string) {

	for _, value := range l.connections {
		log.Println("[Replication master] Replication master: sending update to", value.RemoteAddr())
		value.Write([]byte(s))
	}
}

func (l *listener) readFromConnection(conn net.Conn) {
	tbuf := make([]byte, 81920)

	for {
		n, err := conn.Read(tbuf)
		// Was there an error in reading ?
		if err != nil {
			if err != io.EOF {
				log.Printf("[Replication master] Read error: %s", err)
			}
			break
		}
		fmt.Printf(string(tbuf[0:n]))
	}
	l.removeConnection(conn)
}

func (l *listener) removeConnection(conn net.Conn) {
	log.Println("[Replication master] Connection was closed")
}

// Closes all connected clients
func (l *listener) Close() {
	for _, value := range l.connections {
		value.Close()
	}
}

package replication

import (
	"net"
	"log"
	"bytes"
	"fmt"
	"io"
	"ServiceRegistry"
	"encoding/json"
)

type listener struct {
	connections map[net.Addr]net.Conn
}

// Starts the replication listener, and sends any channel updates down the network connection
func StartListener(sruc chan *ServiceRegistry.ServiceRegistry, sr *ServiceRegistry.ServiceRegistry) {
	l := new(listener)
	l.connections = make(map[net.Addr]net.Conn)
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
		l.handleConnection(conn, sr)
		go l.readFromConnection(conn)
	}

}

func (l *listener) listenForUpdates(sruc chan *ServiceRegistry.ServiceRegistry) {
	var buf bytes.Buffer
	log.Println("[Replication master] Listening for updates...")
	for {
		msg1 := <-sruc
		log.Println("[Replication master] Got an update")
		enc := json.NewEncoder(&buf)
		enc.Encode(msg1)
		l.Write(buf.String())
		buf.Reset()
	}
}

func (l *listener) handleConnection(conn net.Conn, sr *ServiceRegistry.ServiceRegistry) {
	l.connections[conn.RemoteAddr()] = conn
	log.Printf("Got a connection; there are now %d connections\n", len(l.connections))

	// Send initial structure down the wire
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.Encode(sr)
	conn.Write(buf.Bytes())
	buf.Reset()
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
	delete(l.connections, conn.RemoteAddr())
	log.Printf("[Replication master] Connection was closed; there are now %d connections", len(l.connections))
}

//// Closes all connected clients
//func (l *listener) Close() {
//	for _, value := range l.connections {
//		value.Close()
//		delete(l.connections, value.RemoteAddr())
//	}
//}

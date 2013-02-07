package replication

import (
	"net"
	"log"
//	"bytes"
	"fmt"
	"io"
)

type listener struct {
	connections []net.Conn
}

func StartListener() {
	l := new(listener)
	ln, err := net.Listen("tcp", ":8083")
	if err != nil {
		panic(err)
	}
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

func (l *listener) handleConnection(conn net.Conn) {
	log.Println("Got a connection")
	l.connections = append(l.connections, conn)
}

func (l *listener) Write (s string) {
	for _, value := range l.connections {
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
				log.Printf("Read error: %s", err)
			}
			break
		}
		fmt.Println(string(tbuf[0:n]))
	}
	l.removeConnection(conn)
}

func (l *listener) removeConnection(conn net.Conn) {
	log.Println("Connection was closed")
}

func (l *listener) Close() {
	for _,value := range l.connections {
		value.Close()
	}
}

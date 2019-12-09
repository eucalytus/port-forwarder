package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
)

var (
	localPort    int
	remoteServer string
	remotePort   int
)

func main() {
	flag.StringVar(&remoteServer, "remote-server-addr", "0.0.0.0", "remote server address, default is localhost ")
	flag.IntVar(&remotePort, "remote-server-port", 0, "remote server port")
	flag.IntVar(&localPort, "listen", 0, "listen port")
	flag.Parse()
	if len(remoteServer) == 0 || remotePort <= 0 || localPort <= 0 {
		log.Fatal("invalid configuration")
	}

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", localPort))
	if err != nil {
		log.Panic(err)
	}

	for {
		client, err := l.Accept()
		if err != nil {
			log.Panic(err)
		}

		go handleClientRequest(client)
	}
}

func handleClientRequest(client net.Conn) {
	log.Println("handle new connection:", client.RemoteAddr())
	if client == nil {
		return
	}
	defer client.Close()

	server, err := net.Dial("tcp", net.JoinHostPort(remoteServer, fmt.Sprintf("%d", remotePort)))
	if err != nil {
		log.Println(err)
		return
	}
	defer server.Close()
	go io.Copy(server, client)
	io.Copy(client, server)

	log.Println("handle new connection done: ", client.RemoteAddr())
}

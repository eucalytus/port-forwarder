package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
)

var (
	localPort   int
	forwardAddr string
)

func main() {
	flag.StringVar(&forwardAddr, "forward-addr", "", "forward address")
	flag.IntVar(&localPort, "listen", 0, "listen port")
	flag.Parse()
	if len(forwardAddr) == 0 || localPort <= 0 {
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

	server, err := net.Dial("tcp", forwardAddr)
	if err != nil {
		log.Println(err)
		return
	}
	defer server.Close()
	go io.Copy(server, client)
	io.Copy(client, server)

	log.Println("handle new connection done: ", client.RemoteAddr())
}

package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
	"strings"

	"rmda-go/rsocket"
)

var (
	host = flag.String("s", "192.168.96.10", "server address")
	port = flag.Int("p", 8000, "server port")
)

func main() {
	flag.Parse()

	localAddr := fmt.Sprintf("%s:%d", *host, *port)
	ln, err := rsocket.NewListener(localAddr, 128)
	if err != nil {
		log.Fatal(err)
	}

	defer func(ln *rsocket.Listener) {
		if err := ln.Close(); err != nil {
			log.Fatal(err)
		}
	}(ln)

	log.Printf("RDMA server is listening on %s\n", localAddr)
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Printf("failed accepting: %v\n", err)
			return
		}

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer func(conn net.Conn) {
		if err := conn.Close(); err != nil {
			log.Fatal(err)
		}
	}(conn)

	reader := bufio.NewReader(conn)
	for {
		msg, err := reader.ReadString('\n')
		if err != nil {
			log.Println("Client disconnected:", err)
			break
		}

		if msg == "exit" {
			fmt.Println("exited")
			break
		}

		content := strings.TrimSpace(msg)
		log.Printf("Received message from %s: %s\n", conn.RemoteAddr().String(), content)

		response := fmt.Sprintf("received '%s', pong!", content)
		if _, err := conn.Write([]byte(response)); err != nil {
			log.Printf("failed sending response: %v\n", err)
			break
		}
	}
}

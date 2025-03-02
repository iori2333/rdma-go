package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"

	"rmda-go/rsocket"
)

var (
	host  = flag.String("s", "192.168.96.10", "server address")
	local = flag.String("l", "192.168.96.20", "local address")
	port  = flag.Int("p", 8000, "server port")
)

func main() {
	flag.Parse()

	remoteAddr := fmt.Sprintf("%s:%d", *host, *port)
	localAddr := fmt.Sprintf("%s:0", *local)

	conn, err := rsocket.Dial(remoteAddr, rsocket.WithLocal(localAddr))
	if err != nil {
		log.Fatal(err)
	}

	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(conn)
	for {
		fmt.Print("Enter message: ")
		text, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Failed to read input:", err)
			break
		}

		if _, err := writer.WriteString(text); err != nil {
			fmt.Println("Failed to send message:", err)
			break
		}

		if err := writer.Flush(); err != nil {
			fmt.Println("Failed to flush writer:", err)
			break
		}

		if text == "exit\n" {
			break
		}

		response := make([]byte, 1024)
		n, err := conn.Read(response)
		if err != nil {
			fmt.Println("Failed to read response from server:", err)
			break
		}

		fmt.Printf("Server response: %s\n", string(response[:n]))
	}
}

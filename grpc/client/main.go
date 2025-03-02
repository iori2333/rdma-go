package main

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net"
	proto2 "rmda-go/grpc/proto"
	"rmda-go/rsocket"
	"time"
)

func main() {
	insecureOpt := grpc.WithTransportCredentials(insecure.NewCredentials())
	rsocketOpt := grpc.WithContextDialer(func(ctx context.Context, s string) (net.Conn, error) {
		return rsocket.Dial(s, rsocket.WithLocal("192.168.96.20:0"))
	})

	conn, err := grpc.NewClient("192.168.96.10:8000", insecureOpt, rsocketOpt)
	if err != nil {
		log.Fatal("Dial: ", err)
	}
	defer func(conn *grpc.ClientConn) {
		if err := conn.Close(); err != nil {
			panic(err)
		}
	}(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	client := proto2.NewGreeterClient(conn)
	resp, err := client.SayHello(ctx, &proto2.HelloRequest{Name: "Client"})
	if err != nil {
		log.Fatal("Response: ", err)
	}

	log.Printf("Greeting: %s\n", resp.Message)
}

package main

import (
	"context"
	"google.golang.org/grpc"
	"log"
	proto2 "rdma-go/grpc/proto"
	"rdma-go/rsocket"
)

type serverImpl struct {
	proto2.UnimplementedGreeterServer
}

func (s *serverImpl) SayHello(ctx context.Context, in *proto2.HelloRequest) (*proto2.HelloReply, error) {
	log.Printf("invoking SayHello with param: %s\n", in.Name)
	return &proto2.HelloReply{Message: "Hello " + in.Name}, nil
}

func main() {
	server := grpc.NewServer()
	defer server.Stop()

	ln, err := rsocket.NewListener("192.168.96.10:8000", 128)
	if err != nil {
		log.Fatal("create listener: ", err)
	}

	proto2.RegisterGreeterServer(server, &serverImpl{})
	if err := server.Serve(ln); err != nil {
		log.Fatal("serve: ", err)
	}
}

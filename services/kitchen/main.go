package main

import (
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewGRPCClient(addr string) *grpc.ClientConn {
	conn, err := grpc.NewClient(
		addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	log.Println(conn.GetState().String())
	return conn
}

func main() {

	conn := NewGRPCClient(":9000")
	defer conn.Close()

	httpServer := NewHttpServer(":3000", conn)
	log.Fatal(httpServer.Run())
}

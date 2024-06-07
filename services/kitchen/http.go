package main

import (
	kitchen "github.com/shammianand/oms/services/kitchen/handler"
	"google.golang.org/grpc"
	"log"
	"net/http"
)

type httpServer struct {
	addr string
	conn *grpc.ClientConn
}

func NewHttpServer(addr string, conn *grpc.ClientConn) *httpServer {
	return &httpServer{addr: addr, conn: conn}
}

func (s *httpServer) Run() error {
	router := http.NewServeMux()
	kitchenHandler := kitchen.NewHandler(s.conn)
	kitchenHandler.RegisterRoutes(router)

	log.Println("Starting server on", s.addr)
	return http.ListenAndServe(s.addr, router)
}

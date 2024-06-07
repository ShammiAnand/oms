package main

import (
	"log"
	"net/http"

	"github.com/shammianand/oms/services/common/util"
	kitchen "github.com/shammianand/oms/services/kitchen/handler"
	"google.golang.org/grpc"
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
	userStore := kitchen.NewStore()
	kitchenHandler := kitchen.NewHandler(s.conn, userStore)
	kitchenHandler.RegisterRoutes(router)

	log.Println("Starting http server on", s.addr)
	return http.ListenAndServe(s.addr, util.Logging(router))
}

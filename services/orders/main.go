package main

func main() {

	grpcServer := NewGRPCServer(":9000")
	grpcServer.Run()

	// NOTE: order service does not need to have a http transport
	// httpServer := NewHttpServer(":8000")
	// httpServer.Run()
}

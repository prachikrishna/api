package main

import (
	"log"
	"net"

	mysvccore "prachi/core"
	mysvcgrpc "prachi/grpc"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {

	// configure our core service
	wordService := mysvccore.NewService()

	// configure our gRPC service controller
	wordServiceController := NewWordServiceController(wordService)

	// start a gRPC server
	server := grpc.NewServer()
	mysvcgrpc.RegisterWordServiceServer(server, wordServiceController)
	reflection.Register(server)

	con, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		panic(err)
	}

	log.Printf("Starting gRPC user service on %s...\n", con.Addr().String())
	err = server.Serve(con)
	if err != nil {
		panic(err)
	}
}

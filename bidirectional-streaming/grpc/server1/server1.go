package server1

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"

	"google.golang.org/grpc"

	proto "bidirectional-streaming/grpc"
)

// RunServer runs gRPC service to publish Book service
func RunServer(ctx context.Context, API proto.SearchBookServiceServer, port string) error {
	listen, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return err
	}

	// register service
	server := grpc.NewServer()
	proto.RegisterSearchBookServiceServer(server, API)

	// graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			log.Println("shutting down gRPC server...")

			server.GracefulStop()

			<-ctx.Done()
		}
	}()

	// start gRPC server
	log.Println("starting gRPC server...")
	return server.Serve(listen)
}

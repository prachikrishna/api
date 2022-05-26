package main

import (
	"context"
	"fmt"
	"io"
	"log"

	//"net"
	"flag"

	"google.golang.org/grpc/credentials/insecure"

	proto "grpc-bidirectional-bookdetails/grpc"

	"google.golang.org/grpc"
)

func main() {
	// get configuration
	address := flag.String("server", "", "gRPC server in format host:port")
	flag.Parse()

	// Set up a connection to the server.
	conn, err := grpc.Dial(*address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	client := proto.NewBookServiceClient(conn)

	//ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	//defer cancel()

	reqs := []*proto.Request{
		{Id: 1},
		{Id: 2},
	}

	stream, err := client.GetDetails(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	// goroutine to stream outbound requests to server
	go func() {
		for _, req := range reqs {
			if err := stream.Send(req); err != nil {
				log.Fatal(err)
			}
		}
		if err := stream.CloseSend(); err != nil {
			log.Fatal(err)
		}
	}()

	fmt.Println("\nFound BookDetails")
	fmt.Println("-----------------")
	for {
		book, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println(err)
			continue
		}
		fmt.Printf("%s %s\n", book.Book.GetAuthor(), book.Book.GetName())
	}
}

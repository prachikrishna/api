package main

import (
	"context"
	"fmt"
	"io"
	"log"

	//"net"
	"flag"
	"time"

	"google.golang.org/grpc/credentials/insecure"

	pb "grpc-server-streaming/grpc"

	"google.golang.org/grpc"
)

/*const (
	server     = "127.0.0.1"
	serverPort = "9090"
)*/

func printPageItems(client pb.PageInfoServiceClient) {
	curReq := &pb.PageRequest{PageNo: 1, PageSize: "a4"}
	stream, err := client.GetDetails(context.Background(), curReq)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("\nPage Items")
	fmt.Println("-------------")
	for {
		cur, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println(err)
			continue
		}
		fmt.Printf("%s %d %d %d\n", cur.GetTitle(), cur.GetTotalWords(), cur.GetTotalSentences(), cur.GetTotalImages())
	}
}

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

	client := pb.NewPageInfoServiceClient(conn)

	_, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	printPageItems(client)
}

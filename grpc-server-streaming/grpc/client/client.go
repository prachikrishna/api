package main

/*import (
	"context"
	"flag"
	"log"
	"time"

	//"github.com/golang/protobuf/ptypes"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	proto "grpc-crud-unary-bookservice/grpc"
)*/

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

/*func main() {
// get configuration
address := flag.String("server", "", "gRPC server in format host:port")
flag.Parse()

// Set up a connection to the server.
conn, err := grpc.Dial(*address, grpc.WithTransportCredentials(insecure.NewCredentials()))
if err != nil {
	log.Fatalf("did not connect: %v", err)
}
defer conn.Close()

c := proto.NewBookServiceClient(conn)

ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

t := time.Now().In(time.UTC)
//reminder, _ := ptypes.TimestampProto(t)
pfx := t.Format(time.RFC3339Nano)

req2 := proto.ReadRequest{
	Id: id,
}
res2, err := c.ReadBook(ctx, &req2)
if err != nil {
	log.Fatalf("Read failed: %v", err)
}
log.Printf("Read result: <%+v>\n\n", res2)*/

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

	//t := time.Now().In(time.UTC)
	//reminder, _ := ptypes.TimestampProto(t)
	//pfx := t.Format(time.RFC3339Nano)

	/*serverAddr := net.JoinHostPort(server, serverPort)

	// setup insecure connection
	conn, err := grpc.Dial(serverAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}

	//client := pb.NewPageInfoServiceClient(conn)*/

	printPageItems(client)
}

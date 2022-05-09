//PS C:\Users\KIIT\go\src\grpc-crud-unary-bookservice\grpc\client> ./client -server=localhost:9090
//2022/05/05 19:25:02 Create failed: rpc error: code = Unimplemented desc = method CreateBook not implemented
package main

import (
	"context"
	"flag"
	"log"
	"time"

	//"github.com/golang/protobuf/ptypes"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	proto "grpc-crud-unary-bookservice/grpc"
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

	c := proto.NewBookServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	t := time.Now().In(time.UTC)
	//reminder, _ := ptypes.TimestampProto(t)
	pfx := t.Format(time.RFC3339Nano)

	// Call Create
	req1 := proto.CreateRequest{
		Book: &proto.Book{
			Name:   "name (" + pfx + ")",
			Author: "author (" + pfx + ")",
		},
	}
	res1, err := c.CreateBook(ctx, &req1)
	if err != nil {
		log.Fatalf("Create failed: %v", err)
	}
	log.Printf("Create result: <%+v>\n\n", res1)

	id := res1.Id

	// Read
	req2 := proto.ReadRequest{
		Id: id,
	}
	res2, err := c.ReadBook(ctx, &req2)
	if err != nil {
		log.Fatalf("Read failed: %v", err)
	}
	log.Printf("Read result: <%+v>\n\n", res2)

	// Update
	req3 := proto.UpdateRequest{
		Book: &proto.Book{
			Id:     res2.Book.Id,
			Name:   res2.Book.Name,
			Author: res2.Book.Author + " + updated",
		},
	}
	res3, err := c.UpdateBook(ctx, &req3)
	if err != nil {
		log.Fatalf("Update failed: %v", err)
	}
	log.Printf("Update result: <%+v>\n\n", res3)

	// ListAll
	req4 := proto.ListAllRequest{}
	res4, err := c.ListAllBooks(ctx, &req4)
	if err != nil {
		log.Fatalf("ReadAll failed: %v", err)
	}
	log.Printf("ReadAll result: <%+v>\n\n", res4)

	// Delete
	req5 := proto.DeleteRequest{
		Id: id,
	}
	res5, err := c.DeleteBook(ctx, &req5)
	if err != nil {
		log.Fatalf("Delete failed: %v", err)
	}
	log.Printf("Delete result: <%+v>\n\n", res5)
}

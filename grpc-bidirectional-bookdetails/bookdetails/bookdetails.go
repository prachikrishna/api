package bookdetails

import (
	"context"
	"database/sql"
	"fmt"
	"grpc-bidirectional-bookdetails/grpc"
	"io"
	"log"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// bookServiceServer is implementation of grpc.BookServiceServer proto interface
type bookServiceServer struct {
	db *sql.DB
	grpc.UnimplementedBookServiceServer
}

// NewBookServiceServer creates Book service
func NewBookServiceServer(db *sql.DB) grpc.BookServiceServer {
	return &bookServiceServer{db: db}
}

func (s *bookServiceServer) connect(ctx context.Context) (*sql.Conn, error) {
	c, err := s.db.Conn(ctx)
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to connect to database-> "+err.Error())
	}
	return c, nil
}

func (s *bookServiceServer) GetDetails(stream grpc.BookService_GetDetailsServer) error {
	ctx := context.Background()
	c, err := s.connect(ctx)
	if err != nil {
		return err
	}
	defer c.Close()
	for {
		req, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				return nil // we're done
			}
			return err
		}

		rows, err := c.QueryContext(ctx, "SELECT `ID`, `Name`, `Author` FROM books WHERE `ID`=?",
			req.Id)
		if err != nil {
			return status.Error(codes.Unknown, "failed to select from books-> "+err.Error())
		}
		defer rows.Close()

		if !rows.Next() {
			if err := rows.Err(); err != nil {
				return status.Error(codes.Unknown, "failed to retrieve data from books-> "+err.Error())
			}
			return status.Error(codes.NotFound, fmt.Sprintf("Book with ID='%d' is not found",
				req.Id))
		}

		// get book data
		var td grpc.Book
		if err := rows.Scan(&td.Id, &td.Name, &td.Author); err != nil {
			return status.Error(codes.Unknown, "failed to retrieve field values from books row-> "+err.Error())
		}

		if rows.Next() {
			return status.Error(codes.Unknown, fmt.Sprintf("found multiple books rows with ID='%d'",
				req.Id))
		}

		res := stream.Send(&grpc.Response{
			Book: &td,
		})

		if res != nil {
			log.Fatalf("Error when response was sent to the client: %v", res)
		}
	}
}

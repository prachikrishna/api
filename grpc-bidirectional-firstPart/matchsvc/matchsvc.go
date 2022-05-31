package matchsvc

import (
	"context"
	"database/sql"
	"fmt"
	"grpc-bidirectional-firstPart/grpc"
	"io"

	//"log"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// searchBookServiceServer is implementation of grpc.SearchBookServiceServer proto interface
type searchBookServiceServer struct {
	db *sql.DB
	grpc.UnimplementedSearchBookServiceServer
}

// NewBookSearchServiceServer creates SearchBook service
func NewSearchBookServiceServer(db *sql.DB) grpc.SearchBookServiceServer {
	return &searchBookServiceServer{db: db}
}

func (s *searchBookServiceServer) connect(ctx context.Context) (*sql.Conn, error) {
	c, err := s.db.Conn(ctx)
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to connect to database-> "+err.Error())
	}
	return c, nil
}

func (s *searchBookServiceServer) FindMatch(stream grpc.SearchBookService_FindMatchServer) error {
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

		if req.GetName() != "nil" {
			rows, err := c.QueryContext(ctx, "SELECT `ID`, `Name`, `Author` FROM books WHERE `Name`=? LIMIT 3",
				req.GetName())
			if err != nil {
				return status.Error(codes.Unknown, "failed to select from books-> "+err.Error())
			}
			defer rows.Close()
			if !rows.Next() {
				if err := rows.Err(); err != nil {
					return status.Error(codes.Unknown, "failed to retrieve data from books-> "+err.Error())
				}
				return status.Error(codes.NotFound, fmt.Sprintf("Book with Name='%s' is not found",
					req.GetName()))
			}
			// get book data
			var td *grpc.Book
			var books []*grpc.Book

			for rows.Next() {
				if err := rows.Scan(&td.Id, &td.Name, &td.Author); err != nil {
					return status.Error(codes.Unknown, "failed to retrieve field values from books row-> "+err.Error())
				}
				books = append(books, td)
			}

			for _, book := range books {
				res := &grpc.Response{
					Book: book,
				}
				if err := stream.Send(res); err != nil {
					return err
				}

			}
			return nil

		} else if req.GetAuthor() != "nil" {
			rows, err := c.QueryContext(ctx, "SELECT `ID`, `Name`, `Author` FROM books WHERE `Author`=? LIMIT 3",
				req.GetAuthor())
			if err != nil {
				return status.Error(codes.Unknown, "failed to select from books-> "+err.Error())
			}
			defer rows.Close()
			if !rows.Next() {
				if err := rows.Err(); err != nil {
					return status.Error(codes.Unknown, "failed to retrieve data from books-> "+err.Error())
				}
				return status.Error(codes.NotFound, fmt.Sprintf("Books with Author='%s' was not found",
					req.GetAuthor()))
			}
			// get book data
			var td *grpc.Book
			var books []*grpc.Book

			for rows.Next() {
				if err := rows.Scan(&td.Id, &td.Name, &td.Author); err != nil {
					return status.Error(codes.Unknown, "failed to retrieve field values from books row-> "+err.Error())
				}
				books = append(books, td)
			}

			if err := rows.Scan(&td.Id, &td.Name, &td.Author); err != nil {
				return status.Error(codes.Unknown, "failed to retrieve field values from books row-> "+err.Error())
			}

			for _, book := range books {
				res := &grpc.Response{
					Book: book,
				}
				if err := stream.Send(res); err != nil {
					return err
				}

			}
			return nil
		}

	}

}

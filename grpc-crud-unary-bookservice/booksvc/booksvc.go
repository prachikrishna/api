package booksvc

import (
	"context"
	"database/sql"
	"fmt"
	"grpc-crud-unary-bookservice/grpc"

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

// connect returns SQL database connection from the pool
func (s *bookServiceServer) connect(ctx context.Context) (*sql.Conn, error) {
	c, err := s.db.Conn(ctx)
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to connect to database-> "+err.Error())
	}
	return c, nil
}

// Create new book
func (s *bookServiceServer) CreateBook(ctx context.Context, req *grpc.CreateRequest) (*grpc.CreateResponse, error) {

	// get SQL connection from pool
	c, err := s.connect(ctx)
	if err != nil {
		return nil, err
	}
	defer c.Close()

	// insert book entity data
	res, err := c.ExecContext(ctx, "INSERT INTO books(`Name`, `Author`) VALUES(?, ?)",
		req.Book.Name, req.Book.Author)
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to insert into books-> "+err.Error())
	}

	// get ID of created book
	id, err := res.LastInsertId()
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to retrieve id for created book-> "+err.Error())
	}

	return &grpc.CreateResponse{
		Id: id,
	}, nil
}

// Read book
func (s *bookServiceServer) ReadBook(ctx context.Context, req *grpc.ReadRequest) (*grpc.ReadResponse, error) {

	// get SQL connection from pool
	c, err := s.connect(ctx)
	if err != nil {
		return nil, err
	}
	defer c.Close()

	// query book by ID
	rows, err := c.QueryContext(ctx, "SELECT `ID`, `Name`, `Author` FROM books WHERE `ID`=?",
		req.Id)
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to select from books-> "+err.Error())
	}
	defer rows.Close()

	if !rows.Next() {
		if err := rows.Err(); err != nil {
			return nil, status.Error(codes.Unknown, "failed to retrieve data from books-> "+err.Error())
		}
		return nil, status.Error(codes.NotFound, fmt.Sprintf("Book with ID='%d' is not found",
			req.Id))
	}

	// get book data
	var td grpc.Book
	if err := rows.Scan(&td.Id, &td.Name, &td.Author); err != nil {
		return nil, status.Error(codes.Unknown, "failed to retrieve field values from books row-> "+err.Error())
	}

	if rows.Next() {
		return nil, status.Error(codes.Unknown, fmt.Sprintf("found multiple books rows with ID='%d'",
			req.Id))
	}

	return &grpc.ReadResponse{
		Book: &td,
	}, nil

}

// Update book
func (s *bookServiceServer) UpdateBook(ctx context.Context, req *grpc.UpdateRequest) (*grpc.UpdateResponse, error) {

	// get SQL connection from pool
	c, err := s.connect(ctx)
	if err != nil {
		return nil, err
	}
	defer c.Close()

	res, err := c.ExecContext(ctx, "UPDATE books SET `Name`=?, `Author`=? WHERE `ID`=?",
		req.Book.Name, req.Book.Author, req.Book.Id)
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to update books-> "+err.Error())
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to retrieve rows affected value-> "+err.Error())
	}

	if rows == 0 {
		return nil, status.Error(codes.NotFound, fmt.Sprintf("book with ID='%d' is not found",
			req.Book.Id))
	}

	return &grpc.UpdateResponse{
		Updated: rows,
	}, nil
}

// Delete book
func (s *bookServiceServer) DeleteBook(ctx context.Context, req *grpc.DeleteRequest) (*grpc.DeleteResponse, error) {

	// get SQL connection from pool
	c, err := s.connect(ctx)
	if err != nil {
		return nil, err
	}
	defer c.Close()

	res, err := c.ExecContext(ctx, "DELETE FROM books WHERE `ID`=?", req.Id)
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to delete book-> "+err.Error())
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to retrieve rows affected value-> "+err.Error())
	}

	if rows == 0 {
		return nil, status.Error(codes.NotFound, fmt.Sprintf("book with ID='%d' is not found",
			req.Id))
	}

	return &grpc.DeleteResponse{
		Deleted: rows,
	}, nil
}

// Read all books
func (s *bookServiceServer) ListAllBooks(ctx context.Context, req *grpc.ListAllRequest) (*grpc.ListAllResponse, error) {

	// get SQL connection from pool
	c, err := s.connect(ctx)
	if err != nil {
		return nil, err
	}
	defer c.Close()

	// get book list
	rows, err := c.QueryContext(ctx, "SELECT `ID`, `Name`, `Author` FROM books")
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to select from books-> "+err.Error())
	}
	defer rows.Close()
	list := []*grpc.Book{}
	for rows.Next() {
		td := new(grpc.Book)
		if err := rows.Scan(&td.Id, &td.Name, &td.Author); err != nil {
			return nil, status.Error(codes.Unknown, "failed to retrieve field values from books row-> "+err.Error())
		}
		list = append(list, td)
	}

	if err := rows.Err(); err != nil {
		return nil, status.Error(codes.Unknown, "failed to retrieve data from books-> "+err.Error())
	}

	return &grpc.ListAllResponse{
		Books: list,
	}, nil
}

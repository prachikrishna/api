package pageinfosvc

import (
	"context"
	"database/sql"
	"fmt"
	"grpc-server-streaming/grpc"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type pageInfoServiceServer struct {
	db *sql.DB
	grpc.UnimplementedPageInfoServiceServer
}

/*type Page struct{
	Title string
	Total_words int32
	Total_sentences int32
	Total_images int32
}*/

// NewBookServiceServer creates Book service
func NewPageInfoServiceServer(db *sql.DB) grpc.PageInfoServiceServer {
	return &pageInfoServiceServer{db: db}
}

// connect returns SQL database connection from the pool
func (s *pageInfoServiceServer) connect(ctx context.Context) (*sql.Conn, error) {
	c, err := s.db.Conn(ctx)
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to connect to database-> "+err.Error())
	}
	return c, nil
}

func (s *pageInfoServiceServer) GetDetails(req *grpc.PageRequest,
	stream grpc.PageInfoService_GetDetailsServer,
) error {
	ctx := context.Background()
	c, err := s.connect(ctx)
	if err != nil {
		return err
	}
	defer c.Close()

	// query book by ID
	rows, err := c.QueryContext(ctx, "SELECT `Title`, `Total_words`, `Total_sentences`,`Total_images` FROM PageDetails WHERE `Page_no`=?, `Page_size`=?",
		req.PageNo, req.PageSize)
	if err != nil {
		return status.Error(codes.Unknown, "failed to select from books-> "+err.Error())
	}
	defer rows.Close()

	if !rows.Next() {
		if err := rows.Err(); err != nil {
			return status.Error(codes.Unknown, "failed to retrieve data from books-> "+err.Error())
		}
		return status.Error(codes.NotFound, fmt.Sprintf("Book with ID='%d' is not found",
			req.PageNo))
	}

	// get book data
	var td grpc.PageItems
	if err := rows.Scan(&td.Title, &td.TotalWords, &td.TotalSentences, &td.TotalImages); err != nil {
		return status.Error(codes.Unknown, "failed to retrieve field values from books row-> "+err.Error())
	}

	if rows.Next() {
		return status.Error(codes.Unknown, fmt.Sprintf("found multiple books rows with ID='%d'",
			req.PageNo))
	}

	if err := stream.Send(&td); err != nil {
		return err
	}

	/*return &grpc.ReadResponse{
		Book: &td,
	}, nil*/
	return nil
}

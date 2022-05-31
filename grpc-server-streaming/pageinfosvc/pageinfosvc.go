package pageinfosvc

import (
	"context"
	"database/sql"
	"fmt"
	"grpc-server-streaming/grpc"

	//"log"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type pageInfoServiceServer struct {
	db *sql.DB
	grpc.UnimplementedPageInfoServiceServer
}

/*type container struct {
	items []*grpc.PageItems
}*/

/*type Page struct{
	Title string
	Total_words int32
	Total_sentences int32
	Total_images int32
}*/

// NewPageInfoServiceServer creates PageInfo service
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

	/*if req.GetPageNo() > 1 {
		offset_value:=(req.GetPageNo()-1)*req.GetPageSize()
	}*/

	// query items by page number and size
	rows, err := c.QueryContext(ctx, "SELECT `Title`, `Total_words`, `Total_sentences`,`Total_images` FROM PageDetails WHERE `Page_no`=? LIMIT %d OFFSET %d",
		req.PageNo, req.PageSize, (req.GetPageNo()-1)*req.GetPageSize())
	if err != nil {
		return status.Error(codes.Unknown, "failed to select from pagedetails-> "+err.Error())
	}
	defer rows.Close()

	if !rows.Next() {
		if err := rows.Err(); err != nil {
			return status.Error(codes.Unknown, "failed to retrieve data from pagedetails-> "+err.Error())
		}
		return status.Error(codes.NotFound, fmt.Sprintf("Page with page_no='%d'  not found",
			req.PageNo))
	}

	// get data
	var pg []*grpc.PageItems
	var td *grpc.PageItems

	for rows.Next() {

		if err := rows.Scan(&td.Title, &td.TotalWords, &td.TotalSentences, &td.TotalImages); err != nil {
			return status.Error(codes.Unknown, "failed to retrieve field values from pagedetails row-> "+err.Error())
		}
		pg = append(pg, td)

	}

	if rows.Next() {
		return status.Error(codes.Unknown, fmt.Sprintf("found multiple pagedetails rows with page_no='%d'",
			req.PageNo))
	}
	for _, p := range pg {
		res := &grpc.PageItems{
			Title:          p.Title,
			TotalWords:     p.TotalWords,
			TotalSentences: p.TotalSentences,
			TotalImages:    p.TotalImages,
		}
		if err := stream.Send(res); err != nil {
			return err
		}

	}

	/*res := stream.Send(&grpc.PageItems{
		Title:          td.Title,
		TotalWords:     td.TotalWords,
		TotalSentences: td.TotalSentences,
		TotalImages:    td.TotalImages,
	})*/

	/*if res != nil {
		log.Fatalf("Error when response was sent to the client: %v", res)
	}*/

	return nil
}

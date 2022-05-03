package client

import (
	"context"
	"time"

	mysvcgrpc "prachi/grpc"
	"prachi/mysvc"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var defaultRequestTimeout = time.Second * 10

type grpcService struct {
	grpcClient mysvcgrpc.WordServiceClient
}

// NewGRPCService creates a new gRPC word service connection using the specified connection string.
func NewGRPCService(connString string) (mysvc.Service, error) {
	conn, err := grpc.Dial(connString, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return &grpcService{grpcClient: mysvcgrpc.NewWordServiceClient(conn)}, nil
}

func (s *grpcService) GetWCount(text string) (result []*mysvcgrpc.Word, err error) {
	result = []*mysvcgrpc.Word{}

	req := &mysvcgrpc.GetRequest{
		Text: text,
	}

	ctx, cancelFunc := context.WithTimeout(context.Background(), defaultRequestTimeout)
	defer cancelFunc()
	resp, err := s.grpcClient.GetWCount(ctx, req)
	if err != nil {
		return
	}

	result = append(result, resp.GetWords()...)

	return
}

package main

import (
	"context"

	mysvcgrpc "prachi/grpc"
	"prachi/mysvc"
)

// userServiceController implements the gRPC UserServiceServer interface.
type wordServiceController struct {
	wordService mysvc.Service
	mysvcgrpc.UnimplementedWordServiceServer
}

// NewUserServiceController instantiates a new UserServiceServer.
func NewWordServiceController(wordService mysvc.Service) mysvcgrpc.WordServiceServer {
	return &wordServiceController{
		wordService: wordService,
	}
}

func (ctlr *wordServiceController) GetWCount(ctx context.Context, req *mysvcgrpc.GetRequest) (resp *mysvcgrpc.GetResponse, err error) {
	words := ctlr.wordService.GetWCount(req.GetText())
	if err != nil {
		return
	}

	return &mysvcgrpc.GetResponse{
		Words: words,
	}, nil

}

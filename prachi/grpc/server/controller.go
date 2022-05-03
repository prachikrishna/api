package main

import (
	"context"

	mysvcgrpc "prachi/grpc"
	"prachi/mysvc"
)

// wordServiceController implements the gRPC WordServiceServer interface.
type wordServiceController struct {
	wordService mysvc.Service
	mysvcgrpc.UnimplementedWordServiceServer
}

// NewWordServiceController instantiates a new WordServiceServer.
func NewWordServiceController(wordService mysvc.Service) mysvcgrpc.WordServiceServer {
	return &wordServiceController{
		wordService: wordService,
	}
}

func (ctlr *wordServiceController) GetWCount(ctx context.Context, req *mysvcgrpc.GetRequest) (resp *mysvcgrpc.GetResponse, err error) {
	words, err := ctlr.wordService.GetWCount(req.GetText())
	if err != nil {
		return
	}
	//resp = &mysvcgrpc.GetResponse{}
	//resp.Words=append(resp.Words,marshalWord(words))

	return &mysvcgrpc.GetResponse{
		Words: words,
	}, nil

}

/*func marshalWord(*mysvc.WCount) (w *mysvcgrpc.Word) {
	return &mysvcgrpc.Word{
		Count: w.Count,
		Word:  w.Word,
	}
}*/

package main

import (
	"fmt"
	"log"
	"os"

	mysvccore "prachi/core"
	mysvcgrpc "prachi/grpc/client"
	"prachi/mysvc"
	//"github.com/golang/protobuf/ptypes/empty"
	//"github.com/xiam/to"
)

func main() {
	var localService, grpcService mysvc.Service

	localService = mysvccore.NewService()
	grpcService, err := mysvcgrpc.NewGRPCService("0.0.0.0:50051")
	if err != nil {
		log.Printf("error instantiating gRPC service: %v\n", err)
		os.Exit(1)
	}

	var text string
	if text == "empty" {
		log.Printf("invalid input\n")
		os.Exit(1)
	} else if len(text) != 0 {
		localResult, localErr := localService.GetWCount(text)
		if localErr != nil {
			fmt.Printf("localService.GetWCount() returned an error: %v\n", localErr)
		} else {
			fmt.Printf("localService.GetWCount() returned: %+v\n", localResult)
		}
		grpcResult, remoteErr := grpcService.GetWCount(text)
		if remoteErr != nil {
			fmt.Printf("grpcService.GetWCount() returned an error: %v\n", remoteErr)
		} else {
			fmt.Printf("grpcService.GetWCount() returned: %+v\n", grpcResult)
		}
	}
	os.Exit(0)

}

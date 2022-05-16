package main

import (
	"fmt"
	"os"

	"grpc-bidirectional-firstPart/config"
)

func main() {
	if err := config.RunServer(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}

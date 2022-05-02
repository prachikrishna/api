package mysvc

import (
	"errors"
	word_pb "prachi/grpc"
)

// ErrNotFound signifies that a single requested object was not found.
var ErrNotFound = errors.New("not found")

// User is a user business object.
type WCount struct {
	Word  string
	Count uint32
}

// Service defines the interface exposed by this package.
type Service interface {
	GetWCount(text string) []*word_pb.Word
}

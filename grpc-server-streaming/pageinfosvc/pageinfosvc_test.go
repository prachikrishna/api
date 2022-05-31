package pageinfosvc

import (
	"context"
	"reflect"
	"testing"

	//"time"

	//"github.com/golang/protobuf/ptypes"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"

	"grpc-server-streaming/grpc"
)

func makeStreamMock() *StreamMock {
    return &StreamMock{
        ctx:            context.Background(),
        recvToServer:   make(chan *grpc.SumStreamRequest, 10),
        sentFromServer: make(chan *v1.SumStreamResponse, 10),
    }
}
type StreamMock struct {
    grpc.ServerStream
    ctx            context.Context
    recvToServer   chan *v1.SumStreamRequest
    sentFromServer chan *v1.SumStreamResponse
}
func (m *StreamMock) Context() context.Context {
    return m.ctx
}
func (m *StreamMock) Send(resp *v1.SumStreamResponse) error {
    m.sentFromServer <- resp
    return nil
}
func (m *StreamMock) Recv() (*v1.SumStreamRequest, error) {
    req, more := <-m.recvToServer
    if !more {
        return nil, errors.New("empty")
    }
    return req, nil
}
func (m *StreamMock) SendFromClient(req *v1.SumStreamRequest) error{
    m.recvToServer <- req
    return nil
}
func (m *StreamMock) RecvToClient() (*v1.SumStreamResponse, error) {
    response, more := <-m.sentFromServer
    if !more {
        return nil, errors.New("empty")
    }
    return response, nil
}

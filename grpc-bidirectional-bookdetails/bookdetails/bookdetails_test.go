//bookdetails_test.go:96: rpc error: code = Unknown desc = failed to select from books->
//all expectations were already fulfilled, call to Query 'SELECT `ID`, `Name`, `Author` FROM books WHERE `ID`=?' with args [{Name: Ordinal:1 Value:1}] was not expected
package bookdetails

import (
	"context"
	//"reflect"
	"errors"
	"testing"

	//"database/sql"

	//"time"

	//"github.com/golang/protobuf/ptypes"
	pr "google.golang.org/grpc"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"

	"grpc-bidirectional-bookdetails/grpc"
)

func makeStreamMock() *StreamMock {
	return &StreamMock{
		ctx:            context.Background(),
		recvToServer:   make(chan *grpc.Request, 10),
		sentFromServer: make(chan *grpc.Response, 10),
	}
}

type StreamMock struct {
	pr.ServerStream
	ctx            context.Context
	recvToServer   chan *grpc.Request
	sentFromServer chan *grpc.Response
}

func (m *StreamMock) Context() context.Context {
	return m.ctx
}
func (m *StreamMock) Send(resp *grpc.Response) error {
	m.sentFromServer <- resp
	return nil
}
func (m *StreamMock) Recv() (*grpc.Request, error) {
	req, more := <-m.recvToServer
	if !more {
		return nil, errors.New("empty")
	}
	return req, nil
}
func (m *StreamMock) SendFromClient(req *grpc.Request) error {
	m.recvToServer <- req
	return nil
}
func (m *StreamMock) RecvToClient() (*grpc.Response, error) {
	response, more := <-m.sentFromServer
	if !more {
		return nil, errors.New("empty")
	}
	return response, nil
}

func TestGetDetails(t *testing.T) {
	stream := createStream(t)
	err := stream.SendFromClient(&grpc.Request{
		Id: 1,
	})
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	sumStreamResponse, err := stream.RecvToClient()
	if err != nil {
		t.Error(err.Error())
		return
	}
	want := &grpc.Book{
		Id:     1,
		Name:   "name",
		Author: "author",
	}

	if sumStreamResponse.Book != want {
		t.Errorf("expected %v, instead received %v", want, sumStreamResponse)
	}
}
func createStream(t *testing.T) *StreamMock {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	//defer db.Close()
	stream := makeStreamMock()
	go func() {
		bookServiceServer := NewBookServiceServer(db)
		err := bookServiceServer.GetDetails(stream)
		if err != nil {
			t.Errorf(err.Error())
		}
		close(stream.sentFromServer)
		close(stream.recvToServer)
	}()
	return stream
}

/*import (
	"io"
	"testing"

	"github.com/golang/mock/gomock"
	"grpc-bidirectional-bookdetails/grpc"
	//"github.com/toransahu/grpc-eg-go/mock_machine"
)*/

/*------------------------------------2nd attempt-----------------------------------------*/
//panic: runtime error: invalid memory address or nil pointer dereference [recovered]

/*import (
	//"database/sql"
	"io"
	"testing"

	"grpc-bidirectional-bookdetails/grpc"
	"grpc-bidirectional-bookdetails/mock_grpc"

	"gopkg.in/DATA-DOG/go-sqlmock.v1"

	"github.com/golang/mock/gomock"
)

func TestGetDetails(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	s := NewBookServiceServer(db)
	//s := bookServiceServer{}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockServerStream := mock_grpc.NewMockBookService_GetDetailsServer(ctrl)

	mockResults := []*grpc.Response{}
	cRecv1 := mockServerStream.EXPECT().Recv().Return(&grpc.Request{Id: 1}, nil)
	mockServerStream.EXPECT().Recv().Return(nil, io.EOF).After(cRecv1)
	mockServerStream.EXPECT().Send(gomock.Any()).DoAndReturn(
		func(result *grpc.Response) error {
			mockResults = append(mockResults, result)
			return nil
		}).AnyTimes()

	want := &grpc.Book{
		Id:     1,
		Name:   "name",
		Author: "author",
	}

	err1 := s.GetDetails(mockServerStream)
	if err1 != nil {
		t.Errorf("Execute(%v) got unexpected error: %v", mockServerStream, err)
	}

	for _, result := range mockResults {
		got := result.GetBook()
		//want := wants[i]
		if got != want {
			t.Errorf("got %v, want %v", got, want)
		}
	}

}*/

//3rd attempt

/*import (
	//"context"
	//"errors"
	//"reflect"
	"testing"

	//"time"
	"io"

	//"github.com/golang/protobuf/ptypes"
	"grpc-bidirectional-bookdetails/grpc"
	"grpc-bidirectional-bookdetails/mock_grpc"

	"github.com/stretchr/testify/assert"

	"gopkg.in/DATA-DOG/go-sqlmock.v1"

	//"github.com/envoyproxy/go-control-plane/pkg/server/stream/v3"
	"github.com/golang/mock/gomock"
)

func TestGetDetails(t *testing.T) {
	//ctx := context.Background()
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	s := NewBookServiceServer(db)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockServerStream := mock_grpc.NewMockBookService_GetDetailsServer(ctrl)

	//tm := time.Now().In(time.UTC)
	//reminder, _ := ptypes.TimestampProto(tm)

	/*type args struct {
		ctx context.Context
		req *grpc.Request
	}*/
/*tests := []struct {
	name string
	s    grpc.BookServiceServer
	//args    args
	mockServerStream grpc.BookService_GetDetailsServer
	mock             func()
	want             *grpc.Response
	wantErr          bool
}{
	{
		name: "OK",
		s:    s,
		/*args: args{
			ctx: ctx,
			req: &grpc.ReadRequest{
				Id: 1,
			},
		},*/
//mockServerStream: mock_grpc.NewMockBookService_GetDetailsServer(),
/*mock: func() {
			mockResults := []*grpc.Response{}
			cRecv1 := mockServerStream.EXPECT().Recv().Return(&grpc.Request{Id: 1}, nil)
			mockServerStream.EXPECT().Recv().Return(nil, io.EOF).After(cRecv1)
			mockServerStream.EXPECT().Send(gomock.Any()).DoAndReturn(
				func(result *grpc.Response) error {
					mockResults = append(mockResults, result)
					return nil
				}).AnyTimes()

			//rows := sqlmock.NewRows([]string{"ID", "Name", "Author"}).
			//AddRow(1, "name", "author")
			//mock.ExpectQuery("SELECT (.+) FROM books").WithArgs(1).WillReturnRows(rows)
		},
		want: &grpc.Response{
			Book: &grpc.Book{
				Id:     1,
				Name:   "name",
				Author: "author",
			},
		},
	},
}
for _, tt := range tests {
	t.Run(tt.name, func(t *testing.T) {
		tt.mock()
		err := tt.s.GetDetails(mockServerStream)
		if (err != nil) != tt.wantErr {
			t.Errorf("bookServiceServer.Read() error = %v, wantErr %v", err, tt.wantErr)
			return
		}
		assert.NoError(t, err)
		/*for _, result := range mockResults {
			got := result.GetBook()
			//want := wants[i]
			if got != want {
				t.Errorf("got %v, want %v", got, want)
			}
		}

		if err == nil && !reflect.DeepEqual(got, tt.want) {
			t.Errorf("bookServiceServer.Read() = %v, want %v", got, tt.want)
		}*/
//})
//}
//}

//4th attempt

/*import (
	//"database/sql"
	"io"
	"testing"

	"grpc-bidirectional-bookdetails/grpc"
	"grpc-bidirectional-bookdetails/mock_grpc"

	"github.com/golang/mock/gomock"
)

func TestGetDetails(t *testing.T) {
	s := bookServiceServer{}
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("test get details", func(t *testing.T) {
		mockServerStream := mock_grpc.NewMockBookService_GetDetailsServer(ctrl)
		mockResults := []*grpc.Response{}
		cRecv1 := mockServerStream.EXPECT().Recv().Return(&grpc.Request{Id: 1}, nil)
		mockServerStream.EXPECT().Recv().Return(nil, io.EOF).After(cRecv1)
		mockServerStream.EXPECT().Send(gomock.Any()).DoAndReturn(
			func(result *grpc.Response) error {
				mockResults = append(mockResults, result)
				return nil
			}).AnyTimes()

		want := &grpc.Book{
			Id:     1,
			Name:   "name",
			Author: "author",
		}

		err1 := s.GetDetails(mockServerStream)
		if err1 != nil {
			t.Errorf("got unexpected error: %v", err1)
		}

		for _, result := range mockResults {
			got := result.GetBook()
			if got != want {
				t.Errorf("got %v, want %v", got, want)
			}
		}

	})

}*/

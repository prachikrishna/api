package booksvc

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	//"github.com/golang/protobuf/ptypes"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"

	"grpc-crud-unary-bookservice/grpc"
)

func Test_bookServiceServer_Create(t *testing.T) {
	ctx := context.Background()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	s := NewBookServiceServer(db)
	tm := time.Now().In(time.UTC)
	//reminder, _ := ptypes.TimestampProto(tm)

	type args struct {
		ctx context.Context
		req *grpc.CreateRequest
	}
	tests := []struct {
		name    string
		s       grpc.BookServiceServer
		args    args
		mock    func()
		want    *grpc.CreateResponse
		wantErr bool
	}{
		{
			name: "OK",
			s:    s,
			args: args{
				ctx: ctx,
				req: &grpc.CreateRequest{
					Book: &grpc.Book{
						Name:   "name",
						Author: "author",
					},
				},
			},
			mock: func() {
				mock.ExpectExec("INSERT INTO books").WithArgs("name", "author").
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			want: &grpc.CreateResponse{
				Id: 1,
			},
		},
		{
			name: "INSERT failed",
			s:    s,
			args: args{
				ctx: ctx,
				req: &grpc.CreateRequest{
					Book: &grpc.Book{
						Name:   "name",
						Author: "author",
					},
				},
			},
			mock: func() {
				mock.ExpectExec("INSERT INTO books").WithArgs("name", "author", tm).
					WillReturnError(errors.New("INSERT failed"))
			},
			wantErr: true,
		},
		{
			name: "LastInsertId failed",
			s:    s,
			args: args{
				ctx: ctx,
				req: &grpc.CreateRequest{
					Book: &grpc.Book{
						Name:   "name",
						Author: "author",
					},
				},
			},
			mock: func() {
				mock.ExpectExec("INSERT INTO books").WithArgs("name", "author", tm).
					WillReturnResult(sqlmock.NewErrorResult(errors.New("LastInsertId failed")))
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			got, err := tt.s.CreateBook(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("toDoServiceServer.Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("toDoServiceServer.Create() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_bookServiceServer_Read(t *testing.T) {
	ctx := context.Background()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	s := NewBookServiceServer(db)
	//tm := time.Now().In(time.UTC)
	//reminder, _ := ptypes.TimestampProto(tm)

	type args struct {
		ctx context.Context
		req *grpc.ReadRequest
	}
	tests := []struct {
		name    string
		s       grpc.BookServiceServer
		args    args
		mock    func()
		want    *grpc.ReadResponse
		wantErr bool
	}{
		{
			name: "OK",
			s:    s,
			args: args{
				ctx: ctx,
				req: &grpc.ReadRequest{
					Id: 1,
				},
			},
			mock: func() {
				//check whether reminder needs to be present here
				rows := sqlmock.NewRows([]string{"ID", "Name", "Author"}).
					AddRow(1, "name", "author")
				mock.ExpectQuery("SELECT (.+) FROM books").WithArgs(1).WillReturnRows(rows)
			},
			want: &grpc.ReadResponse{
				Book: &grpc.Book{
					Id:     1,
					Name:   "name",
					Author: "author",
				},
			},
		},
		{
			name: "SELECT failed",
			s:    s,
			args: args{
				ctx: ctx,
				req: &grpc.ReadRequest{
					Id: 1,
				},
			},
			mock: func() {
				mock.ExpectQuery("SELECT (.+) FROM books").WithArgs(1).
					WillReturnError(errors.New("SELECT failed"))
			},
			wantErr: true,
		},
		{
			name: "Not found",
			s:    s,
			args: args{
				ctx: ctx,
				req: &grpc.ReadRequest{
					Id: 1,
				},
			},
			mock: func() {
				rows := sqlmock.NewRows([]string{"ID", "Name", "Author"})
				mock.ExpectQuery("SELECT (.+) FROM books").WithArgs(1).WillReturnRows(rows)
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			got, err := tt.s.ReadBook(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("bookServiceServer.Read() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err == nil && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("bookServiceServer.Read() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_toDoServiceServer_Update(t *testing.T) {
	ctx := context.Background()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	s := NewBookServiceServer(db)
	tm := time.Now().In(time.UTC)
	//reminder, _ := ptypes.TimestampProto(tm)

	type args struct {
		ctx context.Context
		req *grpc.UpdateRequest
	}
	tests := []struct {
		name    string
		s       grpc.BookServiceServer
		args    args
		mock    func()
		want    *grpc.UpdateResponse
		wantErr bool
	}{
		{
			name: "OK",
			s:    s,
			args: args{
				ctx: ctx,
				req: &grpc.UpdateRequest{
					Book: &grpc.Book{
						Id:     1,
						Name:   "new name",
						Author: "new author",
					},
				},
			},
			mock: func() {
				mock.ExpectExec("UPDATE books").WithArgs("new name", "new author", 1).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			want: &grpc.UpdateResponse{
				Updated: 1,
			},
		},
		{
			name: "UPDATE failed",
			s:    s,
			args: args{
				ctx: ctx,
				req: &grpc.UpdateRequest{
					Book: &grpc.Book{
						Id:     1,
						Name:   "new name",
						Author: "new author",
					},
				},
			},
			mock: func() {
				mock.ExpectExec("UPDATE Book").WithArgs("new name", "new author", tm, 1).
					WillReturnError(errors.New("UPDATE failed"))
			},
			wantErr: true,
		},
		{
			name: "RowsAffected failed",
			s:    s,
			args: args{
				ctx: ctx,
				req: &grpc.UpdateRequest{
					Book: &grpc.Book{
						Id:     1,
						Name:   "new name",
						Author: "new author",
					},
				},
			},
			mock: func() {
				mock.ExpectExec("UPDATE Book").WithArgs("new name", "new author", tm, 1).
					WillReturnResult(sqlmock.NewErrorResult(errors.New("RowsAffected failed")))
			},
			wantErr: true,
		},
		{
			name: "Not Found",
			s:    s,
			args: args{
				ctx: ctx,
				req: &grpc.UpdateRequest{
					Book: &grpc.Book{
						Id:     1,
						Name:   "new name",
						Author: "new author",
					},
				},
			},
			mock: func() {
				mock.ExpectExec("UPDATE books").WithArgs("new name", "new author", tm, 1).
					WillReturnResult(sqlmock.NewResult(1, 0))
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			got, err := tt.s.UpdateBook(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("bookServiceServer.Update() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("bookServiceServer.Update() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_toDoServiceServer_Delete(t *testing.T) {
	ctx := context.Background()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	s := NewBookServiceServer(db)

	type args struct {
		ctx context.Context
		req *grpc.DeleteRequest
	}
	tests := []struct {
		name    string
		s       grpc.BookServiceServer
		args    args
		mock    func()
		want    *grpc.DeleteResponse
		wantErr bool
	}{
		{
			name: "OK",
			s:    s,
			args: args{
				ctx: ctx,
				req: &grpc.DeleteRequest{
					Id: 1,
				},
			},
			mock: func() {
				mock.ExpectExec("DELETE FROM books").WithArgs(1).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			want: &grpc.DeleteResponse{
				Deleted: 1,
			},
		},
		{
			name: "DELETE failed",
			s:    s,
			args: args{
				ctx: ctx,
				req: &grpc.DeleteRequest{
					Id: 1,
				},
			},
			mock: func() {
				mock.ExpectExec("DELETE FROM books").WithArgs(1).
					WillReturnError(errors.New("DELETE failed"))
			},
			wantErr: true,
		},
		{
			name: "RowsAffected failed",
			s:    s,
			args: args{
				ctx: ctx,
				req: &grpc.DeleteRequest{
					Id: 1,
				},
			},
			mock: func() {
				mock.ExpectExec("DELETE FROM books").WithArgs(1).
					WillReturnResult(sqlmock.NewErrorResult(errors.New("RowsAffected failed")))
			},
			wantErr: true,
		},
		{
			name: "Not Found",
			s:    s,
			args: args{
				ctx: ctx,
				req: &grpc.DeleteRequest{
					Id: 1,
				},
			},
			mock: func() {
				mock.ExpectExec("DELETE FROM books").WithArgs(1).
					WillReturnResult(sqlmock.NewResult(1, 0))
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			got, err := tt.s.DeleteBook(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("bookServiceServer.Delete() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("bookServiceServer.Delete() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_toDoServiceServer_ReadAll(t *testing.T) {
	ctx := context.Background()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	s := NewBookServiceServer(db)

	type args struct {
		ctx context.Context
		req *grpc.ListAllRequest
	}
	tests := []struct {
		name    string
		s       grpc.BookServiceServer
		args    args
		mock    func()
		want    *grpc.ListAllResponse
		wantErr bool
	}{
		{
			name: "OK",
			s:    s,
			args: args{
				ctx: ctx,
				req: &grpc.ListAllRequest{},
			},
			mock: func() {
				rows := sqlmock.NewRows([]string{"ID", "Name", "Author"}).
					AddRow(1, "title 1", "description 1").
					AddRow(2, "title 2", "description 2")
				mock.ExpectQuery("SELECT (.+) FROM books").WillReturnRows(rows)
			},
			want: &grpc.ListAllResponse{
				Books: []*grpc.Book{
					{
						Id:     1,
						Name:   "title 1",
						Author: "description 1",
					},
					{
						Id:     2,
						Name:   "title 2",
						Author: "description 2",
					},
				},
			},
		},
		{
			name: "Empty",
			s:    s,
			args: args{
				ctx: ctx,
				req: &grpc.ListAllRequest{},
			},
			mock: func() {
				rows := sqlmock.NewRows([]string{"ID", "Name", "Author"})
				mock.ExpectQuery("SELECT (.+) FROM books").WillReturnRows(rows)
			},
			want: &grpc.ListAllResponse{
				Books: []*grpc.Book{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			got, err := tt.s.ListAllBooks(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("bookServiceServer.ReadAll() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("bookServiceServer.ReadAll() = %v, want %v", got, tt.want)
			}
		})
	}
}

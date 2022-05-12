package pageinfosvc

import (
	//"context"
	"reflect"
	"testing"

	//"time"

	//"github.com/golang/protobuf/ptypes"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"

	"grpc-server-streaming/grpc"
)

func Test_pageinfoServiceServer_Read(t *testing.T) {
	//ctx := context.Background()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	s := NewPageInfoServiceServer(db)
	//tm := time.Now().In(time.UTC)
	//reminder, _ := ptypes.TimestampProto(tm)

	type args struct {
		req    *grpc.PageRequest
		stream grpc.PageInfoService_GetDetailsServer
	}
	tests := []struct {
		name    string
		s       grpc.PageInfoServiceServer
		args    args
		mock    func()
		want    error
		wantErr bool
	}{
		{
			name: "OK",
			s:    s,
			args: args{
				req: &grpc.PageRequest{
					PageNo:   1,
					PageSize: "a4",
				},
			},
			mock: func() {
				rows := sqlmock.NewRows([]string{"Page_no", "Page_size", "Title", "Total_word", "Total_sentences", "Total_images"}).
					AddRow(1, "a4", "Death on the nile", 1000, 100, 5)
				mock.ExpectQuery("SELECT (.+) FROM pagedetails").WithArgs(1, 2).WillReturnRows(rows)
			},
			/*want: &grpc.PageItems{
				Title:          "Death on the nile",
				TotalWords:     1000,
				TotalSentences: 100,
				TotalImages:    5,
			},*/
			want: nil,
		},
		{
			name: "Not found",
			s:    s,
			args: args{
				req: &grpc.PageRequest{
					PageNo:   1,
					PageSize: "a4",
				},
			},
			mock: func() {
				rows := sqlmock.NewRows([]string{"Page_no", "Page_size", "Title", "Total_word", "Total_sentences", "Total_images"})
				mock.ExpectQuery("SELECT (.+) FROM pagedetails").WithArgs(1, 2).WillReturnRows(rows)
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			got, err := tt.s.GetDetails(tt.args.req, tt.args.stream)
			if (err != nil) != tt.wantErr {
				t.Errorf("pageinfoServiceServer.GetDetails() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err == nil && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("pageinfoServiceServer.GetDetails() = %v, want %v", got, tt.want)
			}
		})
	}
}

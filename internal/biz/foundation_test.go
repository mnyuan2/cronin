package biz

import (
	"context"
	"cron/internal/pb"
	"fmt"
	"testing"
)

func TestFoundationService_ParseProto(t *testing.T) {
	r := &pb.ParseProtoRequest{Proto: `syntax = "proto3";
package merchantpush;
service Merchantpush {
  rpc Echo(EchoRequest)returns(EchoResponse){}
}
message EchoRequest{
  string a = 1;
  int32 b = 2;
  string body = 3;
}
message EchoResponse{
  string a = 1;
  int32 b = 2;
  string body = 3;
}`}

	resp, err := NewDicService(context.Background(), nil).ParseProto(r)
	fmt.Println(resp, err)
}

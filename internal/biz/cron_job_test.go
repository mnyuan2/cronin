package biz

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"net/rpc"
	"testing"
)

// 构建rpc请求
func TestCronJob_Rpc(t *testing.T) {
	type EchoResponse struct {
	}
	type EchoRequest struct {
	}

	cli, err := rpc.DialHTTP("tcp", "127.0.0.1:21014") // 1.14.156.225:21014
	if err != nil {
		panic(err)
	}

	req := &EchoRequest{}
	resp := &EchoResponse{}

	err = cli.Call("/merchantpush.Merchantpush/Echo", req, resp)
	if err != nil {
		panic(fmt.Errorf("调用失败，%w", err))
	}

	fmt.Println(resp)
}

type EchoRequest struct {
}

// grpc必须实现
func (m *EchoRequest) Reset() {

}

// grpc必须实现
func (m *EchoRequest) String() string {
	return ""
}

// grpc必须实现
func (m *EchoRequest) ProtoMessage() {

}

// 构建rpc请求
func TestCronJob_Grpc(t *testing.T) {

	conn, err := grpc.Dial("localhost:21014", grpc.WithTransportCredentials(insecure.NewCredentials())) // 1.14.156.225:21014
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	req := &EchoRequest{}
	resp := &EchoRequest{}

	err = conn.Invoke(context.Background(), "/merchantpush.Merchantpush/Echo", req, resp)
	if err != nil {
		panic(fmt.Errorf("调用失败，%w", err))
	}

	fmt.Println(resp)
}

package biz

import (
	"context"
	"cron/internal/models"
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

// 构建rpc请求
func TestCronJob_Grpc(t *testing.T) {

	conn, err := grpc.Dial("localhost:21014", grpc.WithTransportCredentials(insecure.NewCredentials())) // 1.14.156.225:21014
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	//39;

	req := &models.GrpcRequest{}
	req.SetParam(`{"a":"a","b":1}`) // 这个参数的传递，还要验证一下。
	resp := &models.GrpcRequest{}
	//resp := &models.GrpcResponse{}

	conf := conn.GetMethodConfig("/merchantpush.Merchantpush/Echo")
	fmt.Println(conf)

	err = conn.Invoke(context.Background(), "/merchantpush.Merchantpush/Echo", req, resp)
	if err != nil {
		panic(fmt.Errorf("调用失败，%w", err))
	}

	fmt.Println(resp)
}

package biz

import (
	"context"
	"cron/internal/models"
	"fmt"
	"github.com/axgle/mahonia"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"io/ioutil"
	"net/rpc"
	"os/exec"
	"runtime"
	"strings"
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

// 构建grpc请求
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
	// 还有一个方案：用户输入配置(请求、响应)和请求参数，程序把配置写为文件，这样就有结构体了，就可以在请求时携带上对应的结构体了，就能正常请求了。

	conf := conn.GetMethodConfig("/merchantpush.Merchantpush/Echo")
	fmt.Println(conf)

	err = conn.Invoke(context.Background(), "/merchantpush.Merchantpush/Echo", req, resp)
	if err != nil {
		panic(fmt.Errorf("调用失败，%w", err))
	}

	fmt.Println(resp)
}

// 构建sell执行
func TestCronJob_Cmd(t *testing.T) {
	//name := "curl"
	//arg := "http://baidu.com"

	// 参数：1.命令名称、2.参数；
	//cmd := exec.Command(name, arg) // 命令和参数
	//cmd := exec.Command("sh.exe","-c", "echo abc") // 合并 linux 命令

	data := strings.Split("curl http://175.178.108.84:6123/_cat/indices?v&h=i,tm&s=tm:desc", " ")
	if len(data) < 2 {
		t.Fatal("命令参数不合法")
	}

	cmd := exec.Command(data[0], data[1:]...) // 合并 winds 命令

	// windows 平台执行时，隐藏cmd窗口
	if runtime.GOOS == "windows" {
		//	fmt.Println("windows")
		//	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
		if b, err := cmd.Output(); err != nil {
			t.Fatal("结果获取失败", err)
		} else {
			fmt.Printf("结果 |%s|", str2gbk(b))
			return
		}
	} else {
		//获取输出对象
		stdout, err := cmd.StdoutPipe()
		if err != nil {
			t.Error(err)
		}
		defer stdout.Close()

		if err = cmd.Run(); err != nil {
			t.Fatal("执行失败，", err)
		}

		if b, err := ioutil.ReadAll(stdout); err != nil {
			t.Error("结果获取失败", err)
		} else {
			fmt.Println(string(b))
		}
	}

}

// 字符串GBK2312编码方式解码方法
func str2gbk(text []byte) []byte {

	srcCoder := mahonia.NewDecoder("gbk").ConvertString(string(text))
	return []byte(srcCoder)
}

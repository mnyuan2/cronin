package biz

import (
	"context"
	"cron/internal/models"
	"cron/internal/pb"
	"fmt"
	"github.com/axgle/mahonia"
	"github.com/robfig/cron/v3"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"io/ioutil"
	"log"
	"net/rpc"
	"os/exec"
	"runtime"
	"strings"
	"testing"
	"time"
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

// 构建shell执行
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

// 执行shell脚本
func TestCronJob_Cmd2(t *testing.T) {
	/*
		127 找不到可执行文件；
		126 没有操作权限；
		也可以直接去执行脚本文件，但我认为并不适用当前程序，直接把脚本的文件复制出来好一点。
	*/
	//os.Chdir("e:\\WorkApps\\go\\src\\incron")
	//shell := `./test_shell.sh`
	//	shell := `#!/bin/bash
	//echo 'abc'`
	//	ph, err := exec.LookPath("internal/biz/test_shell.sh")
	//	if err != nil {
	//		t.Fatal("文件错误", err.Error())
	//	}
	//	fmt.Println(ph)

	if runtime.GOOS == "windows" {
		ph := `dir;
echo 258;` // 这里目前的问题是，没有执行第二行
		e := exec.Command("cmd", "/C", ph)
		cmd, err := e.Output()
		if err != nil {
			t.Fatal("命令执行错误：", err.Error())
		}
		srcCoder := mahonia.NewDecoder("gbk").ConvertString(string(cmd))
		fmt.Println("执行结果：", string(srcCoder))
	} else {
		// 对于linux 脚本文件 是支持的
		ph := `#!/bin/bash
dir
echo 258`
		e := exec.Command("sh", "-c", ph) // "/bin/bash"
		cmd, err := e.Output()
		if err != nil {
			t.Fatal("命令执行错误：", err.Error())
		}
		srcCoder := mahonia.NewDecoder("gbk").ConvertString(string(cmd))
		fmt.Println("执行结果：", string(srcCoder))
	}
}

// 执行shell脚本
func TestCronJob_Cmd3(t *testing.T) {
	shell := `#!/bin/bash
list="dev-jaeger-span-2022-11-10
dev-jaeger-span-2022-11-09
dev-jaeger-span-2022-11-08
dev-jaeger-span-2022-11-07
dev-jaeger-span-2022-11-06
dev-jaeger-span-2022-11-05
dev-jaeger-span-2022-11-04
dev-jaeger-span-2022-11-03
dev-jaeger-span-2022-11-02
dev-jaeger-span-2022-11-01"
index=5 # 保留最近的指定数量
echo "-----------------------"
for item in $list; do
 if ((index>0)); then
   ((index--))
   echo "index: $index"
   continue
 fi
 echo "echo: $item" # 执行删除语句
done
echo "任务执行完毕..."`
	//data := strings.Split(shell, " ")
	//shell = "test_shell.sh"

	//os.Chdir("e:\\WorkApps\\go\\src\\incron")
	//ph, err := exec.LookPath("internal/biz/test_shell.sh")
	//if err != nil {
	//	t.Fatal("文件错误", err.Error())
	//}

	cmd, err := exec.Command("sh", "-c", shell).Output()
	if err != nil {
		t.Fatal("命令执行错误：", err.Error())
	}
	for true {
		continue
	}
	fmt.Println("执行结果：", string(cmd))
}

// 执行sql命令
func TestCronJob_Mysql(t *testing.T) {
	conf := &models.CronConfig{}
	r := &pb.CronSql{
		Driver: models.SqlSourceMysql,
		Source: pb.CronSqlSource{
			Id:       0,
			Title:    "zby.dev",
			Hostname: "gz-cdb-6ggn2bux.sql.tencentcdb.com",
			Database: "zhubaoe",
			Username: "root",
			Password: "Zby_123456",
			Port:     "63438",
		},
		ErrAction: models.SqlErrActionProceed,
		Statement: []string{
			"UPDATE cron_log set body='修改3' WHERE id=2247",
			"SELECT 1、",
		},
	}

	ctx := context.Background()
	res, err := NewCronJob(conf).sqlMysql(ctx, r)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(string(res))
}

type J struct {
	cronId cron.EntryID // 任务id
}

func (j *J) Run() {
	e := cronRun.Entry(j.cronId)
	cronRun.Remove(j.cronId)
	if e.ID == j.cronId {
		return
	}

	fmt.Println("任务被执行了", j.cronId)
	//自行移除队列
}

// 特别研究
func TestCronJob_Demo(t *testing.T) {
	/*
			研究一下单次定时任务
				也可以叫做临时任务。
			最好也是把任务创建在表中，会好维护点。
				否则就要放内存，总是不太稳定。
			时间就是标准年月日时分秒了。
		验证没有问题
	*/
	ti := time.Now().Add(time.Second - 10)

	s, err := NewScheduleOnce(ti.Format(time.DateTime))
	if err != nil {
		t.Fatal(err)
	}
	j := &J{}
	j.cronId = cronRun.Schedule(s, j)
	log.Println("等待执行...")
	time.Sleep(time.Minute * 2)
	t.Log("end...")
}

// 字符串GBK2312编码方式解码方法
func str2gbk(text []byte) []byte {

	srcCoder := mahonia.NewDecoder("gbk").ConvertString(string(text))
	return []byte(srcCoder)
}

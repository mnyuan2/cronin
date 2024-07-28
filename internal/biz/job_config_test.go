package biz

import (
	"bytes"
	"context"
	"cron/internal/basic/config"
	"cron/internal/basic/db"
	"cron/internal/models"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/axgle/mahonia"
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/desc/protoparse"
	jsoniter "github.com/json-iterator/go"
	"github.com/robfig/cron/v3"
	"golang.org/x/crypto/ssh"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"io"
	"io/ioutil"
	"log"
	"net/rpc"
	"os/exec"
	"runtime"
	"strings"
	"testing"
	"text/template"
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
// 此方法存在的问题是 请求和响应的结构体必须预先准确定义，否则参数无法传递。
func TestCronJob_Grpc(t *testing.T) {

	conn, err := grpc.Dial("localhost:21014", grpc.WithTransportCredentials(insecure.NewCredentials())) // 1.14.156.225:21014
	if err != nil {
		panic(fmt.Sprintf("无法连接到gRPC服务器：%v", err))
	}
	defer conn.Close()

	req := &models.GrpcRequest{}

	req.SetParam(`{"a":"a","b":1}`) // 这个参数的传递，还要验证一下。
	//reqBytes, err := proto.Marshal(req)
	//if err != nil {
	//	t.Fatal("请求序列化错误，%w", err)
	//}

	resp := &models.GrpcRequest{}

	conf := conn.GetMethodConfig("/merchantpush.Merchantpush/Echo")
	fmt.Println(conf)

	//resp := &models.GrpcResponse{}
	err = conn.Invoke(context.Background(), "/merchantpush.Merchantpush/Echo", req, resp)
	if err != nil {
		panic(fmt.Errorf("调用失败，%w", err))
	}

	fmt.Println(resp)
}

// 构建grpc请求 2版
func TestCronJob_Grpc2(t *testing.T) {
	conf := &models.CronConfig{
		Command: []byte(`{
    "cmd": "",
    "rpc": {
        "addr": "localhost:21014",
        "body": "{\"a\":\"a\",\"b\":2,\"body\":\"中文\"}",
        "proto": "",
        "action": "merchantpush.Merchantpush/Echo",
        "header": null,
        "method": "GRPC"
    },
    "sql": {
        "driver": "mysql",
        "source": {
            "id": 0,
            "port": "",
            "title": "",
            "database": "",
            "hostname": "",
            "password": "",
            "username": ""
        },
        "statement": [],
        "err_action": 1
    },
    "http": {
        "url": "",
        "body": "",
        "header": [
            {
                "key": "",
                "value": ""
            }
        ],
        "method": "GET"
    }
}`),
	}

	//db.New(context.Background()).Write.Where("id=?", 116).Find(conf)

	r := NewJobConfig(conf)
	r.commandParse.Rpc.Proto = `syntax = "proto3";
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
}`

	ctx := context.Background()
	res, err := NewJobConfig(conf).rpcGrpc(ctx, r.commandParse.Rpc)

	fmt.Println(string(res), err)
}

// 解析grpc文件
func Test_GrpcParse(t *testing.T) {
	input := map[string]string{
		"merchantpush.proto": `syntax = "proto3";
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
}
`,
	}
	parser := protoparse.Parser{
		Accessor: func(filename string) (io.ReadCloser, error) {
			f, ok := input[filename]
			if !ok {
				return nil, fmt.Errorf("file not found: %s", filename)
			}
			return io.NopCloser(strings.NewReader(f)), nil
		},
	}
	fileDescriptors, err := parser.ParseFiles("merchantpush.proto")
	if err != nil {
		t.Fatal("文件错误", err.Error())
	}
	// 因为只有一个文件，所以肯定只有一个 fileDescriptor
	fileDescriptor := fileDescriptors[0]
	m := make(map[string]interface{})
	for _, msgDescriptor := range fileDescriptor.GetMessageTypes() {
		m[msgDescriptor.GetName()] = convertMessageToMap(msgDescriptor)
	}
	bs, _ := json.MarshalIndent(m, "", "\t")
	fmt.Println(string(bs))
}

func convertMessageToMap(message *desc.MessageDescriptor) map[string]interface{} {
	m := make(map[string]interface{})
	for _, fieldDescriptor := range message.GetFields() {
		fieldName := fieldDescriptor.GetName()
		if fieldDescriptor.IsRepeated() {
			// 如果是一个数组的话，就返回 nil 吧
			m[fieldName] = nil
			continue
		}
		switch fieldDescriptor.GetType() {
		case descriptor.FieldDescriptorProto_TYPE_MESSAGE:
			m[fieldName] = convertMessageToMap(fieldDescriptor.GetMessageType())
			continue
		}
		m[fieldName] = fieldDescriptor.GetDefaultValue()
	}
	return m
}

// http任务
func TestCronJob_Http(t *testing.T) {
	hader := map[string]string{}
	hader = nil
	if config.MainConf().User.AdminAccount != "" {
		s := base64.StdEncoding.EncodeToString([]byte(config.MainConf().User.AdminAccount + ":" + config.MainConf().User.AdminPassword))
		fmt.Println(s)
		hader = map[string]string{
			"Authorization": "Basic " + s,
		}
	}

	ctx := context.Background()
	job := &JobConfig{}
	resp, err := job.httpRequest(ctx, "POST",
		"http://127.0.0.1:9003/log/del",
		[]byte(fmt.Sprintf(`{"retention":"%s"}`, config.MainConf().Task.LogRetention)),
		hader)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("请求响应", string(resp))
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

// 在远程服务器上执行shell脚本
func TestCronJob_Cmd4(t *testing.T) {
	// SSH连接配置
	config := &ssh.ClientConfig{
		User: "root",
		Auth: []ssh.AuthMethod{
			ssh.Password("xxxxxx"),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	// 连接到远程服务器
	conn, err := ssh.Dial("tcp", "127.0.0.1:22", config)
	if err != nil {
		log.Fatalf("拨号失败: %v", err)
	}
	defer conn.Close()

	// 创建一个新的会话
	session, err := conn.NewSession()
	if err != nil {
		log.Fatalf("创建会话失败: %v", err)
	}
	defer session.Close()

	// 执行Shell脚本
	output, err := session.CombinedOutput("/soft/scripts/menu_script/run.sh asd hjk 999")
	if err != nil {
		log.Fatalf("执行脚本失败: %v", err)
	}

	//session.CombinedOutput()

	// 打印脚本输出
	fmt.Println(string(output))
}

// 执行sql命令
func TestCronJob_Mysql(t *testing.T) {
	conf := &models.CronConfig{}
	//r := &pb.CronSql{
	//	Driver: models.SqlSourceMysql,
	//	Source: &pb.CronSqlSource{
	//		Id:       0,
	//		Title:    "zby.dev",
	//		Hostname: "gz-cdb-6ggn2bux.sql.tencentcdb.com",
	//		Database: "zhubaoe",
	//		Username: "root",
	//		Password: "Zby_123456",
	//		Port:     "63438",
	//	},
	//	ErrAction: models.SqlErrActionProceed,
	//	Statement: []string{
	//		"UPDATE cron_log set body='修改3' WHERE id=2247",
	//		"SELECT 1、",
	//	},
	//}

	db.New(context.Background()).Where("id=?", 114).Find(conf)

	r := NewJobConfig(conf)
	r.Run()

	//ctx := context.Background()
	//res, err := r.sqlMysql(ctx, r.commandParse.Sql)
	//if err != nil {
	//	t.Fatal(err)
	//}
	time.Sleep(time.Minute * 3)
	//fmt.Println(string(res))
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
	ti := time.Now().Add(time.Second * 80)
	time.Sleep(time.Second * 2)

	s, err := NewScheduleOnce(ti.Format(time.DateTime))
	if err != nil {
		t.Fatal(err.Error())
	}
	fmt.Println(s.Next(time.Now()).Format(time.DateTime), "|", ti.Format(time.DateTime), "\n", s.execTime.UnixMilli(), ti.UnixMilli(), "\n", s.execTime.UnixMicro(), ti.UnixMicro(), "\n", s.execTime.UnixNano(), ti.UnixNano())

	//j := &J{}
	//j.cronId = cronRun.Schedule(s, j)
	//log.Println("等待执行...")
	//time.Sleep(time.Minute * 2)
	t.Log("end...")
}

// 字符串GBK2312编码方式解码方法
func str2gbk(text []byte) []byte {

	srcCoder := mahonia.NewDecoder("gbk").ConvertString(string(text))
	return []byte(srcCoder)
}

// 配置解析
func TestConfigParse(t *testing.T) {
	param := `{"a":"A","b":{"b1":"B1","B2":22},"c":3}`

	varParams := map[string]any{}
	if param != "" {
		if er := jsoniter.UnmarshalFromString(param, &varParams); er != nil {
			t.Fatal(er)
		}
	}

	err := NewJobConfig(&models.CronConfig{
		Id:           0,
		Env:          "",
		EntryId:      0,
		Type:         0,
		Name:         "",
		Spec:         "",
		Protocol:     0,
		Command:      []byte(`{"cmd": {"host": {"id": -1, "ip": "", "port": "", "type": "", "user": "", "secret": ""}, "type": "bash", "origin": "local", "statement": {"git": {"ref": "", "path": [""], "owner": "", "link_id": 0, "project": ""}, "type": "", "local": "echo [[.a]] '\\n' [[.b]] [[jsonString2 .b]] [[.b.b1]]", "is_batch": 0}}, "git": {"events": [], "link_id": 0}, "rpc": {"addr": "", "body": "", "proto": "", "action": "", "header": [], "method": "GRPC", "actions": []}, "sql": {"driver": "mysql", "origin": "local", "source": {"id": 0, "port": "", "title": "", "database": "", "hostname": "", "password": "", "username": ""}, "interval": 0, "statement": [], "err_action": 1, "err_action_name": ""}, "http": {"url": "", "body": "", "header": [{"key": "", "value": ""}], "method": "GET"}, "jenkins": {"name": "", "params": [{"key": "", "value": ""}], "source": {"id": 0}}}`),
		Remark:       "",
		Status:       0,
		StatusRemark: "",
		StatusDt:     "",
		UpdateDt:     "",
		CreateDt:     "",
		MsgSet:       nil,
		VarFields:    []byte(`[{"key": "a", "value": "常规参数"}, {"key": "b", "value": "对象,内部至少包含b1"}, {"key": "", "value": ""}]`),
	}).Parse(varParams)

	if err != nil {
		t.Fatal(err)
	}
}

func TestTemplateV2(t *testing.T) {
	input := strReplace("Hello, [[.name]]! Your age is [[.age]] and your sex is [[.sex]].")
	jsonStr := `{"sex": "男", "age": 180, "name": ["title2", "title1"]}`

	// 定义一个 map 用于存储解析后的数据
	data := make(map[string]interface{})
	// 解析 JSON 字符串
	err := json.Unmarshal([]byte(jsonStr), &data)
	if err != nil {
		fmt.Println("Failed to parse JSON:", err)
		return
	}

	// 创建模板
	tmpl := template.New("tmpl")

	// 解析模板字符串
	_tmpl, err := tmpl.Parse(input)
	if err != nil {
		fmt.Println("Failed to parse template:", err)
		return
	}

	// 应用模板到数据
	var buf bytes.Buffer
	err = _tmpl.Execute(&buf, data)
	if err != nil {
		fmt.Println("Failed to execute template:", err)
		return
	}

	// 获取替换后的字符串
	result := buf.String()
	fmt.Println(result)
}

func strReplace(input string) string {
	input = strings.ReplaceAll(input, "[[", "{{")
	input = strings.ReplaceAll(input, "]]", "}}")
	return input
}

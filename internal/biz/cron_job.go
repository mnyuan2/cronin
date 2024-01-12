package biz

import (
	"bytes"
	"context"
	"cron/internal/basic/grpcurl"
	"cron/internal/basic/util"
	"cron/internal/data"
	"cron/internal/models"
	"cron/internal/pb"
	"errors"
	"fmt"
	"github.com/axgle/mahonia"
	"github.com/jhump/protoreflect/grpcreflect"
	jsoniter "github.com/json-iterator/go"
	"github.com/robfig/cron/v3"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

type ScheduleOnce struct {
	s        string // 单次时间标志 YYYY-mm-dd hh:ii:ss
	execTime time.Time
}

func NewScheduleOnce(dateTime string) (m *ScheduleOnce, err error) {
	m = &ScheduleOnce{
		s: dateTime,
	}
	m.execTime, err = time.ParseInLocation(time.DateTime, dateTime, time.Local)
	if err != nil {
		return nil, fmt.Errorf("执行时间格式不规范，%s", err.Error())
	}
	if m.execTime.Unix()-60*60*1 < time.Now().Unix() {
		return nil, fmt.Errorf("执行时间至少间隔1小时以后")
	}

	return m, nil
}

func (m *ScheduleOnce) Next(t time.Time) time.Time {
	if m.execTime.Unix() < t.Unix() {
		return t
	}
	return m.execTime
}

type CronJob struct {
	conf         *models.CronConfig
	commandParse *pb.CronConfigCommand
	ErrorCount   int // 连续错误
}

// 任务执行器
func NewCronJob(conf *models.CronConfig) *CronJob {
	com := &pb.CronConfigCommand{}
	_ = jsoniter.UnmarshalFromString(conf.Command, com)

	return &CronJob{conf: conf, commandParse: com}
}

// 执行任务
func (job *CronJob) Run() {
	var g *models.CronLog
	var res []byte
	var err error
	st := time.Now()
	ctx := context.Background()
	if job.conf.Type == models.TypeOnce { // 单次任务，自行移除
		e := cronRun.Entry(cron.EntryID(job.conf.EntryId))
		cronRun.Remove(cron.EntryID(job.conf.EntryId))
		if e.ID == 0 {
			return
		}
	}

	defer func() {
		if err := util.PanicInfo(recover()); err != "" {
			data.NewCronLogData(ctx).Add(models.NewErrorCronLog(job.conf, fmt.Sprintf("异常：%s", err), "panic", st))
		}
		data.NewCronLogData(ctx).Add(g)
		if job.conf.Type == models.TypeOnce { // 单次执行完毕后，状态也要更新
			job.conf.Status = models.ConfigStatusFinish
			job.conf.EntryId = 0
			data.NewCronConfigData(ctx).ChangeStatus(job.conf, "执行完成")
		}
	}()

	fmt.Println("执行 "+job.conf.GetProtocolName()+" 任务", job.conf.Id, job.conf.Name)
	switch job.conf.Protocol {
	case models.ProtocolHttp:
		res, err = job.httpFunc(ctx)
	case models.ProtocolRpc:
		res, err = job.rpcFunc(ctx)
	case models.ProtocolCmd:
		res, err = job.cmdFunc(ctx)
	case models.ProtocolSql:
		res, err = job.sqlFunc(ctx)
	default:
		err = fmt.Errorf("未支持的protocol=%v", job.conf.Protocol)
	}

	if err != nil {
		g = models.NewErrorCronLog(job.conf, string(res), err.Error(), st)
		job.ErrorCount++
	} else {
		g = models.NewSuccessCronLog(job.conf, string(res), st)
		job.ErrorCount = 0
	}
	// 连续错误达到5次，任务终止。
	if job.ErrorCount >= 5 || job.ErrorCount < 0 || job.conf.Type == models.TypeOnce {
		cronRun.Remove(cron.EntryID(job.conf.EntryId))
		job.conf.Status = models.ConfigStatusError
		job.conf.EntryId = 0
		err := data.NewCronConfigData(ctx).ChangeStatus(job.conf, "执行失败")
		if err != nil {
			log.Println("任务状态写入失败 id", job.conf.Id, err.Error())
		}
	}
}

// http 执行函数
func (job *CronJob) httpFunc(ctx context.Context) (res []byte, err error) {
	header := map[string]string{}
	for _, head := range job.commandParse.Http.Header {
		if head.Key == "" {
			continue
		}
		header[head.Key] = head.Value
	}
	method := models.ProtocolHttpMethodMap()[job.commandParse.Http.Method]
	if method == "" {
		// 任务设置有问题，提出执行队列，记录日志。
		job.ErrorCount = -2
		return nil, fmt.Errorf("未支持的http method，任务已终止。")
	}
	return job.httpRequest(ctx, method, job.commandParse.Http.Url, []byte(job.commandParse.Http.Body), header)
}

// rpc 执行函数
func (job *CronJob) rpcFunc(ctx context.Context) (res []byte, err error) {
	switch job.commandParse.Rpc.Method {
	case "GRPC":
		// 进行grpc处理
		// 目前还存在问题，无法通用性的提交和接收参数！
		return job.rpcGrpc(ctx, job.commandParse.Rpc)
	case "RPC":
		job.ErrorCount = -2
		return nil, fmt.Errorf("未支持的rpc method，任务已终止。")
		// 手头目前没有rpc的服务，不好测试验证。
	default:
		job.ErrorCount = -2
		return nil, fmt.Errorf("未支持的rpc method %s，任务已终止。", job.commandParse.Rpc.Method)
	}
}

// rpc 执行函数
func (job *CronJob) cmdFunc(ctx context.Context) (res []byte, err error) {
	if runtime.GOOS == "windows" {
		data := strings.Split(job.commandParse.Cmd, " ")
		if len(data) < 2 {
			return nil, errors.New("命令参数不合法，已跳过")
		}

		cmd := exec.Command(data[0], data[1:]...) // 合并 winds 命令
		if res, err = cmd.Output(); err != nil {
			return nil, err
		} else {
			srcCoder := mahonia.NewDecoder("gbk").ConvertString(string(res))
			return []byte(srcCoder), nil
		}

		// windows下安装了sh.exe 可以使用；就是git；但要添加环境变量 git/bin
		//cmd := exec.Command("sh.exe","-c", job.commandParse.Cmd)
		//if res, err = cmd.Output(); err != nil{
		//	return nil, err
		//}else {
		//	return res, nil
		//}

	} else {
		cmd := exec.Command("/bin/bash", "-c", job.commandParse.Cmd)
		if res, err = cmd.Output(); err != nil {
			return nil, err
		} else {
			return res, nil
		}
	}
}

// rpc 执行函数
func (job *CronJob) sqlFunc(ctx context.Context) (res []byte, err error) {
	switch job.commandParse.Sql.Driver {
	case models.SqlSourceMysql:
		return job.sqlMysql(ctx, job.commandParse.Sql)
	default:
		job.ErrorCount = -2
		return nil, fmt.Errorf("未支持的sql 驱动 %s，任务已终止。", job.commandParse.Sql.Driver)
	}
}

// http请求
func (job *CronJob) httpRequest(ctx context.Context, method, url string, body []byte, header map[string]string) (resp []byte, err error) {
	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("请求构建失败,%w", err)
	}
	for k, v := range header {
		req.Header.Set(k, v)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求执行失败，%w", err)
	}
	defer res.Body.Close()

	resp, err = io.ReadAll(res.Body)
	if err == nil && res.StatusCode != http.StatusOK {
		err = errors.New(fmt.Sprintf("%v %s", res.StatusCode, http.StatusText(res.StatusCode)))
	}
	return resp, err
}

// grpc调用
func (job *CronJob) rpcGrpc(ctx context.Context, r *pb.CronRpc) (resp []byte, err error) {
	cli, err := grpcurl.BlockingDial(ctx, "tcp", r.Addr, nil)
	if err != nil {
		return nil, fmt.Errorf("拨号目标主机 %s 失败:%w", r.Addr, err)
	}
	// 解析描述文件
	var descSource grpcurl.DescriptorSource
	if r.Proto != "" {
		fds, err := grpcurl.ParseProtoString(r.Proto)
		if err != nil {
			return nil, fmt.Errorf("无法解析给定的proto文件: %w", err)
		}
		descSource, err = grpcurl.DescriptorSourceFromFileDescriptors(fds...)
		if err != nil {
			return nil, fmt.Errorf("proto描述解析错误: %w", err)
		}
	} else { // 大部分服务器是不支持服务端的反射解析的
		md := grpcurl.MetadataFromHeaders(r.Header)
		refCtx := metadata.NewOutgoingContext(ctx, md)
		refClient := grpcreflect.NewClientAuto(refCtx, cli)
		descSource = grpcurl.DescriptorSourceFromServer(ctx, refClient)
	}

	var in io.Reader
	if r.Body == "" {
		in = os.Stdin
	} else {
		in = strings.NewReader(r.Body)
	}

	// 如果不是详细输出，那么还可以在每个消息之间包含记录分隔符，这样输出就可以通过管道输送到另一个grpcurl进程
	// 请求参数处理方法，把原json参数根据描述文件进行了转义处理。
	rf, formatter, err := grpcurl.RequestParserAndFormatter(grpcurl.Format("json"), descSource, in, grpcurl.FormatOptions{
		EmitJSONDefaultFields: true, // 是否json格式
		IncludeTextSeparator:  false,
		AllowUnknownFields:    true,
	})
	if err != nil {
		return nil, fmt.Errorf("请求解析器错误 %w", err)
	}

	h := grpcurl.NewMyEventHandler(formatter)
	// 发起请求
	err = grpcurl.InvokeRPC(ctx, descSource, cli, r.Action, r.Header, h, rf.Next)
	// 处理错误
	if err != nil {
		errStatus, _ := status.FromError(err)
		h.SetStatus(errStatus)
	}
	if h.GetStatus().Code() != codes.OK {
		err = fmt.Errorf("code:%v code_name:%v message:%v", int32(h.GetStatus().Code()), h.GetStatus().Code().String(), h.GetStatus().Message())
	}

	return []byte(h.RespMessages), err
}

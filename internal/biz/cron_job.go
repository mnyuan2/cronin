package biz

import (
	"bytes"
	"context"
	"cron/internal/data"
	"cron/internal/models"
	"cron/internal/pb"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"github.com/robfig/cron/v3"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"io/ioutil"
	"net/http"
	"time"
)

type CronJob struct {
	conf         *models.CronConfig
	commandParse *pb.CronConfigCommand
	cronId       cron.EntryID
	ErrorCount   int // 连续错误
}

// 任务执行器
func NewCronJob(conf *models.CronConfig) *CronJob {
	com := &pb.CronConfigCommand{}
	_ = jsoniter.UnmarshalFromString(conf.Command, com)

	return &CronJob{conf: conf, commandParse: com}
}

// 设置任务执行id
func (job *CronJob) SetCronId(cronId cron.EntryID) {
	job.cronId = cronId
}

// 返回任务执行中的id
func (job *CronJob) GetCronId() cron.EntryID {
	return job.cronId
}

func (job *CronJob) Run() {
	switch job.conf.Protocol {
	case models.ProtocolHttp:
		job.httpFunc()
	case models.ProtocolRpc:
		job.rpcFunc()
	case models.ProtocolCmd:
		job.cmdFunc()
	}
}

// http 执行函数
func (job *CronJob) httpFunc() {

	// 执行请求任务，并记录结果日志
	fmt.Println("执行http 任务")
	/*
		这里有三个点：1.请求类型、2.请求url、3.请求body；
			默认只能说是get请求的一个url
			最好的方案就是前段拼装成一个json
		任务连续失败三次，也应该终止；
		任务执行后，无论成功或失败，都要记录日志。
	*/
	startTime := time.Now()
	ctx := context.Background()
	g := &models.CronLog{}

	switch job.commandParse.Http.Method {
	case http.MethodPost:
		res, err := job.httpPost(ctx, job.commandParse.Http.Url, []byte(job.commandParse.Http.Body), nil)
		if err != nil {
			g = models.NewErrorCronLog(job.conf, err.Error(), startTime)
		} else {
			g = models.NewSuccessCronLog(job.conf, string(res), startTime)
		}

	case http.MethodGet:
		res, err := job.httpGet(ctx, job.commandParse.Http.Url, nil)
		if err != nil {
			g = models.NewErrorCronLog(job.conf, err.Error(), startTime)
		} else {
			g = models.NewSuccessCronLog(job.conf, string(res), startTime)
		}

	default:
		// 任务设置有问题，提出执行队列，记录日志。
		job.ErrorCount = -2
		g = models.NewErrorCronLog(job.conf, "未支持的http method，任务已终止。", startTime)
	}
	if g.Status == models.StatusDisable {
		job.ErrorCount++
	} else {
		job.ErrorCount = 0
	}
	if job.ErrorCount >= 5 || job.ErrorCount < 0 {
		jobList.Delete(job.conf.Id)
		cronRun.Remove(job.cronId)
	}

	data.NewCronLogData(ctx).Add(g)
}

// rpc 执行函数
func (job *CronJob) rpcFunc() {
	startTime := time.Now()
	ctx := context.Background()
	g := &models.CronLog{}

	switch job.commandParse.Rpc.Method {
	case "GRPC":
		// 进行grpc处理
		res, err := job.rpcGrpc(ctx, job.commandParse.Rpc.Addr, job.commandParse.Rpc.Action, job.commandParse.Rpc.Body)
		if err != nil {
			g = models.NewErrorCronLog(job.conf, err.Error(), startTime)
		} else {
			g = models.NewSuccessCronLog(job.conf, string(res), startTime)
		}

	case "RPC":
		job.ErrorCount = -2
		g = models.NewErrorCronLog(job.conf, "未支持的rpc method，任务已终止。", startTime)
		// 手头目前没有rpc的服务，不好测试验证。

	default:
		job.ErrorCount = -2
		g = models.NewErrorCronLog(job.conf, "未支持的rpc method，任务已终止。", startTime)
	}

	if g.Status == models.StatusDisable {
		job.ErrorCount++
	} else {
		job.ErrorCount = 0
	}
	if job.ErrorCount >= 5 || job.ErrorCount < 0 {
		jobList.Delete(job.conf.Id)
		cronRun.Remove(job.cronId)
	}

	data.NewCronLogData(ctx).Add(g)
}

// rpc 执行函数
func (job *CronJob) cmdFunc() {
	// 这个最后兼容
	fmt.Println("执行cmd 任务")
}

// get请求
func (job *CronJob) httpGet(ctx context.Context, url string, header http.Header) (resp []byte, err error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("请求构建失败,%w", err)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求执行失败，%w", err)
	}
	defer res.Body.Close()

	b, _ := ioutil.ReadAll(res.Body)
	return b, nil
}

// post请求
func (job *CronJob) httpPost(ctx context.Context, url string, body []byte, header http.Header) (resp []byte, err error) {
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("请求构建失败,%w", err)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求执行失败，%w", err)
	}
	defer res.Body.Close()

	resp, _ = ioutil.ReadAll(res.Body)
	return resp, nil
}

// grpc调用
func (job *CronJob) rpcGrpc(ctx context.Context, addr, action, param string) (resp []byte, err error) {
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("地址(%s)连接失败,%w", addr, err)
	}
	defer conn.Close()

	req := &models.GrpcRequest{}
	req.SetParam(param)
	res := &models.GrpcRequest{}

	err = conn.Invoke(ctx, action, req, resp)
	if err != nil {
		panic(fmt.Errorf("调用失败，%w", err))
		return nil, fmt.Errorf("%s 调用失败，%w", action, err)
	}

	return []byte(res.String()), nil
}

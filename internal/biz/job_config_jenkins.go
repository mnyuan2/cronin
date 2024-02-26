package biz

import (
	"bytes"
	"context"
	"cron/internal/basic/errs"
	"cron/internal/basic/tracing"
	"cron/internal/basic/util"
	"cron/internal/data"
	"cron/internal/pb"
	"crypto/tls"
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	jsoniter "github.com/json-iterator/go"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"io"
	"net/http"
	"strings"
	"time"
)

// 队列响应
type JenkinsQueueResponse struct {
	Executable *JenkinsQueueExecutable `json:"executable"`
}
type JenkinsQueueExecutable struct {
	Number int32 `json:"number"`
}

// 工作流程 响应
type JenkinsWorkflowResponse struct {
	DisplayName string `json:"displayName"`
	Result      string `json:"result"`
}

// mysql 命令执行
func (job *JobConfig) jenkins(ctx context.Context, r *pb.CronJenkins) (err errs.Errs) {
	ctx, span := job.tracer.Start(ctx, "exec-jenkins")
	defer func() {
		if err != nil {
			span.SetStatus(tracing.StatusError, err.Desc())
			span.AddEvent("执行错误", trace.WithAttributes(attribute.String("error.object", err.Error())))
		} else if er := util.PanicInfo(recover()); er != "" {
			span.SetStatus(tracing.StatusError, "执行异常")
			span.AddEvent("error", trace.WithAttributes(attribute.String("error.panic", er)))
		} else {
			span.SetStatus(tracing.StatusOk, "")
		}
		span.End()
	}()
	span.AddEvent("set", trace.WithAttributes(
		attribute.String("name", r.Name),
	))

	source, er := data.NewCronSettingData(ctx).GetSourceOne(job.conf.Env, r.Source.Id)
	if er != nil {
		return errs.New(er, "连接配置异常")
	}
	s := &pb.SettingSource{}
	if er = jsoniter.UnmarshalFromString(source.Content, s); er != nil {
		return errs.New(er, "连接配置解析异常")
	}

	/*
		1.执行构建
		2.循环查询结果，直到成功或失败
		3.完成进行下一步
	*/

	// 请求构建
	buildRes, er := job.httpJenkins(ctx, s.Jenkins, http.MethodPost, fmt.Sprintf("/job/%s/buildWithParameters", r.Name), r.Params)
	if er != nil {
		return errs.New(er, "构建失败")
	}

	dom, er := goquery.NewDocumentFromReader(bytes.NewReader(buildRes))
	if er != nil {
		return errs.New(er, "构建结果解析错误")
	}

	title := ""
	dom.Find("head title").Each(func(i int, selection *goquery.Selection) {
		title = selection.Text()
	})
	if title != "Error 404 Not Found" {
		return errs.New(errors.New(title), "构建错误")
	}
	// 解析响应并获得队列id
	buildData := map[string]string{}
	dom.Find("body table tr").Each(func(i int, selection *goquery.Selection) {
		th := selection.Find("th").Text()
		buildData[th[:len(th)-1]] = selection.Find("td").Text()
	})
	queueURI, ok := buildData["URI"]
	if !ok {
		return errs.New(errors.New(title), "构建队列错误")
	}
	queueParse := strings.Split(strings.Trim(queueURI, "/"), "/")
	queueId := queueParse[len(queueParse)-1]

	// 查询队列获取任务id
	queueRes, er := job.httpJenkins(ctx, s.Jenkins, http.MethodGet, fmt.Sprintf("/queue/item/%v/api/json", queueId), nil)
	if er != nil {
		return errs.New(er, "队列信息请求错误")
	}
	queueData := &JenkinsQueueResponse{Executable: &JenkinsQueueExecutable{}}
	if er := jsoniter.Unmarshal(queueRes, queueData); er != nil {
		return errs.New(er, "队列结果解析错误")
	}

	// 循环轮询任务，直到成功或失败
	for range time.Tick(time.Second * 5) {
		workflowRes, er := job.httpJenkins(ctx, s.Jenkins, http.MethodGet, fmt.Sprintf("/job/%s/%v/api/json", r.Name, queueData.Executable.Number), nil)
		if er != nil {
			return errs.New(er, "工作流程 请求错误")
		}

		workflowData := &JenkinsWorkflowResponse{}
		if er := jsoniter.Unmarshal(workflowRes, workflowData); er != nil {
			return errs.New(er, "工作流程 结果解析错误")
		}
		if workflowData.Result == "SUCCESS" {
			return nil
		}
		// 失败的标志后面看一下，进行中的标志同理
	}

	return err
}

// http请求
func (job *JobConfig) httpJenkins(ctx context.Context, source *pb.SettingJenkinsSource, method, path string, params []*pb.KvItem) (resp []byte, err errs.Errs) {
	ctx, span := job.tracer.Start(ctx, "http")
	defer func() {
		if err != nil {
			span.SetStatus(tracing.StatusError, err.Desc())
			span.AddEvent("error", trace.WithAttributes(attribute.String("error.object", err.Error())))
		}
		span.End()
	}()
	span.SetAttributes(
		attribute.String("component", "HTTP"),
		attribute.String("method", method),
	)

	// url参数
	paramStr := ""
	if len(params) > 0 {
		ps := []string{}
		for _, pram := range params {
			ps = append(ps, pram.Key+"="+pram.Value)
		}
		paramStr = "?" + strings.Join(ps, "&")
	}

	req, er := http.NewRequest(method, source.Hostname+path+paramStr, nil)
	if er != nil {
		return nil, errs.New(er, "请求构建失败")
	}
	// 用户
	if source.Username != "" {
		req.SetBasicAuth(source.Username, source.Password)
	}

	h, _ := jsoniter.Marshal(req.Header)
	span.AddEvent("", trace.WithAttributes(
		attribute.String("url", req.URL.String()),
		attribute.String("request_header", string(h)),
	))

	// 创建 HTTP 客户端
	client := &http.Client{
		Transport: &http.Transport{
			//MaxIdleConns:    10,
			//MaxConnsPerHost: 10,
			//IdleConnTimeout: 10 * time.Second,
			TLSClientConfig: &tls.Config{
				// 指定不校验 SSL/TLS 证书
				InsecureSkipVerify: true,
			},
		},
		//Timeout:   15 * time.Second,
	}

	res, er := client.Do(req)
	if er != nil {
		return nil, errs.New(er, "请求执行失败")
	}
	defer res.Body.Close()

	resp, er = io.ReadAll(res.Body)
	if er != nil {
		err = errs.New(er, "响应错误")
	}

	h, _ = jsoniter.Marshal(res.Header)
	span.AddEvent("", trace.WithAttributes(
		attribute.Int("status_code", res.StatusCode),
		attribute.String("response_header", string(h)),
		attribute.String("response", string(resp)),
	))
	return resp, err
}

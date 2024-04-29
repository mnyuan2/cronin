package biz

import (
	"bytes"
	"context"
	"cron/internal/basic/enum"
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

// 详情响应
type JenkinsDetailResponse struct {
	Class    string                   `json:"_class"`
	Property []*JenkinsDetailProperty `json:"property"`
}
type JenkinsDetailProperty struct {
	Class string `json:"_class"`
}

// 队列响应
type JenkinsQueueResponse struct {
	Why        *string                 `json:"why"`
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
			go job.messagePush(ctx, enum.StatusDisable, err.Desc(), nil, 0)
		} else if er := util.PanicInfo(recover()); er != "" {
			span.SetStatus(tracing.StatusError, "执行异常")
			span.AddEvent("error", trace.WithAttributes(attribute.String("error.panic", er)))
			go job.messagePush(ctx, enum.StatusDisable, "执行异常", []byte(er), 0)
		} else {
			span.SetStatus(tracing.StatusOk, "")
		}
		span.End()
	}()
	p, _ := jsoniter.MarshalToString(r.Params)
	span.AddEvent("set", trace.WithAttributes(
		attribute.String("name", r.Name),
		attribute.String("params", p),
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

	// 确定构建方式
	buildPath := "build"
	if len(r.Params) > 0 && r.Params[0].Key != "" {
		buildPath = "buildWithParameters"
	} else {
		b, er := job.httpJenkins(ctx, s.Jenkins, http.MethodPost, fmt.Sprintf("/job/%s/api/json?tree=property", r.Name), nil, "get-project")
		if er != nil {
			return errs.New(er, "项目信息有误")
		}
		detail := &JenkinsDetailResponse{}
		if er := jsoniter.Unmarshal(b, detail); er != nil {
			return errs.New(er, "详情解析错误")
		}
		for _, item := range detail.Property {
			if item.Class == "hudson.model.ParametersDefinitionProperty" { // 设置了参数就必须要有参构建，否则会报错
				buildPath = "buildWithParameters"
			}
		}
	}

	// 请求构建
	queueId, er := job.httpJenkins(ctx, s.Jenkins, http.MethodPost, fmt.Sprintf("/job/%s/%s", r.Name, buildPath), r.Params, "exec-build")
	if er != nil {
		return errs.New(er, "构建失败")
	}

	// 查询队列获取任务id；最大尝试次数10
	queueData := &JenkinsQueueResponse{Executable: &JenkinsQueueExecutable{}}
	for index := 0; index < 10; index++ {
		time.Sleep(time.Second * 10)
		queueRes, er := job.httpJenkins(ctx, s.Jenkins, http.MethodGet, fmt.Sprintf("/queue/item/%v/api/json", string(queueId)), nil, "get-queue")
		if er != nil {
			return errs.New(er, "队列信息请求错误")
		}

		if er := jsoniter.Unmarshal(queueRes, queueData); er != nil {
			return errs.New(er, "队列结果解析错误")
		}
		if queueData.Executable.Number > 0 {
			break
		}
	}
	if queueData.Executable.Number == 0 {
		return errs.New(er, "工作流程 编号获取失败")
	}

	// 循环轮询任务，直到成功或失败；这里是真的可能很久。
	for range time.Tick(time.Second * 5) {
		workflowRes, er := job.httpJenkins(ctx, s.Jenkins, http.MethodGet, fmt.Sprintf("/job/%s/%v/api/json", r.Name, queueData.Executable.Number), nil, "find-workflow")
		if er != nil {
			return errs.New(er, "工作流程 请求错误")
		}

		workflowData := &JenkinsWorkflowResponse{}
		if er := jsoniter.Unmarshal(workflowRes, workflowData); er != nil {
			return errs.New(er, "工作流程 结果解析错误")
		}
		if workflowData.Result == "SUCCESS" {
			return nil
		} else if workflowData.Result == "FAILURE" {
			return errs.New(nil, "构建失败 FAILURE")
		}
		// 进行中就延迟后 再次查询检测
	}

	return nil
}

// http请求
func (job *JobConfig) httpJenkins(ctx context.Context, source *pb.SettingJenkinsSource, method, path string, params []*pb.KvItem, operateName string) (resp []byte, err errs.Errs) {
	ctx, span := job.tracer.Start(ctx, operateName)
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
			if pram.Key == "" {
				continue
			}
			ps = append(ps, pram.Key+"="+pram.Value)
		}
		if len(ps) > 0 {
			paramStr = "?" + strings.Join(ps, "&")
		}
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

	if operateName == "exec-build" {
		resp, er = job.jenkinsBuildParse(span, res)
	} else {
		resp, er = io.ReadAll(res.Body)
		if er != nil {
			err = errs.New(er, "响应错误")
		}
	}

	h, _ = jsoniter.Marshal(res.Header)
	span.AddEvent("", trace.WithAttributes(
		attribute.Int("status_code", res.StatusCode),
		attribute.String("response_header", string(h)),
		attribute.String("response", string(resp)),
	))
	return resp, err
}

// 构建解析
func (job *JobConfig) jenkinsBuildParse(span trace.Span, res *http.Response) (resp []byte, err error) {
	if res.StatusCode == http.StatusCreated {
		//res.Header.Get("Location")
		// https://jenkins.xxx.cn/queue/item/34170/
		queueParse := strings.Split(strings.Trim(res.Header.Get("Location"), "/"), "/")
		queueId := queueParse[len(queueParse)-1]

		return []byte(queueId), nil

	} else if res.StatusCode == http.StatusNotFound {
		resp, err = io.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}
		span.AddEvent("", trace.WithAttributes(
			attribute.String("response_body", string(resp)),
		))

		dom, er := goquery.NewDocumentFromReader(bytes.NewReader(resp))
		if er != nil {
			return nil, errs.New(er, "构建结果解析错误")
		}

		title := ""
		dom.Find("head title").Each(func(i int, selection *goquery.Selection) {
			title = selection.Text()
		})
		if title != "Error 404 Not Found" {
			return nil, errs.New(errors.New(title), "构建错误")
		}
		// 解析响应并获得队列id
		buildData := map[string]string{}
		dom.Find("body table tr").Each(func(i int, selection *goquery.Selection) {
			th := selection.Find("th").Text()
			buildData[th[:len(th)-1]] = selection.Find("td").Text()
		})
		queueURI, ok := buildData["URI"]
		if !ok {
			return nil, errs.New(errors.New(title), "构建队列错误")
		}
		queueParse := strings.Split(strings.Trim(queueURI, "/"), "/")
		queueId := queueParse[len(queueParse)-1]

		return []byte(queueId), nil
	}

	return nil, errs.New(nil, "错误")
}

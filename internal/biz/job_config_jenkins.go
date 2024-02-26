package biz

import (
	"context"
	"cron/internal/basic/errs"
	"cron/internal/basic/tracing"
	"cron/internal/basic/util"
	"cron/internal/data"
	"cron/internal/pb"
	"encoding/base64"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

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
	header := map[string]string{
		"Authorization": "Basic " + base64.StdEncoding.EncodeToString([]byte(s.Jenkins.Username+":"+s.Jenkins.Password)),
	}

	// 请求构建
	url := fmt.Sprintf("%s/job/%s/buildWithParameters", s.Jenkins.Hostname, r.Name)
	res1, er := job.httpRequest(ctx, "POST", url, nil, header)
	if er != nil {
		return errs.New(er, "构建失败")
	}
	fmt.Println(string(res1))

	return err
}

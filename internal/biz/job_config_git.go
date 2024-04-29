package biz

import (
	"context"
	"cron/internal/basic/enum"
	"cron/internal/basic/errs"
	"cron/internal/basic/git/gitee"
	"cron/internal/basic/tracing"
	"cron/internal/biz/dtos"
	"cron/internal/data"
	"cron/internal/pb"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

func (job *JobConfig) gitFunc(ctx context.Context, r *pb.CronGit) (resp []byte, err errs.Errs) {
	link, er := data.NewCronSettingData(ctx).GetSourceOne(job.conf.Env, r.LinkId)
	if er != nil {
		return nil, errs.New(er, "链接配置查询错误")
	}
	conf := &pb.SettingSource{
		Git: &pb.SettingGitSource{}, // 通过token来确定调用类型：github、gitee、其它
	}
	if er := jsoniter.UnmarshalFromString(link.Content, conf); er != nil {
		return nil, errs.New(er, "链接配置解析错误")
	}

	api := gitee.NewApiV5(conf.Git)

	for i, e := range r.Events {
		switch e.Id {
		case enum.GitEventPullsMerge:
			resp, err = job.PRMerge(ctx, api, e.PRMerge)
		default:
			return nil, errs.New(nil, fmt.Sprintf("未支持的任务 %v-%v", i, e.Id))
		}
	}

	return resp, err
}

// git 抓取文件数据
func (job *JobConfig) getGitFile(ctx context.Context, r *pb.Git) (flies []*dtos.File, err errs.Errs) {
	link, er := data.NewCronSettingData(ctx).GetSourceOne(job.conf.Env, r.LinkId)
	if er != nil {
		return nil, errs.New(er, "链接配置查询错误")
	}
	conf := &pb.SettingSource{
		Git: &pb.SettingGitSource{},
	}
	if er := jsoniter.UnmarshalFromString(link.Content, conf); er != nil {
		return nil, errs.New(er, "链接配置解析错误")
	}

	api := gitee.NewApiV5(conf.Git)
	flies = []*dtos.File{}
	for _, path := range r.Path {
		if path == "" {
			continue
		}
		file, err := job.gitReposContents(ctx, api, r, path)
		if err != nil {
			return nil, err
		}
		flies = append(flies, &dtos.File{Name: path, Byte: file})
	}

	return flies, nil
}

func (job *JobConfig) gitReposContents(ctx context.Context, api *gitee.ApiV5, r *pb.Git, path string) (file []byte, err errs.Errs) {
	h := gitee.NewHandler(ctx)
	ctx, span := job.tracer.Start(ctx, "repos-contents")
	defer func() {
		if err != nil {
			span.SetStatus(tracing.StatusError, err.Desc())
			span.AddEvent("error", trace.WithAttributes(attribute.String("error.object", err.Error())))
		}
		span.SetAttributes(
			attribute.String("component", "HTTP"),
			attribute.String("method", h.General.Method),
		)
		span.AddEvent("", trace.WithAttributes(
			attribute.String("url", h.General.Url),
			attribute.Int("status_code", h.General.StatusCode),
			attribute.String("response_header", string(h.ResponseHeaderBytes())),
			attribute.String("response_body", string(h.ResponseBody)),
		))
		span.End()
	}()

	res, er := api.ReposContents(h, r.Owner, r.Project, path, r.Ref)
	if er != nil {
		return nil, errs.New(er, "gite文件获取失败")
	}
	span.AddEvent("", trace.WithAttributes(attribute.String("response", string(res))))

	return res, nil
}

// 记录日志
func (job *JobConfig) handlerLog(name string, h *gitee.Handler, err errs.Errs) {
	_, span := job.tracer.Start(h.GetContext(), name, trace.WithTimestamp(h.StartTime()))
	span.SetAttributes(
		attribute.String("component", "HTTP"),
		attribute.String("method", h.General.Method),
	)
	span.AddEvent("", trace.WithAttributes(
		attribute.String("url", h.General.Url),
		attribute.String("body", string(h.RequestBody)),
		attribute.Int("status_code", h.General.StatusCode),
		attribute.String("response_header", string(h.ResponseHeaderBytes())),
		attribute.String("response_body", string(h.ResponseBody)),
	))
	if err != nil {
		span.SetStatus(tracing.StatusError, err.Desc())
		span.AddEvent("error", trace.WithAttributes(attribute.String("error.object", err.Error())))
	}

	span.End()
}

// pr 合并
func (job *JobConfig) PRMerge(ctx context.Context, api *gitee.ApiV5, r *pb.GitEventPRMerge) (resp []byte, err errs.Errs) {
	h := gitee.NewHandler(ctx)
	defer func() {
		job.handlerLog("PRMerge", h, err)
	}()

	request := &gitee.PullsMergeRequest{
		BaseRequest: gitee.BaseRequest{
			Owner: r.Owner,
			Repo:  r.Repo,
		},
		Number:            r.Number,
		MergeMethod:       r.MergeMethod,
		PruneSourceBranch: r.PruneSourceBranch,
		Title:             r.Title,
		Description:       r.Description,
	}
	res, er := api.PullsMerge(h, request)
	if er != nil {
		return nil, errs.New(er, "pr合并失败")
	}
	return res, nil
}

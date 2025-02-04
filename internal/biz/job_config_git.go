package biz

import (
	"context"
	"cron/internal/basic/conv"
	"cron/internal/basic/enum"
	"cron/internal/basic/errs"
	"cron/internal/basic/git"
	"cron/internal/basic/tracing"
	"cron/internal/biz/dtos"
	"cron/internal/data"
	"cron/internal/pb"
	"encoding/base64"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"strconv"
	"strings"
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

	api := git.NewApi(git.Config{
		Driver:      conf.Git.Driver,
		AccessToken: conf.Git.AccessToken,
	})

	for i, e := range r.Events {
		switch e.Id {
		case enum.GitEventPullsDetail:
			resp, err = job.PRDetail(ctx, api, e.PRDetail)
		case enum.GitEventPullsIsMerge:
			resp, err = job.PRIsMerge(ctx, api, e.PRIsMerge)
		case enum.GitEventPullsMerge:
			resp, err = job.PRMerge(ctx, api, e.PRMerge)
		case enum.GitEventFileUpdate:
			resp, err = job.FileUpdate(ctx, api, e.FileUpdate)
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

	api := git.NewApi(git.Config{
		Driver:      conf.Git.Driver,
		AccessToken: conf.Git.AccessToken,
	})
	flies = []*dtos.File{}
	for _, path := range r.Path {
		list := strings.Split(path, ",")
		for _, item := range list {
			item = strings.Trim(strings.TrimSpace(item), "/")
			if item == "" {
				continue
			}
			file, err := job.gitReposContents(ctx, api, r, item)
			if err != nil {
				return nil, err
			}
			flies = append(flies, &dtos.File{Name: item, Byte: file})
		}
	}

	return flies, nil
}

// 获取文件信息
func (job *JobConfig) gitReposContents(ctx context.Context, api git.Api, r *pb.Git, path string) (file []byte, err errs.Errs) {
	h := git.NewHandler(ctx)
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

	res, er := api.FileGet(h, &git.FileGetRequest{
		BaseRequest: git.BaseRequest{
			Owner: r.Owner,
			Repo:  r.Project,
		},
		Path: path,
		Ref:  r.Ref,
	})
	if er != nil {
		return nil, errs.New(er, "gite文件获取失败")
	}
	span.AddEvent("", trace.WithAttributes(attribute.String("response", string(res.Content))))

	return file, nil
}

// 记录日志
func (job *JobConfig) handlerLog(name string, h *git.Handler, err errs.Errs) {
	if h == nil {
		return
	}
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

	span.End(trace.WithTimestamp(h.EndTime()))
}

// pr 列表查询
func (job *JobConfig) PRList(ctx context.Context, api git.Api, r *pb.GetEventPRList) (resp []byte, err errs.Errs) {
	h := git.NewHandler(ctx)
	defer func() {
		job.handlerLog("PRList", h, err)
	}()

	if r.Owner == "" || r.Repo == "" {
		return nil, errs.New(nil, "必填参数不足")
	}

	request := &git.Pulls{
		BaseRequest: git.BaseRequest{
			Owner: r.Owner,
			Repo:  r.Repo,
		},
		State:   r.State,
		Head:    r.Head,
		Base:    r.Base,
		Page:    r.Page,
		PerPage: r.PerPage,
	}
	res, er := api.Pulls(h, request)
	if er != nil {
		return nil, errs.New(er)
	}
	resp, _ = jsoniter.Marshal(res)

	return resp, nil
}

// pr 详情
func (job *JobConfig) PRDetail(ctx context.Context, api git.Api, r *pb.GitEventPRMerge) (resp []byte, err errs.Errs) {
	h := git.NewHandler(ctx)
	defer func() {
		job.handlerLog("PRDetail", h, err)
	}()
	num, er := strconv.Atoi(r.Number)
	if er != nil {
		return nil, errs.New(er, "pr编号输入有误")
	}
	if r.Owner == "" || r.Repo == "" || num == 0 {
		return nil, errs.New(nil, "必填参数不足")
	}

	request := &git.PullsMergeRequest{
		BaseRequest: git.BaseRequest{
			Owner: r.Owner,
			Repo:  r.Repo,
		},
		Number: int32(num),
	}
	res, er := api.PullGet(h, request)
	if er != nil {
		return nil, errs.New(er, "pr详情查询失败")
	}
	resp, _ = jsoniter.Marshal(res)
	return resp, nil
}

// pr 合并校验
func (job *JobConfig) PRIsMerge(ctx context.Context, api git.Api, r *pb.GitEventPRMerge) (resp []byte, err errs.Errs) {
	h := git.NewHandler(ctx)
	defer func() {
		job.handlerLog("PRIsMerge", h, err)
	}()

	num, er := strconv.Atoi(r.Number)
	if er != nil {
		return nil, errs.New(er, "pr编号输入有误")
	}
	if r.Owner == "" || r.Repo == "" || num == 0 {
		return nil, errs.New(nil, "必填参数不足")
	}
	if r.State == "" {
		r.State = "merge"
	} else if r.State != "merge" && r.State != "open" {
		return []byte(""), errs.New(nil, "不支持的合并类型")
	}

	request := &git.PullsMergeRequest{
		BaseRequest: git.BaseRequest{
			Owner: r.Owner,
			Repo:  r.Repo,
		},
		Number: int32(num),
		State:  r.State,
	}
	er = api.PullsIsMerge(h, request)
	if er != nil {
		return []byte(""), errs.New(er)
	}
	return []byte(""), nil
}

// pr 合并
func (job *JobConfig) PRMerge(ctx context.Context, api git.Api, r *pb.GitEventPRMerge) (resp []byte, err errs.Errs) {
	h := git.NewHandler(ctx)
	defer func() {
		job.handlerLog("PRMerge", h, err)
	}()

	num, er := strconv.Atoi(r.Number)
	if er != nil {
		return nil, errs.New(er, "pr编号输入有误")
	}

	request := &git.PullsMergeRequest{
		BaseRequest: git.BaseRequest{
			Owner: r.Owner,
			Repo:  r.Repo,
		},
		Number:            int32(num),
		MergeMethod:       r.MergeMethod,
		PruneSourceBranch: r.PruneSourceBranch,
		Title:             r.Title,
		Description:       r.Description,
	}
	res, er := api.PullMerge(h, request)
	if er != nil {
		return []byte(res.Url), errs.New(er, "pr合并失败")
	}
	return []byte(res.Url), nil
}

// 文件 更新
func (job *JobConfig) FileUpdate(ctx context.Context, api git.Api, r *pb.GitEventFileUpdate) (resp []byte, err errs.Errs) {
	h1 := git.NewHandler(ctx)
	h2 := git.NewHandler(ctx)
	defer func() {
		job.handlerLog("FileGet", h1, err)
		job.handlerLog("FileUpdate", h2, err)
	}()

	// 获取原文件信息
	res1, er := api.FileGet(h1, &git.FileGetRequest{
		BaseRequest: git.BaseRequest{
			Owner: r.Owner,
			Repo:  r.Repo,
		},
		Path: r.Path,
		Ref:  r.Branch,
	})
	if er != nil {
		return nil, errs.New(er, "文件获取错误")
	}

	// 对内容支持模板解析
	inContent, _ := base64.StdEncoding.DecodeString(r.Content)
	p := map[string]any{}
	for k, v := range job.varParams {
		p[k] = v
	}
	p["raw_content"] = res1.Content
	content, er := conv.DefaultStringTemplate().SetParam(p).Execute(inContent)
	if er != nil {
		return nil, errs.New(er, "内容模板错误")
	}

	// 更新文件信息
	res2, er := api.FileUpdate(h2, &git.FileUpdateRequest{
		BaseRequest: git.BaseRequest{
			Owner: r.Owner,
			Repo:  r.Repo,
		},
		Path:    r.Path,
		Content: string(content),
		Sha:     res1.Sha,
		Message: r.Message,
		Branch:  r.Branch,
	})
	if er != nil {
		return nil, errs.New(er, "文件更新错误")
	}
	return []byte(res2.Commit.CommitUrl), nil
}

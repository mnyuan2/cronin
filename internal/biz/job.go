package biz

import (
	"context"
	"cron/internal/basic/auth"
	"cron/internal/basic/errs"
	"cron/internal/models"
	"cron/internal/pb"
	"errors"
	"github.com/robfig/cron/v3"
)

type JobService struct {
	ctx  context.Context
	user *auth.UserToken
}

func NewJobService(ctx context.Context, user *auth.UserToken) *JobService {
	return &JobService{
		ctx:  ctx,
		user: user,
	}
}

// 停止进行中的任务
func (dm *JobService) Stop(r *pb.JobStopRequest) (resp *pb.JobStopReply, err error) {
	if r.RefId == 0 || r.EntryId == 0 {
		return nil, errors.New("参数未传递")
	}
	v := cronRun.Entry(cron.EntryID(r.EntryId))
	switch v.Job.(type) {
	case *JobConfig:
		job := v.Job.(*JobConfig)
		if job.conf.Id != r.RefId {
			return nil, errs.New(nil, "注册任务信息不匹配")
		}
		err = job.Stop(dm.ctx, dm.user, "手动停止")
	case *JobPipeline:
		job := v.Job.(*JobPipeline)
		if job.pipeline.Id != r.RefId {
			return nil, errs.New(nil, "注册任务信息不匹配")
		}
		err = job.conf.Stop(dm.ctx, dm.user, "手动停止")
	case *JobReceive:
		job := v.Job.(*JobReceive)
		if job.conf.conf.EntryId != r.EntryId {
			return nil, errs.New(nil, "注册任务信息不匹配")
		}
		err = job.conf.Stop(dm.ctx, dm.user, "手动停止")
	case nil:
		err = errs.New(nil, "任务未执行")
	default:
		err = errs.New(nil, "任务类型异常，请联系管理员")
	}

	return &pb.JobStopReply{}, err
}

// 日志踪迹
func (dm *JobService) Traces(r *pb.JobTracesRequest) (resp *pb.JobTracesResponse, err error) {
	if r.RefId == 0 || r.EntryId == 0 {
		return nil, errors.New("参数未传递")
	}
	if r.TraceId == "" {
		return nil, errors.New("未指定踪迹id")
	}

	// 树
	resp = &pb.JobTracesResponse{
		List:  []*pb.CronLogTraceItem{},
		Limit: 1000,
	}

	var list []*models.CronLogSpan
	v := cronRun.Entry(cron.EntryID(r.EntryId))
	switch v.Job.(type) {
	case *JobConfig:
		job := v.Job.(*JobConfig)
		list, err = job.Logs(r.TraceId)

	case *JobPipeline:
		job := v.Job.(*JobPipeline)
		list, err = job.GetConf().Logs(r.TraceId)
	case *JobReceive:
		job := v.Job.(*JobReceive)
		list, err = job.GetConf().Logs(r.TraceId)
	default:
		err = errs.New(nil, "任务类型异常，请联系管理员")
	}

	if err != nil && err.Error() == "非执行中任务" {
		re, err := NewCronLogService(dm.ctx, dm.user).Trace(&pb.CronLogTraceRequest{TraceId: r.TraceId})
		if err != nil {
			return nil, err
		}
		resp.List = re.List
		return resp, nil
	}
	if err != nil {
		return nil, err
	}

	resp.Total = len(list)
	if resp.Total == 0 {
		return resp, nil
	}
	tra := &pb.CronLogTraceItem{
		TraceId: list[0].TraceId,
		Spans:   []*pb.CronLogSpan{},
	}
	gs := NewCronLogService(dm.ctx, dm.user)
	for i, item := range list {
		span := gs.toOut(item)
		tra.Spans = append(tra.Spans, span)
		if i > resp.Limit {
			break
		}
	}
	resp.List = append(resp.List, tra)
	return resp, err
}

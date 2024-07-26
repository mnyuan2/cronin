package biz

import (
	"context"
	"cron/internal/basic/auth"
	"cron/internal/basic/errs"
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
		err = job.Stop(dm.ctx, "手动停止")
	case *JobPipeline:
		job := v.Job.(*JobPipeline)
		if job.pipeline.Id != r.RefId {
			return nil, errs.New(nil, "注册任务信息不匹配")
		}
		err = job.conf.Stop(dm.ctx, "手动停止")
	case nil:
		err = errs.New(nil, "任务未执行")
	default:
		err = errs.New(nil, "任务类型异常，请联系管理员")
	}

	return &pb.JobStopReply{}, err
}

package biz

import (
	"context"
	"cron/internal/basic/db"
	"cron/internal/basic/enum"
	"cron/internal/basic/errs"
	"cron/internal/basic/tracing"
	"cron/internal/basic/util"
	"cron/internal/biz/dtos"
	"cron/internal/data"
	"cron/internal/models"
	"fmt"
	"github.com/robfig/cron/v3"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"time"
)

type JobPipeline struct {
	pipeline    *models.CronPipeline
	conf        *JobConfig
	msgSetParse *dtos.MsgSetParse
	tracer      trace.Tracer
}

// 任务执行器
func NewJobPipeline(conf *models.CronPipeline) *JobPipeline {
	job := &JobPipeline{
		pipeline: conf,
	}
	job.conf = NewJobConfig(&models.CronConfig{
		Env:          conf.Env,
		EntryId:      0,
		Type:         0,
		Name:         conf.Name,
		Spec:         "",
		Protocol:     0,
		Command:      nil,
		Remark:       "",
		Status:       0,
		StatusRemark: "",
		StatusDt:     "",
		UpdateDt:     "",
		CreateDt:     "",
		MsgSet:       conf.MsgSet,
	})

	// 日志
	job.tracer = tracing.Tracer(job.pipeline.Env+"-cronin", trace.WithInstrumentationAttributes(
		attribute.String("driver", "mysql"),
		attribute.String("env", job.pipeline.Env),
	))

	return job
}

// 执行任务
func (job *JobPipeline) Run() {
	var err errs.Errs
	//var res []byte
	st := time.Now()
	ctx, span := job.tracer.Start(context.Background(), "job-pipeline", trace.WithAttributes(attribute.Int("ref_id", job.pipeline.Id)))
	defer func() {
		//if res != nil {
		//	span.AddEvent("", trace.WithAttributes(attribute.String("resp", string(res))))
		//}
		if err != nil {
			span.SetStatus(tracing.StatusError, err.Desc())
			span.AddEvent("error", trace.WithAttributes(attribute.String("error.object", err.Error())))
		} else {
			span.SetStatus(tracing.StatusOk, "")
		}
		if er := util.PanicInfo(recover()); er != "" {
			span.SetStatus(tracing.StatusError, "执行异常")
			span.AddEvent("error", trace.WithAttributes(attribute.String("error.panic", er)))
		}
		span.End()
	}()
	span.SetAttributes(
		attribute.String("env", job.pipeline.Env),
		attribute.String("component", "pipeline"),
	)

	e := cronRun.Entry(cron.EntryID(job.pipeline.EntryId))
	cronRun.Remove(cron.EntryID(job.pipeline.EntryId))
	if e.ID == 0 {
		span.SetStatus(tracing.StatusError, "重复执行？")
		return
	}

	// 欢迎语
	job.conf.messagePush(ctx, 0, "开始", nil, 0)
	fmt.Println(ctx, st)

	w := db.NewWhere().Eq("env", job.pipeline.Env).In("id", job.pipeline.ConfigIds)
	list, er := data.NewCronConfigData(ctx).List(w, len(job.pipeline.ConfigIds))
	if er != nil {
		err = errs.New(er, "任务查询错误")
		return
	}
	jobs := []*JobConfig{}
	for _, item := range list {
		if item.Status != models.ConfigStatusActive && item.Status != models.ConfigStatusFinish {
			if job.pipeline.ConfigDisableAction == models.DisableActionStop {
				job.conf.messagePush(ctx, enum.StatusDisable, "任务非激活", []byte(fmt.Sprintf("%s-%s", item.Name, item.GetStatusName())), time.Since(st).Seconds())
				return
			} else if job.pipeline.ConfigDisableAction == models.DisableActionOmit {
				continue
			}
		}
		jobs = append(jobs, NewJobConfig(item))
	}

	for _, j := range jobs {
		_, er := j.run(ctx)
		if er != nil {
			job.conf.messagePush(ctx, enum.StatusDisable, er.Desc()+" 流水线"+job.pipeline.ConfigErrActionName(), []byte(err.Error()), time.Since(st).Seconds())
			// 这里要确认一下是否继续执行下去。
			if job.pipeline.ConfigErrAction == models.ErrActionStop {
				return
			}
		}
	}
	// 结束语
	job.conf.messagePush(ctx, enum.StatusActive, "完成", nil, time.Since(st).Seconds())

}

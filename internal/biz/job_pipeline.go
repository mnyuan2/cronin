package biz

import (
	"context"
	"cron/internal/basic/errs"
	"cron/internal/basic/tracing"
	"cron/internal/basic/util"
	"cron/internal/models"
	"fmt"
	"github.com/robfig/cron/v3"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"time"
)

type JobPipeline struct {
	pipeline   *models.CronPipeline
	ErrorCount int // 连续错误
	tracer     trace.Tracer
}

// 任务执行器
func NewJobPipeline(conf *models.CronPipeline) *JobPipeline {
	job := &JobPipeline{
		pipeline: conf,
	}

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
	var res []byte
	st := time.Now()
	ctx, span := job.tracer.Start(context.Background(), "job-pipeline", trace.WithAttributes(attribute.Int("ref_id", job.pipeline.Id)))
	defer func() {
		if res != nil {
			span.AddEvent("", trace.WithAttributes(attribute.String("resp", string(res))))
		}
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

	fmt.Println(ctx, st)

	//w := db.NewWhere().Eq("env", job.pipeline.Env).In("id", job.pipeline.ConfigIds)
	//list, er := data.NewCronConfigData(ctx).List(w, len(job.pipeline.ConfigIds))
	//if er != nil {
	//	err = errs.New(er, "任务查询错误")
	//	return
	//}
	//jobs := []*JobConfig{}

	//res, err = job.Exec(ctx)
	//if err != nil {
	//	job.ErrorCount++
	//	go job.messagePush(ctx, enum.StatusDisable, err.Desc(), []byte(err.Error()), time.Since(st).Seconds())
	//} else {
	//	job.ErrorCount = 0
	//	go job.messagePush(ctx, enum.StatusActive, "ok", res, time.Since(st).Seconds())
	//}

}

//func (job *JobPipeline) Exec(ctx context.Context) (res []byte, err errs.Errs) {
//	switch job.conf.Protocol {
//	case models.ProtocolHttp:
//		res, err = job.httpFunc(ctx, job.commandParse.Http)
//	case models.ProtocolRpc:
//		res, err = job.rpcFunc(ctx)
//	case models.ProtocolCmd:
//		res, err = job.cmdFunc(ctx)
//	case models.ProtocolSql:
//		err = job.sqlFunc(ctx)
//	default:
//		err = errs.New(nil, fmt.Sprintf("未支持的protocol=%v", job.conf.Protocol))
//	}
//	return res, err
//}

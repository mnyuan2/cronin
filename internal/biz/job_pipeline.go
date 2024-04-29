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
	"cron/internal/pb"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"github.com/robfig/cron/v3"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"log"
	"time"
)

type JobPipeline struct {
	pipeline    *models.CronPipeline
	confs       []*pb.CronConfigListItem // 配置任务集合
	conf        *JobConfig               // 流水线解析后的本身配置
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
		EntryId:      conf.EntryId,
		Type:         0,
		Name:         conf.Name,
		Spec:         conf.Spec,
		Protocol:     99,
		Command:      nil,
		Remark:       "",
		Status:       conf.Status,
		StatusRemark: "",
		StatusDt:     "",
		UpdateDt:     "",
		CreateDt:     "",
		MsgSet:       conf.MsgSet,
	})
	if err := job.parse(conf); err != nil {
		log.Println("流水线配置解析错误", err.Error())
		// ...
	}
	job.conf.Parse(nil)

	// 日志
	job.tracer = tracing.Tracer(job.pipeline.Env+"-cronin", trace.WithInstrumentationAttributes(
		attribute.String("driver", "mysql"),
		attribute.String("env", job.pipeline.Env),
	))

	return job
}

func (job *JobPipeline) parse(conf *models.CronPipeline) error {
	job.confs = []*pb.CronConfigListItem{}
	return jsoniter.Unmarshal(conf.Configs, &job.confs)
}

// 执行任务
func (job *JobPipeline) Run() {
	var err errs.Errs
	//var res []byte
	st := time.Now()
	ctx, span := job.tracer.Start(context.Background(), "job-pipeline", trace.WithAttributes(attribute.Int("ref_id", job.pipeline.Id)))
	defer func() {
		status, remark := 0, ""
		//if res != nil {
		//	span.AddEvent("", trace.WithAttributes(attribute.String("resp", string(res))))
		//}
		if err != nil {
			span.SetStatus(tracing.StatusError, err.Desc())
			span.AddEvent("error", trace.WithAttributes(attribute.String("error.object", err.Error())))
			status, remark = models.ConfigStatusError, err.Desc()
		} else if er := util.PanicInfo(recover()); er != "" {
			span.SetStatus(tracing.StatusError, "执行异常")
			span.AddEvent("error", trace.WithAttributes(attribute.String("error.panic", er)))
			status, remark = models.ConfigStatusError, "执行异常"
		} else {
			span.SetStatus(tracing.StatusOk, "")
			status, remark = models.ConfigStatusFinish, "SUCCESS"
		}

		job.pipeline.Status = status
		if er := data.NewCronPipelineData(ctx).ChangeStatus(job.pipeline, remark); er != nil {
			log.Println(attribute.String("error.object", "完成状态写入失败"+er.Error()))
		}
		span.End()
	}()
	span.SetAttributes(
		attribute.String("env", job.pipeline.Env),
		attribute.String("component", "pipeline"),
	)

	e := cronRun.Entry(cron.EntryID(job.pipeline.EntryId))
	cronRun.Remove(cron.EntryID(job.pipeline.EntryId))
	job.pipeline.EntryId = 0
	if e.ID == 0 {
		span.SetStatus(tracing.StatusError, "重复执行？")
		return
	}

	// 欢迎语
	job.conf.messagePush(ctx, 0, "开始", nil, 0)
	//fmt.Println(ctx, st)
	configIds := []int{}
	if er := jsoniter.Unmarshal(job.pipeline.ConfigIds, &configIds); er != nil {
		err = errs.New(er, "configIds 序列化错误")
		return
	}
	if len(configIds) == 0 {
		err = errs.New(nil, "未配置任务")
		return
	}

	w := db.NewWhere().Eq("env", job.pipeline.Env).In("id", configIds)
	list, er := data.NewCronConfigData(ctx).List(w, len(job.pipeline.ConfigIds))
	if er != nil {
		err = errs.New(er, "任务查询错误")
		return
	}
	listMap := map[int]*models.CronConfig{}
	for _, item := range list {
		listMap[item.Id] = item
	}

	jobs := []*JobConfig{}
	for _, item := range job.confs { // 要保持设置时的顺序关系
		temp, ok := listMap[item.Id]
		if !ok {
			err = errs.New(fmt.Errorf("%v·%v 未找到匹配", item.Id, item.Name), "任务配置异常")
			return
		}

		if item.Status != models.ConfigStatusActive && item.Status != models.ConfigStatusFinish {
			if job.pipeline.ConfigDisableAction == models.DisableActionStop {
				job.conf.messagePush(ctx, enum.StatusDisable, "任务非激活", []byte(fmt.Sprintf("%s-%s", item.Name, temp.GetStatusName())), time.Since(st).Seconds())
				return
			} else if job.pipeline.ConfigDisableAction == models.DisableActionOmit {
				continue
			}
		}
		jobs = append(jobs, NewJobConfig(temp))
	}

	// 参数解析
	varParams := map[string]any{}
	if job.pipeline.VarParams != "" {
		if er := jsoniter.UnmarshalFromString(job.pipeline.VarParams, &varParams); err != nil {
			err = errs.New(er, "参数解析失败")
			return
		}
	}

	for _, j := range jobs {
		_, er := j.Running(ctx, "流水线执行", varParams)
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

func (job *JobPipeline) GetConf() *JobConfig {
	return job.conf
}

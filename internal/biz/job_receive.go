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
	jsoniter "github.com/json-iterator/go"
	"github.com/robfig/cron/v3"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"time"
)

type JobReceive struct {
	set    *models.CronReceive
	params *dtos.ReceiveWebHook

	confs       []*pb.CronConfigListItem // 配置任务集合
	conf        *JobConfig               // 流水线解析后的本身配置
	msgSetParse *dtos.MsgSetParse
	tracer      trace.Tracer
}

// 任务执行器
func NewJobReceive(conf *models.CronReceive, param *dtos.ReceiveWebHook) *JobReceive {
	job := &JobReceive{
		set:    conf,
		params: param,
	}
	job.conf = NewJobConfig(&models.CronConfig{
		Id:  conf.Id,
		Env: conf.Env,
		//Type: models.TypeOnce,
		Name:     conf.Name,
		Protocol: 98,
		Command:  nil,
		Remark:   conf.Remark,
		Status:   conf.Status,
		MsgSet:   conf.MsgSet,
	})
	//if err := job.parse(conf); err != nil {
	//	log.Println("流水线配置解析错误", err.Error())
	//}
	//param, _ := job.conf.ParseParams(nil)
	job.conf.Parse(nil)

	// 日志
	job.tracer = tracing.Tracer(job.set.Env+"-cronin", trace.WithInstrumentationAttributes(
		attribute.String("driver", "mysql"),
		attribute.String("env", job.set.Env),
	))

	return job
}

//func (job *JobReceive) parse(conf *models.CronReceive) error {
//	job.confs = []*pb.CronConfigListItem{}
//	return jsoniter.Unmarshal(conf.RuleConfig, &job.confs)
//}

// 执行任务
func (job *JobReceive) Run() {
	var err errs.Errs
	ctx1, span := job.tracer.Start(context.Background(), "job-receive", trace.WithAttributes(
		attribute.Int("ref_id", job.set.Id),
		attribute.String("env", job.set.Env),
		attribute.String("component", "receive"),
		attribute.String("name", job.set.Name),
	), tracing.Extract(job.params.TraceId))
	defer func() {
		job.conf.isRun = false
		if err != nil {
			span.SetStatus(tracing.StatusError, err.Desc())
			span.AddEvent("error", trace.WithAttributes(attribute.String("error.object", err.Error())))
		} else if er := util.PanicInfo(recover()); er != "" {
			span.SetStatus(tracing.StatusError, "执行异常")
			span.AddEvent("error", trace.WithAttributes(attribute.String("error.panic", er)))
		} else {
			span.SetStatus(tracing.StatusOk, "")
		}
		rawId := cron.EntryID(job.conf.conf.EntryId)
		e := cronRun.Entry(rawId)
		if e.ID == rawId {
			cronRun.Remove(e.ID)
		}
		job.conf.conf.EntryId = 0

		span.End()
	}()
	if job.conf.isRun {
		err = errs.New(nil, "任务正在进行中，跳过")
		return
	}
	ctx, cancel := context.WithCancelCause(ctx1)
	job.conf.ctxCancel = cancel
	job.conf.isRun = true
	job.conf.runTime = time.Now()

	// 匹配任务
	if job.set.RuleConfig == nil {
		return
	}
	rulesRaw := []*pb.ReceiveRuleItem{}
	if er := jsoniter.Unmarshal(job.set.RuleConfig, &rulesRaw); err != nil {
		err = errs.New(er, "规则任务序列化错误")
		return
	}
	confIds := []int{}
	ruleSelected := []*pb.ReceiveRuleItem{}
	for _, rule := range rulesRaw { // 以规则配置的顺序优先
		ll := len(rule.Rule)
		for _, item := range job.params.Dataset {
			rl := 0
			for _, r := range rule.Rule { // 规则需要完全匹配
				//这里应该类似map中去找
				if val, ok := item[r.Key]; ok && val != "" && val == r.Value {
					rl++
				}
			}
			if ll > 0 && ll == rl {
				for _, p := range rule.Param { // 参数匹配替换；空值表示未匹配
					p.Value, _ = item[p.Value]
				}
				ruleSelected = append(ruleSelected, rule)
				confIds = append(confIds, rule.Config.Id)
			}
		}
	}

	// 加载最新任务信息
	w := db.NewWhere().Eq("env", job.set.Env).In("id", confIds)
	list, er := data.NewCronConfigData(ctx).List(w, len(confIds))
	if er != nil {
		err = errs.New(er, "任务查询错误")
		return
	}
	listMap := map[int]*models.CronConfig{}
	for _, item := range list {
		listMap[item.Id] = item
	}

	// 欢迎语
	job.conf.messagePush(ctx, &dtos.MsgPushRequest{
		Status:     0,
		StatusDesc: "开始",
		Args: map[string]any{
			"receive": map[string]any{
				"title": job.params.Title,
				"user":  job.params.User,
			},
		},
	})

	// 执行任务
	for _, rule := range ruleSelected {
		conf, ok := listMap[rule.Config.Id]
		if !ok { // 极端情况下，这里可能存在匹配不上; 说明任务已经被弃用，那就不要执行了。
			continue
		}
		p := map[string]any{}
		for _, item := range rule.Param {
			if item.Value != "" {
				p[item.Key] = item.Value
			}
		}

		j := NewJobConfig(conf)

		_, er := j.Running(ctx, "接收任务", p)
		if er != nil {
			j.messagePush(ctx, &dtos.MsgPushRequest{
				Status:     enum.StatusDisable,
				StatusDesc: er.Desc() + " 接收任务",
				Body:       []byte(err.Error()),
				Duration:   time.Since(job.conf.runTime).Seconds(),
			})
			// 这里要确认一下是否继续执行下去。
			if job.set.ConfigErrAction == models.ErrActionStop {
				return
			}
		}
		if job.set.Interval > 0 {
			if err = j.Sleep(ctx, time.Duration(job.set.Interval)*time.Second); err != nil {
				return
			}
		}
	}
	// 结束语
	job.conf.messagePush(ctx, &dtos.MsgPushRequest{
		Status:     enum.StatusActive,
		StatusDesc: "完成",
		Duration:   time.Since(job.conf.runTime).Seconds(),
		Args: map[string]any{
			"receive": map[string]any{
				"title": job.params.Title,
				"user":  job.params.User,
			},
		},
	})
}

func (job *JobReceive) GetConf() *JobConfig {
	return job.conf
}

func (job *JobReceive) SetEntryId(id int) {
	job.conf.conf.EntryId = id
}

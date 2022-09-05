package biz

import (
	"context"
	"cron/internal/basic/config"
	"cron/internal/basic/conv"
	"cron/internal/basic/db"
	"cron/internal/data"
	"cron/internal/models"
	"cron/internal/pb"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"github.com/robfig/cron/v3"
	"net/http"
	"time"
)

type TaskService struct {
	cron *cron.Cron // 任务计划 组件
	conf config.Main
}

func NewTaskService(conf config.Main) *TaskService {
	return &TaskService{
		cron: cronRun,
		conf: conf,
	}
}

// Init 初始化任务
func (dm *TaskService) Init() (err error) {
	pageSize, total := 500, int64(500)
	cronDb := data.NewCronConfigData(context.Background())
	for page := 1; total >= int64(pageSize*page); page++ {
		list := []*models.CronConfig{}
		w := db.NewWhere().Eq("status", models.StatusActive)
		total, err = cronDb.GetList(w, page, pageSize, &list)
		if err != nil {
			panic(fmt.Sprintf("任务配置读取异常：%s", err.Error()))
		}
		for _, conf := range list {
			dm.Add(conf)
		}
	}

	// 系统内置任务
	dm.Add(dm.sysLogRetentionConf())

	return nil
}

// 添加任务
func (dm *TaskService) Add(conf *models.CronConfig) {
	if conf == nil {
		return
	}
	st := time.Now()
	j := NewCronJob(conf)
	id, err := dm.cron.AddJob(conf.Spec, j)
	if err != nil {
		g := models.NewErrorCronLog(conf, fmt.Sprintf("任务启动失败，%s；", err.Error()), st)
		data.NewCronLogData(context.Background()).Add(g)
		return
	}
	j.SetCronId(id)
	jobList.Store(conf.Id, j)
}

// 删除任务
func (dm *TaskService) Del(conf *models.CronConfig) {
	if temp, ok := jobList.Load(conf.Id); ok == true {
		job := temp.(*CronJob)
		dm.cron.Remove(job.GetCronId())
		jobList.Delete(conf.Id)
	}
}

// sysLogDurationConf 内置任务，日志删除
func (dm *TaskService) sysLogRetentionConf() *models.CronConfig {
	retention := dm.conf.Task.LogRetention
	if retention == "" {
		return nil
	}
	re, err := time.ParseDuration(retention)
	if err != nil {
		panic(fmt.Sprintf("log_retention 日志存续配置有误, %s", err.Error()))
	} else if re.Hours() < 24 {
		panic("log_retention 日志存续不得小于24h")
	}

	var sysLogRetention = &models.CronConfig{
		Id:       -1,
		Name:     "日志留存时间",
		Spec:     "0 0 5 * * *", // 每天5点执行
		Protocol: models.ProtocolHttp,
		Status:   models.StatusActive,
		Remark:   "系统内置任务",
		CreateDt: time.Now().Format(conv.FORMAT_DATETIME),
		UpdateDt: time.Now().Format(conv.FORMAT_DATETIME),
	}
	cmd := &pb.CronConfigCommand{
		Http: struct {
			Method string `json:"method"`
			Url    string `json:"url"`
			Body   string `json:"body"`
		}{
			Method: http.MethodPost,
			Url:    dm.conf.Http.Local() + "/log/del",
			Body:   fmt.Sprintf(`{"retention":"%s"}`, dm.conf.Task.LogRetention),
		},
	}
	sysLogRetention.Command, _ = jsoniter.MarshalToString(cmd)
	return sysLogRetention
}

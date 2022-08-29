package biz

import (
	"context"
	"cron/internal/basic/db"
	"cron/internal/data"
	"cron/internal/models"
	"fmt"
	"github.com/robfig/cron/v3"
)

type TaskService struct {
	cron *cron.Cron // 任务计划 组件
}

func NewTaskService() *TaskService {
	return &TaskService{
		cron: cronRun,
	}
}

func (dm *TaskService) Init() (err error) {
	/*
		1.读取所有已经配置有效的任务；
		2.全部注册到任务服务；
			注册后，会返回一个任务id；
			要维护一个内存队列，配置id对应一个任务id；
				配置任务关闭，就要结束对应的任务id；
		3.对已经注册的任务，要能够进行关闭控制；
	*/
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

	return nil
}

// 添加任务
func (dm *TaskService) Add(conf *models.CronConfig) {
	j := NewCronJob(conf)
	id, err := dm.cron.AddJob(conf.Spec, j)
	if err != nil {
		// 这里记录，不做任何返回(db日志写入失败，就一块要写入到文件了)。
		g := models.NewErrorCronLog(conf, err.Error())
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

package models

import (
	"fmt"
	"github.com/robfig/cron/v3"
)

type CronJob struct {
	conf   *CronConfig
	cronId cron.EntryID
}

// 任务执行器
func NewCronJob(conf *CronConfig) *CronJob {
	return &CronJob{conf: conf}
}

// 设置任务执行id
func (job *CronJob) SetCronId(cronId cron.EntryID) {
	job.cronId = cronId
}

func (job *CronJob) Run() {
	switch job.conf.Protocol {
	case ProtocolHttp:
		job.httpFunc()
	case ProtocolRpc:
		job.rpcFunc()
	case ProtocolCmd:
		job.cmdFunc()
	}
}

// http 执行函数
func (job *CronJob) httpFunc() {
	fmt.Println("执行http 任务")
}

// rpc 执行函数
func (job *CronJob) rpcFunc() {
	fmt.Println("执行rpc 任务")
}

// rpc 执行函数
func (job *CronJob) cmdFunc() {
	fmt.Println("执行cmd 任务")
}

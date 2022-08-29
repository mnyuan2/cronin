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

// 返回任务执行中的id
func (job *CronJob) GetCronId() cron.EntryID {
	return job.cronId
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

	// 执行请求任务，并记录结果日志
	fmt.Println("执行http 任务")
	/*
		这里有三个点：1.请求类型、2.请求url、3.请求body；
			默认只能说是get请求的一个url
			最好的方案就是前段拼装成一个json
	*/

}

// rpc 执行函数
func (job *CronJob) rpcFunc() {
	fmt.Println("执行rpc 任务")
}

// rpc 执行函数
func (job *CronJob) cmdFunc() {
	// 这个最后兼容
	fmt.Println("执行cmd 任务")
}

package biz

import (
	"github.com/robfig/cron/v3"
	"log"
	"os"
	"sync"
)

// 调度器连接
var cronRun *cron.Cron

// 时间解释器 （主要用于验证时间格式）
var secondParser = cron.NewParser(cron.Second | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.DowOptional | cron.Descriptor)

// 执行中的任务列表
var jobList = sync.Map{} //map[int]*models.CronJob{}

// 全局初始化
func init() {
	cronInit()
}

func cronInit() {
	// 这个应该是全局唯一初始化（不能重复初始化）
	cronRun = cron.New(
		// 标准的cron时间语法解析器
		cron.WithSeconds(),
		// 作业 包装器
		cron.WithChain(
			// 异常处理，保证单个任务的panic，不会影响其它任务。
			cron.Recover(cron.DefaultLogger),
			// 如果上一个任务还未执行完成，则跳过该次调度
			cron.SkipIfStillRunning(cron.VerbosePrintfLogger(log.New(os.Stdout, "上一步还未完成, 跳过:", log.LstdFlags))),
		),
	)
	cronRun.Start() // 启动程序；启动之后添加任务也是可以的；
	//cronRun.Stop() // 停止服务
}

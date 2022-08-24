package biz

import (
	"github.com/robfig/cron/v3"
	"log"
	"os"
)

var cronRun *cron.Cron

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
}

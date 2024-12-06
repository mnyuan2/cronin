package biz

import (
	"github.com/robfig/cron/v3"
	"log"
	"os"
)

// 调度器连接
var cronRun *cron.Cron

// 时间解释器 （主要用于验证时间格式）
var secondParser = cron.NewParser(cron.Second | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.DowOptional | cron.Descriptor)

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
			SkipIfStillRunning(cron.VerbosePrintfLogger(log.New(os.Stdout, "上一步还未完成, 跳过:", log.LstdFlags))),
		),
	)
	log.Println("cron 初始化完成")
	cronRun.Start() // 启动程序；启动之后添加任务也是可以的；
	//cronRun.Stop() // 停止服务
}

// SkipIfStillRunning 如果先前的调用仍在运行，跳过对 Job 的调用。它在 Info 级别记录跳转到给定记录器的日志。
func SkipIfStillRunning(logger cron.Logger) cron.JobWrapper {
	return func(j cron.Job) cron.Job {
		var ch = make(chan struct{}, 1)
		ch <- struct{}{}
		return cron.FuncJob(func() {
			select {
			case v := <-ch:
				j.Run()
				ch <- v
			default:
				switch val := j.(type) {
				case *JobConfig:
					logger.Info(val.conf.Name)
				case *JobPipeline:
					logger.Info(val.conf.conf.Name)
				case *JobReceive:
					logger.Info(val.conf.conf.Name)
				default:
					logger.Info("未知任务类型")
				}
			}
		})
	}
}

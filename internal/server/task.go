package server

import (
	"cron/internal/basic/config"
	"cron/internal/biz"
)

// 这里还是要回去看一下参考，看一下对于入参的限制。

func InitTask() {
	task := biz.NewTaskService(config.MainConf())
	task.Init()
	go task.RegisterMonitor()

}

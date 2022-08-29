package server

import "cron/internal/biz"

// 这里还是要回去看一下参考，看一下对于入参的限制。

func InitTask() {
	biz.NewTaskService().Init()
}

package main

import "cron/internal/server"

func main() {
	// 初始化任务

	// 初始化http
	r := server.Init()
	r.Run(":8081")
}

package main

import (
	"cron/internal/server"
	"embed"
)

var (
	//go:embed web
	Resource embed.FS
)

func main() {
	// 初始化任务
	server.InitTask()
	// 初始化http
	r := server.InitHttp(Resource)
	r.Run(":9003")
}

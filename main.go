package main

import (
	"context"
	"cron/internal/basic/config"
	"cron/internal/basic/db"
	"cron/internal/models"
	"cron/internal/server"
	"embed"
	"fmt"
)

var (
	//go:embed web
	Resource        embed.FS
	version         = ""      // 版本号  构建时通过 -ldflags "-X main.version=0.3.4" 进行指定
	isBuildResource = "false" // 是否打包静态资源 构建时通过 -ldflags "-X main.isBuildResource=true" 进行指定
)

func main() {
	config.Version = version
	fmt.Println("版本号", config.Version, isBuildResource)
	// 注册mysql表
	models.AutoMigrate(db.New(context.Background()))
	// 初始化任务
	server.InitTask()
	// 初始化http
	r := server.InitHttp(Resource, isBuildResource == "true")
	r.Run(":" + config.MainConf().Http.Port)
}

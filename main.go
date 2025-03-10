package main

import (
	"context"
	"cron/internal/basic/config"
	"cron/internal/basic/db"
	"cron/internal/basic/tracing"
	"cron/internal/models"
	"cron/internal/server"
	"embed"
	"log"
)

var (
	//go:embed web
	Resource        embed.FS
	version         = "v0.0.0" // 版本号  构建时通过 -ldflags "-X main.version=v0.0.0" 进行指定; 版本号位说明：1.不向下兼容的发布、2.向下兼容功能发布、3.bug修正且向下兼容发布
	isBuildResource = "false"  // 是否打包静态资源 构建时通过 -ldflags "-X main.isBuildResource=true" 进行指定
)

func main() {
	config.Version = version
	log.Println("版本号", config.Version, isBuildResource)
	// 注册mysql表
	models.AutoMigrate(db.New(context.Background()))
	// 日志写入
	go tracing.MysqlCollectorListen()
	// 初始化任务
	server.InitTask()
	// 初始化http
	r := server.InitHttp(Resource, isBuildResource == "true")
	r.Run(":" + config.MainConf().Http.Port)
}

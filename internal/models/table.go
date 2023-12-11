package models

import (
	"cron/internal/basic/db"
	"fmt"
	"time"
)

// 注册表结构
func AutoMigrate(db *db.Database) {
	// 迁移表结构
	err := db.Write.Set("gorm:table_options", "ENGINE=InnoDB CHARSET=utf8mb4").
		AutoMigrate(&CronSetting{}, &CronConfig{}, &CronLog{})
	if err != nil {
		panic(fmt.Sprintf("mysql 表初始化失败，%s", err.Error()))
	}
	// 初始化数据
	err = db.Write.Where("scene='env' and env='system' and title='默认环境'").FirstOrCreate(&CronSetting{
		Scene:    "env",
		Title:    "默认环境",
		Env:      "system",
		CreateDt: time.Now().Format(time.DateTime),
	}).Error
	if err != nil {
		panic(fmt.Sprintf("cron_setting 表默认行数据初始化失败，%s", err.Error()))
	}
}

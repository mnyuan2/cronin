package models

import (
	"cron/internal/basic/db"
	"cron/internal/basic/enum"
	"fmt"
	"time"
)

// 注册表结构
func AutoMigrate(db *db.MyDB) {
	// 迁移表结构
	err := db.Set("gorm:table_options", "ENGINE=InnoDB CHARSET=utf8mb4").
		AutoMigrate(&CronSetting{}, &CronConfig{}, &CronLogSpan{}, &CronUser{})
	if err != nil {
		panic(fmt.Sprintf("mysql 表初始化失败，%s", err.Error()))
	}
	// 初始化数据
	err = db.Where("scene=? and status=?", SceneEnv, enum.StatusActive).FirstOrCreate(&CronSetting{
		Scene:    "env",
		Name:     "public",
		Title:    "public",
		Content:  `{"default":2}`,
		Status:   enum.StatusActive,
		CreateDt: time.Now().Format(time.DateTime),
		UpdateDt: time.Now().Format(time.DateTime),
	}).Error
	if err != nil {
		panic(fmt.Sprintf("cron_setting 表默认行数据初始化失败，%s", err.Error()))
	}
	msg := &CronSetting{}
	err = db.Where("scene=?", SceneMsg).Find(msg).Error
	if err != nil {
		panic(fmt.Sprintf("cron_setting 表默认行数据初始化失败，%s", err.Error()))
	}
	if msg.Id == 0 { // 后期会有多条默认消息模板
		db.CreateInBatches([]*CronSetting{
			{
				Scene: "msg",
				Name:  "",
				Title: "企微xx群",
				Content: `{
	"http":{
		"method":"POST",
		"url":"https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=xx",
		"body":"{
			\"msgtype\": \"text\",
			\"text\": {
				\"content\": \"时间：[[log.create_dt]]\\n任务 [[config.name]]执行[[log.status_name]]了 \\n耗时[[log.duration]]秒\\n响应：[[log.body]]\",
				\"mentioned_mobile_list\": [[user.mobile]]
			}
		}",
		"header":[{"key":"","value":""}]
	}
}`,
				Status:   enum.StatusActive,
				CreateDt: time.Now().Format(time.DateTime),
				UpdateDt: time.Now().Format(time.DateTime),
			},
		}, 10)
	}

}

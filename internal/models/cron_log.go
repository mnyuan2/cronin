package models

import (
	"cron/internal/basic/conv"
	"time"
)

type CronLog struct {
	Id       int    `json:"id"`        // 主键
	ConfId   int    `json:"conf_id"`   // 配置id
	CreateDt string `json:"create_dt"` // 创建时间
	Status   int    `json:"status"`    // 状态：1.错误、2.正常
	Body     string `json:"body"`      // 日志文本
}

var LogStatusMap = map[int]string{
	StatusDisable: "错误",
	StatusActive:  "正常",
}

func NewErrorCronLog(confId int, body string) *CronLog {
	return &CronLog{
		ConfId:   confId,
		CreateDt: time.Now().Format(conv.FORMAT_DATETIME),
		Status:   StatusDisable,
		Body:     body,
	}
}

func NewSuccessCronLog(confId int, body string) *CronLog {
	return &CronLog{
		ConfId:   confId,
		CreateDt: time.Now().Format(conv.FORMAT_DATETIME),
		Status:   StatusActive,
		Body:     body,
	}
}

package models

import (
	"cron/internal/basic/conv"
	jsoniter "github.com/json-iterator/go"
	"time"
)

type CronLog struct {
	Id       int    `json:"id"`        // 主键
	ConfId   int    `json:"conf_id"`   // 配置id
	CreateDt string `json:"create_dt"` // 创建时间
	Status   int    `json:"status"`    // 状态：1.错误、2.正常
	Body     string `json:"body"`      // 日志文本
	Snap     string `json:"snap"`      // 任务快照
}

var LogStatusMap = map[int]string{
	StatusDisable: "错误",
	StatusActive:  "正常",
}

// 新建一个错误日志
func NewErrorCronLog(conf *CronConfig, body string) *CronLog {
	str, _ := jsoniter.MarshalToString(conf)
	return &CronLog{
		ConfId:   conf.Id,
		CreateDt: time.Now().Format(conv.FORMAT_DATETIME),
		Status:   StatusDisable,
		Body:     body,
		Snap:     str,
	}
}

// 新建一个成功日志
func NewSuccessCronLog(conf *CronConfig, body string) *CronLog {
	str, _ := jsoniter.MarshalToString(conf)
	return &CronLog{
		ConfId:   conf.Id,
		CreateDt: time.Now().Format(conv.FORMAT_DATETIME),
		Status:   StatusActive,
		Body:     body,
		Snap:     str,
	}
}

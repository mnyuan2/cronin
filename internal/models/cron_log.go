package models

import (
	"cron/internal/basic/conv"
	jsoniter "github.com/json-iterator/go"
	"time"
)

type CronLog struct {
	Id       int     `json:"id"`        // 主键
	ConfId   int     `json:"conf_id"`   // 配置id
	CreateDt string  `json:"create_dt"` // 完成时间
	Duration float64 `json:"duration"`  // 耗时
	Status   int     `json:"status"`    // 状态：1.错误、2.正常
	Body     string  `json:"body"`      // 日志文本
	Snap     string  `json:"snap"`      // 任务快照
}

var LogStatusMap = map[int]string{
	StatusDisable: "错误",
	StatusActive:  "正常",
}

// 新建一个错误日志
func NewErrorCronLog(conf *CronConfig, body string, startTime time.Time) *CronLog {
	str, _ := jsoniter.MarshalToString(conf)
	t := time.Now()
	return &CronLog{
		ConfId:   conf.Id,
		CreateDt: t.Format(conv.FORMAT_DATETIME),
		Duration: t.Sub(startTime).Seconds(),
		Status:   StatusDisable,
		Body:     body,
		Snap:     str,
	}
}

// 新建一个成功日志
func NewSuccessCronLog(conf *CronConfig, body string, startTime time.Time) *CronLog {
	str, _ := jsoniter.MarshalToString(conf)
	t := time.Now()
	return &CronLog{
		ConfId:   conf.Id,
		CreateDt: t.Format(conv.FORMAT_DATETIME),
		Duration: t.Sub(startTime).Seconds(),
		Status:   StatusActive,
		Body:     body,
		Snap:     str,
	}
}

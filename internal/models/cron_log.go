package models

import (
	"cron/internal/basic/conv"
	"cron/internal/basic/enum"
	jsoniter "github.com/json-iterator/go"
	"time"
)

type CronLog struct {
	Id       int     `json:"id" gorm:"column:id;type:int(11);primary_key;comment:主键;"`
	ConfId   int     `json:"conf_id" gorm:"column:conf_id;type:int(11);index:conf_id;comment:配置id;"`
	CreateDt string  `json:"create_dt" gorm:"column:create_dt;type:datetime;default:null;comment:完成时间;"`
	Duration float64 `json:"duration" gorm:"column:duration;type:double(10,3);default:0;comment:耗时;"`
	Status   int     `json:"status" gorm:"column:status;type:tinyint(2);default:0;comment:状态：1.错误、2.正常;"`
	Body     string  `json:"body" gorm:"column:body;type:text;comment:日志内容;"`
	Snap     string  `json:"snap" gorm:"column:snap;type:text;comment:任务快照;"`
}

var LogStatusMap = map[int]string{
	enum.StatusDisable: "错误",
	enum.StatusActive:  "正常",
}

// 新建一个错误日志
func NewErrorCronLog(conf *CronConfig, body string, startTime time.Time) *CronLog {
	str, _ := jsoniter.MarshalToString(conf)
	t := time.Now()
	return &CronLog{
		ConfId:   conf.Id,
		CreateDt: t.Format(conv.FORMAT_DATETIME),
		Duration: t.Sub(startTime).Seconds(),
		Status:   enum.StatusDisable,
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
		Status:   enum.StatusActive,
		Body:     body,
		Snap:     str,
	}
}

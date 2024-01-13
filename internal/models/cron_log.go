package models

import (
	"cron/internal/basic/conv"
	"cron/internal/basic/enum"
	"cron/internal/basic/errs"
	jsoniter "github.com/json-iterator/go"
	"time"
)

type CronLog struct {
	Id         int     `json:"id" gorm:"column:id;type:int(11);primary_key;comment:主键;"`
	Env        string  `json:"env" gorm:"column:env;type:varchar(32);index:conf_id,priority:10;comment:环境;"`
	ConfId     int     `json:"conf_id" gorm:"column:conf_id;type:int(11);index:conf_id,priority:11;comment:配置id;"`
	CreateDt   string  `json:"create_dt" gorm:"column:create_dt;type:datetime;index:conf_id,priority:12;default:null;comment:完成时间;"`
	Duration   float64 `json:"duration" gorm:"column:duration;type:double(10,3);default:0;comment:耗时;"`
	Status     int     `json:"status" gorm:"column:status;type:tinyint(2);default:0;comment:状态：1.错误、2.正常;"`
	StatusDesc string  `json:"status_desc" gorm:"column:status_desc;type:varchar(255);default:'';comment:错误描述;"`
	Body       string  `json:"body" gorm:"column:body;type:text;comment:日志内容;"`
	Snap       string  `json:"snap" gorm:"column:snap;type:text;comment:任务快照;"`
}

var LogStatusMap = map[int]string{
	enum.StatusDisable: "错误",
	enum.StatusActive:  "正常",
}

// 新建一个错误日志
func NewErrorCronLog(conf *CronConfig, body string, err error, startTime time.Time) *CronLog {
	str, _ := jsoniter.MarshalToString(conf)
	t := time.Now()

	g := &CronLog{
		Env:      conf.Env,
		ConfId:   conf.Id,
		CreateDt: t.Format(conv.FORMAT_DATETIME),
		Duration: t.Sub(startTime).Seconds(),
		Status:   enum.StatusDisable,
		Body:     body,
		Snap:     str,
	}
	e, ok := err.(*errs.Error)
	if ok {
		g.StatusDesc = e.Desc()
		if e.Error() != "" {
			g.Body = e.Error() + "\n" + body
		}
	} else {
		g.StatusDesc = "error"
		g.Body = err.Error() + "\n" + body
	}

	return g
}

// 新建一个成功日志
func NewSuccessCronLog(conf *CronConfig, body string, startTime time.Time) *CronLog {
	str, _ := jsoniter.MarshalToString(conf)
	t := time.Now()
	return &CronLog{
		Env:        conf.Env,
		ConfId:     conf.Id,
		CreateDt:   t.Format(conv.FORMAT_DATETIME),
		Duration:   t.Sub(startTime).Seconds(),
		Status:     enum.StatusActive,
		StatusDesc: "success",
		Body:       body,
		Snap:       str,
	}
}

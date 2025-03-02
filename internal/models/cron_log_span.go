package models

import "cron/internal/basic/enum"

// CronLogSpan 第二版日志 依据规范 OpenTelemetry
type CronLogSpan struct {
	Env          string `json:"env" gorm:"column:env;type:varchar(32);index:log_span_env,priority:10;index:log_span_trace_id,priority:10;comment:环境;"`
	RefId        string `json:"ref_id" gorm:"column:ref_id;type:varchar(32);index:log_span_env,priority:11;comment:引用id;"`
	Timestamp    int64  `json:"timestamp" gorm:"column:timestamp;type:bigint(20);default:0;comment:开始时间us/微秒;"`
	TraceId      string `json:"trace_id" gorm:"column:trace_id;type:varchar(32);default:'';index:log_span_trace_id,priority:11;comment:踪迹id;"`
	SpanId       string `json:"span_id" gorm:"column:span_id;type:varchar(32);default:'';comment:节点id;"`
	ParentSpanId string `json:"parent_span_id" gorm:"column:parent_span_id;type:varchar(32);default:'';comment:父节点id;"`
	Service      string `json:"service" gorm:"column:service;type:varchar(120);default:'';comment:服务名称;"`
	Operation    string `json:"operation" gorm:"column:operation;type:varchar(120);default:'';comment:操作名称;"`
	Duration     int64  `json:"duration" gorm:"column:duration;type:bigint(20);default:0;comment:耗时us/微秒;"`
	Tags         []byte `json:"tags" gorm:"column:tags;type:json;default:null;comment:描述;"`
	//TagsKey      []byte `json:"tags_key" gorm:"column:tags_key;type:json;default:null;comment:描述key数组;"` // Deprecated: 废弃，使用 TagsKV 代替
	//TagsVal      []byte `json:"tags_val" gorm:"column:tags_val;type:json;default:null;comment:描述val数组;"` // Deprecated: 废弃，使用 TagsKV 代替
	TagsKV     []byte `json:"tags_kv" gorm:"column:tags_kv;type:json;default:null;comment:tags描述数组k=v;"`
	Logs       []byte `json:"logs" gorm:"column:logs;type:json;default:null;comment:日志;"`
	Status     int    `json:"status" gorm:"column:status;type:tinyint(2);default:0;comment:状态：0.无、1.错误、2.正常;"`
	StatusDesc string `json:"status_desc" gorm:"column:status_desc;type:varchar(255);default:'';comment:状态描述;"`
}

func (m *CronLogSpan) TableName() string {
	return "cron_log_span"
}

var LogSpanStatusMap = map[int]string{
	enum.StatusEmpty:   "-",
	enum.StatusDisable: "失败",
	enum.StatusActive:  "成功",
}

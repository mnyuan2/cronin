package models

// CronLogSpan 第二版日志 依据规范 OpenTelemetry
type CronLogSpan struct {
	Env          string  `json:"env" gorm:"column:env;type:varchar(32);comment:环境;"`
	Timestamp    string  `json:"timestamp" gorm:"column:timestamp;type:datetime;default:null;comment:开始时间;"`
	TraceId      string  `json:"trace_id" gorm:"column:trace_id;type:varchar(32);default:'';comment:踪迹id;"` // 可以是索引
	SpanId       string  `json:"span_id" gorm:"column:span_id;type:varchar(32);default:'';comment:节点id;"`
	ParentSpanId string  `json:"parent_span_id" gorm:"column:parent_span_id;type:varchar(32);default:'';comment:父节点id;"`
	Service      string  `json:"service" gorm:"column:service;type:varchar(120);default:'';comment:服务名称;"`
	Operation    string  `json:"operation" gorm:"column:operation;type:varchar(120);default:'';comment:操作名称;"`
	Duration     float64 `json:"duration" gorm:"column:duration;type:double(19,0);default:0;comment:耗时/毫秒;"`
	Tags         []byte  `json:"tags" gorm:"column:tags;type:json;default:null;comment:描述;"`
	Logs         []byte  `json:"logs" gorm:"column:logs;type:json;default:null;comment:日志;"`
	Status       int32   `json:"status" gorm:"column:status;type:tinyint(2);default:0;comment:状态：0.无、1.错误、2.正常;"`
}

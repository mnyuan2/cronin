package models

// CronLogSpanIndex 日志指标
type CronLogSpanIndex struct {
	Id        int    `json:"id" gorm:"column:id;type:INTEGER;primary_key;comment:主键;"`
	Env       string `json:"env" gorm:"column:env;type:varchar(32);index:log_span_index,priority:10;comment:环境;"`
	RefId     string `json:"ref_id" gorm:"column:ref_id;type:varchar(32);index:log_span_index,priority:11;comment:引用id;"`
	Operation string `json:"operation" gorm:"column:operation;type:varchar(32);default:'';index:log_span_index,priority:12;comment:操作名称;"`
	Timestamp string `json:"timestamp" gorm:"column:timestamp;type:datetime;default:null;index:log_span_index,priority:13;comment:时间h/小时;"`

	StatusEmptyNum   int   `json:"status_empty_num" gorm:"column:status_empty_num;type:int(11);default:0;comment:状态空数量;"`
	StatusErrorNum   int   `json:"status_error_num" gorm:"column:status_error_num;type:int(11);default:0;comment:状态错误数量;"`
	StatusSuccessNum int   `json:"status_success_num" gorm:"column:status_success_num;type:int(11);default:0;comment:状态成功数量;"`
	DurationMax      int64 `json:"duration_max" gorm:"column:duration_max;type:bigint(20);default:0;comment:最大耗时us/微秒;"`
	DurationAvg      int64 `json:"duration_avg" gorm:"column:duration;type:bigint(20);default:0;comment:平均耗时us/微秒;"`

	TraceIds string `json:"trace_ids" gorm:"column:trace_ids;type:json;default:null;comment:踪迹id集合;"` // mediumtext
	// 平均耗时、最高耗时
}

func (m *CronLogSpanIndex) TableName() string {
	return "cron_log_span_index"
}

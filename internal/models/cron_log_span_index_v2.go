package models

// CronLogSpanIndex 日志指标
type CronLogSpanIndexV2 struct {
	Id        int    `json:"id" gorm:"column:id;type:INTEGER;primary_key;comment:主键;"`
	Env       string `json:"env" gorm:"column:env;type:varchar(32);index:span_index,priority:10;comment:环境;"`
	RefId     string `json:"ref_id" gorm:"column:ref_id;type:varchar(32);index:span_index,priority:11;comment:引用id;"`
	Timestamp string `json:"timestamp" gorm:"column:timestamp;type:datetime;default:null;index:span_index,priority:12;comment:时间h/小时;"`
	Operation string `json:"operation" gorm:"column:operation;type:varchar(32);default:'';comment:操作名称;"`
	Status    int    `json:"status" gorm:"column:status;type:tinyint(2);default:0;comment:状态：0.无、1.错误、2.正常;"`
	Duration  int64  `json:"duration" gorm:"column:duration;type:bigint(20);default:0;comment:耗时us/微秒;"`
	TraceId   string `json:"trace_id" gorm:"column:trace_id;type:varchar(32);default:'';comment:踪迹id;"`
}

func (m *CronLogSpanIndexV2) TableName() string {
	return "cron_log_span_index_v2"
}

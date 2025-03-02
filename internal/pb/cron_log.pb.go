package pb

// 通过配置查询
type CronLogListRequest struct {
	Env   string `json:"env" form:"env"`
	Tags  string `json:"tags" form:"tags"`
	Limit int    `json:"limit" form:"limit"`
}

// 通过配置查询
type CronLogListResponse struct {
	List []*CronLogSpan `json:"list"`
	// 后期可能会做分页
}
type CronLogSpan struct {
	Timestamp    int64             `json:"timestamp"`      // 开始时间
	Duration     int64             `json:"duration"`       // 耗时/秒
	Status       int               `json:"status"`         // 状态：0.无、1.错误、2.正常
	StatusName   string            `json:"status_name"`    //
	StatusDesc   string            `json:"status_desc"`    //
	TraceId      string            `json:"trace_id"`       // 踪迹id
	SpanId       string            `json:"span_id"`        // 节点id
	ParentSpanId string            `json:"parent_span_id"` // 父节点id
	Service      string            `json:"service"`        // 服务名称
	Operation    string            `json:"operation"`      // 操作名称
	Tags         []*CronLogSpanKV  `json:"tags"`           // 描述
	Logs         []*CronLogSpanLog `json:"logs"`           // 日志
}

type CronLogSpanKV struct {
	Key   string               `json:"key"`
	Value *CronLogSpanTagValue `json:"value"`
}
type CronLogSpanTagValue struct {
	Type  string `json:"type"`
	Value any    `json:"value"`
}

type CronLogSpanLog struct {
	Name       string           `json:"name"`
	Timestamp  int64            `json:"timestamp"`
	Attributes []*CronLogSpanKV `json:"attributes"`
}

// 日志踪迹请求
type CronLogTraceRequest struct {
	TraceId string `json:"trace_id" form:"trace_id"`
}

// 日志踪迹响应
type CronLogTraceResponse struct {
	List  []*CronLogTraceItem `json:"list"`
	Limit int                 `json:"limit"`
	Total int                 `json:"total"`
}

type CronLogTraceItem struct {
	TraceId string
	Spans   []*CronLogSpan
}

// 删除日志
type CronLogDelRequest struct {
	Retention string `json:"retention"`
}
type CronLogDelResponse struct {
	Count int `json:"count"`
}

// 变更日志列表
type CronChangeLogListRequest struct {
	Page    int    `form:"page"`
	Size    int    `form:"size"`
	RefType string `form:"ref_type"`
	RefId   int    `form:"ref_id"`
}
type CronChangeLogListResponse struct {
	List []*CronChangeLogItem `json:"list"`
	Page *Page                `json:"page"`
}
type CronChangeLogItem struct {
	Id             int                       `json:"id,omitempty"`
	CreateDt       string                    `json:"create_dt,omitempty"`
	CreateUserId   int                       `json:"create_user_id,omitempty"`
	CreateUserName string                    `json:"create_user_name,omitempty"`
	Type           int                       `json:"type" format:"enum:type"`
	TypeName       string                    `json:"type_name"`
	RefType        string                    `json:"ref_type,omitempty"`
	RefId          int                       `json:"ref_id,omitempty"`
	ContentStr     string                    `json:"-" gorm:"column:content;"`
	Content        []*CronChangeLogItemField `json:"content" gorm:"-"`
}
type CronChangeLogItemField struct {
	Field      string `json:"field"`
	VType      string `json:"v_type"`
	OldVal     any    `json:"old_val"`
	NewVal     any    `json:"new_val"`
	FieldName  string `json:"field_name"`
	OldValName string `json:"old_val_name"`
	NewValName string `json:"new_val_name"`
}

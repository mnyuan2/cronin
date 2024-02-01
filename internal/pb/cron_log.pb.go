package pb

// 通过配置查询
type CronLogListRequest struct {
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
	Status       int32             `json:"status"`         // 状态：0.无、1.错误、2.正常
	StatusName   string            `json:"status_name"`    //
	StatusDesc   string            `json:"status_desc"`    //
	TraceId      string            `json:"trace_id"`       // 踪迹id
	SpanId       string            `json:"span_id"`        // 节点id
	ParentSpanId string            `json:"parent_span_id"` // 父节点id
	Service      string            `json:"service"`        // 服务名称
	Operation    string            `json:"operation"`      // 操作名称
	Tags         []*CronLogSpanTag `json:"tags"`           // 描述
	Logs         []*CronLogSpanLog `json:"logs"`           // 日志
}

type CronLogSpanTag struct {
	Key   string
	Value *CronLogSpanTagValue
}
type CronLogSpanTagValue struct {
	Type  string
	Value any
}

type CronLogSpanLog struct {
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

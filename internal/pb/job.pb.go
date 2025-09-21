package pb

// 用户列表
type JobStopRequest struct {
	RefId   int `json:"ref_id"`
	EntryId int `json:"entry_id"`
}
type JobStopReply struct{}

// 日志详情
type JobTracesRequest struct {
	RefId   int    `json:"ref_id" form:"ref_id"`
	EntryId int    `json:"entry_id" form:"entry_id"`
	TraceId string `json:"trace_id" form:"trace_id"`
}
type JobTracesResponse struct {
	List  []*CronLogTraceItem `json:"list"`
	Limit int                 `json:"limit"`
	Total int                 `json:"total"`
}

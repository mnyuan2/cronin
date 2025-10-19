package dtos

// 执行队列元素
type ExecQueueItem struct {
	RefId    int     `json:"ref_id"`
	RefType  string  `json:"ref_type"`
	EntryId  int     `json:"entry_id"`
	Name     string  `json:"name"`
	Duration float64 `json:"duration"`
	TraceId  string  `json:"trace_id"`
}

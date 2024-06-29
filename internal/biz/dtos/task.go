package dtos

// 执行队列元素
type ExecQueueItem struct {
	RefId    int     `json:"ref_id"`
	EntryId  int     `json:"entryId"`
	Name     string  `json:"name"`
	Duration float64 `json:"duration"`
}

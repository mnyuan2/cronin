package pb

type WorkTableRequest struct {
	Tab       string `json:"tab" form:"tab"`
	SourceIds []int  `json:"source_ids" form:"source_ids[]"`
}
type WorkTableReply struct {
	List []*WorkTableItem `json:"list"`
}

type WorkTableItem struct {
	Env      string `json:"env"`
	EnvTitle string `json:"env_title"`
	Type     string `json:"type"`
	Total    int64  `json:"total"`
}

// 删除任务
type WorkTaskDelRequest struct {
	Retention string `json:"retention"`
}
type WorkTaskDelReply struct {
	Info string `json:"info"`
}

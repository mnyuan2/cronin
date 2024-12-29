package dtos

type ReceiveWebHook struct {
	User    string              `json:"user"` // 操作人
	Type    string              `json:"type"`
	Event   string              `json:"event"`
	Title   string              `json:"title"`    // 标题
	Dataset []map[string]string `json:"dataset"`  // 数据集
	TraceId string              `json:"trace_id"` // 内部占用
}

type ReceiveWebHookData struct {
	Owner   string `json:"owner,omitempty"`
	Repo    string `json:"repo,omitempty"`
	Number  string `json:"number,omitempty"`
	Type    string `json:"type,omitempty"`
	Service string `json:"service,omitempty"`
}

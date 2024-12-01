package dtos

type ReceiveWebHook struct {
	User    string              `json:"user"`
	Type    string              `json:"type"`
	Event   string              `json:"event"`
	Dataset []map[string]string `json:"dataset"`
	TraceId string              `json:"trace_id"`
}

type ReceiveWebHookData struct {
	Owner   string `json:"owner,omitempty"`
	Repo    string `json:"repo,omitempty"`
	Number  string `json:"number,omitempty"`
	Type    string `json:"type,omitempty"`
	Service string `json:"service,omitempty"`
}

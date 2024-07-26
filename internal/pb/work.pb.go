package pb

type WorkTableRequest struct {
	Tab string `json:"tab" form:"tab"`
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

package pb

type WorkTableRequest struct {
	Status int `json:"status" form:"status"`
}
type WorkTableReply struct {
	List []*WorkTableItem `json:"list"`
}

type WorkTableItem struct {
	Env      string `json:"env"`
	EnvTitle string `json:"env_title"`
	JoinType string `json:"join_type"`
	Total    int64  `json:"total"`
}

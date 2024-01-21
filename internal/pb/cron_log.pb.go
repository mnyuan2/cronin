package pb

// 通过配置查询
type CronLogByConfigRequest struct {
	ConfId int `json:"conf_id" form:"conf_id"`
	Limit  int `json:"limit" form:"limit"`
}

// 通过配置查询
type CronLogByConfigResponse struct {
	List []*CronLogItem `json:"list"`
}
type CronLogItem struct {
	Id            int      `json:"id"`
	ConfId        int      `json:"conf_id"`
	CreateDt      string   `json:"create_dt"`
	Duration      float64  `json:"duration"`
	Status        int      `json:"status"`
	StatusName    string   `json:"status_name"`
	StatusDesc    string   `json:"status_desc"`
	Body          string   `json:"body"`
	Snap          string   `json:"snap"`
	MsgStatus     int      `json:"msg_status"`
	MsgStatusName string   `json:"msg_status_name"`
	MsgBody       []string `json:"msg_body"`
}

// 删除日志
type CronLogDelRequest struct {
	Retention string `json:"retention"`
}
type CronLogDelResponse struct {
	Count int `json:"count"`
}

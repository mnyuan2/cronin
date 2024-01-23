package pb

// 告警列表
type UserListRequest struct {
	Page int `form:"page"`
	Size int `form:"size"`
}
type UserListReply struct {
	List []*UserListItem `json:"list"`
	Page *Page           `json:"page"`
}
type UserListItem struct {
	Id       int    `json:"id"`
	Sort     int    `json:"sort"`
	Username string `json:"username"`
	Mobile   string `json:"mobile"`
	CreateDt string `json:"create_dt"`
	UpdateDt string `json:"update_dt"`
}

// 告警设置
type UserSetRequest struct {
	Id       int    `json:"id"`
	Username string `json:"username"`
	Mobile   string `json:"mobile"`
	Sort     int    `json:"sort"`
}
type UserSetReply struct {
	Id int `json:"id"`
}

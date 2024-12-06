package pb

// 用户列表
type TagListRequest struct{}
type TagListReply struct {
	List []*TagListItem `json:"list"`
}
type TagListItem struct {
	Id             int    `json:"id"`
	Name           string `json:"name"`
	Remark         string `json:"remark"`
	CreateDt       string `json:"create_dt"`
	CreateUserId   int    `json:"create_user_id"`
	CreateUserName string `json:"create_user_name"`
	UpdateDt       string `json:"update_dt"`
	UpdateUserId   int    `json:"update_user_id"`
	UpdateUserName string `json:"update_user_name"`
}

// 用户设置
type TagSetRequest struct {
	Id     int    `json:"id"`
	Name   string `json:"name"`
	Remark string `json:"remark"`
}
type TagSetReply struct {
	Id int `json:"id"`
}

// 用户状态
type TagChangeStatusRequest struct {
	Id     int `json:"id"`
	Status int `json:"status"`
}
type TagChangeStatusReply struct {
}

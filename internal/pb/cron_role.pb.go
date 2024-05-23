package pb

// 用户列表
type RoleListRequest struct {
	Page int `form:"page"`
	Size int `form:"size"`
}
type RoleListReply struct {
	List []*RoleListItem `json:"list"`
	Page *Page           `json:"page"`
}
type RoleListItem struct {
	Id         int    `json:"id"`
	Name       string `json:"name"`
	Remark     string `json:"remark"`
	AuthIds    []int  `json:"auth_ids"`
	Status     int    `json:"status"`
	StatusName string `json:"status_name"`
}

// 用户设置
type RoleSetRequest struct {
	Id     int    `json:"id"`
	Name   string `json:"name"`
	Remark string `json:"remark"`
	Status int    `json:"status"`
}
type RoleSetReply struct {
	Id int `json:"id"`
}

// 规则列表
type AuthListRequest struct {
	UserId  int   `json:"user_id"` // 限定用户，不传所有
	RoleIds []int `json:"role_ids"`
}
type AuthListReply struct {
	List []*AuthListItem `json:"list"`
}
type AuthListItem struct {
	Id    int    `json:"id"`
	Name  string `json:"name"`
	Path  string `json:"path"`
	Group string `json:"group"`
	Tag   string `json:"tag"`
}

// 角色规则列表
type RoleAuthSetRequest struct {
	Id      int   `json:"id"`
	AuthIds []int `json:"auth_ids"`
}
type RoleAuthSetReply struct{}

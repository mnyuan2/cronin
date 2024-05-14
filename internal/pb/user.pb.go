package pb

// 用户列表
type UserListRequest struct {
	Page int `form:"page"`
	Size int `form:"size"`
}
type UserListReply struct {
	List []*UserListItem `json:"list"`
	Page *Page           `json:"page"`
}
type UserListItem struct {
	Id         int    `json:"id"`
	Sort       int    `json:"sort"`
	Username   string `json:"username"`
	Mobile     string `json:"mobile"`
	Status     int    `json:"status"`
	StatusName string `json:"status_name"`
	CreateDt   string `json:"create_dt"`
	UpdateDt   string `json:"update_dt"`
}

// 用户设置
type UserSetRequest struct {
	Id       int    `json:"id"`
	Account  string `json:"account"`
	Username string `json:"username"`
	Mobile   string `json:"mobile"`
	Sort     int    `json:"sort"`
	Password string `json:"password"`
}
type UserSetReply struct {
	Id int `json:"id"`
}

// 用户详情
type UserDetailRequest struct {
	Id int `json:"id" form:"id"`
}
type UserDetailReply struct {
	Id         int    `json:"id"`
	Account    string `json:"account"`
	Username   string `json:"username"`
	Mobile     string `json:"mobile"`
	Sort       int    `json:"sort"`
	Status     int    `json:"status"`
	StatusName string `json:"status_name"`
	UpdateDt   string `json:"update_dt,omitempty"`
	CreateDt   string `json:"create_dt,omitempty"`
}

// 用户状态
type UserChangeStatusRequest struct {
	Id     int `json:"id"`
	Status int `json:"status"`
}
type UserChangeStatusReply struct {
}

// 用户登录
type UserLoginRequest struct {
	Account  string `json:"account"`
	Password string `json:"password"`
}
type UserLoginReply struct {
	User  *UserDetailReply `json:"user"`
	Token string           `json:"token"`
	Menus []byte           `json:"menus"` // 权限菜单，后期补充
}

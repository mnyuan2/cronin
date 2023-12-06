package pb

// sql源
type SettingSqlSource struct {
	Hostname string
	Port     string
	Username string
	Password string
}

// 任务列表
type SettingSqlListRequest struct {
	Page int `form:"page"`
	Size int `form:"size"`
}
type SettingSqlListReply struct {
	List []*SettingSqlListItem `json:"list"`
	Page *Page                 `json:"page"`
}
type SettingSqlListItem struct {
	Id       int               `json:"id"`
	Title    string            `json:"title"`
	CreateDt string            `json:"create_dt"`
	UpdateDt string            `json:"update_dt"`
	Source   *SettingSqlSource `json:"source"`
}

// 任务设置
type SettingSqlSetRequest struct {
	Id     int               `json:"id"`
	Title  string            `json:"title"`
	Source *SettingSqlSource `json:"source"`
}
type SettingSqlSetReply struct {
	Id int `json:"id"`
}

type SettingChangeStatusRequest struct {
	Id     int `json:"id"`
	Status int `json:"status"`
}
type SettingChangeStatusReply struct {
}

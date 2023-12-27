package pb

// sql源
type SettingSqlSource struct {
	Hostname string `json:"hostname"`
	Port     string `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	Database string `json:"database"`
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

// sql连接监测
type SettingSqlPingRequest struct {
	*SettingSqlSource
}
type SettingSqlPingReply struct {
}

type SettingChangeStatusRequest struct {
	Id     int `json:"id"`
	Status int `json:"status"`
}
type SettingChangeStatusReply struct {
}

// 环境列表
type SettingEnvListRequest struct {
	Page int `form:"page"`
	Size int `form:"size"`
}
type SettingEnvListReply struct {
	List []*SettingEnvListItem `json:"list"`
	Page *Page                 `json:"page"`
}
type SettingEnvListItem struct {
	Id         int    `json:"id"`
	Key        string `json:"key"`
	Title      string `json:"title"`
	Status     int    `json:"status"`
	StatusName string `json:"status_name"`
	CreateDt   string `json:"create_dt"`
	UpdateDt   string `json:"update_dt"`
}

// 环境设置
type SettingEnvSetRequest struct {
	Id    int    `json:"id"`
	Title string `json:"title"`
	Key   string `json:"key"`
}
type SettingEnvSetReply struct {
	Id int `json:"id"`
}

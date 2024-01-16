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
	Name       string `json:"name"`
	Title      string `json:"title"`
	Status     int    `json:"status"`
	StatusName string `json:"status_name"`
	CreateDt   string `json:"create_dt"`
	UpdateDt   string `json:"update_dt"`
	Default    int    `json:"default"`
}

// 环境设置
type SettingEnvSetRequest struct {
	Id      int    `json:"id"`
	Name    string `json:"name"`
	Title   string `json:"title"`
	Default int    `json:"default"`
}
type SettingEnvSetReply struct {
	Id int `json:"id"`
}

// 环境设置
type SettingEnvDelRequest struct {
	Id int `json:"id"`
}
type SettingEnvDelReply struct {
}

// 消息模板列表
type SettingMessageListRequest struct {
	Page int `form:"page"`
	Size int `form:"size"`
}
type SettingMessageListReply struct {
	List []*SettingMessageListItem `json:"list"`
	Page *Page                     `json:"page"`
}
type SettingMessageListItem struct {
	Id       int       `json:"id"`
	Title    string    `json:"title"`
	Sort     int       `json:"int"`
	Http     *CronHttp `json:"http"`
	CreateDt string    `json:"create_dt"`
	UpdateDt string    `json:"update_dt"`
}

// 消息模板设置
type SettingMessageSetRequest struct {
	Id    int       `json:"id"`
	Title string    `json:"title"`
	Sort  int       `json:"sort"`
	Http  *CronHttp `json:"http"` // 这里应该增加一个层级，后面可能会有别的模板种类支持
}
type SettingMessageSetReply struct {
	Id int `json:"id"`
}

type SettingMessageRunReply struct {
	Result string `json:"result"`
}

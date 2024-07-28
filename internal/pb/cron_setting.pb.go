package pb

type SettingSource struct {
	Sql     *SettingSqlSource     `json:"sql"`
	Jenkins *SettingJenkinsSource `json:"jenkins"`
	Git     *SettingGitSource     `json:"git"`
	Host    *SettingHostSource    `json:"host"`
}

// sql 源
type SettingSqlSource struct {
	Driver   string `json:"driver"`
	Hostname string `json:"hostname"`
	Port     string `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	Database string `json:"database"`
}

// jenkins 源
type SettingJenkinsSource struct {
	Hostname string `json:"hostname"`
	Username string `json:"username"`
	Password string `json:"password"`
}

// git 源
type SettingGitSource struct {
	Type        string `json:"type"`
	AccessToken string `json:"access_token"` // 用户授权码
}

func (m *SettingGitSource) GetAccessToken() string {
	return m.AccessToken
}

// 主机 源
type SettingHostSource struct {
	Id     int    `json:"id"`
	Type   string `json:"type"`
	Ip     string `json:"ip"`
	Port   string `json:"port"`
	User   string `json:"user"`
	Secret string `json:"secret"`
}

// 任务列表
type SettingListRequest struct {
	Type int `form:"type"`
	Page int `form:"page"`
	Size int `form:"size"`
}
type SettingListReply struct {
	List []*SettingListItem `json:"list"`
	Page *Page              `json:"page"`
}
type SettingListItem struct {
	Id       int            `json:"id"`
	Title    string         `json:"title"`
	CreateDt string         `json:"create_dt"`
	UpdateDt string         `json:"update_dt"`
	Type     int            `json:"type"`
	Source   *SettingSource `json:"source"`
}

// 任务设置
type SettingSqlSetRequest struct {
	Id     int            `json:"id"`
	Title  string         `json:"title"`
	Type   int            `json:"type"`
	Source *SettingSource `json:"source"`
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
	Id       int                     `json:"id"`
	Title    string                  `json:"title"`
	Sort     int                     `json:"sort"`
	Template *SettingMessageTemplate `json:"template"`
	CreateDt string                  `json:"create_dt"`
	UpdateDt string                  `json:"update_dt"`
}

// 消息模板设置
type SettingMessageSetRequest struct {
	Id       int                     `json:"id"`
	Title    string                  `json:"title"`
	Sort     int                     `json:"sort"`
	Type     int                     `json:"type"`
	Template *SettingMessageTemplate `json:"template"`
}
type SettingMessageTemplate struct {
	Http *CronHttp `json:"http"`
}

type SettingMessageSetReply struct {
	Id int `json:"id"`
}

type SettingMessageRunReply struct {
	Result string `json:"result"`
}

// 使用习惯设置
type SettingPreferenceSetRequest struct {
	Pipeline *SettingPreferencePipeline `json:"pipeline"`
	Git      *SettingPreferenceGit      `json:"git"`
}
type SettingPreferencePipeline struct {
	Interval            int `json:"interval"`
	ConfigDisableAction int `json:"config_disable_action"`
}
type SettingPreferenceGit struct {
	OwnerRepo []*SettingPreferenceGitOwner `json:"owner_repo"`
	Owner     string                       `json:"owner"`
	Repo      string                       `json:"repo"`
}
type SettingPreferenceGitOwner struct {
	Owner string                      `json:"owner"`
	Repos []*SettingPreferenceGitRepo `json:"repos"`
}
type SettingPreferenceGitRepo struct {
	Name string `json:"name"`
}
type SettingPreferenceSetReply struct{}

// 使用习惯获取
type SettingPreferenceGetRequest struct{}
type SettingPreferenceGetReply struct {
	Pipeline *SettingPreferencePipeline `json:"pipeline"`
	Git      *SettingPreferenceGit      `json:"git"`
}

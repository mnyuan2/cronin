package models

type CronSetting struct {
	Id       int    `json:"id,omitempty"`
	Key      string `json:"key,omitempty"`
	Title    string `json:"title,omitempty"`
	Env      string `json:"env,omitempty"`
	Content  string `json:"content,omitempty"`
	CreateDt string `json:"create_dt,omitempty"`
	UpdateDt string `json:"update_dt,omitempty"`
	Status   int    `json:"status,omitempty"`
}

// sql连接源配置
type SettingSqlSource struct {
	// 主机、端口、用户名、密码
}

const (
	KeySqlSource = "sql_source" // sql数据源配置
)

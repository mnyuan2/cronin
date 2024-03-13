package models

import "cron/internal/basic/enum"

type CronSetting struct {
	Id       int    `json:"id" gorm:"primary_key;column:id;type:int(11);comment:主键;"`
	Scene    string `json:"scene" gorm:"column:scene;type:varchar(255);index:env,priority:12;comment:使用场景;"`
	Name     string `json:"name" gorm:"column:name;type:varchar(64);comment:名称，英文占用;"`
	Title    string `json:"title" gorm:"column:title;type:varchar(255);comment:名称，中文占用;"`
	Env      string `json:"env" gorm:"column:env;type:varchar(32);index:env,priority:11;comment:环境:system.系统信息、其它.业务环境信息;"`
	Content  string `json:"content" gorm:"column:content;type:text;comment:内容;"`
	CreateDt string `json:"create_dt" gorm:"column:create_dt;type:datetime;default:null;comment:创建时间;"`
	UpdateDt string `json:"update_dt" gorm:"column:update_dt;type:datetime;default:null;comment:更新时间;"`
	Status   int    `json:"status" gorm:"column:status;type:tinyint(2);default:2;comment:状态:枚举由业务定义"`
}

// sql连接源配置
type SettingSqlSource struct {
	// 主机、端口、用户名、密码
}

const (
	SceneSqlSource     = "sql_source"     // sql数据源配置
	SceneJenkinsSource = "jenkins_source" // jenkins 连接配置
	SceneGitSource     = "git_source"     // git
	SceneHostSource    = "host_source"    // 主机 连接
	SceneEnv           = "env"            // 程序环境
	SceneMsg           = "msg"            // 消息推送
)

const EnvDefault = "public" // 默认环境 是不可以删除的

// dic 到 source key的转换
var DicToSource = map[int]string{
	enum.DicSqlSource:     SceneSqlSource,
	enum.DicJenkinsSource: SceneJenkinsSource,
	enum.DicGitSource:     SceneGitSource,
	enum.DicHostSource:    SceneHostSource,
}

var SourceToDic = map[string]int{
	SceneSqlSource:     enum.DicSqlSource,
	SceneJenkinsSource: enum.DicJenkinsSource,
	SceneGitSource:     enum.DicGitSource,
	SceneHostSource:    enum.DicHostSource,
}

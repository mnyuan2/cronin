package models

import (
	"cron/internal/basic/enum"
	"net/http"
)

type CronProtocol int
type CronStatus int

const (
	ProtocolHttp = 1 // http
	ProtocolRpc  = 2 // rpc
	ProtocolCmd  = 3 // 命令行 cmd
	ProtocolSql  = 4 // sql 执行
)

var ProtocolMap = map[int]string{
	ProtocolHttp: "http",
	ProtocolRpc:  "rpc",
	ProtocolCmd:  "cmd",
	ProtocolSql:  "sql",
}

const (
	ConfigStatusDisable = 1 // 停用
	ConfigStatusActive  = 2 // 激活
	ConfigStatusFinish  = 3 // 完成
	ConfigStatusError   = 4 // 错误
	ConfigStatusDelete  = 9 // 删除
)

// 通用状态
var ConfigStatusMap = map[int]string{
	ConfigStatusDisable: "停用",
	ConfigStatusActive:  "激活",
	ConfigStatusError:   "错误",
	ConfigStatusFinish:  "完成",
}

const (
	TypeCycle = 1 // 周期
	TypeOnce  = 2 // 单次
)

func ProtocolHttpMethodMap() map[string]string {
	return map[string]string{
		http.MethodGet:  http.MethodGet,
		http.MethodPost: http.MethodPost,
	}
}
func ConfTypeMap() map[int]string {
	return map[int]string{
		TypeCycle: "周期",
		TypeOnce:  "单次",
	}
}

type CronConfig struct {
	Id           int    `json:"id" gorm:"column:id;type:int(11);primary_key;comment:主键;"`
	Env          string `json:"env" gorm:"column:env;type:varchar(32);index:env;comment:环境;"`
	EntryId      int    `json:"entry_id" gorm:"column:entry_id;type:int(11);default:0;comment:执行队列编号;"`
	Type         int    `json:"type" gorm:"column:type;type:tinyint(2);default:1;comment:类型：1.周期任务（默认）、2.单次任务;"`
	Name         string `json:"name" gorm:"column:name;type:varchar(255);default:'';comment:任务名称;"`
	Spec         string `json:"spec" gorm:"column:spec;type:varchar(32);default:'';comment:执行时间 表达式;"`
	Protocol     int    `json:"protocol" gorm:"column:protocol;type:tinyint(2);default:0;comment:协议：1.http、2.grpc、3.系统命令、4.sql执行;"`
	Command      []byte `json:"command" gorm:"column:command;type:json;default:null;comment:命令内容;"`
	Remark       string `json:"remark" gorm:"column:remark;type:varchar(255);comment:备注;"`
	Status       int    `json:"status" gorm:"column:status;type:tinyint(2);default:1;comment:状态：1.停止、2.启用;"`
	StatusRemark string `json:"status_remark" gorm:"column:status_remark;type:varchar(255);comment:状态变更描述;"`
	StatusDt     string `json:"status_dt" gorm:"column:status_dt;type:datetime;default:null;comment:状态变更时间;"`
	UpdateDt     string `json:"update_dt" gorm:"column:update_dt;type:datetime;default:null;comment:更新时间;"`
	CreateDt     string `json:"create_dt" gorm:"column:create_dt;type:datetime;default:null;comment:创建时间;"`
	AlertId      int    `json:"alert_id" gorm:"column:alert_id;type:int(11);default:0;comment:告警配置编号;"`
	AlertUserIds []byte `json:"alert_user_ids" gorm:"column:alert_user_ids;type:json;default:null;comment:告警通知人员编号;"`
}

func (m *CronConfig) GetProtocolName() string {
	return ProtocolMap[m.Protocol]
}

func (m *CronConfig) GetStatusName() string {
	return enum.StatusMap[m.Status]
}

func (m *CronConfig) GetTypeName() string {
	return ConfTypeMap()[m.Type]
}

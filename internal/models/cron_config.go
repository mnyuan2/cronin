package models

import (
	"cron/internal/basic/enum"
	"net/http"
)

type CronProtocol int
type CronStatus int

const (
	ProtocolHttp    = 1 // http
	ProtocolRpc     = 2 // rpc
	ProtocolCmd     = 3 // 命令行 cmd
	ProtocolSql     = 4 // sql 执行
	ProtocolJenkins = 5 // jenkins 构建
	ProtocolGit     = 6 // git api 任务
)

var ProtocolMap = map[int]string{
	ProtocolHttp:    "http",
	ProtocolRpc:     "rpc",
	ProtocolCmd:     "cmd",
	ProtocolSql:     "sql",
	ProtocolJenkins: "jenkins",
	ProtocolGit:     "git",
}

const (
	ConfigStatusDisable = 1 // 草稿
	ConfigStatusAudited = 5 // 待审核
	ConfigStatusReject  = 6 // 驳回
	ConfigStatusActive  = 2 // 激活
	ConfigStatusFinish  = 3 // 完成
	ConfigStatusError   = 4 // 错误
	ConfigStatusClosed  = 8 // 已关闭（预删除）
	ConfigStatusDelete  = 9 // 删除
)

// 通用状态
var ConfigStatusMap = map[int]string{
	ConfigStatusDisable: "草稿",
	ConfigStatusAudited: "待审核",
	ConfigStatusReject:  "驳回",
	ConfigStatusActive:  "激活",
	ConfigStatusFinish:  "完成",
	ConfigStatusError:   "错误",
	ConfigStatusClosed:  "已关闭",
}

const (
	TypeCycle  = 1 // 周期
	TypeOnce   = 2 // 单次
	TypeModule = 5 // 模块
)

// 任务类型
var ConfigTypeMap = map[int]string{
	TypeCycle:  "周期",
	TypeOnce:   "单次",
	TypeModule: "模块",
}

func ProtocolHttpMethodMap() map[string]string {
	return map[string]string{
		http.MethodGet:    http.MethodGet,
		http.MethodPost:   http.MethodPost,
		http.MethodPut:    http.MethodPut,
		http.MethodDelete: http.MethodDelete,
	}
}
func ConfTypeMap() map[int]string {
	return map[int]string{
		TypeCycle: "周期",
		TypeOnce:  "单次",
	}
}

const (
	RetryModeFixed = 1 // 固定间隔
	RetryModeIncr  = 2 // 递增间隔
)

// 重试模式
var RetryModeMap = map[int]string{
	RetryModeFixed: "固定间隔",
	RetryModeIncr:  "递增间隔",
}

type CronConfig struct {
	Id              int    `json:"id" gorm:"column:id;type:INTEGER;primary_key;comment:主键;"`
	Env             string `json:"env" gorm:"column:env;type:varchar(32);index:config_env;comment:环境;"`
	EntryId         int    `json:"entry_id" gorm:"column:entry_id;type:int(11);default:0;comment:执行队列编号;"`
	Type            int    `json:"type" gorm:"column:type;type:tinyint(2);default:1;index:config_env,priority:11;comment:类型：1.周期任务（默认）、2.单次任务;"`
	Name            string `json:"name" gorm:"column:name;type:varchar(255);default:'';comment:任务名称;"`
	Spec            string `json:"spec" gorm:"column:spec;type:varchar(32);default:'';comment:执行时间 表达式;"`
	Protocol        int    `json:"protocol" gorm:"column:protocol;type:tinyint(2);default:0;comment:协议：1.http、2.grpc、3.系统命令、4.sql执行;"`
	Command         []byte `json:"command" gorm:"column:command;type:json;default:null;comment:命令内容;"`
	CommandHash     string `json:"command_hash" gorm:"column:command_hash;type:char(32);default:'';comment:命令内容hash;"`
	AfterTmpl       string `json:"after_tmpl" gorm:"column:after_tmpl;type:varchar(1024);default:'';comment:结束模板;"`
	Remark          string `json:"remark" gorm:"column:remark;type:varchar(255);comment:备注;"`
	Status          int    `json:"status" gorm:"column:status;type:tinyint(2);default:1;comment:状态：1.停止、2.启用、3.完成、4.错误;"`
	StatusRemark    string `json:"status_remark" gorm:"column:status_remark;type:varchar(255);comment:状态变更描述;"`
	StatusDt        string `json:"status_dt" gorm:"column:status_dt;type:datetime;default:null;comment:状态变更时间;"`
	UpdateDt        string `json:"update_dt" gorm:"column:update_dt;type:datetime;default:null;comment:更新时间;"`
	CreateDt        string `json:"create_dt" gorm:"column:create_dt;type:datetime;default:null;comment:创建时间;"`
	MsgSet          []byte `json:"msg_set" gorm:"column:msg_set;type:json;default:null;comment:消息配置详情;"`
	EmptyNotMsg     int    `json:"empty_not_msg" gorm:"column:empty_not_msg;type:tinyint(2);default:2;comment:空结果不发消息：1.是、2.否;"`
	MsgSetHash      string `json:"msg_set_hash" gorm:"column:msg_set_hash;type:char(32);default:'';comment:消息配置hash;"`
	AfterSleep      int    `json:"after_sleep" gorm:"column:after_sleep;type:int(11);default:0;comment:延迟关闭;"`
	ErrRetryNum     int    `json:"err_retry_num" gorm:"column:err_retry_num;type:int(11);default:0;comment:错误重试次数;"`
	ErrRetrySleep   int    `json:"err_retry_sleep" gorm:"column:err_retry_sleep;type:int(11);default:0;comment:错误重试间隔/秒;"`
	ErrRetryMode    int    `json:"err_retry_mode" gorm:"column:err_retry_mode;type:int(11);default:1;comment:错误重试模式：1.固定间隔、2.增长间隔;"`
	VarFields       []byte `json:"var_fields" gorm:"column:var_fields;type:json;default:null;comment:参数变量;"`
	VarFieldsHash   string `json:"var_fields_hash" gorm:"column:var_fields_hash;type:char(32);default:'';comment:参数变量hash;"`
	CreateUserId    int    `json:"create_user_id" gorm:"column:create_user_id;type:int(11);default:0;comment:创建人;"`
	CreateUserName  string `json:"create_user_name" gorm:"column:create_user_name;type:varchar(64);default:'';comment:创建人名称;"`
	AuditUserId     int    `json:"audit_user_id" gorm:"column:audit_user_id;type:int(11);default:0;comment:审核人;"`
	AuditUserName   string `json:"audit_user_name" gorm:"column:audit_user_name;type:varchar(64);default:'';comment:审核人名称;"`
	HandleUserIds   string `json:"handle_user_ids" gorm:"column:handle_user_ids;type:varchar(255);default:'';comment:处理人,多选id逗号分隔;"`
	HandleUserNames string `json:"handle_user_names" gorm:"column:handle_user_names;type:varchar(500);default:'';comment:处理人名称,多选id逗号分隔;"`
	TagIds          string `json:"tag_ids" gorm:"column:tag_ids;type:varchar(255);default:'';comment:标签id,多选逗号分隔;"`
	TagNames        string `json:"tag_names" gorm:"column:tag_names;type:varchar(500);default:'';comment:标签名称,多选逗号分隔;"`
}

func (m *CronConfig) TableName() string {
	return "cron_config"
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

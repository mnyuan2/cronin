package models

const (
	DisableActionStop = 1
	DisableActionOmit = 2
	DisableActionExec = 3

	ErrActionStop = 1
	ErrActionOmit = 2
)

// 停用时行为
var DisableActionMap = map[int]string{
	DisableActionStop: "停止",
	DisableActionOmit: "跳过",
	DisableActionExec: "执行",
}

// 错误时行为
var ErrActionMap = map[int]string{
	ErrActionStop: "停止",
	ErrActionOmit: "跳过",
}

type CronPipeline struct {
	Id                  int    `json:"id" gorm:"column:id;type:INTEGER;primary_key;comment:主键;"`
	Env                 string `json:"env" gorm:"column:env;type:varchar(32);index:pipeline_env;comment:环境;"`
	EntryId             int    `json:"entry_id" gorm:"column:entry_id;type:int(11);default:0;comment:执行队列编号;"`
	Type                int    `json:"type" gorm:"column:type;type:tinyint(2);default:1;comment:类型：1.周期任务（默认）、2.单次任务;"`
	Name                string `json:"name" gorm:"column:name;type:varchar(255);default:'';comment:任务名称;"`
	Spec                string `json:"spec" gorm:"column:spec;type:varchar(32);default:'';comment:执行时间 表达式;"`
	VarParams           string `json:"var_params" gorm:"column:var_params;type:varchar(5000);comment:变量参数;"`
	ConfigIds           []byte `json:"config_ids" gorm:"column:config_ids;type:json;default:null;comment:任务id集合;"`
	Configs             []byte `json:"configs" gorm:"column:configs;type:json;default:null;comment:任务集合;"` // 存储快照用于展示，执行时要取最新的任务信息
	ConfigDisableAction int    `json:"config_disable_action" gorm:"column:config_disable_action;type:tinyint(2);default:1;comment:任务停用行为：1.整体停止、2.忽略跳过、3.执行;"`
	ConfigErrAction     int    `json:"config_err_action" gorm:"column:config_err_action;type:tinyint(2);default:1;comment:任务结果错误行为：1.整体停止、2.忽略跳过"`
	Interval            int    `json:"interval" gorm:"column:interval;type:int(11);default:0;comment:执行间隔;"`
	Remark              string `json:"remark" gorm:"column:remark;type:varchar(255);comment:备注;"`
	Status              int    `json:"status" gorm:"column:status;type:tinyint(2);default:1;comment:状态：1.停止、2.启用、3.完成、4.错误;"`
	StatusRemark        string `json:"status_remark" gorm:"column:status_remark;type:varchar(255);comment:状态变更描述;"`
	StatusDt            string `json:"status_dt" gorm:"column:status_dt;type:datetime;default:null;comment:状态变更时间;"`
	UpdateDt            string `json:"update_dt" gorm:"column:update_dt;type:datetime;default:null;comment:更新时间;"`
	CreateDt            string `json:"create_dt" gorm:"column:create_dt;type:datetime;default:null;comment:创建时间;"`
	MsgSet              []byte `json:"msg_set" gorm:"column:msg_set;type:json;default:null;comment:消息配置详情;"`
	MsgSetHash          string `json:"msg_set_hash" gorm:"column:msg_set_hash;type:char(32);default:'';comment:消息配置hash;"`
	CreateUserId        int    `json:"create_user_id" gorm:"column:create_user_id;type:int(11);default:0;comment:创建人;"`
	CreateUserName      string `json:"create_user_name" gorm:"column:create_user_name;type:varchar(64);default:'';comment:创建人名称;"`
	AuditUserId         int    `json:"audit_user_id" gorm:"column:audit_user_id;type:int(11);default:0;comment:审核人;"`
	AuditUserName       string `json:"audit_user_name" gorm:"column:audit_user_name;type:varchar(64);default:'';comment:审核人名称;"`
	HandleUserIds       string `json:"handle_user_ids" gorm:"column:handle_user_ids;type:varchar(255);default:'';comment:处理人,多选id逗号分隔;"`
	HandleUserNames     string `json:"handle_user_names" gorm:"column:handle_user_names;type:varchar(500);default:'';comment:处理人名称,多选id逗号分隔;"`
}

func (m *CronPipeline) TableName() string {
	return "cron_pipeline"
}

func (m *CronPipeline) ConfigErrActionName() string {
	return ErrActionMap[m.ConfigErrAction]
}

package pb

// 列表
type CronPipelineListRequest struct {
	Type                 int    `form:"type"`
	Page                 int    `form:"page"`
	Size                 int    `form:"size"`
	Status               []int  `form:"status[]"`
	CreateUserIds        []int  `form:"create_user_ids[]"`
	HandleUserIds        []int  `form:"handle_user_ids[]"`
	CreateOrHandleUserId int    `form:"create_or_handle_user_id"`
	Name                 string `form:"name"`
}
type CronPipelineListReply struct {
	List []*CronPipelineListItem `json:"list"`
	Page *Page                   `json:"page"`
}
type CronPipelineListItem struct {
	Id                      int                   `json:"id"`
	EntryId                 int                   `json:"entry_id"`
	Name                    string                `json:"name"`
	Spec                    string                `json:"spec"`
	Remark                  string                `json:"remark"`
	Status                  int                   `json:"status"`
	StatusName              string                `json:"status_name"`
	StatusRemark            string                `json:"status_remark"`
	StatusDt                string                `json:"status_dt"`
	Type                    int                   `json:"type"`
	ConfigDisableAction     int                   `json:"config_disable_action"`
	ConfigDisableActionName string                `json:"config_disable_action_name"`
	ConfigErrAction         int                   `json:"config_err_action"`
	Interval                int                   `json:"interval"`
	TopNumber               int                   `json:"top_number"`       // 最近执行次数（最大5次）
	TopErrorNumber          int                   `json:"top_error_number"` // 最近执行次数中，失败的次数
	UpdateDt                string                `json:"update_dt"`
	VarParams               string                `json:"var_params"`
	ConfigIds               []int                 `json:"config_ids" gorm:"-"`
	Configs                 []*CronConfigListItem `json:"configs" gorm:"-"`
	MsgSet                  []*CronMsgSet         `json:"msg_set" gorm:"-"`
	ConfigIdsStr            []byte                `json:"-" gorm:"column:config_ids;"`
	ConfigsStr              []byte                `json:"-" gorm:"column:configs;"`
	MsgSetStr               []byte                `json:"-" gorm:"column:msg_set;"`
	HandleUserStr           []byte                `json:"-" gorm:"column:handle_user_ids;"`
	CreateUserId            int                   `json:"create_user_id"`
	CreateUserName          string                `json:"create_user_name" gorm:"-"`
	HandleUserIds           []int                 `json:"handle_user_ids" gorm:"-"` // 处理人
}

// 流水线设置
type CronPipelineSetRequest struct {
	Id                  int                   `json:"id"`                    // 主键
	Name                string                `json:"name"`                  // 任务名称
	Type                int                   `json:"type"`                  // 类型
	Spec                string                `json:"spec"`                  // 执行时间表达式
	ConfigIds           []int                 `json:"config_ids"`            // 任务id集合
	Configs             []*CronConfigListItem `json:"configs"`               // 任务集合
	ConfigDisableAction int                   `json:"config_disable_action"` //
	ConfigErrAction     int                   `json:"config_err_action"`     //
	Interval            int                   `json:"interval"`              // 执行间隔
	Status              int                   `json:"status"`                // 状态
	Remark              string                `json:"remark"`                // 备注
	VarParams           string                `json:"var_params"`            // 变量参数
	MsgSet              []*CronMsgSet         `json:"msg_set"`               // 消息设置
}
type CronPipelineSetReply struct {
	Id int `json:"id"`
}

// 流水线执行一下
type CronPipelineRunRequest struct {
	Id                  int                   `json:"id"`                    // 主键
	Name                string                `json:"name"`                  // 任务名称
	Type                int                   `json:"type"`                  // 类型
	Spec                string                `json:"spec"`                  // 执行时间表达式
	ConfigIds           []int                 `json:"config_ids"`            // 任务id集合
	Configs             []*CronConfigListItem `json:"configs"`               // 任务集合
	ConfigDisableAction int                   `json:"config_disable_action"` //
	ConfigErrAction     int                   `json:"config_err_action"`     //
	Interval            int                   `json:"interval"`              // 执行间隔
	Status              int                   `json:"status"`                // 状态
	Remark              string                `json:"remark"`                // 备注
	VarParams           string                `json:"var_params"`            // 变量参数
	MsgSet              []*CronMsgSet         `json:"msg_set"`               // 消息设置
}
type CronPipelineRunReply struct {
	Result string `json:"result"`
}

// 流水线详情
type CronPipelineDetailRequest struct {
	Id int `json:"id" form:"id"`
}
type CronPipelineDetailReply struct {
	Id                      int                   `json:"id"`
	EntryId                 int                   `json:"entry_id"`
	Name                    string                `json:"name"`
	Spec                    string                `json:"spec"`
	Remark                  string                `json:"remark"`
	Status                  int                   `json:"status"`
	StatusName              string                `json:"status_name"`
	StatusRemark            string                `json:"status_remark"`
	StatusDt                string                `json:"status_dt"`
	Type                    int                   `json:"type"`
	ConfigDisableAction     int                   `json:"config_disable_action"`
	ConfigDisableActionName string                `json:"config_disable_action_name"`
	ConfigErrAction         int                   `json:"config_err_action"`
	ConfigErrActionName     string                `json:"config_err_action_name"`
	Interval                int                   `json:"interval"`
	TopNumber               int                   `json:"top_number"`       // 最近执行次数（最大5次）
	TopErrorNumber          int                   `json:"top_error_number"` // 最近执行次数中，失败的次数
	UpdateDt                string                `json:"update_dt"`
	CreateDt                string                `json:"create_dt"`
	VarParams               string                `json:"var_params"`
	ConfigIds               []int                 `json:"config_ids"`
	Configs                 []*CronConfigListItem `json:"configs"`
	MsgSet                  []*CronMsgSet         `json:"msg_set"`
	CreateUserId            int                   `json:"create_user_id"`
	CreateUserName          string                `json:"create_user_name"`
	AuditUserId             int                   `json:"audit_user_id"`
	AuditUserName           string                `json:"audit_user_name"`
	HandleUserIds           []int                 `json:"handle_user_ids"` // 处理人
	HandleUserNames         string                `json:"handle_user_names"`
}

// 状态变更
type CronPipelineChangeStatusRequest struct {
	Id            int    `json:"id"`
	Status        int    `json:"status"` // 状态
	StatusRemark  string `json:"status_remark"`
	HandleUserIds []int  `json:"handle_user_ids"`
}

type CronPipelineChangeStatusReply struct{}

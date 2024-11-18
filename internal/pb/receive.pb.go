package pb

// 用户列表
type ReceiveListRequest struct {
	Page                 int    `form:"page"`
	Size                 int    `form:"size"`
	Status               []int  `form:"status[]"`
	CreateUserIds        []int  `form:"create_user_ids[]"`
	HandleUserIds        []int  `form:"handle_user_ids[]"`
	CreateOrHandleUserId int    `form:"create_or_handle_user_id"`
	Name                 string `form:"name"`
}
type ReceiveListReply struct {
	List []*ReceiveListItem `json:"list"`
	Page *Page              `json:"page"`
}
type ReceiveListItem struct {
	Id             int    `json:"id"`
	Name           string `json:"name"`
	Remark         string `json:"remark"`
	Status         int    `json:"status"`
	StatusName     string `json:"status_name"`
	StatusRemark   string `json:"status_remark"`
	StatusDt       string `json:"status_dt"`
	TopNumber      int    `json:"top_number"`       // 最近执行次数（最大5次）
	TopErrorNumber int    `json:"top_error_number"` // 最近执行次数中，失败的次数
	UpdateDt       string `json:"update_dt"`
	ConfigIds      []int  `json:"config_ids" gorm:"-"`
	ConfigIdsStr   []byte `json:"-" gorm:"column:config_ids;"`
	HandleUserStr  []byte `json:"-" gorm:"column:handle_user_ids;"`
	CreateUserId   int    `json:"create_user_id"`
	CreateUserName string `json:"create_user_name" gorm:"-"`
	HandleUserIds  []int  `json:"handle_user_ids" gorm:"-"` // 处理人
}

// 用户设置
type ReceiveSetRequest struct {
	Id                  int                `json:"id"`                    // 主键
	Name                string             `json:"name"`                  // 任务名称
	ConfigIds           []int              `json:"config_ids"`            // 任务id集合
	RuleConfig          []*ReceiveRuleItem `json:"rule_config"`           // 任务集合
	ConfigDisableAction int                `json:"config_disable_action"` //
	ConfigErrAction     int                `json:"config_err_action"`     //
	Interval            int                `json:"interval"`              // 执行间隔
	Status              int                `json:"status"`                // 状态
	Remark              string             `json:"remark"`                // 备注
	ReceiveTmpl         string             `json:"receive_tmpl"`          // 接收模板
	MsgSet              []*CronMsgSet      `json:"msg_set"`               // 消息设置
}
type ReceiveRuleItem struct {
	Rule   []*KvItem           `json:"rule"`   // 匹配规则
	Param  []*KvItem           `json:"param"`  // 参数关联
	Config *CronConfigListItem `json:"config"` // 任务描述
}
type ReceiveSetReply struct {
	Id int `json:"id"`
}

// 接收配置 详情
type ReceiveDetailRequest struct {
	Id int `json:"id" form:"id"`
}

// 接收配置 响应
type ReceiveDetailReply struct {
	Id                      int                `json:"id"`
	Name                    string             `json:"name"`
	Remark                  string             `json:"remark"`
	Status                  int                `json:"status"`
	StatusName              string             `json:"status_name"`
	StatusRemark            string             `json:"status_remark"`
	StatusDt                string             `json:"status_dt"`
	ConfigDisableAction     int                `json:"config_disable_action"`
	ConfigDisableActionName string             `json:"config_disable_action_name"`
	ConfigErrAction         int                `json:"config_err_action"`
	ConfigErrActionName     string             `json:"config_err_action_name"`
	Interval                int                `json:"interval"`
	TopNumber               int                `json:"top_number"`
	TopErrorNumber          int                `json:"top_error_number"`
	UpdateDt                string             `json:"update_dt"`
	CreateDt                string             `json:"create_dt"`
	ReceiveTmpl             string             `json:"receive_tmpl"`
	ConfigIds               []int              `json:"config_ids"`
	RuleConfig              []*ReceiveRuleItem `json:"rule_config"`
	MsgSet                  []*CronMsgSet      `json:"msg_set"`
	CreateUserId            int                `json:"create_user_id"`
	AuditUserId             int                `json:"audit_user_id"`
	HandleUserIds           []int              `json:"handle_user_ids"`
}

// 接收配置 状态
type ReceiveChangeStatusRequest struct {
	Id            int    `json:"id"`
	Status        int    `json:"status"` // 状态
	StatusRemark  string `json:"status_remark"`
	HandleUserIds []int  `json:"handle_user_ids"`
}
type ReceiveChangeStatusReply struct {
}

// 接收配置 触发钩子 请求
type ReceiveWebhookRequest struct {
	Id   int    `json:"id"`
	Body []byte `json:"body"`
}

// 接收配置 触发钩子 响应
type ReceiveWebhookReply struct {
}

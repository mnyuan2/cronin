package pb

// 列表
type CronPipelineListRequest struct {
	Type int `form:"type"`
	Page int `form:"page"`
	Size int `form:"size"`
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
	TopNumber               int                   `json:"top_number"`       // 最近执行次数（最大5次）
	TopErrorNumber          int                   `json:"top_error_number"` // 最近执行次数中，失败的次数
	UpdateDt                string                `json:"update_dt"`
	ConfigIds               []int                 `json:"config_ids" gorm:"-"`
	Configs                 []*CronConfigListItem `json:"configs" gorm:"-"`
	MsgSet                  []*CronMsgSet         `json:"msg_set" gorm:"-"`
	ConfigIdsStr            []byte                `json:"-" gorm:"column:config_ids;"`
	ConfigsStr              []byte                `json:"-" gorm:"column:configs;"`
	MsgSetStr               []byte                `json:"-" gorm:"column:msg_set;"`
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
	Status              int                   `json:"status"`                // 状态
	Remark              string                `json:"remark"`                // 备注
	MsgSet              []*CronMsgSet         `json:"msg_set"`               // 消息设置
}
type CronPipelineSetReply struct {
	Id int `json:"id"`
}

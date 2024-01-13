package pb

type KvItem struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// 任务列表
type CronConfigListRequest struct {
	Type int `form:"type"`
	Page int `form:"page"`
	Size int `form:"size"`
}
type CronConfigListReply struct {
	List []*CronConfigListItem `json:"list"`
	Page *Page                 `json:"page"`
}
type CronConfigListItem struct {
	Id             int                `json:"id"`
	EntryId        int                `json:"entry_id"`
	Name           string             `json:"name"`
	Spec           string             `json:"spec"`
	Protocol       int                `json:"protocol"`
	ProtocolName   string             `json:"protocol_name"`
	Remark         string             `json:"remark"`
	Status         int                `json:"status"`
	StatusName     string             `json:"status_name"`
	StatusRemark   string             `json:"status_remark"`
	StatusDt       string             `json:"status_dt"`
	Type           int                `json:"type"`
	TopNumber      int                `json:"top_number"`       // 最近执行次数（最大5次）
	TopErrorNumber int                `json:"top_error_number"` // 最近执行次数中，失败的次数
	UpdateDt       string             `json:"update_dt"`
	Command        *CronConfigCommand `json:"command" gorm:"-"`
	CommandStr     string             `json:"-" gorm:"column:command;"` // 这里只能读取字符串后，载入到结构体
}

// 任务设置
type CronConfigSetRequest struct {
	Id       int                `json:"id,omitempty"`       // 主键
	Name     string             `json:"name,omitempty"`     // 任务名称
	Type     int                `json:"type"`               // 类型
	Spec     string             `json:"spec"`               // 执行时间表达式
	Protocol int                `json:"protocol,omitempty"` // 协议：1.http、2.grpc、3.系统命令
	Command  *CronConfigCommand `json:"command,omitempty"`  // 命令
	Status   int                `json:"status"`             // 状态
	Remark   string             `json:"remark"`
}
type CronConfigSetResponse struct {
	Id int `json:"id"`
}
type CronConfigCommand struct {
	Http *CronHttp `json:"http"`
	Rpc  *CronRpc  `json:"rpc"`
	Cmd  string    `json:"cmd"`
	Sql  *CronSql  `json:"sql"`
}

type CronHttp struct {
	Method string    `json:"method"`
	Url    string    `json:"url"`
	Body   string    `json:"body"`
	Header []*KvItem `json:"header"`
}
type CronRpc struct {
	Proto   string   `json:"proto"`   // proto定义文件类容
	Method  string   `json:"method"`  // 执行类型：rpc、grpc
	Addr    string   `json:"addr"`    // 地址，包含端口
	Action  string   `json:"action"`  // 方法
	Actions []string `json:"actions"` // 方法集合，辅助
	Header  []string `json:"header"`  // 请求头
	Body    string   `json:"body"`    // 请求参数
}

// sql任务配置
type CronSql struct {
	Driver    string         `json:"driver"`     // 驱动，默认mysql
	Source    *CronSqlSource `json:"source"`     // 具体链接配置
	ErrAction int            `json:"err_action"` // 错误后行为
	Statement []string       `json:"statement"`  // sql语句多条
}

// CronSqlSource sql任务 来源配置
type CronSqlSource struct {
	Id       int    `json:"id"`
	Title    string `json:"title"`
	Hostname string `json:"hostname"`
	Database string `json:"database"`
	Username string `json:"username"`
	Password string `json:"password"`
	Port     string `json:"port"`
}

// 已注册列表
type CronConfigRegisterListRequest struct{}
type CronConfigRegisterListResponse struct {
	List []*CronConfigListItem `json:"list"`
}

// 任务设置
type CronConfigRunRequest struct {
	Name     string             `json:"name,omitempty"`     // 任务名称
	Type     int                `json:"type"`               // 类型
	Spec     string             `json:"spec"`               // 执行时间表达式
	Protocol int                `json:"protocol,omitempty"` // 协议：1.http、2.grpc、3.系统命令
	Command  *CronConfigCommand `json:"command,omitempty"`  // 命令
	Status   int                `json:"status"`             // 状态
	Remark   string             `json:"remark"`
}
type CronConfigRunResponse struct {
	Result string `json:"result"`
}

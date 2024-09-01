package pb

type KvItem struct {
	Key    string `json:"key"`
	Value  string `json:"value"`
	Remark string `json:"remark"`
}

// 任务语句
type CronStatement struct {
	Type    string `json:"type"`
	Local   string `json:"local"`    // 本地输入
	Git     *Git   `json:"git"`      // git输入
	IsBatch int    `json:"is_batch"` // 是否批量解析
}

type Git struct {
	LinkId  int      `json:"link_id"` // 连接配置id
	Owner   string   `json:"owner"`   // 仓库所属空间
	Project string   `json:"project"` // 仓库项目
	Path    []string `json:"path"`    // 文件的路径
	Ref     string   `json:"ref"`     // 分支、tag或commit。默认: 仓库的默认分支(通常是master)
}

type GitEvent struct {
	Id      int              `json:"id"`       // 事件id
	PRMerge *GitEventPRMerge `json:"pr_merge"` // pr合并内容
}

type GitEventPRMerge struct {
	Owner string `json:"owner"` // 空间地址
	Repo  string `json:"repo"`  // 项目名称（仓库路径）
	// 第几个PR，即本仓库PR的序数
	Number string `json:"number"`
	// 可选。合并PR的方法，merge（合并所有提交）、squash（扁平化分支合并）和rebase（变基并合并）。默认为merge。
	MergeMethod string `json:"merge_method"`
	// 可选。合并PR后是否删除源分支，默认false（不删除）
	PruneSourceBranch bool `json:"prune_source_branch"`
	// 可选。合并标题，默认为PR的标题
	Title string `json:"title"`
	// 可选。合并描述，默认为 "Merge pull request !{pr_id} from {author}/{source_branch}"，与页面显示的默认一致。
	Description string `json:"description"`
}

// 任务列表
type CronConfigListRequest struct {
	Ids                  []int  `form:"ids[]"`
	Type                 int    `form:"type"`
	Page                 int    `form:"page"`
	Size                 int    `form:"size"`
	Protocol             []int  `form:"protocol[]"`
	Status               []int  `form:"status[]"`
	CreateUserIds        []int  `form:"create_user_ids[]"`
	HandleUserIds        []int  `form:"handle_user_ids[]"`
	CreateOrHandleUserId int    `form:"create_or_handle_user_id"`
	Name                 string `form:"name"`
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
	TypeName       string             `json:"type_name"`
	TopNumber      int                `json:"top_number"`       // 最近执行次数（最大5次）
	TopErrorNumber int                `json:"top_error_number"` // 最近执行次数中，失败的次数
	UpdateDt       string             `json:"update_dt"`
	AfterTmpl      string             `json:"after_tmpl"`          // 结果模板
	VarFields      []*KvItem          `json:"var_fields" gorm:"-"` // 定义变量参数
	Command        *CronConfigCommand `json:"command" gorm:"-"`
	MsgSet         []*CronMsgSet      `json:"msg_set" gorm:"-"`
	VarFieldsStr   []byte             `json:"-" gorm:"column:var_fields;"`
	CommandStr     []byte             `json:"-" gorm:"column:command;"` // 这里只能读取字符串后，载入到结构体
	MsgSetStr      []byte             `json:"-" gorm:"column:msg_set;"`
	HandleUserStr  []byte             `json:"-" gorm:"column:handle_user_ids;"`
	CreateUserId   int                `json:"create_user_id"`
	CreateUserName string             `json:"create_user_name" gorm:"-"`
	HandleUserIds  []int              `json:"handle_user_ids" gorm:"-"` // 处理人
}

type CronConfigDetailRequest struct {
	Id int `json:"id" form:"id"`
}
type CronConfigDetailReply struct {
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
	TypeName       string             `json:"type_name"`
	TopNumber      int                `json:"top_number"`       // 最近执行次数（最大5次）
	TopErrorNumber int                `json:"top_error_number"` // 最近执行次数中，失败的次数
	UpdateDt       string             `json:"update_dt"`
	CreateDt       string             `json:"create_dt"`
	AfterTmpl      string             `json:"after_tmpl"` // 结果模板
	VarFields      []*KvItem          `json:"var_fields"` // 定义变量参数
	Command        *CronConfigCommand `json:"command"`
	MsgSet         []*CronMsgSet      `json:"msg_set"`
	CreateUserId   int                `json:"create_user_id"`
	AuditUserId    int                `json:"audit_user_id"`
	HandleUserIds  []int              `json:"handle_user_ids"` // 处理人
}

// 任务设置
type CronConfigSetRequest struct {
	Id            int                `json:"id,omitempty"`       // 主键
	Name          string             `json:"name,omitempty"`     // 任务名称
	Type          int                `json:"type"`               // 类型
	Spec          string             `json:"spec"`               // 执行时间表达式
	Protocol      int                `json:"protocol,omitempty"` // 协议：1.http、2.grpc、3.系统命令
	AfterTmpl     string             `json:"after_tmpl"`         // 结果模板
	VarFields     []*KvItem          `json:"var_fields"`         // 定义变量参数
	Command       *CronConfigCommand `json:"command,omitempty"`  // 命令
	Status        int                `json:"status"`             // 状态
	StatusRemark  string             `json:"status_remark"`      // （审核时）状态备注
	HandleUserIds []int              `json:"handle_user_ids"`    // 处理人
	Remark        string             `json:"remark"`             // 备注
	MsgSet        []*CronMsgSet      `json:"msg_set"`            // 消息设置
}
type CronConfigSetResponse struct {
	Id int `json:"id"`
}

type CronMsgSet struct {
	MsgId         int   `json:"msg_id"`
	Status        int   `json:"status"`
	NotifyUserIds []int `json:"notify_user_ids"`
}
type CronConfigCommand struct {
	Http    *CronHttp    `json:"http"`
	Rpc     *CronRpc     `json:"rpc"`
	Cmd     *CronCmd     `json:"cmd"`
	Sql     *CronSql     `json:"sql"`
	Jenkins *CronJenkins `json:"jenkins"`
	Git     *CronGit     `json:"git"`
}

type CronCmd struct {
	Host      *SettingHostSource `json:"host"`
	Type      string             `json:"type"`
	Origin    string             `json:"origin"` // 语句来源
	Statement *CronStatement     `json:"statement"`
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
	Driver        string           `json:"driver"`          // 驱动，默认mysql
	Source        *CronSqlSource   `json:"source"`          // 具体链接配置
	ErrAction     int              `json:"err_action"`      // 错误后行为
	ErrActionName string           `json:"err_action_name"` // 错误后行为名称
	Interval      int64            `json:"interval"`        // 执行间隔
	Origin        string           `json:"origin"`          // 语句来源
	Statement     []*CronStatement `json:"statement"`
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

type CronJenkinsSource struct {
	Id int `json:"id"`
}

type CronJenkins struct {
	Source *CronJenkinsSource `json:"source"` // 具体链接配置
	Name   string             `json:"name"`   // 项目名称
	Params []*KvItem          `json:"params"` // 参数
}

type CronGit struct {
	LinkId int         `json:"link_id"`
	Events []*GitEvent `json:"events"`
}

// 已注册列表
type CronConfigRegisterListRequest struct{}
type CronConfigRegisterListResponse struct {
	List []*CronConfigListItem `json:"list"`
}

// 任务设置
type CronConfigRunRequest struct {
	Id        int                `json:"id"`                 // 任务编号
	Name      string             `json:"name,omitempty"`     // 任务名称
	Type      int                `json:"type"`               // 类型
	Spec      string             `json:"spec"`               // 执行时间表达式
	Protocol  int                `json:"protocol,omitempty"` // 协议：1.http、2.grpc、3.系统命令
	Command   *CronConfigCommand `json:"command,omitempty"`  // 命令
	Status    int                `json:"status"`             // 状态
	Remark    string             `json:"remark"`
	AfterTmpl string             `json:"after_tmpl"`          // 结果模板
	VarFields []*KvItem          `json:"var_fields" gorm:"-"` // 定义变量参数
	MsgSet    []*CronMsgSet      `json:"msg_set"`             // 消息设置
}
type CronConfigRunResponse struct {
	Result string `json:"result"`
}

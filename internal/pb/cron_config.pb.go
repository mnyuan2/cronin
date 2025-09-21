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
	Id         int                 `json:"id"`                    // 事件id
	PRDetail   *GitEventPRMerge    `json:"pr_detail,omitempty"`   // pr详情
	PRIsMerge  *GitEventPRMerge    `json:"pr_is_merge,omitempty"` // pr是否合并
	PRMerge    *GitEventPRMerge    `json:"pr_merge"`              // pr合并内容
	FileUpdate *GitEventFileUpdate `json:"file_update,omitempty"` // 文件更新
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
	// 状态
	State string `json:"state"`
}
type GitEventFileUpdate struct {
	Owner   string `json:"owner"`   // 空间地址
	Repo    string `json:"repo"`    // 项目名称（仓库路径）
	Path    string `json:"path"`    // 文件路径
	Content string `json:"content"` // 文件内容
	Message string `json:"message"` // 提交描述
	Branch  string `json:"branch"`  // 分支名称
}

type GetEventPRList struct {
	Owner   string `json:"owner"`    // 空间地址
	Repo    string `json:"repo"`     // 项目名称（仓库路径）
	State   string `json:"state"`    // 可选。Pull Request 状态: open、closed、merged、all
	Head    string `json:"head"`     // 可选。Pull Request 提交的源分支。格式：branch 或者：username:branch
	Base    string `json:"base"`     // 可选。Pull Request 提交目标分支的名称。
	Page    int    `json:"page"`     // 当前的页码
	PerPage int    `json:"per_page"` // 每页的数量，最大为 100
}

// 任务列表
type CronConfigListRequest struct {
	IsExecRatio          int    `json:"is_exec_ratio"` // 执行率
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
	TagIds               []int  `json:"tag_ids" form:"tag_ids[]"`
	SourceIds            []int  `json:"source_ids" form:"source_ids[]"`
	Env                  string `form:"env"`
}
type CronConfigListReply struct {
	List []*CronConfigListItem `json:"list"`
	Page *Page                 `json:"page"`
}
type CronConfigListItem struct {
	Id             int       `json:"id"`
	Env            string    `json:"env"`
	Name           string    `json:"name"`
	Spec           string    `json:"spec"`
	Protocol       int       `json:"protocol"`
	ProtocolName   string    `json:"protocol_name"`
	Remark         string    `json:"remark"`
	Status         int       `json:"status"`
	StatusName     string    `json:"status_name"`
	StatusRemark   string    `json:"status_remark"`
	StatusDt       string    `json:"status_dt"`
	Type           int       `json:"type"`
	TypeName       string    `json:"type_name"`
	TopNumber      int       `json:"top_number"`       // 最近执行次数（最大5次）
	TopErrorNumber int       `json:"top_error_number"` // 最近执行次数中，失败的次数
	UpdateDt       string    `json:"update_dt"`
	VarFields      []*KvItem `json:"var_fields" gorm:"-"` // 定义变量参数
	VarFieldsStr   []byte    `json:"-" gorm:"column:var_fields;"`
	HandleUserStr  []byte    `json:"-" gorm:"column:handle_user_ids;"`
	TagIdsStr      []byte    `json:"-" gorm:"column:tag_ids"`
	CreateUserId   int       `json:"create_user_id"`
	CreateUserName string    `json:"create_user_name" gorm:"-"`
	HandleUserIds  []int     `json:"handle_user_ids" gorm:"-"` // 处理人
	TagIds         []int     `json:"tag_ids" gorm:"-"`         //
	TagNames       string    `json:"tag_names"`                //
}

// 任务匹配列表
type CronMatchListRequest struct {
	Search     []*CronMatchListSearchItem `json:"search"`
	SearchText string                     `json:"search_text"`
}
type CronMatchListSearchItem struct {
	Type  string   `json:"type"`
	Value []string `json:"value"`
}
type CronMatchListReply struct {
	List      []*CronConfigListItem `json:"list"`
	VarParams map[string]string     `json:"var_params"` // 这里还要返回pr的变量实现，还有就是任务中所有包含的变量。
}

type CronConfigDetailRequest struct {
	Id        int    `json:"id" form:"id"`
	VarParams string `json:"var_params" form:"var_params"`
}
type CronConfigDetailReply struct {
	Id               int                `json:"id"`
	Env              string             `json:"env"`
	EntryId          int                `json:"entry_id"`
	Name             string             `json:"name"`
	Spec             string             `json:"spec"`
	Protocol         int                `json:"protocol"`
	ProtocolName     string             `json:"protocol_name"`
	Remark           string             `json:"remark"`
	Status           int                `json:"status"`
	StatusName       string             `json:"status_name"`
	StatusRemark     string             `json:"status_remark"`
	StatusDt         string             `json:"status_dt"`
	Type             int                `json:"type"`
	TypeName         string             `json:"type_name"`
	TopNumber        int                `json:"top_number"`       // 最近执行次数（最大5次）
	TopErrorNumber   int                `json:"top_error_number"` // 最近执行次数中，失败的次数
	UpdateDt         string             `json:"update_dt"`
	CreateDt         string             `json:"create_dt"`
	AfterTmpl        string             `json:"after_tmpl"` // 结果模板
	VarFields        []*KvItem          `json:"var_fields"` // 定义变量参数
	Command          *CronConfigCommand `json:"command"`
	MsgSet           []*CronMsgSet      `json:"msg_set"`
	EmptyNotMsg      int                `json:"empty_not_msg"` // 空结果不发送消息：2.空结果发送消息（默认）、1.空结果不发送消息
	AfterSleep       int                `json:"after_sleep"`
	ErrRetryNum      int                `json:"err_retry_num"`
	ErrRetrySleep    int                `json:"err_retry_sleep"` // 错误重试间隔
	ErrRetryMode     int                `json:"err_retry_mode"`  // 错误重试模式：1.固定间隔、2.增长间隔
	ErrRetryModeName string             `json:"err_retry_mode_name"`
	CreateUserId     int                `json:"create_user_id"`
	CreateUserName   string             `json:"create_user_name"`
	AuditUserId      int                `json:"audit_user_id"`
	AuditUserName    string             `json:"audit_user_name"`
	HandleUserIds    []int              `json:"handle_user_ids"` // 处理人
	TagIds           []int              `json:"tag_ids"`
	TagNames         string             `json:"tag_names"`
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
	TagIds        []int              `json:"tag_ids"`            // 标签
	MsgSet        []*CronMsgSet      `json:"msg_set"`            // 消息设置
	EmptyNotMsg   int                `json:"empty_not_msg"`      // 空结果不发消息
	AfterSleep    int                `json:"after_sleep"`        // 延迟关闭
	ErrRetryNum   int                `json:"err_retry_num"`      // 错误重试次数
	ErrRetrySleep int                `json:"err_retry_sleep"`    // 错误重试间隔
	ErrRetryMode  int                `json:"err_retry_mode"`     // 错误重试模式：1.固定间隔、2.增长间隔
}
type CronConfigSetResponse struct {
	Id int `json:"id"`
}

type CronMsgSet struct {
	MsgId         int   `json:"msg_id"`
	Status        []int `json:"status"`
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
	Method  string    `json:"method"`
	Url     string    `json:"url"`
	Body    string    `json:"body"`
	Header  []*KvItem `json:"header"`
	Timeout int       `json:"timeout"`
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
	GitSourceId   int              `json:"git_source_id"`   // git资源id
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
type CronJenkinsParamsGroup struct {
	EnableRule string    `json:"enable_rule"`
	Params     []*KvItem `json:"params"`
}

type CronJenkins struct {
	Source      *CronJenkinsSource        `json:"source"`       // 具体链接配置
	Name        string                    `json:"name"`         // 项目名称
	ParamsMode  int                       `json:"params_mode"`  // 参数模式: 1.参数、2.参数组
	Params      []*KvItem                 `json:"params"`       // 参数
	ParamsGroup []*CronJenkinsParamsGroup `json:"params_group"` // 参数组
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
	Id          int                `json:"id"`                 // 任务编号
	Name        string             `json:"name,omitempty"`     // 任务名称
	Type        int                `json:"type"`               // 类型
	Spec        string             `json:"spec"`               // 执行时间表达式
	Protocol    int                `json:"protocol,omitempty"` // 协议：1.http、2.grpc、3.系统命令
	Command     *CronConfigCommand `json:"command,omitempty"`  // 命令
	Status      int                `json:"status"`             // 状态
	Remark      string             `json:"remark"`
	AfterTmpl   string             `json:"after_tmpl"`          // 结果模板
	VarFields   []*KvItem          `json:"var_fields" gorm:"-"` // 定义变量参数
	MsgSet      []*CronMsgSet      `json:"msg_set"`             // 消息设置
	EmptyNotMsg int                `json:"empty_not_msg"`
	AfterSleep  int                `json:"after_sleep"`
}
type CronConfigRunResponse struct {
	Result string `json:"result"`
}

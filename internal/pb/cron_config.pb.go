package pb

import (
	"encoding/json"
	"strconv"
)

type KvItem struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// 任务语句
type CronStatement struct {
	Type  string `json:"type"`
	Local string `json:"local"` // 本地输入
	Git   *Git   `json:"git"`   // git输入
}

type Git struct {
	LinkId  int      `json:"link_id"` // 连接配置id
	Owner   string   `json:"owner"`   // 仓库所属空间
	Project string   `json:"project"` // 仓库项目
	Path    []string `json:"path"`    // 文件的路径
	Ref     string   `json:"ref"`     // 分支、tag或commit。默认: 仓库的默认分支(通常是master)
}

// 任务列表
type CronConfigListRequest struct {
	Ids  []int `form:"ids[]"`
	Type int   `form:"type"`
	Page int   `form:"page"`
	Size int   `form:"size"`
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
	Command        *CronConfigCommand `json:"command" gorm:"-"`
	MsgSet         []*CronMsgSet      `json:"msg_set" gorm:"-"`
	CommandStr     []byte             `json:"-" gorm:"column:command;"` // 这里只能读取字符串后，载入到结构体
	MsgSetStr      []byte             `json:"-" gorm:"column:msg_set;"`
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
	Remark   string             `json:"remark"`             // 备注
	MsgSet   []*CronMsgSet      `json:"msg_set"`            // 消息设置
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

// 已注册列表
type CronConfigRegisterListRequest struct{}
type CronConfigRegisterListResponse struct {
	List []*CronConfigListItem `json:"list"`
}

// 任务设置
type CronConfigRunRequest struct {
	Id       int                `json:"id"`                 // 任务编号
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

// UnmarshalJSON 转换前端status字符串到int
func (c *CronConfigSetRequest) UnmarshalJSON(data []byte) error {
	// 定义一个临时结构体用于解析 JSON
	type Alias CronConfigSetRequest
	aux := &struct {
		Status string `json:"status"`
		*Alias
	}{
		Alias: (*Alias)(c),
	}

	// 解析 JSON 数据到临时结构体
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	// 将字符串转换为 int
	status, err := strconv.Atoi(aux.Status)
	if err != nil {
		return err
	}
	c.Status = status

	return nil
}

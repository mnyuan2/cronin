package models

type CronProtocol int
type CronStatus int

const (
	ProtocolHttp = 1 // http
	ProtocolRpc  = 2 // rpc
	ProtocolCmd  = 3 // 命令行 cmd
)

var ProtocolMap = map[int]string{
	ProtocolHttp: "http",
	ProtocolRpc:  "rpc",
	ProtocolCmd:  "cmd",
}

const (
	StatusDisable = 1 // 停用
	StatusActive  = 2 // 激活
)

var ConfStatusMap = map[int]string{
	StatusDisable: "停用",
	StatusActive:  "激活",
}

const (
	TypeCycle = 1 // 周期
	TypeOnce  = 2 // 单次
)

var ConfTypeMap = map[int]string{
	TypeCycle: "周期",
	TypeOnce:  "单次",
}

type CronConfig struct {
	Id       int    `json:"id,omitempty"`        // 主键
	EntryId  int    `json:"entry_id,omitempty"`  // 执行队列编号
	Name     string `json:"name,omitempty"`      // 任务名称
	Spec     string `json:"spec"`                // 执行时间 表达式
	Protocol int    `json:"protocol,omitempty"`  // 协议：1.http、2.grpc、3.系统命令
	Command  string `json:"command,omitempty"`   // 命令
	Status   int    `json:"status,omitempty"`    // 状态：1.停止、2.启用
	Type     int    `json:"type"`                // 类型：1.周期任务（默认）、2.单次任务
	Remark   string `json:"remark"`              // 备注
	CreateDt string `json:"create_dt,omitempty"` // 创建时间
	UpdateDt string `json:"update_dt,omitempty"` // 更新时间
}

func (m *CronConfig) GetProtocolName() string {
	return ProtocolMap[m.Protocol]
}

func (m *CronConfig) GetStatusName() string {
	return ConfStatusMap[m.Status]
}

func (m *CronConfig) GetTypeName() string {
	return ConfTypeMap[m.Type]
}

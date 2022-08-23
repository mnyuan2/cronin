package models

type CronProtocol int
type CronStatus int

const (
	ProtocolHttp CronProtocol = 1 // http
	ProtocolRpc                   // rpc
	ProtocolCmd                   // 命令行 cmd
)
const (
	StatusDisable CronStatus = 1 // 停用
	StatusActive                 // 激活
)

type CronConfig struct {
	Id       int          `json:"id,omitempty"`        // 主键
	Name     string       `json:"name,omitempty"`      // 任务名称
	Spec     string       `json:"spec"`                // 执行时间 表达式
	Protocol CronProtocol `json:"protocol,omitempty"`  // 协议：1.http、2.grpc、3.系统命令
	Command  string       `json:"command,omitempty"`   // 命令
	Status   CronStatus   `json:"status,omitempty"`    // 状态：1.停止、2.启用
	Remark   string       `json:"remark"`              // 备注
	CreateDt string       `json:"create_dt,omitempty"` // 创建时间
	UpdateDt string       `json:"update_dt,omitempty"` // 更新时间
}

package enum

const (
	DicEnv              = 2    // 程序环境
	DicMsg              = 3    // 消息模板
	DicUser             = 4    // 用户
	DicRole             = 5    // 角色列表
	DicSqlSource        = 11   // sql 资源连接
	DicJenkinsSource    = 12   // jenkins 资源连接
	DicGitSource        = 13   // git 连接
	DicHostSource       = 14   // 主机 连接
	DicCmdType          = 1001 // 命令行类型
	DicGitEvent         = 1002 // git 事件
	DicSqlDriver        = 1003 // sql 驱动
	DicConfigStatus     = 1004 // 任务状态
	DicProtocolType     = 1005 // 协议类型
	DicReceiveDataField = 1006 // 接收字段
)

const (
	StatusDisable = 1 // 停用
	StatusActive  = 2 // 激活
	StatusDelete  = 9 // 删除
)

// 通用状态
var StatusMap = map[int]string{
	StatusDisable: "停用",
	StatusActive:  "激活",
}

const (
	BoolYes = 1 // 是
	BoolNot = 2 // 否
)

var BoolMap = map[int]string{
	BoolYes: "是",
	BoolNot: "否",
}

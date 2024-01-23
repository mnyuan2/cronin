package enum

const (
	DicSqlSource = 1 // sql连接源
	DicEnv       = 2 // 程序环境
	DicMsg       = 3 // 消息模板
	DicUser      = 4 // 用户
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

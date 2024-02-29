package enum

const (
	DicEnv           = 2  // 程序环境
	DicMsg           = 3  // 消息模板
	DicUser          = 4  // 用户
	DicSqlSource     = 11 // sql 资源连接
	DicJenkinsSource = 12 // jenkins 资源连接
	DicGitSource     = 13 // git 连接
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

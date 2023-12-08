package models

const (
	SqlErrActionAbort   = 1 // 终止
	SqlErrActionProceed = 2 // 继续
)

var SqlErrActionMap = map[int]string{
	SqlErrActionAbort:   "错误终止任务",
	SqlErrActionProceed: "错误跳过继续",
}

// sql驱动
const (
	SqlSourceMysql = "mysql"
)

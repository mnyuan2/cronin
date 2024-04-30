package enum

const (
	SqlStatementSourceLocal = "local" // sql语句来源 本地输入
	SqlStatementSourceGit   = "git"   // 远程 git
)

const (
	SqlDriverMysql      = "mysql"
	SqlDriverClickhouse = "clickhouse"
)

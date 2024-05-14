package models

// 任务公共表
type CronTable struct {
	Id       int    `json:"id" gorm:"column:id;type:int(11);primary_key;comment:主键;"`
	Env      string `json:"env" gorm:"column:env;type:varchar(32);index:env;comment:环境;"`
	JoinType string `json:"join_type" gorm:"column:join_type;type:varchar(32);comment:连接表;"`
	JoinId   int    `json:"join_id" gorm:"column:env;type:varchar(32);index:env;comment:连接id;"`
}

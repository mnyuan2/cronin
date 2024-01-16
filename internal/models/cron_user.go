package models

// 用户/人员
type CronUser struct {
	Id       int    `json:"id" gorm:"column:id;type:int(11);primary_key;comment:主键;"`
	Username string `json:"user_name" gorm:"column:username;type:varchar(255);default:'';comment:用户名;"`
	Mobile   string `json:"mobile" gorm:"column:mobile;type:varchar(24);default:'';comment:手机号;"`
	Sort     int    `json:"sort" gorm:"column:sort;type:int(11);default:1;comment:主键;"`
	UpdateDt string `json:"update_dt" gorm:"column:update_dt;type:datetime;default:null;comment:更新时间;"`
	CreateDt string `json:"create_dt" gorm:"column:create_dt;type:datetime;default:null;comment:创建时间;"`
}

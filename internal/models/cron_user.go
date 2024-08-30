package models

// 用户/人员
type CronUser struct {
	Id       int    `json:"id" gorm:"column:id;type:int(11);primary_key;comment:主键;"`
	Account  string `json:"account"  gorm:"column:account;type:varchar(64);default:'';comment:账号;"`
	Username string `json:"username" gorm:"column:username;type:varchar(64);default:'';comment:用户名;"`
	Mobile   string `json:"mobile" gorm:"column:mobile;type:varchar(24);default:'';comment:手机号;"`
	Sort     int    `json:"sort" gorm:"column:sort;type:int(11);default:1;comment:序号;"`
	Password string `json:"password" gorm:"column:password;type:varchar(64);default:'';comment:密码;"`
	Status   int    `json:"status" gorm:"column:status;type:tinyint(2);default:2;comment:状态：1.停止、2.启用、9.删除;"`
	UpdateDt string `json:"update_dt" gorm:"column:update_dt;type:datetime;default:null;comment:更新时间;"`
	CreateDt string `json:"create_dt" gorm:"column:create_dt;type:datetime;default:null;comment:创建时间;"`
	RoleIds  string `json:"role_ids" gorm:"column:role_ids;type:varchar(255);default:'';comment:拥有角色;"`
}

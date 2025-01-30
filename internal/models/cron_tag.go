package models

type CronTag struct {
	Id             int    `json:"id" gorm:"column:id;type:INTEGER;primary_key;comment:主键;"`
	Name           string `json:"name" gorm:"column:name;type:varchar(64);comment:名称;"`
	Remark         string `json:"remark" gorm:"column:remark;type:varchar(255);comment:描述;"`
	Status         int    `json:"status" gorm:"column:status;type:tinyint(2);default:2;comment:状态：9.删除、2.启用;"`
	CreateDt       string `json:"create_dt" gorm:"column:create_dt;type:datetime;default:null;comment:创建时间;"`
	CreateUserId   int    `json:"create_user_id" gorm:"column:create_user_id;type:int(11);default:0;comment:创建人;"`
	CreateUserName string `json:"create_user_name" gorm:"column:create_user_name;type:varchar(64);default:'';comment:创建人名称;"`
	UpdateDt       string `json:"update_dt" gorm:"column:update_dt;type:datetime;default:null;comment:更新时间;"`
	UpdateUserId   int    `json:"update_user_id" gorm:"column:update_user_id;type:int(11);default:0;comment:修改人;"`
	UpdateUserName string `json:"update_user_name" gorm:"column:update_user_name;type:varchar(64);default:'';comment:修改人名称;"`
}

func (m *CronTag) TableName() string {
	return "cron_tag"
}

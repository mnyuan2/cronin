package models

// 权限角色
type CronAuthRole struct {
	Id      int    `json:"id" gorm:"column:id;type:int(11);primary_key;comment:主键;"`
	Name    string `json:"name" gorm:"column:name;type:varchar(64);comment:角色名称;"`
	Remark  string `json:"remark" gorm:"column:remark;type:varchar(255);comment:描述;"`
	AuthIds string `json:"auth_ids" gorm:"column:auth_ids;type:text;comment:权限节点集合;"`
	Status  int    `json:"status" gorm:"column:status;type:tinyint(2);default:2;comment:状态：1.停止、2.启用;"`
}

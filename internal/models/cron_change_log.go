package models

const (
	LogTypeCreate    = 1 // 新增数据
	LogTypeUpdateDiy = 2 // 手动更新
	LogTypeUpdateSys = 3 // 系统更新
)

var LogTypeMap = map[int]string{
	LogTypeCreate:    "新增数据",
	LogTypeUpdateDiy: "手动更新",
	LogTypeUpdateSys: "系统更新",
}

// 变更日志
type CronChangeLog struct {
	Id             int    `json:"id" gorm:"column:id;type:int(11);primary_key;comment:主键;"`
	CreateDt       string `json:"create_dt" gorm:"column:create_dt;type:datetime;default:null;comment:创建时间;"`
	CreateUserId   int    `json:"create_user_id" gorm:"column:create_user_id;type:int(11);default:0;comment:创建人;"`
	CreateUserName string `json:"create_user_name" gorm:"column:create_user_name;type:varchar(64);default:'';comment:创建人名称;"`
	Type           int    `json:"type" gorm:"column:type;type:tinyint(2);default:0;comment:操作类型：1.新增、2.手动更新、3.自动更新;"`
	RefType        string `json:"ref_type" gorm:"column:ref_type;type:varchar(16);index:ref_type;comment:引用数据类型;"`
	RefId          int    `json:"ref_id" gorm:"column:ref_id;type:int(11);index:ref_type,ref_id:11;comment:引用数据id;"`
	Content        string `json:"content" gorm:"column:content;type:text;comment:变更内容;"`
}

// 变更字段
type ChangeLogField struct {
	Field      string `json:"field"`
	VType      string `json:"v_type"`
	OldVal     any    `json:"old_val"`
	NewVal     any    `json:"new_val"`
	FieldName  string `json:"field_name"`
	OldValName string `json:"old_val_name"`
	NewValName string `json:"new_val_name"`
}

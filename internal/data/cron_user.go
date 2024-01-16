package data

import (
	"context"
	"cron/internal/basic/db"
	"cron/internal/models"
)

type CronUserData struct {
	db        *db.MyDB
	tableName string
}

func NewCronUserData(ctx context.Context) *CronUserData {
	return &CronUserData{
		db:        db.New(ctx),
		tableName: "cron_user",
	}
}

// 查询列表数据
func (m *CronUserData) GetList(where *db.Where, page, size int, list interface{}) (total int64, err error) {
	str, args := where.Build()

	return m.db.Paginate(list, page, size, m.tableName, "*", "sort asc,id desc", str, args...)
}

func (m *CronUserData) Set(data *models.CronUser) error {
	if data.Id > 0 {
		return m.db.Where("id=?", data.Id).Omit("id", "create_dt").Updates(data).Error
	} else {
		return m.db.Create(data).Error
	}
}

//func (m *CronUserData) ChangeStatus(data *models.CronConfig, remark string) error {
//	data.UpdateDt = time.Now().Format(conv.FORMAT_DATETIME)
//	data.StatusDt = data.UpdateDt
//	data.StatusRemark = remark
//	return m.db.Where("id=?", data.Id).Select("status", "status_remark", "status_dt", "update_dt", "entry_id").Updates(data).Error
//}

func (m *CronUserData) GetOne(Id int) (data *models.CronUser, err error) {
	data = &models.CronUser{}
	return data, m.db.Where("id=?", Id).Take(data).Error
}

package data

import (
	"context"
	"cron/internal/basic/conv"
	"cron/internal/basic/db"
	"cron/internal/models"
	"time"
)

type CronConfigData struct {
	db        *db.Database
	tableName string
}

func NewCronConfigData(ctx context.Context) *CronConfigData {
	return &CronConfigData{
		db:        db.New(ctx),
		tableName: "cron_config",
	}
}

func (m *CronConfigData) GetList(where *db.Where, page, size int, list interface{}) (total int64, err error) {
	str, args := where.Build()

	return m.db.Read.Paginate(list, page, size, m.tableName, "*", "id", str, args...)
}

func (m *CronConfigData) Set(data *models.CronConfig) error {
	data.UpdateDt = time.Now().Format(conv.FORMAT_DATETIME)
	if data.Id > 0 {
		return m.db.Write.Where("id=?", data.Id).Omit("status", "status_remark", "status_dt", "entry_id").Updates(data).Error
	} else {
		data.CreateDt = time.Now().Format(conv.FORMAT_DATETIME)
		return m.db.Write.Omit("status_dt").Create(data).Error
	}
}

func (m *CronConfigData) ChangeStatus(data *models.CronConfig, remark string) error {
	data.UpdateDt = time.Now().Format(conv.FORMAT_DATETIME)
	data.StatusDt = data.UpdateDt
	data.StatusRemark = remark
	return m.db.Write.Where("id=?", data.Id).Select("status", "status_remark", "status_dt", "update_dt", "entry_id").Updates(data).Error
}

func (m *CronConfigData) GetOne(Id int) (data *models.CronConfig, err error) {
	data = &models.CronConfig{}
	return data, m.db.Read.Where("id=?", Id).Take(data).Error
}

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
		return m.db.Write.Where("id=?", data.Id).Updates(data).Error
	} else {
		data.CreateDt = time.Now().Format(conv.FORMAT_DATETIME)
		return m.db.Write.Create(data).Error
	}
}

func (m *CronConfigData) GetOne(Id int) (data *models.CronConfig, err error) {
	data = &models.CronConfig{}
	return data, m.db.Read.Where("id=?", Id).Take(data).Error
}

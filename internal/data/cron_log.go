package data

import (
	"context"
	"cron/internal/basic/db"
	"cron/internal/models"
)

type CronLogData struct {
	db        *db.Database
	tableName string
}

func NewCronLogData(ctx context.Context) *CronLogData {
	return &CronLogData{
		db:        db.New(ctx),
		tableName: "cron_log",
	}
}

// 添加数据
func (m *CronLogData) Add(data *models.CronLog) error {
	return m.db.Write.Create(data).Error
}

// 查询列表数据
func (m *CronLogData) GetList(where *db.Where, page, size int, list interface{}) (total int64, err error) {
	str, args := where.Build()

	return m.db.Read.Paginate(list, page, size, m.tableName, "*", "id desc", str, args...)
}

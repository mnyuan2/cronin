package data

import (
	"context"
	"cron/internal/basic/db"
	"cron/internal/models"
)

type CronTableData struct {
	db        *db.MyDB
	tableName string
}

func NewCronTableData(ctx context.Context) *CronTableData {
	return &CronTableData{
		db:        db.New(ctx),
		tableName: "cron_table",
	}
}

func (m CronTableData) Add(data *models.CronTable) error {
	return m.db.Create(data).Error
}

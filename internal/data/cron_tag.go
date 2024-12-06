package data

import (
	"context"
	"cron/internal/basic/db"
	"cron/internal/models"
)

type CronTagData struct {
	db        *db.MyDB
	tableName string
}

func NewCronTagData(ctx context.Context) *CronTagData {
	return &CronTagData{
		db:        db.New(ctx),
		tableName: "tag",
	}
}

func (m *CronTagData) List(where string, args ...any) ([]*models.CronTag, error) {
	list := []*models.CronTag{}
	return list, m.db.Where(where, args...).Find(&list).Error
}

func (m *CronTagData) GetOne(where string, args ...any) (*models.CronTag, error) {
	one := &models.CronTag{}
	return one, m.db.Where(where, args...).Take(one).Error
}

func (m *CronTagData) Set(data *models.CronTag) error {
	if data.Id > 0 {
		return m.db.Where("id=?", data.Id).
			Select("name", "remark", "update_dt", "update_user_id", "update_user_name").
			Updates(data).Error
	} else {
		return m.db.Create(data).Error
	}
}

func (m *CronTagData) ChangeStatus(data *models.CronTag) error {
	return m.db.Where("id=?", data.Id).Select("status", "update_dt", "update_user_id", "update_user_name").Updates(data).Error
}

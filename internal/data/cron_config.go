package data

import (
	"context"
	"cron/internal/basic/conv"
	"cron/internal/basic/db"
	"cron/internal/models"
	"time"
)

type CronConfigData struct {
	ctx       context.Context
	db        *db.MyDB
	tableName string
}

func NewCronConfigData(ctx context.Context) *CronConfigData {
	return &CronConfigData{
		ctx:       ctx,
		db:        db.New(ctx),
		tableName: "cron_config",
	}
}

func (m *CronConfigData) List(where *db.Where, size int) (list []*models.CronConfig, err error) {
	w, args := where.Build()
	list = []*models.CronConfig{}
	err = m.db.Where(w, args...).Limit(size).Find(&list).Error
	return list, err
}

func (m *CronConfigData) ListPage(where *db.Where, page, size int, list interface{}) (total int64, err error) {
	str, args := where.Build()

	return m.db.Paginate(list, page, size, m.tableName, "*", "id desc", str, args...)
}

func (m *CronConfigData) Set(data *models.CronConfig) error {
	data.UpdateDt = time.Now().Format(conv.FORMAT_DATETIME)
	if data.Id > 0 {
		return m.db.Where("id=?", data.Id).
			Omit("id", "entry_id", "env").
			Updates(data).Error
	} else {
		data.CreateDt = time.Now().Format(conv.FORMAT_DATETIME)
		return m.db.Omit("status_dt").Create(data).Error
	}
}

func (m *CronConfigData) ChangeStatus(data *models.CronConfig, remark string) error {
	data.UpdateDt = time.Now().Format(conv.FORMAT_DATETIME)
	data.StatusDt = data.UpdateDt
	data.StatusRemark = remark
	return m.db.Where("id=?", data.Id).Select("status", "status_remark", "status_dt", "update_dt", "entry_id", "handle_user_ids").Updates(data).Error
}

func (m *CronConfigData) SetEntryId(data *models.CronConfig) error {
	return m.db.Where("id=?", data.Id).Select("entry_id").Updates(data).Error
}

func (m *CronConfigData) GetOne(env string, Id int) (data *models.CronConfig, err error) {
	data = &models.CronConfig{}
	return data, m.db.Where("env=? and id=?", env, Id).Take(data).Error
}

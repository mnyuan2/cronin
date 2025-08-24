package data

import (
	"context"
	"cron/internal/basic/conv"
	"cron/internal/basic/db"
	"cron/internal/models"
	"errors"
	"fmt"
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
			Select("name", "spec", "protocol", "command", "remark", "update_dt", "type", "msg_set", "empty_not_msg", "var_fields",
				"after_tmpl", "command_hash", "msg_set_hash", "var_fields_hash", "after_sleep", "err_retry_num", "err_retry_sleep", "err_retry_mode",
				"audit_user_id", "audit_user_name", "tag_ids", "tag_names", "status", "status_remark", "status_dt", "source_ids").
			Updates(data).Error
	} else {
		data.CreateDt = time.Now().Format(conv.FORMAT_DATETIME)
		return m.db.Create(data).Error
	}
}

func (m *CronConfigData) ChangeStatus(data *models.CronConfig, remark string) error {
	data.UpdateDt = time.Now().Format(conv.FORMAT_DATETIME)
	data.StatusDt = data.UpdateDt
	data.StatusRemark = remark
	return m.db.Where("id=?", data.Id).
		Select("status", "status_remark", "status_dt", "update_dt", "entry_id", "handle_user_ids", "handle_user_names", "audit_user_id", "audit_user_name").
		Updates(data).Error
}

func (m *CronConfigData) SetEntryId(data *models.CronConfig) error {
	return m.db.Where("id=?", data.Id).Select("entry_id").Updates(data).Error
}

func (m *CronConfigData) GetOne(env string, Id int) (data *models.CronConfig, err error) {
	data = &models.CronConfig{}
	return data, m.db.Where("env=? and id=?", env, Id).Take(data).Error
}

// Del 删除
func (m *CronConfigData) Del(where *db.Where) (count int, err error) {
	if where.Len() == 0 {
		return 0, errors.New("未指定 config 删除条件")
	}
	count = 0
	w, args := where.Build()
	err = m.db.Model(&models.CronConfig{}).Where(w, args...).Select("count(*)").Find(&count).Error
	if err != nil {
		return 0, err
	}
	if count == 0 {
		return count, nil
	}

	err = m.db.Where(w, args...).Delete(&models.CronConfig{}).Error
	if err != nil {
		return 0, fmt.Errorf("删除失败，%w", err)
	}
	return count, nil
}

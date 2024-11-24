package data

import (
	"context"
	"cron/internal/basic/conv"
	"cron/internal/basic/db"
	"cron/internal/models"
	"time"
)

type CronPipelineData struct {
	db        *db.MyDB
	tableName string
}

func NewCronPipelineData(ctx context.Context) *CronPipelineData {
	return &CronPipelineData{
		db:        db.New(ctx),
		tableName: "cron_pipeline",
	}
}

func (m *CronPipelineData) ListPage(where *db.Where, page, size int, list interface{}) (total int64, err error) {
	str, args := where.Build()

	return m.db.Paginate(list, page, size, m.tableName, "*", "id desc", str, args...)
}

func (m *CronPipelineData) Set(data *models.CronPipeline) error {
	data.UpdateDt = time.Now().Format(conv.FORMAT_DATETIME)
	if data.Id > 0 {
		return m.db.Where("id=?", data.Id).
			Omit("entry_id", "env").
			Select("type", "name", "spec", "config_ids", "configs", "config_disable_action", "config_err_action", "remark", "update_dt", "msg_set",
				"var_params", "interval", "msg_set_hash", "audit_user_id", "audit_user_name").
			Updates(data).Error
	} else {
		data.CreateDt = time.Now().Format(conv.FORMAT_DATETIME)
		data.StatusDt = data.CreateDt
		data.StatusRemark = "新增"
		return m.db.Create(data).Error
	}
}

func (m *CronPipelineData) ChangeStatus(data *models.CronPipeline, remark string) error {
	data.UpdateDt = time.Now().Format(conv.FORMAT_DATETIME)
	data.StatusDt = data.UpdateDt
	data.StatusRemark = remark
	return m.db.Where("id=?", data.Id).
		Select("status", "status_remark", "status_dt", "update_dt", "entry_id", "handle_user_ids", "handle_user_names", "audit_user_id", "audit_user_name").
		Updates(data).Error
}

func (m *CronPipelineData) SetEntryId(data *models.CronPipeline) error {
	return m.db.Where("id=?", data.Id).Select("entry_id").Updates(data).Error
}

func (m *CronPipelineData) GetOne(env string, Id int) (data *models.CronPipeline, err error) {
	data = &models.CronPipeline{}
	return data, m.db.Where("env=? and id=?", env, Id).Take(data).Error
}

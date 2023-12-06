package data

import (
	"context"
	"cron/internal/basic/db"
	"cron/internal/models"
)

type CronSettingData struct {
	db        *db.Database
	tableName string
}

func NewCronSettingData(ctx context.Context) *CronSettingData {
	return &CronSettingData{
		db:        db.New(ctx),
		tableName: "cron_setting",
	}
}

// 列表查询
func (m *CronSettingData) GetList(scene string, page, size int, list interface{}) (total int64, err error) {
	str, args := db.NewWhere().Eq("scene", scene, db.RequiredOption()).Build()
	total, err = m.db.Read.Paginate(list, page, size, m.tableName, "*", "update_dt,id desc", str, args...)

	return total, err
}

// 获得单条配置
func (m *CronSettingData) GetOne(where *db.Where) (one *models.CronSetting, err error) {
	one = &models.CronSetting{}
	w, args := where.Build()

	return one, m.db.Write.Where(w, args...).Take(one).Error
}

// 设置
func (m *CronSettingData) Set(one *models.CronSetting) error {
	if one.Id > 0 {
		return m.db.Write.Where("id=?", one.Id).Omit("create_dt", "key", "env", "status").Updates(one).Error
	} else {
		return m.db.Write.Create(one).Error
	}
}

// 删除
func (m *CronSettingData) Del(key string, id int) error {
	one := &models.CronSetting{}
	return m.db.Write.Where("key=? and id=?", key, id).Delete(one).Error
}

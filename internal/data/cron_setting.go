package data

import (
	"context"
	"cron/internal/basic/db"
	"cron/internal/basic/enum"
	"cron/internal/models"
)

type CronSettingData struct {
	db        *db.MyDB
	tableName string
}

func NewCronSettingData(ctx context.Context) *CronSettingData {
	return &CronSettingData{
		db:        db.New(ctx),
		tableName: "cron_setting",
	}
}

// 列表查询
func (m *CronSettingData) GetList(scene string, env string, page, size int, list interface{}) (total int64, err error) {
	str, args := db.NewWhere().
		Eq("scene", scene, db.RequiredOption()).
		Eq("env", env, db.RequiredOption()).
		Build()
	total, err = m.db.Paginate(list, page, size, m.tableName, "*", "update_dt desc,id desc", str, args...)

	return total, err
}

// 获得单条配置
func (m *CronSettingData) GetOne(where *db.Where) (one *models.CronSetting, err error) {
	one = &models.CronSetting{}
	w, args := where.Build()

	return one, m.db.Where(w, args...).Take(one).Error
}

// 设置
func (m *CronSettingData) Set(one *models.CronSetting) error {
	if one.Id > 0 {
		return m.db.Where("id=?", one.Id).Omit("create_dt", "scene", "env", "status").Updates(one).Error
	} else {
		return m.db.Create(one).Error
	}
}

// 设置
func (m *CronSettingData) ChangeStatus(one *models.CronSetting) error {
	return m.db.Where("id=?", one.Id).Select("status", "update_dt").Updates(one).Error
}

// 删除
func (m *CronSettingData) Del(scene, env string, id int) error {
	one := &models.CronSetting{}
	return m.db.Where("scene=? and env=? and id=?", scene, env, id).Delete(one).Error
}

// 获得sql连接源
func (m *CronSettingData) GetSqlSourceOne(env string, id int) (one *models.CronSetting, err error) {
	w := db.NewWhere().
		Eq("scene", models.SceneSqlSource).
		Eq("env", env, db.RequiredOption()).
		Eq("id", id, db.RequiredOption()).
		Eq("status", enum.StatusActive)
	return m.GetOne(w)
}

// 获得env信息
func (m *CronSettingData) GetEnvOne(id int) (one *models.CronSetting, err error) {
	w := db.NewWhere().Eq("scene", models.SceneEnv).Eq("id", id, db.RequiredOption())
	return m.GetOne(w)
}

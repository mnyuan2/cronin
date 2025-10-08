package data

import (
	"context"
	"cron/internal/basic/conv"
	"cron/internal/basic/db"
	"cron/internal/basic/enum"
	"cron/internal/basic/errs"
	"cron/internal/models"
	jsoniter "github.com/json-iterator/go"
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

// 列表查询(这个应该叫getListPage)
func (m *CronSettingData) GetList(scene string, env string, page, size int, list interface{}) (total int64, err error) {
	str, args := db.NewWhere().
		Eq("scene", scene, db.RequiredOption()).
		Eq("env", env, db.RequiredOption()).
		Build()
	total, err = m.db.Paginate(list, page, size, m.tableName, "*", "update_dt desc,id desc", str, args...)

	return total, err
}

// 获得多条数据
func (m *CronSettingData) Gets(where *db.Where) (list []*models.CronSetting, err error) {
	list = []*models.CronSetting{}
	w, args := where.Build()

	return list, m.db.Where(w, args...).Find(&list).Error
}

// 获得单条配置
func (m *CronSettingData) GetOne(where *db.Where) (one *models.CronSetting, err error) {
	one = &models.CronSetting{}
	w, args := where.Build()

	return one, m.db.Where(w, args...).Find(one).Error
}

// 设置
func (m *CronSettingData) Set(one *models.CronSetting) error {
	if one.Id > 0 {
		return m.db.Where("id=?", one.Id).Omit("create_dt", "scene", "status").Updates(one).Error
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

// 获得连接源
func (m *CronSettingData) GetSourceOne(env string, id int) (one *models.CronSetting, err error) {
	w := db.NewWhere().
		FindInSet("env", env).
		Eq("id", id, db.RequiredOption()).
		Eq("status", enum.StatusActive)
	return m.GetOne(w)
}

// 获得资源列表
func (m *CronSettingData) GetSourceList(scene string) (list []*models.CronSetting, err error) {
	return list, m.db.Where("scene=?", scene).Find(&list).Error
}

// 获得环境列表
func (m *CronSettingData) GetEnvList() (list []*models.CronSetting, err error) {
	w := db.NewWhere().Eq("scene", models.SceneEnv)
	list = []*models.CronSetting{}
	where, args := w.Build()

	return list, m.db.Where(where, args...).Find(&list).Error
}

// 获得env信息
func (m *CronSettingData) GetEnvOne(id int) (one *models.CronSetting, err error) {
	w := db.NewWhere().Eq("scene", models.SceneEnv).Eq("id", id, db.RequiredOption())
	return m.GetOne(w)
}

// 获得env信息
func (m *CronSettingData) GetMessageOne(id int) (one *models.CronSetting, err error) {
	w := db.NewWhere().Eq("scene", models.SceneMsg).Eq("id", id, db.RequiredOption())
	return m.GetOne(w)
}

// 全局变量列表
func (m *CronSettingData) GetGlobalVariateList() (list []*models.CronSetting, err error) {
	list = []*models.CronSetting{}
	return list, m.db.Where("scene=?", models.SceneGlobalVar).Find(&list).Error
}

// 全局变量设置
func (m *CronSettingData) SetGlobalVariate(one *models.CronSetting) error {
	one.UpdateDt = conv.TimeNew().String()
	if one.Id > 0 {
		row := &models.CronSetting{}
		m.db.Where("scene=? and name=? and id!=?", models.SceneGlobalVar, one.Name, one.Id).Find(row)
		if row.Id > 0 {
			return errs.New(nil, "名称已存在")
		}
		return m.db.Where("id=?", one.Id).Omit("id", "scene", "env", "create_dt", "status").Updates(one).Error
	} else {
		one.Status = enum.StatusDisable
		one.CreateDt = one.UpdateDt
		one.Env = ""
		one.Scene = models.SceneGlobalVar
		row := &models.CronSetting{}
		m.db.Where("scene=? and name=?", models.SceneGlobalVar, one.Name).Find(row)
		if row.Id > 0 {
			return errs.New(nil, "名称已存在")
		}
		return m.db.Create(one).Error
	}
}

// 全局变量 状态设置
func (m *CronSettingData) ChangeGlobalVariateStatus(one *models.CronSetting) error {
	row := &models.CronSetting{}
	m.db.Where("scene=? and id=?", models.SceneGlobalVar, one.Id).Find(row)
	if row.Id == 0 {
		return errs.New(nil, "数据不存在")
	}
	if one.Status == enum.StatusDelete {
		if row.Status != enum.StatusDisable {
			return errs.New(nil, "数据激活中，删除失败")
		}
		return m.db.Where("id=?", row.Id).Delete(row).Error
	}
	row.UpdateDt = conv.TimeNew().String()
	row.Status = one.Status
	err := m.db.Where("id=? and status!=?", one.Id, one.Status).Select("status", "update_dt").Updates(row).Error
	if err != nil {
		return err
	}
	one.Name = row.Name
	one.Content = row.Content
	return err
}

// 模板列表
func (m *CronSettingData) GetTemplateList(w *db.Where) (list []*models.CronSetting, err error) {
	list = []*models.CronSetting{}
	if w == nil {
		w = db.NewWhere()
	}
	where, args := w.Eq("scene", models.Template).Build()
	return list, m.db.Where(where, args...).Find(&list).Error
}

// 模板设置
func (m *CronSettingData) SetTemplate(one *models.CronSetting) error {
	one.UpdateDt = conv.TimeNew().String()

	row := &models.CronSetting{}
	m.db.Where("scene=? and name=? and id=?", models.Template, one.Name, one.Id).Find(row)
	if row.Id == 0 {
		return errs.New(nil, "模板不存在")
	}

	if err := m.db.Where("id=?", one.Id).Select("content", "update_dt").Updates(one).Error; err != nil {
		return err
	}
	return nil
}

// 查询单个模板
func (m *CronSettingData) GetTemplateOne(name string) (data *models.TemplateConfig, err error) {
	w := db.NewWhere().Eq("scene", models.Template).Eq("name", name)
	one, err := m.GetOne(w)
	if err != nil || one.Id == 0 {
		return nil, errs.New(err, "模板不存在")
	}
	data = &models.TemplateConfig{}
	err = jsoniter.UnmarshalFromString(one.Content, data)
	return data, err
}

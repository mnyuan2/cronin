package data

import (
	"context"
	"cron/internal/basic/db"
	"cron/internal/basic/enum"
	"cron/internal/models"
	"errors"
)

type CronUserData struct {
	db        *db.MyDB
	tableName string
}

func NewCronUserData(ctx context.Context) *CronUserData {
	return &CronUserData{
		db:        db.New(ctx),
		tableName: "cron_user",
	}
}

// 查询列表数据
func (m *CronUserData) GetListPage(where *db.Where, page, size int, list interface{}) (total int64, err error) {
	str, args := where.Build()

	return m.db.Paginate(list, page, size, m.tableName, "*", "sort asc,id desc", str, args...)
}

// 查询列表数据
func (m *CronUserData) GetList(where *db.Where) (list []*models.CronUser, err error) {
	list = []*models.CronUser{}
	w, args := where.Build()

	return list, m.db.Where(w, args...).Find(&list).Error
}

func (m *CronUserData) Set(data *models.CronUser) error {
	if data.Id > 0 {
		return m.db.Where("id=?", data.Id).Omit("id", "create_dt", "account").Updates(data).Error
	} else {
		return m.db.Create(data).Error
	}
}

func (m *CronUserData) ChangePassword(data *models.CronUser) error {
	return m.db.Where("id=?", data.Id).Select("password", "update_dt").Updates(data).Error
}

func (m *CronUserData) ChangeStatus(data *models.CronUser) error {
	return m.db.Where("id=?", data.Id).Select("status", "update_dt").Updates(data).Error
}

func (m *CronUserData) ChangeAccount(data *models.CronUser) error {
	if data.Account != "" {
		old := &models.CronUser{}
		m.db.Where("id!=? and account=? and status != ?", data.Id, data.Account, enum.StatusDelete).Select("id").Find(old)
		if old.Id > 0 {
			return errors.New("账户已经存在")
		}
	}

	return m.db.Where("id=?", data.Id).Select("account", "update_dt").Updates(data).Error
}

func (m *CronUserData) GetOne(Id int) (data *models.CronUser, err error) {
	data = &models.CronUser{}
	return data, m.db.Where("id=?", Id).Take(data).Error
}

func (m *CronUserData) Login(account, password string) (data *models.CronUser, err error) {
	data = &models.CronUser{}
	return data, m.db.Where("account=? and password=?", account, password).Find(data).Error
}

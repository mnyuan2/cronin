package data

import (
	"context"
	"cron/internal/basic/db"
	"cron/internal/models"
)

type CronAuthRoleData struct {
	ctx       context.Context
	db        *db.MyDB
	tableName string
}

func NewCronAuthRoleData(ctx context.Context) *CronAuthRoleData {
	return &CronAuthRoleData{
		ctx:       ctx,
		db:        db.New(ctx),
		tableName: "cron_auth_role",
	}
}

func (m *CronAuthRoleData) Set(one *models.CronAuthRole) error {
	if one.Id > 0 {
		return m.db.Where("id=?", one.Id).Select("name", "remark").Updates(one).Error
	} else {
		return m.db.Create(one).Error
	}
}

func (m *CronAuthRoleData) SetAuthIds(one *models.CronAuthRole) error {
	return m.db.Where("id=?", one.Id).Select("auth_ids").Updates(one).Error
}

func (m *CronAuthRoleData) GetOne(Id int) (data *models.CronAuthRole, err error) {
	data = &models.CronAuthRole{}
	return data, m.db.Where("id=?", Id).Take(data).Error
}

func (m *CronAuthRoleData) GetList(where *db.Where) (list []*models.CronAuthRole, err error) {
	list = []*models.CronAuthRole{}
	w, args := where.Build()

	return list, m.db.Where(w, args...).Find(&list).Error
}

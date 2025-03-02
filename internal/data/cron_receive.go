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

type CronReceiveData struct {
	db        *db.MyDB
	tableName string
}

func NewCronReceiveData(ctx context.Context) *CronReceiveData {
	return &CronReceiveData{
		db:        db.New(ctx),
		tableName: "cron_receive",
	}
}

// 查询列表数据
func (m *CronReceiveData) ListPage(where *db.Where, page, size int, list interface{}) (total int64, err error) {
	str, args := where.Build()

	return m.db.Paginate(list, page, size, m.tableName, "*", "id desc", str, args...)
}

// 查询列表数据
func (m *CronReceiveData) GetList(where *db.Where) (list []*models.CronReceive, err error) {
	list = []*models.CronReceive{}
	w, args := where.Build()

	return list, m.db.Where(w, args...).Find(&list).Error
}

func (m *CronReceiveData) Set(data *models.CronReceive) error {
	data.UpdateDt = time.Now().Format(conv.FORMAT_DATETIME)
	if data.Id > 0 {
		return m.db.Where("id=?", data.Id).Omit("entry_id", "env").Updates(data).Error
	} else {
		data.CreateDt = time.Now().Format(conv.FORMAT_DATETIME)
		return m.db.Create(data).Error
	}
}

func (m *CronReceiveData) ChangeStatus(data *models.CronReceive, remark string) error {
	data.UpdateDt = time.Now().Format(conv.FORMAT_DATETIME)
	data.StatusDt = data.UpdateDt
	data.StatusRemark = remark
	return m.db.Where("id=?", data.Id).Select("status", "status_remark", "status_dt", "update_dt",
		"handle_user_ids", "handle_user_names", "audit_user_id", "audit_user_name").
		Updates(data).Error
}

func (m *CronReceiveData) GetOne(Id int) (data *models.CronReceive, err error) {
	data = &models.CronReceive{}
	return data, m.db.Where("id=?", Id).Find(data).Error
}

// Del 删除
func (m *CronReceiveData) Del(where *db.Where) (count int, err error) {
	if where.Len() == 0 {
		return 0, errors.New("未指定 receive 删除条件")
	}
	count = 0
	w, args := where.Build()
	err = m.db.Model(&models.CronReceive{}).Where(w, args...).Select("count(*)").Find(&count).Error
	if err != nil {
		return 0, err
	}
	if count == 0 {
		return count, nil
	}

	err = m.db.Where(w, args...).Delete(&models.CronReceive{}).Error
	if err != nil {
		return 0, fmt.Errorf("删除失败，%w", err)
	}
	return count, nil
}

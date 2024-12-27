package data

import (
	"context"
	"cron/internal/basic/db"
	"cron/internal/models"
	"fmt"
)

type CronLogSpanData struct {
	db        *db.MyDB
	tableName string
}

func NewCronLogSpanData(ctx context.Context) *CronLogSpanData {
	return &CronLogSpanData{
		db:        db.New(ctx),
		tableName: "cron_log_span",
	}
}

// ListPage 查询列表数据
func (m *CronLogSpanData) ListPage(where *db.Where, page, size int, list interface{}) (total int64, err error) {
	str, args := where.Build()

	return m.db.Paginate(list, page, size, m.tableName, "*", "timestamp desc", str, args...)
}

// List 获得列表数据
func (m *CronLogSpanData) List(where *db.Where, size int, field string) (list []*models.CronLogSpan, err error) {
	w, args := where.Build()
	list = []*models.CronLogSpan{}
	err = m.db.Where(w, args...).Limit(size).Select(field).Order("timestamp desc,span_id").Find(&list).Error
	return list, err
}

// Del 删除
func (m *CronLogSpanData) Del(where *db.Where) (count int, err error) {
	count = 0
	w, args := where.Build()
	err = m.db.Model(&models.CronLogSpan{}).Where(w, args...).Select("count(*)").Find(&count).Error
	if err != nil {
		return 0, err
	}
	if count == 0 {
		return count, nil
	}

	err = m.db.Where(w, args...).Delete(&models.CronLogSpan{}).Error
	if err != nil {
		return 0, fmt.Errorf("删除失败，%w", err)
	}
	m.db.Where(w, args...).Delete(&models.CronLogSpanIndex{})

	return count, nil
}

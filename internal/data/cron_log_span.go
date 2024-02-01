package data

import (
	"context"
	"cron/internal/basic/db"
	"cron/internal/models"
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

// 查询列表数据
func (m *CronLogSpanData) ListPage(where *db.Where, page, size int, list interface{}) (total int64, err error) {
	str, args := where.Build()

	return m.db.Paginate(list, page, size, m.tableName, "*", "timestamp desc", str, args...)
}

// 获得列表数据
func (m *CronLogSpanData) List(where *db.Where, size int) (list []*models.CronLogSpan, err error) {
	w, args := where.Build()
	list = []*models.CronLogSpan{}
	err = m.db.Where(w, args...).Limit(size).Order("timestamp asc").Find(&list).Error
	return list, err
}

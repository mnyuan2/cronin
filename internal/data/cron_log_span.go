package data

import (
	"context"
	"cron/internal/basic/db"
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

	return m.db.Paginate(list, page, size, m.tableName, "*", "id desc", str, args...)
}

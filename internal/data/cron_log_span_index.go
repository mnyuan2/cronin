package data

import (
	"context"
	"cron/internal/basic/db"
	"cron/internal/models"
)

type CronLogSpanIndexData struct {
	db        *db.MyDB
	tableName string
}

func NewCronLogSpanIndexData(ctx context.Context) *CronLogSpanIndexData {
	return &CronLogSpanIndexData{
		db:        db.New(ctx),
		tableName: "cron_log_span_index",
	}
}

// 合计指标
func (m *CronLogSpanIndexData) SumIndex(w *db.Where) []*models.CronLogSpanIndex {
	where, args := w.Build()
	list := []*models.CronLogSpanIndex{}
	if m.db.GetDriver() == db.DriverMysql {
		m.db.Model(&models.CronLogSpan{}).Where(where, args...).
			Select("FROM_UNIXTIME(LEFT(`timestamp`,10),'%Y-%m-%d %H:00:00') timestamp",
				"env",
				"ref_id",
				"operation",
				"sum(status=0) status_empty_num", "sum(status=1) status_error_num", "sum(status=2) status_success_num",
				"max(duration) duration_max", "round(avg(duration)) duration_avg").
			Group("FROM_UNIXTIME(LEFT(`timestamp`,10),'%Y-%m-%d %H'), env, operation, ref_id").Scan(&list)
	} else if m.db.GetDriver() == db.DriverSqlite {
		m.db.Model(&models.CronLogSpan{}).Where(where, args...).
			Select("strftime('%Y-%m-%d %H:00:00', leftstr(`timestamp`, 10), 'unixepoch') AS timestamp",
				"env",
				"ref_id",
				"operation",
				"SUM(CASE WHEN status = 0 THEN 1 ELSE 0 END) AS status_empty_num",
				"SUM(CASE WHEN status = 1 THEN 1 ELSE 0 END) AS status_error_num",
				"SUM(CASE WHEN status = 2 THEN 1 ELSE 0 END) AS status_success_num",
				"MAX(duration) AS duration_max",
				"ROUND(AVG(duration)) AS duration_avg").
			Group("strftime('%Y-%m-%d %H', leftstr(`timestamp`, 10), 'unixepoch'), env, operation, ref_id").Scan(&list)
	} else {

	}

	return list
}

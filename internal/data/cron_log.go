package data

import (
	"context"
	"cron/internal/basic/conv"
	"cron/internal/basic/db"
	"cron/internal/models"
	"fmt"
	"strings"
	"time"
)

type SumConfTop struct {
	ConfId      int `json:"conf_id"`
	TotalNumber int `json:"total_number"`
	ErrorNumber int `json:"error_number"`
}

type CronLogData struct {
	db        *db.MyDB
	tableName string
}

func NewCronLogData(ctx context.Context) *CronLogData {
	return &CronLogData{
		db:        db.New(ctx),
		tableName: "cron_log",
	}
}

// 添加数据
func (m *CronLogData) Add(data *models.CronLog) error {
	return m.db.Create(data).Error
}

// 查询列表数据
func (m *CronLogData) GetList(where *db.Where, page, size int, list interface{}) (total int64, err error) {
	str, args := where.Build()

	return m.db.Paginate(list, page, size, m.tableName, "*", "id desc", str, args...)
}

// 统计配置置顶的错误数
func (m *CronLogData) SumConfTopError(env string, confId []int, startTime, endTime time.Time, component string) (list map[int]*SumConfTop, err error) {
	w, args := db.NewWhere().
		Eq("a.env", env).
		In("a.ref_id", confId).
		Gte("a.timestamp", startTime.UnixMicro()).Lte("a.timestamp", endTime.UnixMicro()).
		Like("tags_kv", fmt.Sprintf("component=%v", component)).
		Build()
	sql := strings.Replace(`SELECT
	a.ref_id conf_id, count(*) total_number, sum(a.status=1) error_number
FROM
	cron_log_span as a
%WHERE
GROUP BY a.ref_id;`, "%WHERE", "WHERE "+w, 1)
	temps := []*SumConfTop{}
	list = map[int]*SumConfTop{}
	err = m.db.Raw(sql, args...).Find(&temps).Error
	if err != nil {
		return list, err
	}

	for _, temp := range temps {
		list[temp.ConfId] = temp
	}
	return list, nil
}

// 批量删除
func (m *CronLogData) DelBatch(end time.Time) (count int, err error) {
	count = 0
	endDate := end.Format(conv.FORMAT_DATETIME)
	err = m.db.Model(&models.CronLog{}).Where("create_dt <= ?", endDate).Select("count(*)").Find(&count).Error
	if err != nil {
		return 0, err
	}
	if count == 0 {
		return count, nil
	}

	err = m.db.Where("create_dt <= ?", endDate).Delete(&models.CronLog{}).Error
	if err != nil {
		return 0, fmt.Errorf("删除失败，%w", err)
	}
	return count, nil
}

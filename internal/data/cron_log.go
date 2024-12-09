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
func (m *CronLogData) SumConfTopError(w *db.Where) (list map[int]*SumConfTop, err error) {
	where, args := w.Build()
	sql := strings.Replace(`SELECT
	ref_id conf_id,
	sum(status_empty_num) status_empty_num,
	sum(status_error_num) status_error_num,
	sum(status_success_num) status_success_num
FROM cron_log_span_index 
%WHERE
GROUP BY ref_id;`, "%WHERE", "WHERE "+where, 1)
	temps := []*models.CronLogSpanIndex{}
	list = map[int]*SumConfTop{}
	err = m.db.Raw(sql, args...).Find(&temps).Error
	if err != nil {
		return list, err
	}

	for _, temp := range temps {
		id, _ := conv.Ints().Parse(temp.RefId)
		list[id] = &SumConfTop{
			ConfId:      id,
			TotalNumber: temp.StatusEmptyNum + temp.StatusErrorNum + temp.StatusSuccessNum,
			ErrorNumber: temp.StatusErrorNum,
		}
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

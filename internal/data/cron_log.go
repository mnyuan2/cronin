package data

import (
	"context"
	"cron/internal/basic/conv"
	"cron/internal/basic/db"
	"cron/internal/models"
	"time"
)

type SumConfTop struct {
	ConfId      int `json:"conf_id"`
	TotalNumber int `json:"total_number"`
	ErrorNumber int `json:"error_number"`
}

type CronLogData struct {
	db        *db.Database
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
	return m.db.Write.Create(data).Error
}

// 查询列表数据
func (m *CronLogData) GetList(where *db.Where, page, size int, list interface{}) (total int64, err error) {
	str, args := where.Build()

	return m.db.Read.Paginate(list, page, size, m.tableName, "*", "id desc", str, args...)
}

// 统计配置置顶的错误数
func (m *CronLogData) SumConfTopError(confId []int, startTime, endTime time.Time, maxNumber int) (list map[int]*SumConfTop, err error) {
	sql := `SELECT 
	t1.conf_id, count(t1.id) total_number, sum(t1.status=?) error_number
FROM 
	cron_log t1 
WHERE 
	?>(SELECT count(*) FROM cron_log WHERE t1.conf_id=conf_id and t1.id<id) and t1.conf_id in(?) and create_dt between ? and ? GROUP BY t1.conf_id`

	temps := []*SumConfTop{}
	list = map[int]*SumConfTop{}
	err = m.db.Read.Raw(sql, models.StatusDisable, maxNumber, confId, startTime.Format(conv.FORMAT_DATETIME), endTime.Format(conv.FORMAT_DATETIME)).Take(&temps).Error
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

	err = m.db.Write.Model(&models.CronLog{}).Where("create_dt <= ?", end.Format(conv.FORMAT_DATETIME)).Take(&count).Error
	if err != nil {
		return 0, err
	}

	err = m.db.Write.Where("create_dt <= ?", end.Format(conv.FORMAT_DATETIME)).Delete(&models.CronLog{}).Error
	return 0, err
}

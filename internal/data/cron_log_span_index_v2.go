package data

import (
	"context"
	"cron/internal/basic/db"
	"cron/internal/models"
	"fmt"
)

type TraceIdGroupItem struct {
	RefId     string `json:"ref_id"`
	RefName   string `json:"ref_name"`
	Operation string `json:"operation"`
	TraceId   string `json:"trace_id"`
	TotalNum  int    `json:"total_num"`
	ErrorNum  int    `json:"error_num"`
}

type CronLogSpanIndexV2Data struct {
	db        *db.MyDB
	tableName string
}

func NewCronLogSpanIndexV2Data(ctx context.Context) *CronLogSpanIndexV2Data {
	return &CronLogSpanIndexV2Data{
		db:        db.New(ctx),
		tableName: "cron_log_span_index",
	}
}

// Del 删除
func (m *CronLogSpanIndexV2Data) Del(where *db.Where) (count int, err error) {
	count = 0
	w, args := where.Build()
	err = m.db.Model(&models.CronLogSpanIndexV2{}).Where(w, args...).Select("count(*)").Find(&count).Error
	if err != nil {
		return 0, err
	}
	if count == 0 {
		return count, nil
	}

	err = m.db.Where(w, args...).Delete(&models.CronLogSpanIndexV2{}).Error
	if err != nil {
		return 0, fmt.Errorf("删除失败，%w", err)
	}

	return count, nil
}

func (m *CronLogSpanIndexV2Data) GetTraceIds(where *db.Where, limit int) (ids []string, err error) {
	list := []string{}
	w, args := where.Build()
	err = m.db.Model(&models.CronLogSpanIndexV2{}).Where(w, args...).Limit(limit).Select("trace_id").Order("timestamp desc").Scan(&list).Error
	if err != nil {
		return nil, err
	}
	return list, nil
}

// 统计配置置顶的错误数
func (m *CronLogSpanIndexV2Data) SumStatus(w *db.Where) (list map[int]*SumStatus, err error) {
	where, args := w.Build()
	temps := []*SumStatus{}
	if m.db.GetDriver() == db.DriverMysql {
		err = m.db.Model(&models.CronLogSpanIndexV2{}).Where(where, args...).
			Select(
				"ref_id",
				"count(1) total_number",
				"sum(status=1) error_number").
			Group("ref_id").Scan(&temps).Error
	} else if m.db.GetDriver() == db.DriverSqlite {
		err = m.db.Model(&models.CronLogSpanIndexV2{}).Where(where, args...).
			Select(
				"ref_id",
				"count(1) AS total_number",
				"SUM(CASE WHEN status = 1 THEN 1 ELSE 0 END) AS error_number").
			Group("ref_id").Scan(&temps).Error
	} else {
		err = fmt.Errorf("未支持的数据库类型")
	}
	if err != nil {
		return list, err
	}

	list = map[int]*SumStatus{}
	for _, temp := range temps {
		list[temp.RefId] = temp
	}
	return list, nil
}

func (m *CronLogSpanIndexV2Data) TraceIdGroup(where *db.Where) (list map[string][]*TraceIdGroupItem, err error) {
	w, args := where.Build()
	rows := []*TraceIdGroupItem{}
	err = m.db.Model(&models.CronLogSpanIndexV2{}).Where(w, args...).
		Select(
			"trace_id",
			"ref_id",
			"ref_name",
			"operation",
			"count(1) total_num",
			"sum(status=1) error_num",
		).Group("trace_id, operation").Order("timestamp,span_id asc").Scan(&rows).Error
	if err != nil {
		return nil, err
	}
	list = map[string][]*TraceIdGroupItem{}
	for _, row := range rows {
		list[row.TraceId] = append(list[row.TraceId], row)
	}
	return list, nil
}

func (m *CronLogSpanIndexV2Data) List(where *db.Where, page, limit int) (total int64, list []*models.CronLogSpanIndexV2, err error) {
	w, args := where.Build()
	err = m.db.Model(&models.CronLogSpanIndexV2{}).Where(w, args...).Select("COUNT(DISTINCT trace_id)").Scan(&total).Error
	if err != nil {
		return 0, nil, err
	}
	if total == 0 {
		return total, []*models.CronLogSpanIndexV2{}, nil
	}
	err = m.db.Model(&models.CronLogSpanIndexV2{}).
		Where(w, args...).Limit(limit).Offset((page - 1) * limit).
		Group("trace_id").
		Order("timestamp desc,span_id asc").
		Find(&list).Error
	if err != nil {
		return 0, nil, err
	}
	return total, list, nil
}

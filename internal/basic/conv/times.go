package conv

import (
	"time"
)

// mysql写入日期格式
const (
	FORMAT_DATETIME = "2006-01-02 15:04:05"
	FORMAT_DATE     = "2006-01-02"
	NULL_DATETIME   = "0000-00-00 00:00:00" // 空时间格式
)

type Time struct {
	time.Time
}

func TimeNew(ti ...time.Time) *Time {
	t := &Time{}
	switch len(ti) {
	case 0:
		t.Time = time.Now()
	case 1:
		t.Time = ti[0]
	default:
		panic("too many parameters")
	}
	return t
}

// 通过周期载入时间
func NewTimeFormCycle(str string) (t *Time, err error) {
	t = &Time{}
	t.Time, err = time.Parse("200601", str)
	return t, err
}

// 时间运算
func (t *Time) AddDate(years int, months int, days int) *Time {
	return &Time{t.Time.AddDate(years, months, days)}
}

// 截止日期
func (t *Time) EndDateTime() string {
	return t.Format(FORMAT_DATE) + " 23:59:59"
}

// 开始日期
func (t *Time) StartDateTime() string {
	return t.Format(FORMAT_DATE) + " 00:00:00"
}

// 获得周期
// 默认当前月，如果有指定月，以指定月为准。
func (t *Time) Cycle() string {
	return t.Time.Format("200601")
}

func (t *Time) String() string {
	return t.Format(FORMAT_DATETIME)
}

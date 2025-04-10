package biz

import (
	"context"
	"cron/internal/basic/db"
	"cron/internal/data"
	"fmt"
	"testing"
	"time"
)

func TestCronLogTop(t *testing.T) {
	endTime := time.Now()
	startTime := time.Now().Add(-time.Hour * 24 * 7) // 取七天前
	w2 := db.NewWhere().
		Eq("env", "public").
		Eq("operation", "job-task").
		In("ref_id", []int{1, 4}).
		Between("timestamp", startTime.UnixMicro(), endTime.UnixMicro())
	list, err := data.NewCronLogSpanIndexV2Data(context.Background()).SumStatus(w2)
	fmt.Println(list, err)
}

func TestCronLogIndex(t *testing.T) {
	endTime := time.Now()
	startTime := time.Now().Add(-(time.Hour * 3))
	w2 := db.NewWhere().Between("timestamp", startTime.UnixMicro(), endTime.UnixMicro())
	list := data.NewCronLogSpanIndexData(context.Background()).SumIndex(w2)
	for i, item := range list {
		fmt.Println(i, item)
	}
}

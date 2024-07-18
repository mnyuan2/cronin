package biz

import (
	"context"
	"cron/internal/data"
	"fmt"
	"testing"
	"time"
)

func TestCronLogTop(t *testing.T) {
	endTime := time.Now()
	startTime := time.Now().Add(-time.Hour * 24 * 7) // 取七天前
	list, err := data.NewCronLogData(context.Background()).SumConfTopError("public", []int{1, 4}, startTime, endTime, "config")
	fmt.Println(list, err)
}

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
	list, err := data.NewCronLogData(context.Background()).SumConfTopError([]int{1, 4}, startTime, endTime, 3)
	fmt.Println(list, err)
}

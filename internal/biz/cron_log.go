package biz

import (
	"context"
	"cron/internal/basic/db"
	"cron/internal/data"
	"cron/internal/models"
	"cron/internal/pb"
	"fmt"
	"time"
)

type CronLogService struct {
}

func NewCronLogService() *CronLogService {
	return &CronLogService{}
}

// 通过配置查询日志
func (dm *CronLogService) ByConfig(ctx context.Context, r *pb.CronLogByConfigRequest) (resp *pb.CronLogByConfigResponse, err error) {
	w := db.NewWhere().Eq("conf_id", r.ConfId, db.RequiredOption())
	resp = &pb.CronLogByConfigResponse{List: []*pb.CronLogItem{}}

	_, err = data.NewCronLogData(ctx).GetList(w, 1, r.Limit, &resp.List)
	for _, item := range resp.List {
		item.StatusName = models.LogStatusMap[item.Status]
	}

	return resp, err
}

// 通过配置查询日志
func (dm *CronLogService) Del(ctx context.Context, r *pb.CronLogDelRequest) (resp *pb.CronLogDelResponse, err error) {
	if r.Retention == "" {
		return nil, fmt.Errorf("retention 参数为必须")
	}

	re, err := time.ParseDuration(r.Retention)
	if err != nil {
		return nil, fmt.Errorf("retention 参数有误, %s", err.Error())
	} else if re.Hours() < 24 {
		return nil, fmt.Errorf("retention 参数不得小于24h")
	}
	end := time.Now().Add(-re)
	resp = &pb.CronLogDelResponse{}
	resp.Count, err = data.NewCronLogData(ctx).DelBatch(end)

	return resp, err
}

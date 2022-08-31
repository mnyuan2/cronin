package biz

import (
	"context"
	"cron/internal/basic/db"
	"cron/internal/data"
	"cron/internal/models"
	"cron/internal/pb"
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

package biz

import (
	"context"
	"cron/internal/basic/db"
	"cron/internal/data"
	"cron/internal/models"
	"cron/internal/pb"
)

type CronConfigService struct {
}

func NewCronConfigService() *CronConfigService {
	return &CronConfigService{}
}

// 任务配置列表
func (dm *CronConfigService) List(ctx context.Context, r *pb.CronConfigListRequest) (resp *pb.CronConfigListReply, err error) {
	w := db.NewWhere()
	// 构建查询条件

	resp = &pb.CronConfigListReply{
		List: []*pb.CronConfigListItem{},
		Page: &pb.Page{
			Page: 1,
			Size: 10,
		},
	}
	resp.Page.Total, err = data.NewCronConfigData(ctx).GetList(w, 1, 10, &resp.List)

	return resp, err
}

func (dm *CronConfigService) Get() {

}

// 任务配置
func (dm *CronConfigService) Set(ctx context.Context, r *pb.CronConfigSetRequest) (resp *pb.CronConfigSetResponse, err error) {
	d := &models.CronConfig{
		Id:       r.Id,
		Name:     r.Name,
		Spec:     r.Spec,
		Protocol: models.CronProtocol(r.Protocol),
		Command:  r.Command,
		Status:   models.StatusDisable,
		Remark:   r.Remark,
	}

	err = data.NewCronConfigData(ctx).Set(d)
	if err != nil {
		return nil, err
	}
	return &pb.CronConfigSetResponse{
		Id: d.Id,
	}, err
}

func (dm *CronConfigService) Del() {

}

package biz

import (
	"context"
	"cron/internal/basic/db"
	"cron/internal/data"
	"cron/internal/models"
	"cron/internal/pb"
	"fmt"
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
	resp.Page.Total, err = data.NewCronConfigData(ctx).GetList(w, 1, 20, &resp.List)
	for _, item := range resp.List {
		item.StatusName = models.ConfStatusMap[item.Status]
		item.ProtocolName = models.ProtocolMap[item.Protocol]
	}

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
		Protocol: r.Protocol,
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

// 编辑任务
func (dm *CronConfigService) Edit(ctx context.Context, r *pb.CronConfigSetRequest) (resp *pb.CronConfigSetResponse, err error) {
	if r.Id <= 0 {
		return nil, fmt.Errorf("未指定配置")
	}

	da := data.NewCronConfigData(ctx)
	conf, err := da.GetOne(r.Id)
	if err != nil {
		return nil, err
	}

	if _, ok := models.ConfStatusMap[r.Status]; ok {
		//if conf.Status == models.StatusActive {
		//
		//}
		conf.Status = r.Status
		// 启用 到 停用 要关闭执行中的对应任务；
		// 停用 到 启用 要把任务注册；
		// 新旧状态一致，就不用附加操作了。
	}

	if err = da.Set(conf); err != nil {
		// 前面操作了任务，这里失败了；要将任务进行反向操作（回滚）（并附带两条对应日志）
		return nil, err
	}

	return &pb.CronConfigSetResponse{
		Id: conf.Id,
	}, err
}

func (dm *CronConfigService) Del() {

}

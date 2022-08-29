package biz

import (
	"context"
	"cron/internal/basic/db"
	"cron/internal/data"
	"cron/internal/models"
	"cron/internal/pb"
	"fmt"
	jsoniter "github.com/json-iterator/go"
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
		item.Command = &pb.CronConfigCommand{}
		jsoniter.UnmarshalFromString(item.CommandStr, item.Command)
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
		Status:   models.StatusDisable,
		Remark:   r.Remark,
	}
	d.Command, _ = jsoniter.MarshalToString(r.Command)
	if _, err = secondParser.Parse(d.Spec); err != nil {
		return nil, fmt.Errorf("时间格式不规范，%s", err.Error())
	}
	if r.Protocol == models.ProtocolHttp {
		// 这里要校验一下协议的规范性；
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
	if conf.Status == models.StatusActive {
		return nil, fmt.Errorf("请先停用任务后编辑")
	}
	if r.Name != "" {
		conf.Name = r.Name
	}
	if r.Spec != "" {
		conf.Spec = r.Spec
	}
	if r.Remark != "" {
		conf.Remark = r.Remark
	}
	if _, err = secondParser.Parse(conf.Spec); err != nil {
		return nil, fmt.Errorf("时间格式不规范，%s", err.Error())
	}

	if err = da.Set(conf); err != nil {
		// 前面操作了任务，这里失败了；要将任务进行反向操作（回滚）（并附带两条对应日志）
		return nil, err
	}

	return &pb.CronConfigSetResponse{
		Id: conf.Id,
	}, err
}

// 任务状态变更
func (dm *CronConfigService) ChangeStatus(ctx context.Context, r *pb.CronConfigSetRequest) (resp *pb.CronConfigSetResponse, err error) {
	// 同一个任务，这里要加请求锁
	da := data.NewCronConfigData(ctx)
	conf, err := da.GetOne(r.Id)
	if err != nil {
		return nil, err
	}
	if conf.Status == r.Status {
		return nil, fmt.Errorf("状态相等")
	}
	if _, ok := models.ConfStatusMap[r.Status]; !ok {
		return nil, fmt.Errorf("错误状态请求")
	}

	if conf.Status == models.StatusActive && r.Status == models.StatusDisable { // 启用 到 停用 要关闭执行中的对应任务；
		NewTaskService().Del(conf)
	} else if conf.Status == models.StatusDisable && r.Status == models.StatusActive { // 停用 到 启用 要把任务注册；
		NewTaskService().Add(conf)
	}

	conf.Status = r.Status
	if err = da.Set(conf); err != nil {
		// 前面操作了任务，这里失败了；要将任务进行反向操作（回滚）（并附带两条对应日志）
		return nil, err
	}
	return &pb.CronConfigSetResponse{
		Id: conf.Id,
	}, nil
}

func (dm *CronConfigService) Del() {

}

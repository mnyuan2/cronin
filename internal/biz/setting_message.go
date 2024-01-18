package biz

import (
	"context"
	"cron/internal/basic/conv"
	"cron/internal/basic/db"
	"cron/internal/basic/enum"
	"cron/internal/biz/dtos"
	"cron/internal/data"
	"cron/internal/models"
	"cron/internal/pb"
	jsoniter "github.com/json-iterator/go"
)

type SettingMessageService struct {
	db  *db.MyDB
	ctx context.Context
	//user *auth.UserToken
}

func NewSettingMessageService(ctx context.Context) *SettingMessageService {
	return &SettingMessageService{
		ctx: ctx,
	}
}

// 任务配置列表
func (dm *SettingMessageService) List(r *pb.SettingMessageListRequest) (resp *pb.SettingMessageListReply, err error) {
	// 构建查询条件
	if r.Page <= 1 {
		r.Page = 1
	}
	if r.Size <= 10 {
		r.Size = 20
	}
	resp = &pb.SettingMessageListReply{
		Page: &pb.Page{
			Page: r.Page,
			Size: r.Size,
		},
	}
	list := []*models.CronSetting{}
	resp.Page.Total, err = data.NewCronSettingData(dm.ctx).GetList(models.SceneMsg, "", r.Page, r.Size, &list)
	if err != nil {
		return nil, err
	}

	// 格式化
	resp.List = make([]*pb.SettingMessageListItem, len(list))
	for i, item := range list {
		data := &pb.SettingMessageListItem{
			Id:    item.Id,
			Title: item.Title,
			//Sort:  item.Status,
			Template: &pb.SettingMessageTemplate{
				Http: &pb.CronHttp{},
			},
			UpdateDt: item.UpdateDt,
			CreateDt: item.CreateDt,
		}
		jsoniter.UnmarshalFromString(item.Content, data.Template)
		resp.List[i] = data
	}

	return resp, err
}

// 设置源
func (dm *SettingMessageService) Set(r *pb.SettingMessageSetRequest) (resp *pb.SettingMessageSetReply, err error) {
	if err = dtos.CheckHttp(r.Template.Http); err != nil {
		return nil, err
	}

	one := &models.CronSetting{}
	_data := data.NewCronSettingData(dm.ctx)
	ti := conv.TimeNew()
	oldSource := &pb.SettingSqlSource{}
	// 分为新增和编辑
	if r.Id > 0 {
		one, err = _data.GetMessageOne(r.Id)
		if err != nil {
			return nil, err
		}
		jsoniter.UnmarshalFromString(one.Content, oldSource)
	} else {
		one.Scene = models.SceneMsg
		one.Status = enum.StatusActive
		one.CreateDt = ti.String()
	}

	one.UpdateDt = ti.String()
	one.Title = r.Title
	one.Content, err = jsoniter.MarshalToString(r.Template)
	if err != nil {
		return nil, err
	}
	// 执行写入
	err = _data.Set(one)
	if err != nil {
		return nil, err
	}
	return &pb.SettingMessageSetReply{
		Id: one.Id,
	}, err
}

// 运行一下 模板消息
func (dm *SettingMessageService) Run(r *pb.SettingMessageSetRequest) (resp *pb.SettingMessageRunReply, err error) {
	if err = dtos.CheckHttp(r.Template.Http); err != nil {
		return nil, err
	}

	res, err := NewCronConfigService(dm.ctx, nil).
		Run(&pb.CronConfigRunRequest{
			Protocol: models.ProtocolHttp,
			Command:  &pb.CronConfigCommand{Http: r.Template.Http},
		})
	if err != nil {
		return nil, err
	}

	return &pb.SettingMessageRunReply{Result: res.Result}, nil
}

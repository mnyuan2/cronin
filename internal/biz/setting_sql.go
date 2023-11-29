package biz

import (
	"context"
	"cron/internal/basic/conv"
	"cron/internal/basic/db"
	"cron/internal/data"
	"cron/internal/models"
	"cron/internal/pb"
	"errors"
	jsoniter "github.com/json-iterator/go"
)

type SettingSqlService struct {
}

func NewSettingSqlService() *SettingSqlService {
	return &SettingSqlService{}
}

// 任务配置列表
func (dm *SettingSqlService) List(ctx context.Context, r *pb.SettingSqlListRequest) (resp *pb.SettingSqlListReply, err error) {
	// 构建查询条件
	if r.Page <= 1 {
		r.Page = 1
	}
	if r.Size <= 10 {
		r.Size = 20
	}
	resp = &pb.SettingSqlListReply{
		List: []*pb.SettingSqlListItem{},
		Page: &pb.Page{
			Page: r.Page,
			Size: r.Size,
		},
	}
	list := []*models.CronSetting{}
	resp.Page.Total, err = data.NewCronSettingData(ctx).GetList(models.KeySqlSource, r.Page, r.Size, &list)
	if err != nil {
		return nil, err
	}

	// 格式化
	for _, item := range list {
		data := &pb.SettingSqlListItem{
			Id:       item.Id,
			Title:    item.Title,
			CreateDt: item.CreateDt,
			UpdateDt: item.UpdateDt,
			Source:   &pb.SettingSqlSource{},
		}
		jsoniter.UnmarshalFromString(item.Content, data.Source)
	}

	return resp, err
}

// 设置源
func (dm *SettingSqlService) Set(ctx context.Context, r *pb.SettingSqlSetRequest) (resp *pb.SettingSqlSetReply, err error) {
	one := &models.CronSetting{}
	_data := data.NewCronSettingData(ctx)
	ti := conv.TimeNew()
	// 分为新增和编辑
	if r.Id > 0 {
		w := db.NewWhere().Eq("key", models.KeySqlSource).Eq("id", r.Id).Eq("status", models.StatusActive)
		one, err = _data.GetOne(w)
		if err != nil {
			return nil, err
		}
	} else {
		one.Key = models.KeySqlSource
		one.Status = models.StatusActive
		one.CreateDt = ti.String()
	}

	one.UpdateDt = ti.String()
	one.Title = r.Title
	one.Content, err = jsoniter.MarshalToString(r.Source)
	if err != nil {
		return nil, err
	}
	// 执行写入
	err = _data.Set(one)
	if err != nil {
		return nil, err
	}
	return &pb.SettingSqlSetReply{
		Id: one.Id,
	}, err
}

// 删除源
func (dm *SettingSqlService) ChangeStatus(ctx context.Context, r *pb.SettingChangeStatusRequest) (resp *pb.SettingChangeStatusReply, err error) {
	// 同一个任务，这里要加请求锁
	_data := data.NewCronSettingData(ctx)
	w := db.NewWhere().Eq("key", models.KeySqlSource).Eq("id", r.Id, db.RequiredOption())
	one, err := _data.GetOne(w)
	if err != nil {
		return nil, err
	}
	if one.Id <= 0 {
		return nil, errors.New("操作数据不存在")
	}
	// 目前仅支持删除
	if r.Status != models.StatusDelete {
		return nil, errors.New("不支持的状态操作")
	}

	err = _data.Del(one.Key, one.Id)
	return &pb.SettingChangeStatusReply{}, err
}

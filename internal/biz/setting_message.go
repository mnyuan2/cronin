package biz

import (
	"bytes"
	"context"
	"cron/internal/basic/conv"
	"cron/internal/basic/db"
	"cron/internal/basic/enum"
	"cron/internal/basic/errs"
	"cron/internal/biz/dtos"
	"cron/internal/data"
	"cron/internal/models"
	"cron/internal/pb"
	"encoding/json"
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
				Http: &pb.CronHttp{
					Header: []*pb.KvItem{},
				},
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
	// 方案1固定测试值、方案2随机测试值
	args := map[string]string{
		"env":                  "测试环境",
		"config.name":          "xx任务",
		"config.protocol_name": "sql脚本",
		"log.status_name":      "成功",
		"log.status_desc":      "success",
		"log.body":             "xxxxxxxxxxxxxx\nyyyyyyyyyyyyyy",
		"log.duration":         "3.2s",
		"log.create_dt":        "2023-01-01 11:12:59",
		"user.username":        "管理员,大王",
		"user.mobile":          "13118265689,12345678910",
	}
	b, _ := json.Marshal(r.Template)
	for k, v := range args {
		b = bytes.Replace(b, []byte("[["+k+"]]"), []byte(v), -1)
	}
	temp := &pb.SettingMessageTemplate{Http: &pb.CronHttp{}}
	if err = jsoniter.Unmarshal(b, temp); err != nil {
		return nil, errs.New(err, "解析错误")
	}

	res, err := NewCronConfigService(dm.ctx, nil).
		Run(&pb.CronConfigRunRequest{
			Protocol: models.ProtocolHttp,
			Command:  &pb.CronConfigCommand{Http: temp.Http},
		})
	if err != nil {
		return nil, err
	}

	return &pb.SettingMessageRunReply{Result: res.Result}, nil
}

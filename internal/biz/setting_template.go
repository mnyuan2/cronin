package biz

import (
	"context"
	"cron/internal/basic/db"
	"cron/internal/data"
	"cron/internal/models"
	"cron/internal/pb"
	"errors"
	jsoniter "github.com/json-iterator/go"
)

type Template struct {
	Default int `json:"default"`
}

// 环境设置
type TemplateService struct {
	ctx context.Context
	db  *db.MyDB
}

func NewTemplateService(ctx context.Context) *TemplateService {
	return &TemplateService{
		ctx: ctx,
	}
}

// 任务配置列表
func (dm *TemplateService) List(r *pb.TemplateListRequest) (resp *pb.TemplateListReply, err error) {

	list, err := data.NewCronSettingData(dm.ctx).GetTemplateList(db.NewWhere().Eq("name", r.Name))
	if err != nil {
		return nil, err
	}

	// 格式化
	resp = &pb.TemplateListReply{
		List: make([]*pb.TemplateListItem, len(list)),
	}
	for i, row := range list {
		conf := &models.TemplateConfig{}
		jsoniter.UnmarshalFromString(row.Content, conf)

		item := &pb.TemplateListItem{
			Id:       row.Id,
			Name:     row.Name,
			Title:    row.Title,
			Temp:     conf.Temp,
			Hint:     conf.Hint,
			UpdateDt: row.UpdateDt,
		}
		resp.List[i] = item
	}

	return resp, err
}

// 设置环境
func (dm *TemplateService) Set(r *pb.TemplateSetRequest) (resp *pb.TemplateSetReply, err error) {
	// 校验
	if r.Id <= 0 || r.Name == "" {
		return nil, errors.New("非法请求")
	}
	if len(r.Hint) > 500 {
		return nil, errors.New("提示信息长度不能超过500个字符")
	}

	// 格式化
	conf := &models.TemplateConfig{
		Temp: r.Temp,
		Hint: r.Hint,
	}
	str, err := jsoniter.MarshalToString(conf)
	if err != nil {
		return nil, err
	}
	one := &models.CronSetting{
		Id:      r.Id,
		Name:    r.Name,
		Content: str,
	}
	err = data.NewCronSettingData(dm.ctx).SetTemplate(one)
	if err != nil {
		return nil, err
	}

	return &pb.TemplateSetReply{}, nil
}

package biz

import (
	"context"
	"cron/internal/basic/conv"
	"cron/internal/basic/db"
	"cron/internal/basic/enum"
	"cron/internal/basic/errs"
	"cron/internal/data"
	"cron/internal/models"
	"cron/internal/pb"
	jsoniter "github.com/json-iterator/go"
	"log"
	"strings"
)

// 环境设置
type SettingService struct {
	ctx context.Context
	db  *db.MyDB
}

func NewSettingService(ctx context.Context) *SettingService {
	return &SettingService{
		ctx: ctx,
	}
}

// 偏好设置
func (dm *SettingService) PreferenceSet(r *pb.SettingPreferenceSetRequest) (resp *pb.SettingPreferenceSetReply, err error) {
	if r.Pipeline == nil {
		return nil, errs.New(nil, "流水线 配置异常")
	}
	if r.Git == nil || r.Git.OwnerRepo == nil {
		return nil, errs.New(nil, "git选项 配置异常")
	}

	_data := data.NewCronSettingData(dm.ctx)
	ti := conv.TimeNew()

	one, err := _data.GetOne(db.NewWhere().Eq("scene", models.ScenePreference))
	if err != nil {
		return nil, err
	}
	one.Content, _ = jsoniter.MarshalToString(r)

	// 分为新增和编辑
	if one.Id <= 0 {
		one.Scene = models.ScenePreference
		one.Status = enum.StatusActive
		one.CreateDt = ti.String()
		one.Title = "偏好设置"
	}
	one.UpdateDt = ti.String()

	// 执行写入
	err = _data.Set(one)
	if err != nil {
		return nil, err
	}
	return &pb.SettingPreferenceSetReply{}, err
}

// 偏好查看
func (dm *SettingService) PreferenceGet(r *pb.SettingPreferenceGetRequest) (resp *pb.SettingPreferenceGetReply, err error) {

	resp = &pb.SettingPreferenceGetReply{
		Pipeline: &pb.SettingPreferencePipeline{},
		Git:      &pb.SettingPreferenceGit{OwnerRepo: []*pb.SettingPreferenceGitOwner{}},
		Other: &pb.SettingPreferenceOther{
			ConfigSelectType: models.TypeCycle,
		},
	}

	da := data.NewCronSettingData(dm.ctx)
	one, err := da.GetOne(db.NewWhere().Eq("scene", models.ScenePreference))
	if err != nil {
		return nil, err
	}
	if one.Content != "" {
		if err := jsoniter.UnmarshalFromString(one.Content, &resp); err != nil {
			log.Println("偏好查询解析错误，", err.Error(), " --> ", one.Content)
		}
	} else { // 组织默认值
		resp.Pipeline.ConfigDisableAction = 1
	}

	return resp, err
}

// 全局变量 列表
func (dm *SettingService) GlobalVariateList(r *pb.GlobalVariateListRequest) (resp *pb.GlobalVariateListReply, err error) {
	list, err := data.NewCronSettingData(dm.ctx).GetGlobalVariateList()
	if err != nil {
		return nil, err
	}
	resp = &pb.GlobalVariateListReply{
		List: make([]*pb.GlobalVariateListItem, len(list)),
	}
	for i, v := range list {
		resp.List[i] = &pb.GlobalVariateListItem{
			Id:         v.Id,
			Name:       v.Name,
			Value:      v.Content,
			Remark:     v.Title,
			Status:     v.Status,
			StatusName: enum.StatusMap[v.Status],
			Register:   enum.BoolNot,
		}
		if val := globalVariateList.Get(v.Name); val != nil {
			resp.List[i].Register = enum.BoolYes
		}
	}
	return resp, err
}

// 全局变量设置
func (dm *SettingService) GlobalVariateSet(r *pb.GlobalVariateSetRequest) (resp *pb.GlobalVariateSetReply, err error) {
	one := &models.CronSetting{
		Id:      r.Id,
		Name:    strings.TrimSpace(r.Name),
		Title:   r.Remark,
		Content: r.Value,
	}

	if r.Name == "" {
		return nil, errs.New(nil, "名称不能为空")
	}
	err = data.NewCronSettingData(dm.ctx).SetGlobalVariate(one)
	if err != nil {
		return nil, err
	}
	return &pb.GlobalVariateSetReply{}, err
}

// 全局变量 状态设置
func (dm *SettingService) GlobalVariateChangeStatus(r *pb.GlobalVariateSetRequest) (resp *pb.GlobalVariateSetReply, err error) {
	one := &models.CronSetting{
		Id:     r.Id,
		Status: r.Status,
	}
	if r.Id <= 0 {
		return nil, errs.New(nil, "未指定数据")
	}
	if enum.StatusDisable != r.Status && enum.StatusActive != r.Status && enum.StatusDelete != r.Status {
		return nil, errs.New(nil, "状态值错误")
	}

	err = data.NewCronSettingData(dm.ctx).ChangeGlobalVariateStatus(one)
	if err != nil {
		return nil, err
	}
	if enum.StatusActive == r.Status {
		globalVariateList.Set(one.Name, one.Content)
	} else if enum.StatusDisable == r.Status {
		globalVariateList.Del(one.Name)
	}

	return &pb.GlobalVariateSetReply{}, err
}

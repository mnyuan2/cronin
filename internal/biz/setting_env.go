package biz

import (
	"context"
	"cron/internal/basic/conv"
	"cron/internal/basic/db"
	"cron/internal/basic/enum"
	"cron/internal/data"
	"cron/internal/models"
	"cron/internal/pb"
	"errors"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"gorm.io/gorm"
	"time"
)

type EnvContent struct {
	Default int `json:"default"`
}

// 环境设置
type SettingEnvService struct {
	ctx context.Context
	db  *db.MyDB
}

func NewSettingEnvService(ctx context.Context) *SettingEnvService {
	return &SettingEnvService{
		ctx: ctx,
	}
}

// 任务配置列表
func (dm *SettingEnvService) List(r *pb.SettingEnvListRequest) (resp *pb.SettingEnvListReply, err error) {
	// 构建查询条件
	if r.Page <= 1 {
		r.Page = 1
	}
	if r.Size <= 10 {
		r.Size = 20
	}
	resp = &pb.SettingEnvListReply{
		Page: &pb.Page{
			Page: r.Page,
			Size: r.Size,
		},
	}
	list := []*models.CronSetting{}
	resp.Page.Total, err = data.NewCronSettingData(dm.ctx).GetList(models.SceneEnv, "", r.Page, r.Size, &list)
	if err != nil {
		return nil, err
	}

	// 格式化
	resp.List = make([]*pb.SettingEnvListItem, len(list))
	for i, item := range list {
		data := &pb.SettingEnvListItem{
			Id:         item.Id,
			Name:       item.Name,
			Title:      item.Title,
			CreateDt:   item.CreateDt,
			UpdateDt:   item.UpdateDt,
			Status:     item.Status,
			StatusName: enum.StatusMap[item.Status],
		}
		if item.Content != "" {
			con := &EnvContent{}
			jsoniter.UnmarshalFromString(item.Content, con)
			if con.Default == enum.StatusActive {
				data.Default = enum.StatusActive
			} else {
				data.Default = enum.StatusDisable
			}
		}
		resp.List[i] = data
	}

	return resp, err
}

// 设置环境
func (dm *SettingEnvService) Set(r *pb.SettingEnvSetRequest) (resp *pb.SettingEnvSetReply, err error) {
	// 校验
	if r.Name == "" {
		return nil, errors.New("key 为必填")
	}
	if len(r.Name) >= 32 {
		return nil, errors.New("key 长度不得超过32字符")
	}
	if !conv.NewStr().IsLettersAndNumbers(r.Name) {
		return nil, errors.New("key 只能使用字母或数字")
	}
	if r.Title == "" {
		return nil, errors.New("名称 为必填")
	}
	if len(r.Title) >= 210 {
		return nil, errors.New("名称 不得超过210个字符")
	}

	one := &models.CronSetting{}
	_data := data.NewCronSettingData(dm.ctx)
	ti := conv.TimeNew()
	// 分为新增和编辑
	if r.Id > 0 {
		one, err = _data.GetEnvOne(r.Id)
		if err != nil {
			return nil, err
		}
		if r.Name != one.Title { // key 不可以更改，因为可能有关联数据
			return nil, errors.New("key 不可以更改")
		}
	} else {
		one.Scene = models.SceneEnv
		one.Status = enum.StatusActive
		one.CreateDt = ti.String()
		one.Content = "{}"
		one.Name = r.Name
		old, _ := _data.GetOne(db.NewWhere().Eq("scene", models.SceneEnv).Eq("name", r.Name))
		if old.Id > 0 {
			return nil, errors.New("key 已经存在，请更换")
		}
	}

	one.UpdateDt = ti.String()
	one.Title = r.Title

	// 执行写入
	err = _data.Set(one)
	if err != nil {
		return nil, err
	}
	return &pb.SettingEnvSetReply{
		Id: one.Id,
	}, err
}

// 设置源
func (dm *SettingEnvService) SetContent(r *pb.SettingEnvSetRequest) (resp *pb.SettingEnvSetReply, err error) {
	if r.Id <= 0 {
		return nil, errors.New("操作数据未指定")
	}
	if r.Default != enum.StatusActive {
		return nil, errors.New("参数未传递")
	}

	// 查找当前数据并处理
	one := &models.CronSetting{}
	_data := data.NewCronSettingData(dm.ctx)
	ti := conv.TimeNew()
	one, err = _data.GetEnvOne(r.Id)
	if err != nil {
		return nil, err
	}
	oneCon := &EnvContent{}
	jsoniter.UnmarshalFromString(one.Content, oneCon)
	if oneCon.Default == r.Default {
		return nil, errors.New("数据无需操作")
	}
	oneCon.Default = r.Default
	one.Content, _ = jsoniter.MarshalToString(oneCon)
	one.UpdateDt = ti.String()

	// 查找旧数据并处理
	oldOne, err := _data.GetOne(db.NewWhere().Eq("scene", models.SceneEnv).Raw("json_contains(content, '2', '$.default')"))
	oldCon := &EnvContent{}
	jsoniter.UnmarshalFromString(one.Content, oldCon)
	oldCon.Default = enum.StatusDisable
	oldOne.Content, _ = jsoniter.MarshalToString(oldCon)
	oldOne.UpdateDt = ti.String()

	// 执行写入
	db.New(dm.ctx).Transaction(func(tx *gorm.DB) error {
		if err = tx.Select("content", "update_dt").Updates(one).Error; err != nil {
			return err
		}
		if err = tx.Select("content", "update_dt").Updates(oldOne).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &pb.SettingEnvSetReply{
		Id: one.Id,
	}, err
}

// 停用环境
func (dm *SettingEnvService) ChangeStatus(r *pb.SettingChangeStatusRequest) (resp *pb.SettingChangeStatusReply, err error) {
	dm.db = db.New(dm.ctx)

	_data := data.NewCronSettingData(dm.ctx)
	one, err := _data.GetEnvOne(r.Id)
	if err != nil {
		return nil, err
	}
	if one.Id <= 0 {
		return nil, errors.New("操作数据不存在")
	}
	if _, ok := enum.StatusMap[r.Status]; !ok {
		return nil, errors.New("不支持的状态操作")
	}
	if one.Status == r.Status {
		return nil, fmt.Errorf("已是目标状态")
	}
	if r.Status == enum.StatusDisable { // 停用时，不得有进行中的任务
		total := int64(0)
		dm.db.Raw("SELECT count(*) FROM cron_config WHERE env=? and `status`=?", one.Name, models.ConfigStatusActive).Find(&total)
		if total > 0 {
			return nil, fmt.Errorf("环境下存在进行中任务，停用失败！")
		}
		oneCon := &EnvContent{}
		jsoniter.UnmarshalFromString(one.Content, oneCon)
		if oneCon.Default == enum.StatusActive {
			return nil, fmt.Errorf("请先取消默认！")
		}
	}

	one.Status = r.Status
	one.UpdateDt = time.Now().Format(time.DateTime)
	err = _data.ChangeStatus(one)

	return &pb.SettingChangeStatusReply{}, err
}

// 删除日志
func (dm *SettingEnvService) Del(r *pb.SettingEnvDelRequest) (resp *pb.SettingEnvDelReply, err error) {
	if r.Id <= 0 {
		return nil, fmt.Errorf("参数未传递")
	}
	_data := data.NewCronSettingData(dm.ctx)
	one, err := _data.GetEnvOne(r.Id)
	if err != nil {
		return nil, err
	}
	if one.Id <= 0 {
		return nil, errors.New("操作数据不存在")
	}
	if one.Status != enum.StatusDisable {
		return nil, errors.New("请先停用环境") // 可以考虑把所有任务也清空（但是也有风险）。
	}

	err = _data.Del(one.Scene, one.Env, one.Id)
	return resp, err
}

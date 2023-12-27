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
	"strings"
)

// 环境设置
type SettingEnvService struct {
	db *db.Database
}

func NewSettingEnvService() *SettingEnvService {
	return &SettingEnvService{}
}

// 任务配置列表
func (dm *SettingEnvService) List(ctx context.Context, r *pb.SettingEnvListRequest) (resp *pb.SettingEnvListReply, err error) {
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
	resp.Page.Total, err = data.NewCronSettingData(ctx).GetList(models.SceneEnv, r.Page, r.Size, &list)
	if err != nil {
		return nil, err
	}

	// 格式化
	resp.List = make([]*pb.SettingEnvListItem, len(list))
	for i, item := range list {
		data := &pb.SettingEnvListItem{
			Id:         item.Id,
			Key:        item.Title,
			Title:      item.Content,
			CreateDt:   item.CreateDt,
			UpdateDt:   item.UpdateDt,
			Status:     item.Status,
			StatusName: enum.StatusMap[item.Status],
		}
		resp.List[i] = data
	}

	return resp, err
}

// 设置源
func (dm *SettingEnvService) Set(ctx context.Context, r *pb.SettingEnvSetRequest) (resp *pb.SettingEnvSetReply, err error) {
	// 校验
	if r.Key == "" {
		return nil, errors.New("key 为必填")
	}
	if len(r.Key) >= 32 {
		return nil, errors.New("key 长度不得超过32字符")
	}
	if !conv.NewStr().IsLettersAndNumbers(r.Key) {
		return nil, errors.New("key 只能使用字母或数字")
	}
	if r.Title == "" {
		return nil, errors.New("名称 为必填")
	}

	one := &models.CronSetting{}
	_data := data.NewCronSettingData(ctx)
	ti := conv.TimeNew()
	// 分为新增和编辑
	if r.Id > 0 {
		one, err = _data.GetEnvOne(r.Id)
		if err != nil {
			return nil, err
		}
		one.Env = models.EnvDefault
	} else {
		one.Scene = models.SceneEnv
		one.Status = enum.StatusActive
		one.CreateDt = ti.String()
	}

	one.UpdateDt = ti.String()
	one.Title = r.Key
	one.Content = r.Title

	// 执行写入
	err = _data.Set(one)
	if err != nil {
		return nil, err
	}
	return &pb.SettingEnvSetReply{
		Id: one.Id,
	}, err
}

// 删除源
func (dm *SettingEnvService) ChangeStatus(ctx context.Context, r *pb.SettingChangeStatusRequest) (resp *pb.SettingChangeStatusReply, err error) {
	dm.db = db.New(ctx)

	// 同一个任务，这里要加请求锁
	_data := data.NewCronSettingData(ctx)
	one, err := _data.GetSqlSourceOne(r.Id)
	if err != nil {
		return nil, err
	}
	if one.Id <= 0 {
		return nil, errors.New("操作数据不存在")
	}
	// 目前仅支持删除
	if r.Status != enum.StatusDelete {
		return nil, errors.New("不支持的状态操作")
	}
	// 这里还是要做是否使用的检测；
	// 如果使用未启用就联动置空（也不能删除，要么删除任务或者改任务），如果使用并启用禁止删除；
	// 如果没有试用就直接删除。
	list := []string{}
	err = dm.db.Write.Raw(fmt.Sprintf("SELECT `name` FROM `cron_config` WHERE protocol=%v and JSON_CONTAINS(command, '%v', '$.sql.source.id') = 1", models.ProtocolSql, one.Id)).
		Scan(&list).Error
	if err != nil {
		return nil, fmt.Errorf("任务检测错误，%w", err)
	}
	if len(list) > 0 {
		return nil, fmt.Errorf("任务 %s 已使用连接，删除失败！", strings.Join(list, "、"))
	}

	err = _data.Del(one.Scene, one.Id)
	return &pb.SettingChangeStatusReply{}, err
}

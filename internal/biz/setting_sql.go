package biz

import (
	"context"
	"cron/internal/basic/auth"
	"cron/internal/basic/config"
	"cron/internal/basic/conv"
	"cron/internal/basic/db"
	"cron/internal/basic/enum"
	"cron/internal/data"
	"cron/internal/models"
	"cron/internal/pb"
	"errors"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"strings"
)

type SettingSqlService struct {
	db   *db.Database
	ctx  context.Context
	user *auth.UserToken
}

func NewSettingSqlService(ctx context.Context, user *auth.UserToken) *SettingSqlService {
	return &SettingSqlService{
		ctx:  ctx,
		user: user,
	}
}

// 任务配置列表
func (dm *SettingSqlService) List(r *pb.SettingSqlListRequest) (resp *pb.SettingSqlListReply, err error) {
	// 构建查询条件
	if r.Page <= 1 {
		r.Page = 1
	}
	if r.Size <= 10 {
		r.Size = 20
	}
	resp = &pb.SettingSqlListReply{
		Page: &pb.Page{
			Page: r.Page,
			Size: r.Size,
		},
	}
	list := []*models.CronSetting{}
	resp.Page.Total, err = data.NewCronSettingData(dm.ctx).GetList(models.SceneSqlSource, dm.user.Env, r.Page, r.Size, &list)
	if err != nil {
		return nil, err
	}

	// 格式化
	resp.List = make([]*pb.SettingSqlListItem, len(list))
	for i, item := range list {
		data := &pb.SettingSqlListItem{
			Id:       item.Id,
			Title:    item.Title,
			CreateDt: item.CreateDt,
			UpdateDt: item.UpdateDt,
			Source:   &pb.SettingSqlSource{},
		}
		jsoniter.UnmarshalFromString(item.Content, data.Source)
		resp.List[i] = data
	}

	return resp, err
}

// 设置源
func (dm *SettingSqlService) Set(r *pb.SettingSqlSetRequest) (resp *pb.SettingSqlSetReply, err error) {
	one := &models.CronSetting{}
	_data := data.NewCronSettingData(dm.ctx)
	ti := conv.TimeNew()
	oldSource := &pb.SettingSqlSource{}
	// 分为新增和编辑
	if r.Id > 0 {
		one, err = _data.GetSqlSourceOne(dm.user.Env, r.Id)
		if err != nil {
			return nil, err
		}
		jsoniter.UnmarshalFromString(one.Content, oldSource)
	} else {
		one.Scene = models.SceneSqlSource
		one.Env = dm.user.Env
		one.Status = enum.StatusActive
		one.CreateDt = ti.String()
	}

	// 提交密码与旧密码不一致就加密
	if r.Source.Password != "" && r.Source.Password != oldSource.Password {
		r.Source.Password, err = models.SqlSourceEncrypt(r.Source.Password)
		if err != nil {
			return nil, fmt.Errorf("加密失败，%w", err)
		}
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

func (dm *SettingSqlService) Ping(r *pb.SettingSqlPingRequest) (resp *pb.SettingSqlPingReply, err error) {
	password, err := models.SqlSourceDecode(r.Password)
	if err != nil {
		return nil, fmt.Errorf("密码异常,%w", err)
	}
	conf := config.DataBaseConf{
		Source: fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=false&loc=Local",
			r.Username, password, r.Hostname, r.Port, r.Database),
	}
	err = db.Conn(conf).Error
	if err != nil {
		return nil, err
	}
	return &pb.SettingSqlPingReply{}, nil
}

// 删除源
func (dm *SettingSqlService) ChangeStatus(r *pb.SettingChangeStatusRequest) (resp *pb.SettingChangeStatusReply, err error) {
	dm.db = db.New(dm.ctx)

	// 同一个任务，这里要加请求锁
	_data := data.NewCronSettingData(dm.ctx)
	one, err := _data.GetSqlSourceOne(dm.user.Env, r.Id)
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

	err = _data.Del(one.Scene, dm.user.Env, one.Id)
	return &pb.SettingChangeStatusReply{}, err
}

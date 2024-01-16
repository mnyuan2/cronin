package biz

import (
	"context"
	"cron/internal/basic/auth"
	"cron/internal/basic/config"
	"cron/internal/basic/conv"
	"cron/internal/basic/db"
	"cron/internal/basic/enum"
	"cron/internal/basic/grpcurl"
	"cron/internal/models"
	"cron/internal/pb"
	"errors"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"runtime"
	"strings"
)

type FoundationService struct {
	ctx  context.Context
	user *auth.UserToken
}

type DicGetItem struct {
	// 键
	Id int32 `json:"id"`
	// 备用键
	Key string `json:"key"`
	// 值
	Name string `json:"name"`
	// 其它数据，用于业务放关联操作
	Extend string `json:"extend"`
}

func NewDicService(ctx context.Context, user *auth.UserToken) *FoundationService {
	return &FoundationService{
		ctx:  ctx,
		user: user,
	}
}

// 获得枚举
func (dm *FoundationService) DicGets(r *pb.DicGetsRequest) (resp *pb.DicGetsReply, err error) {
	types := []int{}
	err = conv.NewStr().Slice(r.Types, &types)
	if err != nil {
		return nil, err
	}

	resp = &pb.DicGetsReply{
		Maps: map[int]*pb.DicGetsList{},
	}
	for _, t := range types {
		list := &pb.DicGetsList{}
		if t <= 1000 {
			list.List, err = dm.getDb(t)
			if err != nil {
				return nil, err
			}
		} else {
			list.List, err = dm.getEnum(t)
			if err != nil {
				return nil, err
			}
		}

		resp.Maps[t] = list
	}

	return resp, err
}

// 通过数据库获取
func (dm *FoundationService) getDb(t int) ([]*pb.DicGetItem, error) {
	_sql := ""
	w := db.NewWhere()

	switch t {
	case enum.DicSqlSource:
		_sql = "SELECT id,title as name FROM `cron_setting` %WHERE ORDER BY update_dt,id desc"
		w.Eq("scene", models.SceneSqlSource).Eq("status", enum.StatusActive).Eq("env", dm.user.Env, db.RequiredOption())
	case enum.DicEnv:
		_sql = "SELECT id,name 'key', title name, content extend FROM `cron_setting` %WHERE ORDER BY update_dt,id desc"
		w.Eq("scene", models.SceneEnv).Eq("status", enum.StatusActive)
	case enum.DicMsg:
		_sql = "SELECT id, title name FROM `cron_setting` %WHERE ORDER BY update_dt,id desc"
		w.Eq("scene", models.SceneMsg).Eq("status", enum.StatusActive)
	case enum.DicUser:
		_sql = "SELECT id, username name FROM `cron_user` ORDER BY sort asc,id desc"
	}

	items := []*pb.DicGetItem{}
	if _sql != "" {
		temp := []*DicGetItem{}
		where, args := w.Build()
		_sql = strings.Replace(_sql, "%WHERE", "WHERE "+where, 1)

		err := db.New(dm.ctx).Raw(_sql, args...).Scan(&temp).Error
		if err != nil {
			return nil, err
		}
		for _, v := range temp {
			item := &pb.DicGetItem{
				Id:     v.Id,
				Key:    v.Key,
				Name:   v.Name,
				Extend: &pb.DicExtendItem{},
			}
			if v.Extend != "" {
				if err = jsoniter.Unmarshal([]byte(v.Extend), item.Extend); err != nil {
					return nil, err
				}
			}
			items = append(items, item)
		}
	}

	return items, nil
}

// 通过枚举获取
func (dm *FoundationService) getEnum(t int) ([]*pb.DicGetItem, error) {
	// 待完善
	return nil, nil
}

func (dm *FoundationService) SystemInfo(r *pb.SystemInfoRequest) (resp *pb.SystemInfoReply, err error) {
	resp = &pb.SystemInfoReply{
		Version: config.Version,
		CmdName: "sh",
	}
	// 根据运行环境确认cmd的类型
	if runtime.GOOS == "windows" {
		resp.CmdName = "cmd"
	}
	// 查默认环境
	envData := &DicGetItem{}
	err = db.New(dm.ctx).Raw(`SELECT id, name as 'key', title as name FROM cron_setting WHERE scene='env' and status=2 and json_contains(content, '2', '$.default');`).Scan(&envData).Error
	if err != nil {
		return resp, fmt.Errorf("环境查询异常,%w", err)
	}
	if envData.Id == 0 {
		return resp, errors.New("未设置运行环境")
	}
	if envData.Name == "" {
		return resp, errors.New("运行环境设置异常")
	}
	resp.Env = envData.Key
	resp.EnvName = envData.Name

	return resp, nil
}

func (dm *FoundationService) ParseProto(r *pb.ParseProtoRequest) (resp *pb.ParseProtoReply, err error) {
	fds, err := grpcurl.ParseProtoString(r.Proto)
	if err != nil {
		return nil, fmt.Errorf("无法解析给定的proto文件: %w", err)
	}

	resp = &pb.ParseProtoReply{Actions: grpcurl.ParseProtoMethods(fds)}

	return resp, nil
}

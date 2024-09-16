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
	"time"
)

type FoundationService struct {
	ctx  context.Context
	user *auth.UserToken
}

type DicGetItem struct {
	// 键
	Id int `json:"id"`
	// 备用键
	Key string `json:"key"`
	// 值
	Name string `json:"name"`
	// 其它数据，用于业务放关联操作
	Extend string `json:"extend"`
}

//var enumSource = map[int][]*pb.DicGetItem{
//	enum.DicCmdType: {
//		{Key: "cmd", Name: "cmd"},
//		{Key: "bash", Name: "bash"},
//		{Key: "sh", Name: "sh"},
//	},
//}

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
	items := []*pb.DicGetItem{}
	switch t {
	case enum.DicSqlSource:
		_sql = "SELECT id,title as name, concat('{\"driver\":',content->'$.sql.driver','}') extend  FROM `cron_setting` %WHERE ORDER BY update_dt,id desc"
		w.Eq("scene", models.SceneSqlSource).Eq("status", enum.StatusActive).Eq("env", dm.user.Env, db.RequiredOption())
	case enum.DicJenkinsSource:
		_sql = "SELECT id,title as name FROM `cron_setting` %WHERE ORDER BY update_dt,id desc"
		w.Eq("scene", models.SceneJenkinsSource).Eq("status", enum.StatusActive).Eq("env", dm.user.Env, db.RequiredOption())
	case enum.DicGitSource:
		_sql = "SELECT id,title as name FROM `cron_setting` %WHERE ORDER BY update_dt,id desc"
		w.Eq("scene", models.SceneGitSource).Eq("status", enum.StatusActive).Eq("env", dm.user.Env, db.RequiredOption())
	case enum.DicHostSource:
		_sql = "SELECT id,title as name FROM `cron_setting` %WHERE ORDER BY update_dt,id desc"
		w.Eq("scene", models.SceneHostSource).Eq("status", enum.StatusActive).Eq("env", dm.user.Env, db.RequiredOption())
		items = append(items, &pb.DicGetItem{
			Id:     -1,
			Name:   "本机",
			Extend: &pb.DicExtendItem{},
		})
	case enum.DicEnv:
		_sql = "SELECT id,name 'key', title name, content extend FROM `cron_setting` %WHERE ORDER BY update_dt,id desc"
		w.Eq("scene", models.SceneEnv).Eq("status", enum.StatusActive)
	case enum.DicMsg:
		_sql = "SELECT id, title name FROM `cron_setting` %WHERE ORDER BY update_dt,id desc"
		w.Eq("scene", models.SceneMsg).Eq("status", enum.StatusActive)
	case enum.DicUser:
		_sql = "SELECT id, username name FROM `cron_user` ORDER BY sort asc,id desc"
	case enum.DicRole:
		_sql = "SELECT id, name FROM cron_auth_role %WHERE "
		w.Eq("status", enum.StatusActive)
	}

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
	items := []*pb.DicGetItem{}

	switch t {
	case enum.DicCmdType:
		if runtime.GOOS == "windows" {
			items = append(items, &pb.DicGetItem{Key: "cmd", Name: "cmd"})
		} else {
			items = append(items, &pb.DicGetItem{Key: "bash", Name: "bash"},
				&pb.DicGetItem{Key: "sh", Name: "sh"})
		}
	case enum.DicGitEvent:
		items = []*pb.DicGetItem{
			{Id: enum.GitEventPullsCreate, Name: enum.GitEventMap[enum.GitEventPullsCreate]},
			{Id: enum.GitEventPullsMerge, Name: enum.GitEventMap[enum.GitEventPullsMerge]},
			{Id: enum.GitEventFileUpdate, Name: enum.GitEventMap[enum.GitEventFileUpdate]},
		}
	case enum.DicSqlDriver:
		items = []*pb.DicGetItem{
			{Key: enum.SqlDriverMysql, Name: enum.SqlDriverMysql},
			{Key: enum.SqlDriverClickhouse, Name: enum.SqlDriverClickhouse},
		}
	case enum.DicConfigStatus:
		items = []*pb.DicGetItem{
			{Id: models.ConfigStatusDisable, Name: models.ConfigStatusMap[models.ConfigStatusDisable]},
			{Id: models.ConfigStatusAudited, Name: models.ConfigStatusMap[models.ConfigStatusAudited]},
			{Id: models.ConfigStatusReject, Name: models.ConfigStatusMap[models.ConfigStatusReject]},
			{Id: models.ConfigStatusActive, Name: models.ConfigStatusMap[models.ConfigStatusActive]},
			{Id: models.ConfigStatusFinish, Name: models.ConfigStatusMap[models.ConfigStatusFinish]},
			{Id: models.ConfigStatusError, Name: models.ConfigStatusMap[models.ConfigStatusError]},
		}
	case enum.DicProtocolType:
		items = []*pb.DicGetItem{
			{Id: models.ProtocolHttp, Name: models.ProtocolMap[models.ProtocolHttp]},
			{Id: models.ProtocolRpc, Name: models.ProtocolMap[models.ProtocolRpc]},
			{Id: models.ProtocolCmd, Name: models.ProtocolMap[models.ProtocolCmd]},
			{Id: models.ProtocolSql, Name: models.ProtocolMap[models.ProtocolSql]},
			{Id: models.ProtocolJenkins, Name: models.ProtocolMap[models.ProtocolJenkins]},
			{Id: models.ProtocolGit, Name: models.ProtocolMap[models.ProtocolGit]},
		}
	}

	return items, nil
}

func (dm *FoundationService) SystemInfo(r *pb.SystemInfoRequest) (resp *pb.SystemInfoReply, err error) {
	resp = &pb.SystemInfoReply{
		Version:     config.Version,
		CmdName:     "sh",
		CurrentDate: time.Now().Format(time.RFC3339),
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

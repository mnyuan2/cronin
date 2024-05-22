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
)

type RoleService struct {
	db  *db.MyDB
	ctx context.Context
	//user *auth.UserToken
}

func NewRoleService(ctx context.Context) *RoleService {
	return &RoleService{
		ctx: ctx,
	}
}

// 任务配置列表
func (dm *RoleService) List(r *pb.RoleListRequest) (resp *pb.RoleListReply, err error) {
	resp = &pb.RoleListReply{
		List: make([]*pb.RoleListItem, 0),
	}

	list, err := data.NewCronAuthRoleData(dm.ctx).GetList(db.NewWhere())
	if err != nil {
		return nil, err
	}
	for _, item := range list {
		row := &pb.RoleListItem{
			Id:         item.Id,
			Name:       item.Name,
			Remark:     item.Remark,
			AuthIds:    []int{},
			Status:     item.Status,
			StatusName: enum.StatusMap[item.Status],
		}
		conv.NewStr().Slice(item.AuthIds, &row.AuthIds)
		resp.List = append(resp.List, row)
	}

	return resp, err
}

// 任务配置列表
func (dm *RoleService) AuthList(r *pb.AuthListRequest) (resp *pb.AuthListReply, err error) {
	resp = &pb.AuthListReply{
		List: make([]*pb.AuthListItem, 0),
	}

	limited, is := map[int]int{}, false
	if len(r.RoleIds) > 0 {
		list, err := data.NewCronAuthRoleData(dm.ctx).GetList(db.NewWhere().In("id", r.RoleIds))
		if err != nil {
			return nil, err
		}
		if len(list) != len(r.RoleIds) {
			return nil, errors.New("角色信息有误")
		}
		for _, item := range list {
			temp := []int{}
			conv.NewStr().Slice(item.AuthIds, &temp)
			for _, i := range temp {
				limited[i] = i
			}
		}
		is = true
	}

	list := data.NewAuthData().List()
	for _, item := range list {
		if item.Type != data.AuthTypeGrant {
			continue // 不需要授权的忽略
		}
		if is {
			if _, ok := limited[item.Id]; !ok {
				continue
			}
		}
		resp.List = append(resp.List, &pb.AuthListItem{
			Id:    item.Id,
			Name:  item.Title,
			Path:  item.Path,
			Group: item.Group,
		})
	}

	return resp, err
}

// 设置角色
func (dm *RoleService) AuthSet(r *pb.RoleAuthSetRequest) (resp *pb.RoleAuthSetReply, err error) {
	if r.Id <= 0 {
		return nil, errors.New("未指定角色")
	}

	one := &models.CronAuthRole{}
	_data := data.NewCronAuthRoleData(dm.ctx)
	// 分为新增和编辑
	if one, err = _data.GetOne(r.Id); err != nil {
		return nil, err
	}

	one.AuthIds, _ = conv.Int64s().Join(r.AuthIds)

	// 执行写入
	err = _data.SetAuthIds(one)
	if err != nil {
		return nil, err
	}
	return &pb.RoleAuthSetReply{}, nil
}

// 设置角色
func (dm *RoleService) Set(r *pb.RoleSetRequest) (resp *pb.RoleSetReply, err error) {
	if r.Name == "" {
		return nil, errors.New("名称不得为空")
	}

	one := &models.CronAuthRole{}
	_data := data.NewCronAuthRoleData(dm.ctx)
	// 分为新增和编辑
	if r.Id > 0 {
		if one, err = _data.GetOne(r.Id); err != nil {
			return nil, err
		}
	} else {
		one.Status = enum.StatusActive
	}

	one.Name = r.Name
	one.Remark = r.Remark
	// 执行写入
	err = _data.Set(one)
	if err != nil {
		return nil, err
	}
	return &pb.RoleSetReply{
		Id: one.Id,
	}, err
}

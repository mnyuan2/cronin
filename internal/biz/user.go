package biz

import (
	"context"
	"cron/internal/basic/conv"
	"cron/internal/basic/db"
	"cron/internal/data"
	"cron/internal/models"
	"cron/internal/pb"
	"errors"
)

type UserService struct {
	db  *db.MyDB
	ctx context.Context
	//user *auth.UserToken
}

func NewUserService(ctx context.Context) *UserService {
	return &UserService{
		ctx: ctx,
	}
}

// 任务配置列表
func (dm *UserService) List(r *pb.UserListRequest) (resp *pb.UserListReply, err error) {
	if r.Page <= 1 {
		r.Page = 1
	}
	if r.Size <= 10 {
		r.Size = 20
	}
	resp = &pb.UserListReply{
		Page: &pb.Page{
			Page: r.Page,
			Size: r.Size,
		},
		List: make([]*pb.UserListItem, 0),
	}
	w := db.NewWhere()

	resp.Page.Total, err = data.NewCronUserData(dm.ctx).GetListPage(w, r.Page, r.Size, &resp.List)
	if err != nil {
		return nil, err
	}

	return resp, err
}

// 设置源
func (dm *UserService) Set(r *pb.UserSetRequest) (resp *pb.UserSetReply, err error) {
	if r.Username == "" {
		return nil, errors.New("名称不得为空")
	}

	one := &models.CronUser{}
	_data := data.NewCronUserData(dm.ctx)
	ti := conv.TimeNew()
	// 分为新增和编辑
	if r.Id > 0 {
		if one, err = _data.GetOne(r.Id); err != nil {
			return nil, err
		}
	} else {
		one.CreateDt = ti.String()
	}

	one.UpdateDt = ti.String()
	one.Username = r.Username
	one.Mobile = r.Mobile
	one.Sort = r.Sort

	// 执行写入
	err = _data.Set(one)
	if err != nil {
		return nil, err
	}
	return &pb.UserSetReply{
		Id: one.Id,
	}, err
}

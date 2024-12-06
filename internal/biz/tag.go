package biz

import (
	"context"
	"cron/internal/basic/auth"
	"cron/internal/basic/conv"
	"cron/internal/basic/db"
	"cron/internal/basic/enum"
	"cron/internal/data"
	"cron/internal/models"
	"cron/internal/pb"
	"errors"
	"fmt"
)

type CronTagService struct {
	ctx  context.Context
	user *auth.UserToken
}

func NewCronTagService(ctx context.Context, user *auth.UserToken) *CronTagService {
	return &CronTagService{
		ctx:  ctx,
		user: user,
	}
}

// List 列表
func (dm *CronTagService) List(r *pb.TagListRequest) (resp *pb.TagListReply, err error) {
	resp = &pb.TagListReply{List: []*pb.TagListItem{}}

	err = db.New(dm.ctx).Model(&models.CronTag{}).Where("status=?", models.ConfigStatusActive).Scan(&resp.List).Error
	if err != nil {
		return nil, err
	}
	return resp, err
}

// 设置环境
func (dm *CronTagService) Set(r *pb.TagSetRequest) (resp *pb.SettingEnvSetReply, err error) {
	if r.Name == "" {
		return nil, errors.New("名称 为必填")
	}
	if len(r.Name) > 32 {
		return nil, errors.New("名称 长度不得超过32字符")
	}
	//if !conv.NewStr().IsLettersAndNumbers(r.Name) {
	//	return nil, errors.New("key 只能使用字母或数字")
	//}

	one := &models.CronTag{}
	_data := data.NewCronTagData(dm.ctx)
	ti := conv.TimeNew()
	// 分为新增和编辑
	if r.Id > 0 {
		one, err = _data.GetOne("id=?", r.Id)
		if err != nil {
			return nil, err
		}
	} else {
		one.Status = enum.StatusActive
		one.CreateDt = ti.String()
		one.CreateUserId = dm.user.UserId
		one.CreateUserName = dm.user.UserName
	}
	if r.Name != one.Name {
		_, err := _data.GetOne("name=? and status=?", r.Name, models.ConfigStatusActive)
		if err == nil {
			return nil, errors.New("标签名称 已经存在，请更换")
		}
	}

	one.Name = r.Name
	one.Remark = r.Remark
	one.UpdateDt = ti.String()
	one.UpdateUserId = dm.user.UserId
	one.UpdateUserName = dm.user.UserName

	err = _data.Set(one)
	if err != nil {
		return nil, err
	}
	return &pb.SettingEnvSetReply{
		Id: one.Id,
	}, err
}

// Del 删除
func (dm *CronTagService) ChangeStatus(r *pb.TagChangeStatusRequest) (resp *pb.TagChangeStatusReply, err error) {
	if r.Status != models.ConfigStatusClosed {
		return nil, fmt.Errorf("不支持的状态操作")
	}

	_data := data.NewCronTagData(dm.ctx)
	one, err := _data.GetOne("id=? and status=?", r.Id, models.ConfigStatusActive)
	if err != nil {
		return nil, err
	}

	one.Status = r.Status
	one.UpdateDt = conv.TimeNew().String()
	one.UpdateUserId = dm.user.UserId
	one.UpdateUserName = dm.user.UserName
	if err = _data.ChangeStatus(one); err != nil {
		return nil, err
	}
	// 已经使用的标记要进行移除，

	return &pb.TagChangeStatusReply{}, err
}

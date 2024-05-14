package biz

import (
	"context"
	"cron/internal/basic/auth"
	"cron/internal/basic/conv"
	"cron/internal/basic/db"
	"cron/internal/basic/enum"
	"cron/internal/basic/errs"
	"cron/internal/data"
	"cron/internal/models"
	"cron/internal/pb"
	"errors"
	"fmt"
	"strings"
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
	for _, item := range resp.List {
		item.StatusName = enum.StatusMap[item.Status]
	}

	return resp, err
}

// 设置源
func (dm *UserService) Set(r *pb.UserSetRequest) (resp *pb.UserSetReply, err error) {
	if r.Username == "" {
		return nil, errors.New("名称不得为空")
	}
	r.Account = strings.ToUpper(strings.TrimSpace(r.Account))

	one := &models.CronUser{}
	_data := data.NewCronUserData(dm.ctx)
	ti := conv.TimeNew()
	// 分为新增和编辑
	if r.Id > 0 {
		if one, err = _data.GetOne(r.Id); err != nil {
			return nil, err
		}
	} else {
		if r.Password == "" {
			return nil, errors.New("密码不得为空")
		}
		if one.Password, err = models.SqlSourceEncrypt(r.Password); err != nil {
			return nil, fmt.Errorf("密码异常，%w", err)
		}
		one.Account = r.Account
		one.Status = enum.StatusActive
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

// 修改密码
func (dm UserService) ChangePassword(r *pb.UserSetRequest) (resp *pb.UserSetReply, err error) {
	if r.Id <= 0 {
		return nil, errors.New("用户不得为空")
	}
	if r.Password == "" {
		return nil, errors.New("密码不得为空")
	}

	_data := data.NewCronUserData(dm.ctx)
	one, err := _data.GetOne(r.Id)
	if err != nil {
		return nil, err
	}

	if one.Password, err = models.SqlSourceEncrypt(r.Password); err != nil {
		return nil, fmt.Errorf("密码异常，%w", err)
	}
	one.UpdateDt = conv.TimeNew().String()

	// 执行写入
	err = _data.ChangePassword(one)
	if err != nil {
		return nil, err
	}
	return &pb.UserSetReply{}, nil
}

// 设置源
func (dm *UserService) ChangeStatus(r *pb.UserChangeStatusRequest) (resp *pb.UserChangeStatusReply, err error) {
	if r.Id <= 0 {
		return nil, errors.New("用户不得为空")
	}
	if _, ok := enum.StatusMap[r.Status]; !ok {
		return nil, errors.New("状态不合法")
	}

	_data := data.NewCronUserData(dm.ctx)
	ti := conv.TimeNew()
	one, err := _data.GetOne(r.Id)
	if err != nil {
		return nil, err
	}
	if one.Status == r.Status {
		return &pb.UserChangeStatusReply{}, nil
	}

	one.Status = r.Status
	one.UpdateDt = ti.String()

	// 执行写入
	err = _data.ChangeStatus(one)
	if err != nil {
		return nil, err
	}
	return &pb.UserChangeStatusReply{}, nil
}

// 设置账号
func (dm *UserService) ChangeAccount(r *pb.UserSetRequest) (resp *pb.UserSetReply, err error) {
	if r.Id <= 0 {
		return nil, errors.New("用户不得为空")
	}

	_data := data.NewCronUserData(dm.ctx)
	ti := conv.TimeNew()
	one, err := _data.GetOne(r.Id)
	if err != nil {
		return nil, err
	}
	newAccount := strings.ToUpper(strings.TrimSpace(r.Account))
	if newAccount == one.Account {
		return &pb.UserSetReply{}, nil
	}
	if !conv.NewStr().IsChinese(r.Account) {
		return nil, errors.New("可输入字母、数字、符号")
	}

	one.Account = newAccount
	one.UpdateDt = ti.String()

	// 执行写入
	err = _data.ChangeAccount(one)
	if err != nil {
		return nil, err
	}
	return &pb.UserSetReply{}, nil
}

// Detail 用户详情
func (dm *UserService) Detail(r *pb.UserDetailRequest) (resp *pb.UserDetailReply, err error) {
	if r.Id <= 0 {
		return nil, errs.New(nil, errs.ParamNotFound)
	}

	user, err := data.NewCronUserData(dm.ctx).GetOne(r.Id)
	if err != nil {
		return nil, err
	}
	if user.Status == enum.StatusDelete {
		return nil, errors.New("数据已被删除")
	}

	resp = &pb.UserDetailReply{
		Id:         user.Id,
		Account:    user.Account,
		Username:   user.Username,
		Mobile:     user.Mobile,
		Sort:       user.Sort,
		Status:     user.Status,
		StatusName: enum.StatusMap[user.Status],
		UpdateDt:   user.UpdateDt,
		CreateDt:   user.CreateDt,
	}
	return resp, err
}

// Login 用户详情
func (dm *UserService) Login(r *pb.UserLoginRequest) (resp *pb.UserLoginReply, err error) {
	if r.Account == "" || r.Password == "" {
		return nil, errors.New("账号密码不得为空")
	}
	password, err := models.SqlSourceEncrypt(r.Password)
	if err != nil {
		return nil, errors.New("密码异常")
	}

	user, err := data.NewCronUserData(dm.ctx).Login(r.Account, password)
	if err != nil {
		return nil, fmt.Errorf("登录错误，%w", err)
	}
	if user.Id <= 0 {
		return nil, errors.New("账号或密码错误")
	}
	if user.Status != enum.StatusActive {
		return nil, errors.New("账号异常，请联系管理员")
	}

	resp = &pb.UserLoginReply{
		User: &pb.UserDetailReply{
			Id:         user.Id,
			Username:   user.Username,
			Mobile:     user.Mobile,
			Sort:       user.Sort,
			Status:     user.Status,
			StatusName: enum.StatusMap[user.Status],
			UpdateDt:   user.UpdateDt,
			CreateDt:   user.CreateDt,
		},
	}
	if resp.Token, err = auth.GenToken(user.Id, user.Username); err != nil {
		return nil, err
	}
	return resp, err
}

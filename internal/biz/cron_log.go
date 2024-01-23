package biz

import (
	"context"
	"cron/internal/basic/auth"
	"cron/internal/basic/db"
	"cron/internal/data"
	"cron/internal/models"
	"cron/internal/pb"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"time"
)

type CronLogService struct {
	ctx  context.Context
	user *auth.UserToken
}

func NewCronLogService(ctx context.Context, user *auth.UserToken) *CronLogService {
	return &CronLogService{
		ctx:  ctx,
		user: user,
	}
}

// 通过配置查询日志
func (dm *CronLogService) ByConfig(r *pb.CronLogByConfigRequest) (resp *pb.CronLogByConfigResponse, err error) {
	env := dm.user.Env
	if r.ConfId <= 0 {
		env = ""
	}
	w := db.NewWhere().
		Eq("conf_id", r.ConfId, db.RequiredOption()).
		Eq("env", env, db.RequiredOption())
	list := []*models.CronLog{}
	_, err = data.NewCronLogData(dm.ctx).GetList(w, 1, r.Limit, &list)
	resp = &pb.CronLogByConfigResponse{List: make([]*pb.CronLogItem, len(list))}
	for i, one := range list {
		item := &pb.CronLogItem{
			Id:            one.Id,
			ConfId:        one.ConfId,
			CreateDt:      one.CreateDt,
			Duration:      one.Duration,
			Status:        one.Status,
			StatusName:    models.LogStatusMap[one.Status],
			StatusDesc:    one.StatusDesc,
			Body:          one.Body,
			Snap:          one.Snap,
			MsgStatus:     one.MsgStatus,
			MsgStatusName: models.LogStatusMap[one.Status],
			MsgBody:       []string{},
		}
		if item.MsgStatusName == "" {
			item.MsgStatusName = "无"
		}
		jsoniter.UnmarshalFromString(one.MsgBody, &item.MsgBody)

		resp.List[i] = item
	}

	return resp, err
}

// 删除日志
func (dm *CronLogService) Del(r *pb.CronLogDelRequest) (resp *pb.CronLogDelResponse, err error) {
	if r.Retention == "" {
		return nil, fmt.Errorf("retention 参数为必须")
	}

	re, err := time.ParseDuration(r.Retention)
	if err != nil {
		return nil, fmt.Errorf("retention 参数有误, %s", err.Error())
	} else if re.Hours() < 24 {
		return nil, fmt.Errorf("retention 参数不得小于24h")
	}
	end := time.Now().Add(-re)
	resp = &pb.CronLogDelResponse{}
	resp.Count, err = data.NewCronLogData(dm.ctx).DelBatch(end)

	return resp, err
}

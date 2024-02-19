package biz

import (
	"context"
	"cron/internal/basic/auth"
	"cron/internal/basic/db"
	"cron/internal/basic/enum"
	"cron/internal/data"
	"cron/internal/models"
	"cron/internal/pb"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"time"
)

type CronPipelineService struct {
	ctx  context.Context
	user *auth.UserToken
}

func NewCronPipelineService(ctx context.Context, user *auth.UserToken) *CronPipelineService {
	return &CronPipelineService{
		ctx:  ctx,
		user: user,
	}
}

// 任务配置列表
func (dm *CronPipelineService) List(r *pb.CronPipelineListRequest) (resp *pb.CronPipelineListReply, err error) {
	if r.Type == 0 {
		r.Type = models.TypeOnce
	}
	w := db.NewWhere().Eq("type", r.Type).Eq("env", dm.user.Env, db.RequiredOption())
	// 构建查询条件
	if r.Page <= 1 {
		r.Page = 1
	}
	if r.Size <= 10 {
		r.Size = 10
	}
	resp = &pb.CronPipelineListReply{
		List: []*pb.CronPipelineListItem{},
		Page: &pb.Page{
			Page: r.Page,
			Size: r.Size,
		},
	}
	resp.Page.Total, err = data.NewCronPipelineData(dm.ctx).GetList(w, r.Page, r.Size, &resp.List)
	topList := map[int]*data.SumConfTop{}
	if len(resp.List) > 0 {
		endTime := time.Now()
		startTime := time.Now().Add(-time.Hour * 24 * 7) // 取七天前
		ids := make([]int, len(resp.List))
		for i, temp := range resp.List {
			ids[i] = temp.Id
		}
		topList, _ = data.NewCronLogData(dm.ctx).SumConfTopError(dm.user.Env, ids, startTime, endTime, 7)
	}

	for _, item := range resp.List {
		item.ConfigIds = []int{}
		item.MsgSet = []*pb.CronMsgSet{}
		item.StatusName = models.ConfigStatusMap[item.Status]
		jsoniter.Unmarshal(item.ConfigIdsStr, item.ConfigIds)
		jsoniter.Unmarshal(item.MsgSetStr, &item.MsgSet)
		if top, ok := topList[item.Id]; ok {
			item.TopNumber = top.TotalNumber
			item.TopErrorNumber = top.ErrorNumber
		}
	}

	return resp, err
}

// 任务配置
func (dm *CronPipelineService) Set(r *pb.CronPipelineSetRequest) (resp *pb.CronPipelineSetReply, err error) {

	d := &models.CronPipeline{}
	if r.Id > 0 {
		da := data.NewCronPipelineData(dm.ctx)
		d, err = da.GetOne(dm.user.Env, r.Id)
		if err != nil {
			return nil, err
		}
		if d.Status == enum.StatusActive {
			return nil, fmt.Errorf("请先停用任务后编辑")
		}
	} else {
		d.Status = enum.StatusDisable
		d.Env = dm.user.Env
	}

	d.Name = r.Name
	d.Spec = r.Spec
	d.Remark = r.Remark
	d.ConfigIds, _ = jsoniter.Marshal(r.ConfigIds)
	d.MsgSet, _ = jsoniter.Marshal(r.MsgSet)
	d.Type = r.Type
	if r.Type == models.TypeCycle {
		if _, err = secondParser.Parse(d.Spec); err != nil {
			return nil, fmt.Errorf("时间格式不规范，%s", err.Error())
		}
	} else if r.Type == models.TypeOnce {
		if _, err = NewScheduleOnce(d.Spec); err != nil {
			return nil, err
		}
	} else {
		return nil, fmt.Errorf("类型输入有误")
	}

	for _, id := range r.ConfigIds {
		fmt.Println("校验任务", id) // 晚一点补充
	}

	for i, msg := range r.MsgSet {
		if msg.MsgId == 0 {
			return nil, fmt.Errorf("推送%v未设置消息模板", i)
		}
	}

	err = data.NewCronPipelineData(dm.ctx).Set(d)
	if err != nil {
		return nil, err
	}
	return &pb.CronPipelineSetReply{
		Id: d.Id,
	}, err
}

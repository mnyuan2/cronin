package biz

import (
	"context"
	"cron/internal/basic/auth"
	"cron/internal/basic/config"
	"cron/internal/basic/db"
	"cron/internal/basic/enum"
	"cron/internal/basic/errs"
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
	w := db.NewWhere().Eq("type", r.Type).Eq("env", dm.user.Env, db.RequiredOption()).Eq("status", r.Status)
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
	resp.Page.Total, err = data.NewCronPipelineData(dm.ctx).ListPage(w, r.Page, r.Size, &resp.List)
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
		item.Configs = []*pb.CronConfigListItem{}
		item.MsgSet = []*pb.CronMsgSet{}
		item.StatusName = models.ConfigStatusMap[item.Status]
		item.ConfigDisableActionName = models.DisableActionMap[item.ConfigDisableAction]
		if err = jsoniter.Unmarshal(item.ConfigIdsStr, &item.ConfigIds); err != nil {
			fmt.Println("	", err.Error())
		}
		if err = jsoniter.Unmarshal(item.ConfigsStr, &item.Configs); err != nil {
			fmt.Println("	", err.Error())
		}
		if err = jsoniter.Unmarshal(item.MsgSetStr, &item.MsgSet); err != nil {
			fmt.Println("	", err.Error())
		}
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

	if r.VarParams != "" {
		prams := map[string]any{}
		if err = jsoniter.UnmarshalFromString(r.VarParams, &prams); err != nil {
			return nil, fmt.Errorf("变量参数序列化错误，%w", err)
		}
	}

	d.Name = r.Name
	d.Spec = r.Spec
	d.Remark = r.Remark
	d.VarParams = r.VarParams
	d.ConfigIds, _ = jsoniter.Marshal(r.ConfigIds)
	d.Configs, _ = jsoniter.Marshal(r.Configs)
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
	if _, ok := models.DisableActionMap[r.ConfigDisableAction]; !ok {
		return nil, errs.New(nil, "任务停用行为未正确设置")
	}
	if _, ok := models.ErrActionMap[r.ConfigErrAction]; !ok {
		return nil, errs.New(nil, "任务错误行为未正确设置")
	}
	d.ConfigDisableAction = r.ConfigDisableAction
	d.ConfigErrAction = r.ConfigErrAction

	err = data.NewCronPipelineData(dm.ctx).Set(d)
	if err != nil {
		return nil, err
	}
	return &pb.CronPipelineSetReply{
		Id: d.Id,
	}, err
}

// 任务状态变更
func (dm *CronPipelineService) ChangeStatus(r *pb.CronPipelineChangeStatusRequest) (resp *pb.CronPipelineChangeStatusReply, err error) {
	// 同一个任务，这里要加请求锁
	da := data.NewCronPipelineData(dm.ctx)
	conf, err := da.GetOne(dm.user.Env, r.Id)
	if err != nil {
		return nil, err
	}
	if conf.Status == r.Status {
		return nil, fmt.Errorf("状态相等")
	}

	if conf.Status == models.ConfigStatusActive && r.Status == models.ConfigStatusDisable { // 启用 到 停用 要关闭执行中的对应任务；
		NewTaskService(config.MainConf()).DelPipeline(conf)
		conf.EntryId = 0
	} else if conf.Status != models.ConfigStatusActive && r.Status == models.ConfigStatusActive { // 停用 到 启用 要把任务注册；
		if err = NewTaskService(config.MainConf()).AddPipeline(conf); err != nil {
			return nil, err
		}
	} else {
		return nil, fmt.Errorf("错误状态请求")
	}

	conf.Status = r.Status
	if err = da.ChangeStatus(conf, "视图操作"+models.ConfigStatusMap[r.Status]); err != nil {
		// 前面操作了任务，这里失败了；要将任务进行反向操作（回滚）（并附带两条对应日志）
		return nil, err
	}
	return &pb.CronPipelineChangeStatusReply{}, nil
}

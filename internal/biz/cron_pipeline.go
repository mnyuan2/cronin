package biz

import (
	"context"
	"cron/internal/basic/auth"
	"cron/internal/basic/config"
	"cron/internal/basic/conv"
	"cron/internal/basic/db"
	"cron/internal/basic/errs"
	"cron/internal/data"
	"cron/internal/models"
	"cron/internal/pb"
	"crypto/md5"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"log"
	"strings"
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
	w := db.NewWhere().
		Eq("type", r.Type).
		Eq("env", dm.user.Env, db.RequiredOption()).
		In("status", r.Status).
		In("create_user_id", r.CreateUserIds).
		FindInSet("handle_user_ids", r.HandleUserIds).
		Like("name", r.Name)
	if r.CreateOrHandleUserId > 0 {
		w.Raw("(create_user_id IN (?) OR FIND_IN_SET(?,handle_user_ids))", r.CreateOrHandleUserId, r.CreateOrHandleUserId)
	}
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
	topList := map[int]*data.SumStatus{}
	if len(resp.List) > 0 {
		endTime := time.Now()
		startTime := time.Now().Add(-time.Hour * 24 * 7) // 取七天前
		ids := make([]int, len(resp.List))
		for i, temp := range resp.List {
			ids[i] = temp.Id
		}
		w2 := db.NewWhere().
			Eq("env", dm.user.Env).
			Eq("operation", "job-pipeline").
			In("ref_id", ids).
			Between("timestamp", startTime.UnixMicro(), endTime.UnixMicro())
		topList, _ = data.NewCronLogSpanIndexV2Data(dm.ctx).SumStatus(w2)
	}

	for _, item := range resp.List {
		item.ConfigIds = []int{}
		item.Configs = []*pb.CronConfigListItem{}
		item.MsgSet = []*pb.CronMsgSet{}
		item.StatusName = models.ConfigStatusMap[item.Status]
		item.ConfigDisableActionName = models.DisableActionMap[item.ConfigDisableAction]
		if item.ConfigIdsStr != nil {
			if err = jsoniter.Unmarshal(item.ConfigIdsStr, &item.ConfigIds); err != nil {
				fmt.Println("	", err.Error())
			}
		}
		if item.ConfigsStr != nil {
			if err = jsoniter.Unmarshal(item.ConfigsStr, &item.Configs); err != nil {
				fmt.Println("	", err.Error())
			}
		}
		if item.MsgSetStr != nil {
			if err = jsoniter.Unmarshal(item.MsgSetStr, &item.MsgSet); err != nil {
				fmt.Println("	", err.Error())
			}
		}
		if item.HandleUserStr != nil {
			conv.NewStr().Slice(string(item.HandleUserStr), &item.HandleUserIds)
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
	g := data.NewChangeLogHandle(dm.user)
	d := &models.CronPipeline{}
	if r.Id > 0 {
		da := data.NewCronPipelineData(dm.ctx)
		d, err = da.GetOne(dm.user.Env, r.Id)
		if err != nil {
			return nil, err
		}
		if d.Status == models.ConfigStatusActive {
			return nil, fmt.Errorf("请先停用任务后编辑")
		}
		if d.Status == models.ConfigStatusClosed {
			return nil, fmt.Errorf("请先恢复任务后编辑")
		}
		g.SetType(models.LogTypeUpdateDiy).OldPipeline(*d)
	} else {
		g.SetType(models.LogTypeCreate).OldPipeline(*d)
		d.Status = models.ConfigStatusDisable
		d.Env = dm.user.Env
		d.CreateUserId = dm.user.UserId
		d.CreateUserName = dm.user.UserName
		d.StatusDt = time.Now().Format(time.DateTime)
		d.StatusRemark = "新增"
	}

	if r.VarParams != "" {
		prams := map[string]any{}
		if err = jsoniter.UnmarshalFromString(r.VarParams, &prams); err != nil {
			return nil, fmt.Errorf("变量参数序列化错误，%w", err)
		}
	}

	d.AuditUserId = 0
	d.AuditUserName = ""
	d.Name = r.Name
	d.Spec = r.Spec
	d.Interval = r.Interval
	d.Remark = r.Remark
	d.VarParams = r.VarParams
	d.ConfigIds, _ = jsoniter.Marshal(r.ConfigIds)
	d.Configs, _ = jsoniter.Marshal(r.Configs)
	d.MsgSet, _ = jsoniter.Marshal(r.MsgSet)
	d.MsgSetHash = fmt.Sprintf("%x", md5.Sum(d.MsgSet))
	d.Type = r.Type
	if d.Status != models.ConfigStatusDisable { // 编辑后，单子都是草稿
		d.Status = models.ConfigStatusDisable
		d.StatusRemark = "编辑"
		d.StatusDt = time.Now().Format(time.DateTime)
	}
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
		if len(msg.Status) == 0 {
			return nil, fmt.Errorf("推送%v未设置消息状态", i)
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
	err = data.NewCronChangeLogData(dm.ctx).Write(g.NewPipeline(*d))
	if err != nil {
		log.Println("变更日志写入错误", err.Error())
	}

	return &pb.CronPipelineSetReply{
		Id: d.Id,
	}, err
}

// 任务执行
func (dm *CronPipelineService) Run(r *pb.CronPipelineRunRequest) (resp *pb.CronPipelineRunReply, err error) {
	//conf := &models.CronPipeline{
	//	Id:        r.Id,
	//	Env:       dm.user.Env,
	//	Name:      r.Name,
	//	Type:      r.Type,
	//	Protocol:  r.Protocol,
	//	AfterTmpl: r.AfterTmpl,
	//}
	//conf.Command, err = jsoniter.Marshal(r.Command)
	//if err != nil {
	//	return nil, err
	//}
	//if r.MsgSet != nil {
	//	if conf.MsgSet, err = jsoniter.Marshal(r.MsgSet); err != nil {
	//		return nil, errs.New(err, "消息设置序列化错误")
	//	}
	//}
	//if r.VarFields != nil {
	//	if conf.VarFields, err = jsoniter.Marshal(r.VarFields); err != nil {
	//		return nil, errs.New(err, "字段设置序列化错误")
	//	}
	//}
	//res, err := NewJobPipeline(conf).Running(dm.ctx, "手动执行", map[string]any{})
	//if err != nil {
	return nil, err
	//}
	//return &pb.CronPipelineRunReply{Result: string(res)}, nil
}

// 流水线详情
func (dm *CronPipelineService) Detail(r *pb.CronPipelineDetailRequest) (resp *pb.CronPipelineDetailReply, err error) {
	if r.Id == 0 {
		return nil, errs.New(nil, "参数未传递")
	}

	one := &models.CronPipeline{}
	if r.Id < 0 {
		list := cronRun.Entries()
		for _, v := range list {
			c, ok := v.Job.(*JobPipeline)
			if !ok {
				continue
			}
			if c.conf.conf.Id == r.Id {
				one = c.pipeline
			}
		}
	} else {
		one, err = data.NewCronPipelineData(dm.ctx).GetOne(dm.user.Env, r.Id)
		if err != nil {
			return nil, err
		}
	}
	if one.Id == 0 {
		return nil, errs.New(nil, "未找到任务信息")
	}

	resp = &pb.CronPipelineDetailReply{
		ConfigIds:     []int{},
		Configs:       []*pb.CronConfigListItem{},
		MsgSet:        []*pb.CronMsgSet{},
		HandleUserIds: []int{},
	}
	err = conv.NewMapper().Exclude("ConfigIds", "Configs", "MsgSet", "HandleUserIds").Map(one, resp)
	if err != nil {
		return nil, errs.New(err, "系统错误")
	}

	resp.StatusName = models.ConfigStatusMap[resp.Status]
	resp.ConfigDisableActionName = models.DisableActionMap[resp.ConfigDisableAction]
	resp.ConfigErrActionName = models.ErrActionMap[resp.ConfigErrAction]
	if one.ConfigIds != nil {
		if err = jsoniter.Unmarshal(one.ConfigIds, &resp.ConfigIds); err != nil {
			fmt.Println("	", err.Error())
		}
	}
	if one.Configs != nil {
		if err = jsoniter.Unmarshal(one.Configs, &resp.Configs); err != nil {
			fmt.Println("	", err.Error())
		}
		// 加载最新任务信息
		confs, _ := NewCronConfigService(dm.ctx, dm.user).List(&pb.CronConfigListRequest{
			Ids:  resp.ConfigIds,
			Page: 1,
			Size: len(resp.Configs),
		})
		if confs != nil {
			for _, item := range confs.List {
				for _, c := range resp.Configs {
					if item.Id == c.Id {
						c.Name = item.Name
						c.TypeName = item.TypeName
						c.Status = item.Status
						c.StatusName = item.StatusName
						c.ProtocolName = item.ProtocolName
					}
				}
			}
		}
	}
	if resp.MsgSet != nil {
		if err = jsoniter.Unmarshal(one.MsgSet, &resp.MsgSet); err != nil {
			fmt.Println("	", err.Error())
		}
	}
	if one.HandleUserIds != "" {
		conv.NewStr().Slice(one.HandleUserIds, &resp.HandleUserIds)
	}

	return resp, err
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
	g := data.NewChangeLogHandle(dm.user).SetType(models.LogTypeUpdateDiy).OldPipeline(*conf)
	// 校验处理人
	if len(r.HandleUserIds) > 0 {
		users, err := data.NewCronUserData(dm.ctx).GetList(db.NewWhere().In("id", r.HandleUserIds))
		if err != nil {
			return nil, fmt.Errorf("审核人信息有误")
		}
		if len(users) != len(r.HandleUserIds) {
			return nil, fmt.Errorf("审核人信息有误！")
		}
		ids := make([]int, len(users))
		names := make([]string, len(users))
		for i, user := range users {
			ids[i] = user.Id
			names[i] = user.Username
		}
		conf.HandleUserIds, _ = conv.Int64s().Join(ids)
		conf.HandleUserNames = strings.Join(names, ",")
	} else {
		conf.HandleUserIds = ""
		conf.HandleUserNames = ""
	}

	switch r.Status {
	case models.ConfigStatusAudited: // 待审核
		if conf.Status != models.ConfigStatusDisable && conf.Status != models.ConfigStatusReject && conf.Status != models.ConfigStatusFinish && conf.Status != models.ConfigStatusError {
			return nil, fmt.Errorf("错误状态请求")
		}
		conf.AuditUserId = 0
		conf.AuditUserName = ""
	case models.ConfigStatusDisable: // 草稿、停用
		NewTaskService(config.MainConf()).DelPipeline(conf)
		conf.EntryId = 0

	case models.ConfigStatusActive: // 激活、通过
		if conf.Status != models.ConfigStatusDisable &&
			conf.Status != models.ConfigStatusAudited &&
			conf.Status != models.ConfigStatusReject &&
			conf.Status != models.ConfigStatusFinish &&
			conf.Status != models.ConfigStatusError {
			return nil, fmt.Errorf("不支持的状态变更操作")
		}
		conf.AuditUserId = dm.user.UserId
		conf.AuditUserName = dm.user.UserName
		if err = NewTaskService(config.MainConf()).AddPipeline(conf); err != nil {
			return nil, err
		}

	case models.ConfigStatusReject: // 驳回
		if conf.Status != models.ConfigStatusAudited {
			return nil, fmt.Errorf("不支持的状态变更操作")
		}
		conf.AuditUserId = dm.user.UserId
		conf.AuditUserName = dm.user.UserName

	case models.ConfigStatusClosed:
		if conf.Status != models.ConfigStatusDisable &&
			conf.Status != models.ConfigStatusReject &&
			conf.Status != models.ConfigStatusFinish &&
			conf.Status != models.ConfigStatusError {
			return nil, fmt.Errorf("不支持的状态变更操作")
		}
		conf.AuditUserId = dm.user.UserId
		conf.AuditUserName = dm.user.UserName
	default:
		return nil, fmt.Errorf("错误状态请求")
	}
	statusRemark := "视图操作" + models.ConfigStatusMap[r.Status]
	if r.StatusRemark != "" {
		statusRemark = r.StatusRemark
	}

	conf.Status = r.Status
	if err = da.ChangeStatus(conf, statusRemark); err != nil {
		// 前面操作了任务，这里失败了；要将任务进行反向操作（回滚）（并附带两条对应日志）
		return nil, err
	}
	err = data.NewCronChangeLogData(dm.ctx).Write(g.NewPipeline(*conf))
	if err != nil {
		log.Println("变更日志写入错误", err.Error())
	}
	return &pb.CronPipelineChangeStatusReply{}, nil
}

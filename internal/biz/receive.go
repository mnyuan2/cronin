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
	"crypto/md5"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"log"
	"strings"
	"time"
)

type ReceiveService struct {
	ctx  context.Context
	user *auth.UserToken
}

func NewReceiveService(ctx context.Context, user *auth.UserToken) *ReceiveService {
	return &ReceiveService{
		ctx:  ctx,
		user: user,
	}
}

// 列表
func (dm *ReceiveService) List(r *pb.ReceiveListRequest) (resp *pb.ReceiveListReply, err error) {
	w := db.NewWhere().
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
	resp = &pb.ReceiveListReply{
		List: []*pb.ReceiveListItem{},
		Page: &pb.Page{
			Page: r.Page,
			Size: r.Size,
		},
	}
	resp.Page.Total, err = data.NewCronReceiveData(dm.ctx).ListPage(w, r.Page, r.Size, &resp.List)
	topList := map[int]*data.SumConfTop{}
	if len(resp.List) > 0 {
		endTime := time.Now()
		startTime := time.Now().Add(-time.Hour * 24 * 7) // 取七天前
		ids := make([]int, len(resp.List))
		for i, temp := range resp.List {
			ids[i] = temp.Id
		}
		topList, _ = data.NewCronLogData(dm.ctx).SumConfTopError(dm.user.Env, ids, startTime, endTime, "receive")
	}

	for _, item := range resp.List {
		item.ConfigIds = []int{}
		item.StatusName = models.ConfigStatusMap[item.Status]
		if item.ConfigIdsStr != nil {
			if err = jsoniter.Unmarshal(item.ConfigIdsStr, &item.ConfigIds); err != nil {
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

// 设置
func (dm *ReceiveService) Set(r *pb.ReceiveSetRequest) (resp *pb.ReceiveSetReply, err error) {
	g := data.NewChangeLogHandle(dm.user)
	d := &models.CronReceive{}
	if r.Id > 0 {
		da := data.NewCronReceiveData(dm.ctx)
		d, err = da.GetOne(r.Id)
		if err != nil {
			return nil, err
		}
		if d.Status == enum.StatusActive {
			return nil, fmt.Errorf("请先停用任务后编辑")
		}
		g.SetType(models.LogTypeUpdateDiy).OldReceive(*d)
	} else {
		g.SetType(models.LogTypeCreate).OldReceive(*d)
		d.Status = enum.StatusDisable
		d.Env = dm.user.Env
	}

	if r.ReceiveTmpl == "" {
		return nil, fmt.Errorf("接收模板不得为空")
	}

	d.Name = r.Name
	d.Interval = r.Interval
	d.Remark = r.Remark
	d.ReceiveTmpl = r.ReceiveTmpl
	d.ConfigIds, _ = jsoniter.Marshal(r.ConfigIds)
	d.RuleConfig, _ = jsoniter.Marshal(r.RuleConfig)
	d.MsgSet, _ = jsoniter.Marshal(r.MsgSet)
	d.MsgSetHash = fmt.Sprintf("%x", md5.Sum(d.MsgSet))
	d.RuleConfigHash = fmt.Sprintf("%x", md5.Sum(d.RuleConfig))
	if d.Status != enum.StatusDisable { // 编辑后，单子都是草稿
		d.Status = enum.StatusDisable
		d.StatusRemark = "编辑"
		d.StatusDt = time.Now().Format(time.DateTime)
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

	err = data.NewCronReceiveData(dm.ctx).Set(d)
	if err != nil {
		return nil, err
	}
	err = data.NewCronChangeLogData(dm.ctx).Write(g.NewReceive(*d))
	if err != nil {
		log.Println("变更日志写入错误", err.Error())
	}

	return &pb.ReceiveSetReply{
		Id: d.Id,
	}, err
}

// 状态变更
func (dm *ReceiveService) ChangeStatus(r *pb.ReceiveChangeStatusRequest) (resp *pb.ReceiveChangeStatusReply, err error) {
	// 同一个任务，这里要加请求锁
	da := data.NewCronReceiveData(dm.ctx)
	conf, err := da.GetOne(r.Id)
	if err != nil {
		return nil, err
	}
	if conf.Status == r.Status {
		return nil, fmt.Errorf("状态相等")
	}
	g := data.NewChangeLogHandle(dm.user).SetType(models.LogTypeUpdateDiy).OldReceive(*conf)
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
	case models.ConfigStatusDisable: // 草稿、停用
		//NewTaskService(config.MainConf()).DelPipeline(conf)
		//conf.EntryId = 0

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
		//if err = NewTaskService(config.MainConf()).AddPipeline(conf); err != nil {
		//	return nil, err
		//}

	case models.ConfigStatusReject: // 驳回
		if conf.Status != models.ConfigStatusAudited {
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
	err = data.NewCronChangeLogData(dm.ctx).Write(g.NewReceive(*conf))
	if err != nil {
		log.Println("变更日志写入错误", err.Error())
	}
	return &pb.ReceiveChangeStatusReply{}, nil
}

// Detail 详情
func (dm *ReceiveService) Detail(r *pb.ReceiveDetailRequest) (resp *pb.ReceiveDetailReply, err error) {
	if r.Id == 0 {
		return nil, errs.New(nil, "参数未传递")
	}

	one := &models.CronReceive{}
	//if r.Id < 0 {
	//	list := cronRun.Entries()
	//	for _, v := range list {
	//		c, ok := v.Job.(*JobPipeline)
	//		if !ok {
	//			continue
	//		}
	//		if c.conf.conf.Id == r.Id {
	//			one = c.pipeline
	//		}
	//	}
	//} else {
	one, err = data.NewCronReceiveData(dm.ctx).GetOne(r.Id)
	if err != nil {
		return nil, err
	}
	//}
	if one.Id == 0 {
		return nil, errs.New(nil, "未找到任务信息")
	}

	resp = &pb.ReceiveDetailReply{}
	err = conv.NewMapper().Exclude("ConfigIds", "RuleConfig", "MsgSet", "HandleUserIds").Map(one, resp)
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
	if one.RuleConfig != nil {
		if err = jsoniter.Unmarshal(one.RuleConfig, &resp.RuleConfig); err != nil {
			fmt.Println("	", err.Error())
		}
		// 加载最新任务信息
		confs, _ := NewCronConfigService(dm.ctx, dm.user).List(&pb.CronConfigListRequest{
			Ids:  resp.ConfigIds,
			Page: 1,
			Size: len(resp.RuleConfig),
		})
		if confs != nil {
			for _, item := range confs.List {
				for _, c := range resp.RuleConfig {
					if item.Id == c.Config.Id {
						c.Config.Name = item.Name
						c.Config.TypeName = item.TypeName
						c.Config.Status = item.Status
						c.Config.StatusName = item.StatusName
						c.Config.ProtocolName = item.ProtocolName
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

// 接收钩子
func (dm *ReceiveService) Webhook(r *pb.ReceiveWebhookRequest) (resp *pb.ReceiveWebhookReply, err error) {
	/*
		目前的目标就是接收第三方消息，启动对应任务执行。
		但是每一个第三方的消息样式存在不同且无法定制，这个也是要解决的问题。
		简单方案是，针对每一个第三方一个独立接口。
			这样灵活性太低了，可以写一个预解析方案，由用户自己去处理。
	*/
	//if r.Id <= 0 {
	//	return nil, errors.New("用户不得为空")
	//}

	//_data := data.NewCronReceiveData(dm.ctx)
	//one, err := _data.GetOne(r.Id)
	//if err != nil {
	//	return nil, err
	//}

	return &pb.ReceiveWebhookReply{}, nil
}

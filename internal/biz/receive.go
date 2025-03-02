package biz

import (
	"bytes"
	"context"
	"cron/internal/basic/auth"
	"cron/internal/basic/config"
	"cron/internal/basic/conv"
	"cron/internal/basic/db"
	"cron/internal/basic/errs"
	"cron/internal/basic/tracing"
	"cron/internal/biz/dtos"
	"cron/internal/data"
	"cron/internal/models"
	"cron/internal/pb"
	"crypto/md5"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"log"
	"regexp"
	"strconv"
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
		w2 := db.NewWhere().
			Eq("env", dm.user.Env).
			Eq("operation", "job-receive").
			In("ref_id", ids).
			Between("timestamp", startTime.Format(time.DateTime), endTime.Format(time.DateTime))
		topList, _ = data.NewCronLogData(dm.ctx).SumConfTopError(w2)
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
	da := data.NewCronReceiveData(dm.ctx)
	if r.Id > 0 {
		d, err = da.GetOne(r.Id)
		if err != nil {
			return nil, err
		}
		if d.Status == models.ConfigStatusActive {
			return nil, fmt.Errorf("请先停用任务后编辑")
		}
		if d.Status == models.ConfigStatusClosed {
			return nil, fmt.Errorf("请先恢复任务后编辑")
		}
		g.SetType(models.LogTypeUpdateDiy).OldReceive(*d)
	} else {
		g.SetType(models.LogTypeCreate).OldReceive(*d)
		d.Status = models.ConfigStatusDisable
		d.Env = dm.user.Env
		d.CreateUserId = dm.user.UserId
		d.CreateUserName = dm.user.UserName
		d.StatusDt = time.Now().Format(time.DateTime)
		d.StatusRemark = "新增"
	}

	if r.ReceiveTmpl == "" {
		return nil, fmt.Errorf("接收模板不得为空")
	}
	if r.Alias != "" {
		if ok := regexp.MustCompile(`^[A-Za-z0-9\-_]+$`).MatchString(r.Alias); !ok {
			return nil, fmt.Errorf("别名为：字母、数字、-_ 的组合")
		}
		if _, err := strconv.Atoi(r.Alias); err == nil {
			return nil, fmt.Errorf("别名不能为纯数字")
		}
	}
	if r.MsgSet == nil {
		r.MsgSet = []*pb.CronMsgSet{}
	}

	d.AuditUserId = 0
	d.AuditUserName = ""
	d.Name = r.Name
	d.Alias = r.Alias
	d.Interval = r.Interval
	d.Remark = r.Remark
	d.ReceiveTmpl = r.ReceiveTmpl
	d.ConfigIds, _ = jsoniter.Marshal(r.ConfigIds)
	d.RuleConfig, _ = jsoniter.Marshal(r.RuleConfig)
	d.MsgSet, _ = jsoniter.Marshal(r.MsgSet)
	d.MsgSetHash = fmt.Sprintf("%x", md5.Sum(d.MsgSet))
	d.RuleConfigHash = fmt.Sprintf("%x", md5.Sum(d.RuleConfig))
	if d.Status != models.ConfigStatusDisable { // 编辑后，单子都是草稿
		d.Status = models.ConfigStatusDisable
		d.StatusRemark = "编辑"
		d.StatusDt = time.Now().Format(time.DateTime)
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

	err = da.Set(d)
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
		if conf.ConfigIds == nil || len(conf.ConfigIds) == 0 {
			return nil, fmt.Errorf("请至少指定一个匹配任务")
		}
		if conf.Alias != "" { // 唯一性校验
			list, _ := da.GetList(db.NewWhere().Eq("alias", conf.Alias).Neq("id", r.Id).Eq("status", models.ConfigStatusActive))
			if len(list) > 0 {
				return nil, fmt.Errorf("别名已被（%s.%s）使用，请更换", list[0].Env, list[0].Name)
			}
		}
		conf.AuditUserId = 0
		conf.AuditUserName = ""
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
		if conf.Alias != "" { // 唯一性校验
			list, _ := da.GetList(db.NewWhere().Eq("alias", conf.Alias).Neq("id", r.Id).Eq("status", models.ConfigStatusActive))
			if len(list) > 0 {
				return nil, fmt.Errorf("别名已被（%s.%s）使用，请更换", list[0].Env, list[0].Name)
			}
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

	resp = &pb.ReceiveDetailReply{
		MsgSet:        []*pb.CronMsgSet{},
		HandleUserIds: []int{},
	}
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
	fmt.Println(string(one.MsgSet), []byte("null"))
	if one.MsgSet != nil && string(one.MsgSet) != "null" {
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
	if r.Id <= 0 && r.Alias == "" {
		return nil, errs.New(nil, pb.ParamNotFound, "未指定 key")
	}

	list, err := data.NewCronReceiveData(dm.ctx).GetList(db.NewWhere().Eq("id", r.Id).Eq("alias", r.Alias))
	if err != nil {
		return nil, err
	}
	if len(list) == 0 || list[0].Status != models.ConfigStatusActive {
		return nil, errs.New(nil, pb.ParamNotFound, "请确认接收规则是否存在且激活")
	}
	one := list[0]

	// 这里不友好，下期优化
	tra := tracing.Tracer(one.Env+"-cronin", trace.WithInstrumentationAttributes(
		attribute.String("driver", "mysql"),
		attribute.String("env", one.Env),
	))
	_, s := tra.Start(dm.ctx, "receive/webhook")
	defer func() {
		if err != nil {
			s.SetStatus(tracing.StatusError, err.Error())
		} else {
			s.SetStatus(tracing.StatusOk, "")
		}
		if resp != nil {
			s.AddEvent("response", trace.WithAttributes(attribute.String("message", resp.Message)))
		}
		s.End()
	}()
	s.SetAttributes(
		attribute.Int("ref_id", one.Id),
		attribute.String("component", "HTTP"),
		attribute.String("method", "POST"),
		attribute.String("ref_type", "receive"),
	)
	s.AddEvent("request", trace.WithAttributes(attribute.String("request_body", string(r.Body))))

	// 解析接收
	b, err := conv.DefaultStringTemplate().SetParam(map[string]any{
		"request_body": string(r.Body),
	}).Execute([]byte(one.ReceiveTmpl))
	if err != nil {
		return nil, errs.New(err, pb.OperationFailure, "接收模板解析失败")
	}
	b = bytes.TrimSpace(b)
	resp = &pb.ReceiveWebhookReply{}
	if b == nil {
		resp.Message = "解析空忽略"
		return resp, nil
	}
	s.AddEvent("process", trace.WithAttributes(attribute.String("tmpl_result", string(b))))

	param := &dtos.ReceiveWebHook{}
	if err = jsoniter.Unmarshal(b, param); err != nil {
		return nil, errs.New(err, pb.OperationFailure, "接收模板解析结果错误")
	}

	if len(param.Dataset) > 0 {
		param.TraceId = tracing.Inject(s)
		id, err := NewTaskService(config.MainConf()).AddReceive(one, param)
		if err != nil {
			return nil, err
		}
		resp.Message = "task_id=" + strconv.Itoa(int(id))
		s.AddEvent("process", trace.WithAttributes(attribute.Int("entry_id", int(id))))
	}

	return resp, nil
}

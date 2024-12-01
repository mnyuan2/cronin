package biz

import (
	"context"
	"cron/internal/basic/auth"
	"cron/internal/basic/config"
	"cron/internal/basic/conv"
	"cron/internal/basic/db"
	"cron/internal/basic/enum"
	"cron/internal/basic/errs"
	"cron/internal/biz/dtos"
	"cron/internal/data"
	"cron/internal/models"
	"cron/internal/pb"
	"crypto/md5"
	"encoding/base64"
	"errors"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"log"
	"strings"
	"time"
)

type CronConfigService struct {
	ctx  context.Context
	user *auth.UserToken
}

func NewCronConfigService(ctx context.Context, user *auth.UserToken) *CronConfigService {
	return &CronConfigService{
		ctx:  ctx,
		user: user,
	}
}

// 任务配置列表
func (dm *CronConfigService) List(r *pb.CronConfigListRequest) (resp *pb.CronConfigListReply, err error) {
	w := db.NewWhere().
		Eq("type", r.Type).
		Eq("env", dm.user.Env, db.RequiredOption()).
		In("id", r.Ids).
		In("protocol", r.Protocol).
		In("status", r.Status).
		In("create_user_id", r.CreateUserIds).
		FindInSet("handle_user_ids", r.HandleUserIds).
		Like("name", r.Name)
	if r.CreateOrHandleUserId > 0 {
		w.Raw("(create_user_id IN (?) OR FIND_IN_SET(?,handle_user_ids))", r.CreateOrHandleUserId, r.CreateOrHandleUserId)
	}
	if w.Len() == 0 {
		return nil, errs.New(nil, "未指定查询条件")
	}
	// 构建查询条件
	if r.Page <= 1 {
		r.Page = 1
	}
	if r.Size <= 10 {
		r.Size = 10
	}
	resp = &pb.CronConfigListReply{
		List: []*pb.CronConfigListItem{},
		Page: &pb.Page{
			Page: r.Page,
			Size: r.Size,
		},
	}
	resp.Page.Total, err = data.NewCronConfigData(dm.ctx).ListPage(w, r.Page, r.Size, &resp.List)
	topList := map[int]*data.SumConfTop{}
	if len(resp.List) > 0 {
		endTime := time.Now()
		startTime := time.Now().Add(-time.Hour * 24 * 7) // 取七天前
		ids := make([]int, len(resp.List))
		for i, temp := range resp.List {
			ids[i] = temp.Id
		}
		topList, _ = data.NewCronLogData(dm.ctx).SumConfTopError(dm.user.Env, ids, startTime, endTime, "config")
	}

	dicUser, err := NewDicService(dm.ctx, dm.user).getDb(enum.DicUser)
	if err != nil {
		return nil, err
	}
	dicUserMap := dtos.DicToMap(dicUser)
	for _, item := range resp.List {
		item.TypeName = models.ConfigTypeMap[item.Type]
		item.StatusName = models.ConfigStatusMap[item.Status]
		item.ProtocolName = models.ProtocolMap[item.Protocol]
		item.CreateUserName = dicUserMap[item.CreateUserId]
		item.HandleUserIds = []int{}
		if item.VarFieldsStr != nil {
			_ = jsoniter.Unmarshal(item.VarFieldsStr, &item.VarFields)
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

// 任务详情
func (dm *CronConfigService) Detail(r *pb.CronConfigDetailRequest) (resp *pb.CronConfigDetailReply, err error) {
	if r.Id == 0 {
		return nil, errs.New(nil, "参数未传递")
	}

	one := &models.CronConfig{}
	if r.Id < 0 { // 兼容系统任务
		list := cronRun.Entries()
		for _, v := range list {
			c, ok := v.Job.(*JobConfig)
			if !ok {
				continue
			}
			if c.conf.Id == r.Id {
				one = c.conf
			}
		}
	} else {
		one, err = data.NewCronConfigData(dm.ctx).GetOne(dm.user.Env, r.Id)
	}
	if one.Id == 0 {
		return nil, errs.New(nil, "未找到任务信息")
	}

	resp = &pb.CronConfigDetailReply{
		VarFields: []*pb.KvItem{},
		Command: &pb.CronConfigCommand{
			Http:    &pb.CronHttp{Header: []*pb.KvItem{}},
			Rpc:     &pb.CronRpc{Actions: []string{}},
			Cmd:     &pb.CronCmd{Statement: &pb.CronStatement{Git: &pb.Git{Path: []string{}}}, Host: &pb.SettingHostSource{}},
			Sql:     &pb.CronSql{Statement: []*pb.CronStatement{}, Source: &pb.CronSqlSource{}},
			Jenkins: &pb.CronJenkins{Source: &pb.CronJenkinsSource{}, Params: []*pb.KvItem{}},
			Git:     &pb.CronGit{Events: []*pb.GitEvent{}},
		},
		MsgSet:        []*pb.CronMsgSet{},
		TypeName:      models.ConfigTypeMap[one.Type],
		StatusName:    models.ConfigStatusMap[one.Status],
		ProtocolName:  models.ProtocolMap[one.Protocol],
		HandleUserIds: []int{},
	}
	err = conv.NewMapper().Exclude("VarFields", "Command", "MsgSet", "HandleUserIds").Map(one, resp)
	if err != nil {
		return nil, errs.New(err, "系统错误")
	}

	if one.VarFields != nil {
		if er := jsoniter.Unmarshal(one.VarFields, &resp.VarFields); er != nil {
			return nil, errs.New(er, "var_fields 解析错误")
		}
	}
	// 流水线预览任务时，需要将参数传递进来解析看最终效果。
	if r.VarParams != "" {
		params := map[string]any{}
		if err := jsoniter.UnmarshalFromString(r.VarParams, &params); err != nil {
			return nil, errs.New(err, "参数实现错误")
		}
		newParams := map[string]any{}
		for _, item := range resp.VarFields {
			if val, ok := params[item.Key]; ok {
				newParams[item.Key] = val
				item.Value = fmt.Sprintf("%v", val)
			} else {
				newParams[item.Key] = item.Value
			}
		}
		cmd, err := conv.DefaultStringTemplate().SetParam(newParams).Execute(one.Command)
		if err != nil {
			return nil, errs.New(err, "模板错误.")
		}
		one.Command = cmd
	}

	if er := jsoniter.Unmarshal(one.Command, resp.Command); er != nil {
		return nil, errs.New(er, "command 解析错误")
	}
	for _, item := range resp.Command.Git.Events {
		if item.FileUpdate != nil {
			if item.FileUpdate.Content != "" {
				content, _ := base64.StdEncoding.DecodeString(item.FileUpdate.Content)
				item.FileUpdate.Content = string(content)
			}
		}
	}
	if one.MsgSet != nil {
		if er := jsoniter.Unmarshal(one.MsgSet, &resp.MsgSet); er != nil {
			return nil, errs.New(er, "msg_set 解析错误")
		}
	}
	if one.HandleUserIds != "" {
		conv.NewStr().Slice(one.HandleUserIds, &resp.HandleUserIds)
	}

	return resp, err
}

// 已注册任务列表
func (dm *CronConfigService) RegisterList(r *pb.CronConfigRegisterListRequest) (resp *pb.CronConfigRegisterListResponse, err error) {

	list := cronRun.Entries()
	resp = &pb.CronConfigRegisterListResponse{List: []*pb.CronConfigListItem{}}
	for _, v := range list {
		c, ok := v.Job.(*JobConfig)
		if !ok {
			c2, ok := v.Job.(*JobPipeline)
			if !ok {
				resp.List = append(resp.List, &pb.CronConfigListItem{
					Id: 0,
					//EntryId:  int(v.ID),
					Name:     "未识别注册任务",
					UpdateDt: v.Next.Format(time.DateTime),
				})
				continue
			}
			c = c2.GetConf()
		}
		if c.conf.Id > 0 && dm.user.Env != c.conf.Env {
			continue
		}
		param, _ := c.ParseParams(nil)
		c.Parse(param)
		conf := c.conf
		next := ""
		if s, err := secondParser.Parse(conf.Spec); err == nil {
			next = s.Next(time.Now()).Format(conv.FORMAT_DATETIME)
		}
		resp.List = append(resp.List, &pb.CronConfigListItem{
			Id: conf.Id,
			//EntryId:      int(v.ID),
			Name:         conf.Name,
			Spec:         conf.Spec,
			Protocol:     conf.Protocol,
			ProtocolName: conf.GetProtocolName(),
			Remark:       conf.Remark,
			Status:       conf.Status,
			StatusName:   conf.GetStatusName(),
			UpdateDt:     next, // 下一次时间
		})
	}

	return resp, err
}

// 任务配置
func (dm *CronConfigService) Set(r *pb.CronConfigSetRequest) (resp *pb.CronConfigSetResponse, err error) {
	g := data.NewChangeLogHandle(dm.user)
	d := &models.CronConfig{}
	if r.Id > 0 {
		da := data.NewCronConfigData(dm.ctx)
		d, err = da.GetOne(dm.user.Env, r.Id)
		if err != nil {
			return nil, err
		}
		if d.Status == models.ConfigStatusActive {
			return nil, fmt.Errorf("请先停用任务后编辑")
		}
		g.SetType(models.LogTypeUpdateDiy).OldConfig(*d)
	} else {
		g.SetType(models.LogTypeCreate).OldConfig(*d)
		d.Env = dm.user.Env
		d.CreateUserId = dm.user.UserId
		d.CreateUserName = dm.user.UserName
	}

	if r.Type == models.TypeCycle {
		if _, err = secondParser.Parse(r.Spec); err != nil {
			return nil, fmt.Errorf("时间格式不规范，%s", err.Error())
		}
	} else if r.Type == models.TypeOnce {
		if r.Spec != "" {
			_, err = time.ParseInLocation(time.DateTime, r.Spec, time.Local)
			if err != nil {
				return nil, errs.New(err, "执行时间格式不规范")
			}
		}
	} else if r.Type == models.TypeModule {
		//
	} else {
		return nil, fmt.Errorf("类型输入有误")
	}

	if r.Protocol == models.ProtocolHttp {
		if err := dtos.CheckHttp(r.Command.Http); err != nil {
			return nil, err
		}
	} else if r.Protocol == models.ProtocolRpc {
		if err := dtos.CheckRPC(r.Command.Rpc); err != nil {
			return nil, err
		}

	} else if r.Protocol == models.ProtocolCmd {
		if err := dtos.CheckCmd(r.Command.Cmd); err != nil {
			return nil, err
		}
		if r.Command.Cmd.Host.Id > 0 {
			if one, _ := data.NewCronSettingData(dm.ctx).GetSourceOne(dm.user.Env, r.Command.Cmd.Host.Id); one.Id == 0 || one.Scene != models.SceneHostSource {
				return nil, errors.New("sql 连接 配置有误，请确认")
			}
		}
	} else if r.Protocol == models.ProtocolSql {
		if err := dtos.CheckSql(r.Command.Sql); err != nil {
			return nil, err
		}
		if one, _ := data.NewCronSettingData(dm.ctx).GetSourceOne(dm.user.Env, r.Command.Sql.Source.Id); one.Id == 0 || one.Scene != models.SceneSqlSource {
			return nil, errors.New("sql 连接 配置有误，请确认")
		} else if set, err := dtos.ParseSource(one); err != nil || set.Sql.Driver != r.Command.Sql.Driver {
			return nil, errors.New("sql 连接 驱动有误，请确认")
		}
	} else if r.Protocol == models.ProtocolJenkins {
		if err := dtos.CheckJenkins(r.Command.Jenkins); err != nil {
			return nil, err
		}
		if one, _ := data.NewCronSettingData(dm.ctx).GetSourceOne(dm.user.Env, r.Command.Jenkins.Source.Id); one.Id == 0 || one.Scene != models.SceneJenkinsSource {
			return nil, errors.New("jenkins 连接 配置有误，请确认")
		}
	} else if r.Protocol == models.ProtocolGit {
		if err := dtos.CheckGit(r.Command.Git); err != nil {
			return nil, err
		}
		if one, _ := data.NewCronSettingData(dm.ctx).GetSourceOne(dm.user.Env, r.Command.Git.LinkId); one.Id == 0 || one.Scene != models.SceneGitSource {
			return nil, errors.New("git 连接 配置有误，请确认")
		}
	}
	if r.ErrRetryNum > 30 {
		return nil, errors.New("最大重试次数不得超过30")
	} else if r.ErrRetryNum < 0 {
		return nil, errors.New("最大重试次数不得小于0")
	}
	pl := len(r.VarFields)
	for i, param := range r.VarFields {
		if param.Key == "" && i < (pl-1) {
			return nil, fmt.Errorf("变量参数 %v 名称不得为空", i+1)
		} else if strings.Contains(param.Key, ".") {
			return nil, fmt.Errorf("参数key %v 不得包含.点符合", i+1)
		}
	}
	for i, msg := range r.MsgSet {
		if msg.MsgId == 0 {
			return nil, fmt.Errorf("推送%v未设置消息模板", i)
		}
	}

	if d.Status != models.ConfigStatusDisable { // 编辑后，单子都是草稿
		d.Status = models.ConfigStatusDisable
		d.StatusRemark = "编辑"
		d.StatusDt = time.Now().Format(time.DateTime)
	}
	d.AuditUserId = 0
	d.AuditUserName = ""
	d.Name = r.Name
	d.Spec = r.Spec
	d.Protocol = r.Protocol
	d.Remark = r.Remark
	d.Type = r.Type
	d.AfterTmpl = r.AfterTmpl
	d.AfterSleep = r.AfterSleep
	d.ErrRetryNum = r.ErrRetryNum
	d.VarFields, _ = jsoniter.Marshal(r.VarFields)
	d.Command, _ = jsoniter.Marshal(r.Command)
	d.MsgSet, _ = jsoniter.Marshal(r.MsgSet)
	d.VarFieldsHash = fmt.Sprintf("%x", md5.Sum(d.VarFields))
	d.CommandHash = fmt.Sprintf("%x", md5.Sum(d.Command))
	d.MsgSetHash = fmt.Sprintf("%x", md5.Sum(d.MsgSet))
	err = data.NewCronConfigData(dm.ctx).Set(d)
	if err != nil {
		return nil, err
	}

	err = data.NewCronChangeLogData(dm.ctx).Write(g.NewConfig(*d))
	if err != nil {
		log.Println("变更日志写入错误", err.Error())
	}

	return &pb.CronConfigSetResponse{
		Id: d.Id,
	}, nil
}

// 任务状态变更
func (dm *CronConfigService) ChangeStatus(r *pb.CronConfigSetRequest) (resp *pb.CronConfigSetResponse, err error) {
	// 同一个任务，这里要加请求锁
	da := data.NewCronConfigData(dm.ctx)
	conf, err := da.GetOne(dm.user.Env, r.Id)
	if err != nil {
		return nil, err
	}
	if conf.Status == r.Status {
		return nil, fmt.Errorf("状态相等")
	}
	g := data.NewChangeLogHandle(dm.user).SetType(models.LogTypeUpdateDiy).OldConfig(*conf)
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
		if r.Type == models.TypeOnce {
			if _, err = NewScheduleOnce(r.Spec); err != nil {
				return nil, err
			}
		}
		conf.AuditUserId = 0
		conf.AuditUserName = ""
	case models.ConfigStatusDisable: // 草稿、停用
		NewTaskService(config.MainConf()).Del(conf)
		conf.EntryId = 0
	case models.ConfigStatusActive: // 激活、通过
		if conf.Status != models.ConfigStatusDisable &&
			conf.Status != models.ConfigStatusAudited &&
			conf.Status != models.ConfigStatusReject &&
			conf.Status != models.ConfigStatusFinish &&
			conf.Status != models.ConfigStatusError {
			return nil, fmt.Errorf("不支持的状态变更操作")
		}
		if conf.Type != models.TypeModule {
			if conf.Status != models.ConfigStatusActive { // 停用 到 启用 要把任务注册；
				if err = NewTaskService(config.MainConf()).AddConfig(conf); err != nil {
					return nil, err
				}
			}
		}
		conf.AuditUserId = dm.user.UserId
		conf.AuditUserName = dm.user.UserName
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
	err = data.NewCronChangeLogData(dm.ctx).Write(g.NewConfig(*conf))
	if err != nil {
		log.Println("变更日志写入错误", err.Error())
	}
	return &pb.CronConfigSetResponse{
		Id: conf.Id,
	}, nil
}

func (dm *CronConfigService) Del() {

}

// 任务执行
func (dm *CronConfigService) Run(r *pb.CronConfigRunRequest) (resp *pb.CronConfigRunResponse, err error) {
	for _, item := range r.Command.Git.Events {
		if item.FileUpdate == nil || item.FileUpdate.Content == "" {
			continue
		}
		item.FileUpdate.Content = base64.StdEncoding.EncodeToString([]byte(item.FileUpdate.Content))
	}

	conf := &models.CronConfig{
		Id:         r.Id,
		Env:        dm.user.Env,
		Name:       r.Name,
		Type:       r.Type,
		Protocol:   r.Protocol,
		AfterTmpl:  r.AfterTmpl,
		AfterSleep: r.AfterSleep,
	}
	conf.Command, err = jsoniter.Marshal(r.Command)
	if err != nil {
		return nil, err
	}

	if r.MsgSet != nil {
		if conf.MsgSet, err = jsoniter.Marshal(r.MsgSet); err != nil {
			return nil, errs.New(err, "消息设置序列化错误")
		}
	}
	if r.VarFields != nil {
		if conf.VarFields, err = jsoniter.Marshal(r.VarFields); err != nil {
			return nil, errs.New(err, "字段设置序列化错误")
		}
	}
	res, err := NewJobConfig(conf).Running(dm.ctx, "手动执行", map[string]any{})
	if err != nil {
		return nil, err
	}
	return &pb.CronConfigRunResponse{Result: string(res)}, nil
}

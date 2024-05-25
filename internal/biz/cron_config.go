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
		Eq("status", r.Status)
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
		topList, _ = data.NewCronLogData(dm.ctx).SumConfTopError(dm.user.Env, ids, startTime, endTime, 7)
	}

	dicUser, err := NewDicService(dm.ctx, dm.user).getDb(enum.DicUser)
	if err != nil {
		return nil, err
	}
	dicUserMap := dtos.DicToMap(dicUser)
	for _, item := range resp.List {
		item.Command = &pb.CronConfigCommand{
			Http:    &pb.CronHttp{Header: []*pb.KvItem{}},
			Rpc:     &pb.CronRpc{Actions: []string{}},
			Cmd:     &pb.CronCmd{Statement: &pb.CronStatement{Git: &pb.Git{Path: []string{}}}, Host: &pb.SettingHostSource{}},
			Sql:     &pb.CronSql{Statement: []*pb.CronStatement{}, Source: &pb.CronSqlSource{}},
			Jenkins: &pb.CronJenkins{Source: &pb.CronJenkinsSource{}, Params: []*pb.KvItem{}},
			Git:     &pb.CronGit{Events: []*pb.GitEvent{}},
		}
		item.MsgSet = []*pb.CronMsgSet{}
		item.TypeName = models.ConfigTypeMap[item.Type]
		item.StatusName = models.ConfigStatusMap[item.Status]
		item.ProtocolName = models.ProtocolMap[item.Protocol]
		item.CreateUserName = dicUserMap[item.CreateUserId]
		item.StatusUserName = dicUserMap[item.StatusUserId]
		if er := jsoniter.Unmarshal(item.CommandStr, item.Command); er != nil {
			log.Println("	command 解析错误", item.Id, er.Error())
		}
		if item.MsgSetStr != nil {
			if er := jsoniter.Unmarshal(item.MsgSetStr, &item.MsgSet); er != nil {
				log.Println("	msg_set 解析错误", item.Id, er.Error())
			}
		}
		if item.VarFieldsStr != nil {
			_ = jsoniter.Unmarshal(item.VarFieldsStr, &item.VarFields)
		}
		if top, ok := topList[item.Id]; ok {
			item.TopNumber = top.TotalNumber
			item.TopErrorNumber = top.ErrorNumber
		}
	}

	return resp, err
}

func (dm *CronConfigService) Get() {

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
					Id:       0,
					EntryId:  int(v.ID),
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
			Id:           conf.Id,
			EntryId:      int(v.ID),
			Name:         conf.Name,
			Spec:         conf.Spec,
			Protocol:     conf.Protocol,
			ProtocolName: conf.GetProtocolName(),
			Remark:       conf.Remark,
			Status:       conf.Status,
			StatusName:   conf.GetStatusName(),
			UpdateDt:     next, // 下一次时间
			Command:      c.commandParse,
			MsgSet:       c.msgSetParse.Set,
		})
	}

	return resp, err
}

// 任务配置
func (dm *CronConfigService) Set(r *pb.CronConfigSetRequest) (resp *pb.CronConfigSetResponse, err error) {

	d := &models.CronConfig{}
	if r.Id > 0 {
		da := data.NewCronConfigData(dm.ctx)
		d, err = da.GetOne(dm.user.Env, r.Id)
		if err != nil {
			return nil, err
		}
		if d.Status == enum.StatusActive {
			return nil, fmt.Errorf("请先停用任务后编辑")
		}
	} else {
		d.Env = dm.user.Env
		d.CreateUserId = dm.user.UserId
	}

	if r.Type == models.TypeCycle {
		if _, err = secondParser.Parse(r.Spec); err != nil {
			return nil, fmt.Errorf("时间格式不规范，%s", err.Error())
		}
	} else if r.Type == models.TypeOnce {
		if _, err = NewScheduleOnce(r.Spec); err != nil {
			return nil, err
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

	if d.Status != enum.StatusDisable { // 编辑后，单子都是草稿
		d.Status = enum.StatusDisable
		d.StatusRemark = "编辑"
		d.StatusDt = time.Now().Format(time.DateTime)
	}
	d.Name = r.Name
	d.Spec = r.Spec
	d.Protocol = r.Protocol
	d.Remark = r.Remark
	d.Type = r.Type
	d.StatusUserId = dm.user.UserId
	d.AfterTmpl = r.AfterTmpl
	d.VarFields, _ = jsoniter.Marshal(r.VarFields)
	d.Command, _ = jsoniter.Marshal(r.Command)
	d.MsgSet, _ = jsoniter.Marshal(r.MsgSet)
	err = data.NewCronConfigData(dm.ctx).Set(d)
	if err != nil {
		return nil, err
	}
	return &pb.CronConfigSetResponse{
		Id: d.Id,
	}, err
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
	conf.StatusUserId = dm.user.UserId
	switch r.Status {
	case models.ConfigStatusAudited: // 待审核
		if conf.Status != models.ConfigStatusDisable && conf.Status != models.ConfigStatusReject {
			return nil, fmt.Errorf("错误状态请求")
		}
		conf.StatusUserId = r.AuditorUserId // 待审核可以指定操作人
	case models.ConfigStatusDisable: // 草稿、停用
		if conf.Type != models.TypeModule {
			if conf.Status == models.ConfigStatusActive { // 启用 到 停用 要关闭执行中的对应任务；
				NewTaskService(config.MainConf()).Del(conf)
				conf.EntryId = 0
			}
		}
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
	case models.ConfigStatusReject: // 驳回
		if conf.Status != models.ConfigStatusAudited {
			return nil, fmt.Errorf("不支持的状态变更操作")
		}
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
	return &pb.CronConfigSetResponse{
		Id: conf.Id,
	}, nil
}

func (dm *CronConfigService) Del() {

}

// 任务执行
func (dm *CronConfigService) Run(r *pb.CronConfigRunRequest) (resp *pb.CronConfigRunResponse, err error) {
	conf := &models.CronConfig{
		Id:        r.Id,
		Env:       dm.user.Env,
		Name:      r.Name,
		Type:      r.Type,
		Protocol:  r.Protocol,
		AfterTmpl: r.AfterTmpl,
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

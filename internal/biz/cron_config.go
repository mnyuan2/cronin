package biz

import (
	"context"
	"cron/internal/basic/auth"
	"cron/internal/basic/config"
	"cron/internal/basic/conv"
	"cron/internal/basic/db"
	"cron/internal/basic/enum"
	"cron/internal/biz/dtos"
	"cron/internal/data"
	"cron/internal/models"
	"cron/internal/pb"
	"errors"
	"fmt"
	jsoniter "github.com/json-iterator/go"
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
	if r.Type == 0 {
		r.Type = models.TypeCycle
	}
	w := db.NewWhere().Eq("type", r.Type).Eq("env", dm.user.Env, db.RequiredOption())
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
	resp.Page.Total, err = data.NewCronConfigData(dm.ctx).GetList(w, r.Page, r.Size, &resp.List)
	topList := map[int]*data.SumConfTop{}
	if len(resp.List) > 0 {
		endTime := time.Now()
		startTime := time.Now().Add(-time.Hour * 24 * 7) // 取七天前
		ids := make([]int, len(resp.List))
		for i, temp := range resp.List {
			ids[i] = temp.Id
		}
		topList, _ = data.NewCronLogData(dm.ctx).SumConfTopError(dm.user.Env, ids, startTime, endTime, 5)
	}

	for _, item := range resp.List {
		item.Command = &pb.CronConfigCommand{Http: &pb.CronHttp{Header: []*pb.KvItem{}}, Rpc: &pb.CronRpc{Actions: []string{}}, Sql: &pb.CronSql{}}
		item.StatusName = models.ConfigStatusMap[item.Status]
		item.ProtocolName = models.ProtocolMap[item.Protocol]
		jsoniter.UnmarshalFromString(item.CommandStr, item.Command)
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
		c, ok := v.Job.(*CronJob)
		if !ok {
			resp.List = append(resp.List, &pb.CronConfigListItem{
				Id:       int(v.ID),
				Name:     "未识别注册任务",
				UpdateDt: v.Next.Format(time.DateTime),
			})
			continue
		}
		if c.conf.Id > 0 && dm.user.Env != c.conf.Env {
			continue
		}
		conf := c.conf
		next := ""
		if s, err := secondParser.Parse(conf.Spec); err == nil {
			next = s.Next(time.Now()).Format(conv.FORMAT_DATETIME)
		}
		resp.List = append(resp.List, &pb.CronConfigListItem{
			Id:           conf.Id,
			Name:         conf.Name,
			Spec:         conf.Spec,
			Protocol:     conf.Protocol,
			ProtocolName: conf.GetProtocolName(),
			Remark:       conf.Remark,
			Status:       conf.Status,
			StatusName:   conf.GetStatusName(),
			UpdateDt:     next, // 下一次时间
			Command:      c.commandParse,
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
		d.Status = enum.StatusDisable
		d.Env = dm.user.Env
	}

	d.Name = r.Name
	d.Spec = r.Spec
	d.Protocol = r.Protocol
	d.Remark = r.Remark
	d.Command, _ = jsoniter.Marshal(r.Command)
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

	if r.Protocol == models.ProtocolHttp {
		if err := dtos.CheckHttp(r.Command.Http); err != nil {
			return nil, err
		}
	} else if r.Protocol == models.ProtocolRpc {
		if err := dtos.CheckRPC(r.Command.Rpc); err != nil {
			return nil, err
		}

	} else if r.Protocol == models.ProtocolCmd {
		if r.Command.Cmd == "" {
			return nil, fmt.Errorf("请输入 cmd 命令类容")
		}
	} else if r.Protocol == models.ProtocolSql {
		if err := dtos.CheckSql(r.Command.Sql); err != nil {
			return nil, err
		}
		if one, _ := data.NewCronSettingData(dm.ctx).GetSqlSourceOne(dm.user.Env, r.Command.Sql.Source.Id); one.Id == 0 {
			return nil, errors.New("sql 连接 配置有误，请确认")
		}
	}

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
	if conf.Status == models.ConfigStatusActive && r.Status == models.ConfigStatusDisable { // 启用 到 停用 要关闭执行中的对应任务；
		NewTaskService(config.MainConf()).Del(conf)
		conf.EntryId = 0
	} else if conf.Status != models.ConfigStatusActive && r.Status == models.ConfigStatusActive { // 停用 到 启用 要把任务注册；
		if err = NewTaskService(config.MainConf()).Add(conf); err != nil {
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
	return &pb.CronConfigSetResponse{
		Id: conf.Id,
	}, nil
}

func (dm *CronConfigService) Del() {

}

// 任务执行
func (dm *CronConfigService) Run(r *pb.CronConfigRunRequest) (resp *pb.CronConfigRunResponse, err error) {
	conf := &models.CronConfig{
		Type:     r.Type,
		Protocol: r.Protocol,
	}
	conf.Command, err = jsoniter.Marshal(r.Command)
	if err != nil {
		return nil, err
	}

	res, err := NewCronJob(conf).Exec(dm.ctx)
	if err != nil {
		return nil, err
	}
	return &pb.CronConfigRunResponse{Result: string(res)}, nil
}

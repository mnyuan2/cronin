package biz

import (
	"context"
	"cron/internal/basic/auth"
	"cron/internal/basic/config"
	"cron/internal/basic/conv"
	"cron/internal/basic/db"
	"cron/internal/basic/enum"
	"cron/internal/basic/errs"
	"cron/internal/basic/git"
	"cron/internal/basic/host"
	"cron/internal/data"
	"cron/internal/models"
	"cron/internal/pb"
	"crypto/tls"
	"errors"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"io"
	"net/http"
	"strings"
)

type SettingSqlService struct {
	db   *db.MyDB
	ctx  context.Context
	user *auth.UserToken
}

func NewSettingSqlService(ctx context.Context, user *auth.UserToken) *SettingSqlService {
	return &SettingSqlService{
		ctx:  ctx,
		user: user,
	}
}

// 任务配置列表
func (dm *SettingSqlService) List(r *pb.SettingListRequest) (resp *pb.SettingListReply, err error) {
	// 构建查询条件
	if r.Page <= 1 {
		r.Page = 1
	}
	if r.Size <= 10 {
		r.Size = 20
	}
	resp = &pb.SettingListReply{
		Page: &pb.Page{
			Page: r.Page,
			Size: r.Size,
		},
	}
	source, ok := models.DicToSource[r.Type]
	if !ok {
		return nil, errs.New(nil, "type参数错误")
	}

	list := []*models.CronSetting{}
	resp.Page.Total, err = data.NewCronSettingData(dm.ctx).GetList(source, dm.user.Env, r.Page, r.Size, &list)
	if err != nil {
		return nil, err
	}

	// 格式化
	resp.List = make([]*pb.SettingListItem, len(list))
	for i, item := range list {
		data := &pb.SettingListItem{
			Id:       item.Id,
			Title:    item.Title,
			CreateDt: item.CreateDt,
			UpdateDt: item.UpdateDt,
			Type:     models.SourceToDic[item.Scene],
			Source: &pb.SettingSource{
				Sql:     &pb.SettingSqlSource{},
				Jenkins: &pb.SettingJenkinsSource{},
				Git:     &pb.SettingGitSource{},
				Host:    &pb.SettingHostSource{},
			},
		}
		jsoniter.UnmarshalFromString(item.Content, data.Source)
		resp.List[i] = data
	}

	return resp, err
}

// 设置源
func (dm *SettingSqlService) Set(r *pb.SettingSqlSetRequest) (resp *pb.SettingSqlSetReply, err error) {
	source, ok := models.DicToSource[r.Type]
	if !ok {
		return nil, errs.New(nil, "type参数错误")
	}

	one := &models.CronSetting{}
	_data := data.NewCronSettingData(dm.ctx)
	ti := conv.TimeNew()
	oldSource := &pb.SettingSqlSource{}
	// 分为新增和编辑
	if r.Id > 0 {
		one, err = _data.GetSourceOne(dm.user.Env, r.Id)
		if err != nil {
			return nil, err
		}
		if one.Scene != source {
			return nil, errs.New(nil, "数据场景不一致")
		}
		jsoniter.UnmarshalFromString(one.Content, oldSource)
	} else {
		one.Scene = source
		one.Env = dm.user.Env
		one.Status = enum.StatusActive
		one.CreateDt = ti.String()
	}

	switch r.Type {
	case enum.DicSqlSource:
		// 提交密码与旧密码不一致就加密
		if r.Source.Sql.Password != "" && r.Source.Sql.Password != oldSource.Password {
			r.Source.Sql.Password, err = models.SqlSourceEncrypt(r.Source.Sql.Password)
			if err != nil {
				return nil, fmt.Errorf("加密失败，%w", err)
			}
		}
		if _, ok := enum.SqlDriverMap[r.Source.Sql.Driver]; !ok {
			if err != nil {
				return nil, fmt.Errorf("sql 驱动有误")
			}
		}
	case enum.DicJenkinsSource:
		r.Source.Jenkins.Hostname = strings.Trim(r.Source.Jenkins.Hostname, "/")
	case enum.DicGitSource:
		if r.Source.Git.Type != "gitee" && r.Source.Git.Type != "github" {
			return nil, fmt.Errorf("git 类型错误 %s", r.Source.Git.Type)
		}
	case enum.DicHostSource:
		if r.Source.Host.Ip == "" {
			return nil, errs.New(nil, "ip地址必填")
		}
		if r.Source.Host.Port == "" {
			return nil, errs.New(nil, "端口必填")
		}
	default:
		return nil, errs.New(nil, "type参数错误")
	}

	one.UpdateDt = ti.String()
	one.Title = r.Title
	one.Content, err = jsoniter.MarshalToString(r.Source)
	if err != nil {
		return nil, err
	}
	// 执行写入
	err = _data.Set(one)
	if err != nil {
		return nil, err
	}
	return &pb.SettingSqlSetReply{
		Id: one.Id,
	}, err
}

func (dm *SettingSqlService) Ping(r *pb.SettingSqlSetRequest) (resp *pb.SettingSqlSetReply, err error) {

	switch r.Type {
	case enum.DicSqlSource:
		password, err := models.SqlSourceDecode(r.Source.Sql.Password)
		if err != nil {
			return nil, fmt.Errorf("密码异常,%w", err)
		}
		conf := &config.MysqlSource{
			Hostname: r.Source.Sql.Hostname,
			Database: r.Source.Sql.Database,
			Username: r.Source.Sql.Username,
			Password: password,
			Port:     r.Source.Sql.Port,
		}
		if r.Source.Sql.Driver == enum.SqlDriverClickhouse {
			err = db.ConnClickhouse(conf).Error
		} else if r.Source.Sql.Driver == enum.SqlDriverMysql {
			err = db.Conn(conf).Error
		}
		if err != nil {
			return nil, err
		}

	case enum.DicJenkinsSource:
		r.Source.Jenkins.Hostname = strings.Trim(r.Source.Jenkins.Hostname, "/")
		req, er := http.NewRequest("GET", r.Source.Jenkins.Hostname+"/api/json", nil)
		if er != nil {
			return nil, errs.New(er, "请求构建失败")
		}
		if r.Source.Jenkins.Username != "" {
			req.SetBasicAuth(r.Source.Jenkins.Username, r.Source.Jenkins.Password)
		}
		// 创建 HTTP 客户端
		client := &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
			},
		}
		res, er := client.Do(req)
		if er != nil {
			return nil, errs.New(er, "请求执行失败")
		}
		defer res.Body.Close()

		if res.StatusCode != 200 {
			b, _ := io.ReadAll(res.Body)
			return nil, errs.New(fmt.Errorf(string(b)), "链接失败")
		}

	case enum.DicGitSource:
		cli := git.NewApi(git.Config{
			Type:        r.Source.Git.Type,
			AccessToken: r.Source.Git.AccessToken,
		})
		if cli == nil {
			return nil, fmt.Errorf("git 类型错误 %s", r.Source.Git.Type)
		}
		_, er := cli.User(git.NewHandler(dm.ctx))
		if err != nil {
			return nil, errs.New(er, "链接失败")
		}
	case enum.DicHostSource:
		if r.Source.Host.Ip == "" {
			return nil, errs.New(nil, "主机ip为必须")
		}
		if r.Source.Host.Port == "" {
			return nil, errs.New(nil, "主机端口为必须")
		}
		res, er := host.NewHost(&host.Config{
			Ip:     r.Source.Host.Ip,
			Port:   r.Source.Host.Port,
			User:   r.Source.Host.User,
			Secret: r.Source.Host.Secret,
		}).RemoteExec("echo ok")
		if er != nil {
			return nil, er
		}
		if strings.TrimSpace(string(res)) != "ok" {
			return nil, errs.New(nil, string(res))
		}
	default:
		return nil, errs.New(nil, "type参数错误")
	}

	return &pb.SettingSqlSetReply{}, nil
}

// 删除源
func (dm *SettingSqlService) ChangeStatus(r *pb.SettingChangeStatusRequest) (resp *pb.SettingChangeStatusReply, err error) {
	dm.db = db.New(dm.ctx)

	// 同一个任务，这里要加请求锁
	_data := data.NewCronSettingData(dm.ctx)
	one, err := _data.GetSourceOne(dm.user.Env, r.Id)
	if err != nil {
		return nil, err
	}
	if one.Id <= 0 {
		return nil, errors.New("操作数据不存在")
	}
	// 目前仅支持删除
	if r.Status != enum.StatusDelete {
		return nil, errors.New("不支持的状态操作")
	}
	// 这里还是要做是否使用的检测；
	// 如果使用未启用就联动置空（也不能删除，要么删除任务或者改任务），如果使用并启用禁止删除；
	// 如果没有试用就直接删除。
	list := []string{}
	err = dm.db.Raw(fmt.Sprintf("SELECT `name` FROM `cron_config` WHERE protocol=%v and JSON_CONTAINS(command, '%v', '$.sql.source.id') = 1", models.ProtocolSql, one.Id)).
		Scan(&list).Error
	if err != nil {
		return nil, fmt.Errorf("任务检测错误，%w", err)
	}
	if len(list) > 0 {
		return nil, fmt.Errorf("任务 %s 已使用连接，删除失败！", strings.Join(list, "、"))
	}

	err = _data.Del(one.Scene, dm.user.Env, one.Id)
	return &pb.SettingChangeStatusReply{}, err
}

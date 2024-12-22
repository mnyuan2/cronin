package data

import (
	"context"
	"cron/internal/basic/conv"
	"cron/internal/basic/enum"
	"cron/internal/basic/errs"
	"cron/internal/basic/git"
	"cron/internal/models"
	"cron/internal/pb"
	jsoniter "github.com/json-iterator/go"
)

type PRListItem struct {
	Number  int    `json:"number"`
	State   string `json:"state"`
	HeadRef string `json:"head_ref"`
	BaseRef string `json:"base_ref"`
}

type ConfigData struct {
	c            *models.CronConfig
	commandParse *pb.CronConfigCommand
	err          errs.Errs
}

func NewConfigData(c *models.CronConfig) *ConfigData {
	return (&ConfigData{
		c: c,
	}).Parse(nil)
}

func (m *ConfigData) Error() error {
	return m.err
}

// 解析
func (m *ConfigData) Parse(params map[string]any) *ConfigData {
	// 进行模板替换
	cmd, err := conv.DefaultStringTemplate().SetParam(params).Execute(m.c.Command)
	if err != nil {
		m.err = errs.New(err, "模板错误")
		return m
	}

	m.commandParse = &pb.CronConfigCommand{}
	if cmd != nil {
		if err := jsoniter.Unmarshal(cmd, m.commandParse); err != nil {
			m.err = errs.New(err, "配置解析错误")
			return m
		}
	}
	m.err = nil
	return m
}

// pr 列表查询（非job使用）
func (m *ConfigData) PRList(ctx context.Context, r *pb.GetEventPRList) (resp []*PRListItem, err errs.Errs) {
	if m.err != nil {
		return nil, m.err
	}
	if m.c.Protocol != models.ProtocolGit {
		return nil, errs.New(nil, "非git任务，操作失败")
	}
	if m.commandParse.Git.LinkId == 0 {
		return nil, errs.New(nil, "未设置链接")
	}

	// 查找到第一次有效的pr合并任务
	for _, e := range m.commandParse.Git.Events {
		if e.Id != enum.GitEventPullsMerge || e.PRMerge.Repo == "" || e.PRMerge.Number == "" || e.PRMerge.MergeMethod == "" {
			continue
		}
		r.Owner = e.PRMerge.Owner
		r.Repo = e.PRMerge.Repo
		break
	}

	link, er := NewCronSettingData(ctx).GetSourceOne(m.c.Env, m.commandParse.Git.LinkId)
	if er != nil {
		return nil, errs.New(er, "链接配置查询错误")
	}
	conf := &pb.SettingSource{
		Git: &pb.SettingGitSource{},
	}
	if er := jsoniter.UnmarshalFromString(link.Content, conf); er != nil {
		return nil, errs.New(er, "链接配置解析错误")
	}
	api := git.NewApi(git.Config{
		Type:        conf.Git.Type,
		AccessToken: conf.Git.AccessToken,
	})
	h := git.NewHandler(ctx)

	res, er := api.Pulls(h, &git.Pulls{
		BaseRequest: git.BaseRequest{
			Owner: r.Owner,
			Repo:  r.Repo,
		},
		State:   r.State,
		Head:    r.Head,
		Base:    r.Base,
		Page:    r.Page,
		PerPage: r.PerPage,
	})
	if er != nil {
		return nil, errs.New(er, "列表请求失败")
	}
	resp = make([]*PRListItem, len(res.List))
	for i, item := range res.List {
		resp[i] = &PRListItem{
			Number:  item.Number,
			State:   item.State,
			HeadRef: item.HeadRefName,
			BaseRef: item.BaseRefName,
		}
	}
	return resp, nil
}

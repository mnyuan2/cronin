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
	"time"
)

type WorkService struct {
	db   *db.MyDB
	ctx  context.Context
	user *auth.UserToken
}

func NewWorkService(ctx context.Context, user *auth.UserToken) *WorkService {
	return &WorkService{
		ctx:  ctx,
		user: user,
		db:   db.New(ctx),
	}
}

// 工作表格
func (dm *WorkService) Table(r *pb.WorkTableRequest) (resp *pb.WorkTableReply, err error) {
	/*
		查询所有任务，
			目前仅做待审核的。
	*/
	resp = &pb.WorkTableReply{
		List: []*pb.WorkTableItem{},
	}
	//if _, ok := models.ConfigStatusMap[r.Status]; !ok {
	//	return nil, errors.New("不支持的状态请求")
	//}
	w := db.NewWhere()
	switch r.Tab {
	case "todo":
		w.FindInSet("handle_user_ids", dm.user.UserId).
			In("status", []int{models.ConfigStatusAudited})
	case "created":
		w.Eq("create_user_id", dm.user.UserId)
	case "draft":
		w.Sub(func(sub *db.Where) {
			sub.Eq("create_user_id", dm.user.UserId).
				FindInSet("handle_user_ids", dm.user.UserId, db.OrOption())
		}).Eq("status", models.ConfigStatusDisable)
	default:
		return resp, nil
	}
	where, args := w.Build()

	sql := fmt.Sprintf(`SELECT COUNT(*) total, 'config' type, env FROM cron_config where %s GROUP BY env
UNION ALL
SELECT COUNT(*) total, 'pipeline' join_type, env FROM cron_pipeline where %s GROUP BY env`, where, where)
	args = append(args, args...)
	list := []*pb.WorkTableItem{}
	dm.db.Raw(sql, args...).Scan(&list)
	if len(list) == 0 {
		return resp, nil
	}
	envs, err := NewDicService(dm.ctx, dm.user).getDb(enum.DicEnv)
	if err != nil {
		return nil, err
	}

	listMap := map[string]*pb.WorkTableItem{}
	for _, val := range list {
		listMap[val.Env+"|"+val.Type] = val
	}
	for _, env := range envs {
		if item, ok := listMap[env.Key+"|config"]; ok {
			item.EnvTitle = env.Name
			resp.List = append(resp.List, item)
		}
		if item, ok := listMap[env.Key+"|pipeline"]; ok {
			item.EnvTitle = env.Name
			resp.List = append(resp.List, item)
		}
	}
	return resp, nil
}

// TaskDel 删除任务
func (dm *WorkService) TaskDel(r *pb.WorkTaskDelRequest) (resp *pb.WorkTaskDelReply, err error) {
	if r.Retention == "" {
		return nil, fmt.Errorf("retention 参数为必须")
	}

	re, err := time.ParseDuration(r.Retention)
	if err != nil {
		return nil, fmt.Errorf("retention 参数有误, %s", err.Error())
	} else if re.Hours() < 24 {
		return nil, fmt.Errorf("retention 参数不得小于24h")
	}
	end := time.Now().Add(-re)
	resp = &pb.WorkTaskDelReply{}
	w := db.NewWhere().Lte("status_dt", end.Format(time.DateTime)).Eq("status", models.ConfigStatusClosed)
	// 删除 config
	countConf, err := data.NewCronConfigData(dm.ctx).Del(w)
	if err != nil {
		return nil, fmt.Errorf("config 删除错误，%w", err)
	}
	resp.Info += fmt.Sprintf(" config delete count %v;", countConf)
	// 删除 pipeline
	countLine, err := data.NewCronPipelineData(dm.ctx).Del(w)
	if err != nil {
		return nil, fmt.Errorf("pipeline 删除错误，%w", err)
	}
	resp.Info += fmt.Sprintf(" pipeline delete count %v;", countLine)
	// 删除 receive
	countRece, err := data.NewCronReceiveData(dm.ctx).Del(w)
	if err != nil {
		return nil, fmt.Errorf("receive 删除错误，%w", err)
	}
	resp.Info += fmt.Sprintf(" receive delete count %v;", countRece)

	return resp, err
}

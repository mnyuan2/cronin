package biz

import (
	"context"
	"cron/internal/basic/auth"
	"cron/internal/basic/db"
	"cron/internal/data"
	"cron/internal/models"
	"cron/internal/pb"
	jsoniter "github.com/json-iterator/go"
)

type CronChangeLogService struct {
	ctx  context.Context
	user *auth.UserToken
}

func NewCronChangeLogService(ctx context.Context, user *auth.UserToken) *CronChangeLogService {
	return &CronChangeLogService{
		ctx:  ctx,
		user: user,
	}
}

// List 列表
func (dm *CronChangeLogService) List(r *pb.CronChangeLogListRequest) (resp *pb.CronChangeLogListResponse, err error) {
	w := db.NewWhere().
		Eq("ref_type", r.RefType, db.RequiredOption()).
		Eq("ref_id", r.RefId, db.RequiredOption())

	resp = &pb.CronChangeLogListResponse{
		List: []*pb.CronChangeLogItem{},
		Page: &pb.Page{
			Page: r.Page,
			Size: r.Size,
		},
	}
	resp.Page.Total, err = data.NewCronChangeLogData(dm.ctx).ListPage(w, r.Page, r.Size, &resp.List)

	for _, item := range resp.List {
		item.Content = []*pb.CronChangeLogItemField{}
		jsoniter.UnmarshalFromString(item.ContentStr, &item.Content)
		item.TypeName = models.LogTypeMap[item.Type]
	}

	return resp, err
}

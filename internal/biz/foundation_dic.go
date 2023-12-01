package biz

import (
	"context"
	"cron/internal/basic/conv"
	"cron/internal/basic/db"
	"cron/internal/basic/enum"
	"cron/internal/pb"
	jsoniter "github.com/json-iterator/go"
	"strings"
)

type DicService struct {
	ctx context.Context
}

type DicGetItem struct {
	// 键
	Id int32 `json:"id"`
	// 值
	Name string `json:"name"`
	// 其它数据，用于业务放关联操作
	Extend string `json:"extend"`
}

func NewDicService() *DicService {
	return &DicService{}
}

// 获得枚举
func (dm *DicService) Gets(ctx context.Context, r *pb.DicGetsRequest) (resp *pb.DicGetsReply, err error) {
	dm.ctx = ctx
	types := []int{}
	err = conv.NewStr().Slice(r.Types, types)
	if err != nil {
		return nil, err
	}

	resp = &pb.DicGetsReply{
		Maps: map[int]*pb.DicGetsList{},
	}
	for _, t := range types {
		list := &pb.DicGetsList{}
		if t <= 1000 {
			list.List, err = dm.getDb(t)
			if err != nil {
				return nil, err
			}
		} else {
			list.List, err = dm.getEnum(t)
			if err != nil {
				return nil, err
			}
		}

		resp.Maps[t] = list
	}

	return resp, err
}

// 通过数据库获取
func (dm *DicService) getDb(t int) ([]*pb.DicGetItem, error) {
	_sql := ""
	w := db.NewWhere()

	switch t {
	case enum.DicSqlSource:
		_sql = ""
	}

	items := []*pb.DicGetItem{}
	if _sql != "" {
		temp := []*DicGetItem{}
		where, args := w.Build()
		_sql = strings.Replace(_sql, "%WHERE", "WHERE "+where, 1)

		err := db.New(dm.ctx).Read.Raw(_sql, args...).Scan(&temp).Error
		if err != nil {
			return nil, err
		}
		for _, v := range temp {
			item := &pb.DicGetItem{
				Id:     v.Id,
				Name:   v.Name,
				Extend: &pb.DicExtendItem{},
			}
			if v.Extend != "" {
				if err = jsoniter.Unmarshal([]byte(v.Extend), item.Extend); err != nil {
					return nil, err
				}
			}
			items = append(items, item)
		}
	}

	return items, nil
}

// 通过枚举获取
func (dm *DicService) getEnum(t int) ([]*pb.DicGetItem, error) {
	// 待完善
	return nil, nil
}

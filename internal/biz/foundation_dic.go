package biz

import (
	"context"
	"cron/internal/basic/conv"
	"cron/internal/pb"
)

type DicService struct {
}

func NewDicService() *DicService {
	return &DicService{}
}

// 获得枚举
func (dm *DicService) Gets(ctx context.Context, r *pb.DicGetsRequest) (resp *pb.DicGetsReply, err error) {
	types := []int{}
	err = conv.NewStr().Slice(r.Types, types)
	if err != nil {
		return nil, err
	}

	return nil, err
}

func (dm *DicService) get() {

}

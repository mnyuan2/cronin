package dtos

import "cron/internal/pb"

func DicToMap(in []*pb.DicGetItem) map[int]string {
	out := map[int]string{}
	for _, item := range in {
		out[item.Id] = item.Name
	}
	return out
}

type DicGetRequest struct {
	Type int
	Env  string
}

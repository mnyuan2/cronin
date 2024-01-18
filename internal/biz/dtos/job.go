package dtos

import "cron/internal/pb"

// 消息设置解析
type MsgSetParse struct {
	MsgIds        []int
	NotifyUserIds []int
	Set           []*pb.CronMsgSet
}

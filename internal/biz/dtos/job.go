package dtos

import "cron/internal/pb"

// 消息设置解析
type MsgSetParse struct {
	MsgIds        []int
	StatusIn      map[int]any
	NotifyUserIds []int
	Set           []*pb.CronMsgSet
}

// 文件信息
type File struct {
	Name string
	Byte []byte
}

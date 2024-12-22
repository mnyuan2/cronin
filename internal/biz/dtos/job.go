package dtos

import "cron/internal/pb"

// 消息设置解析
type MsgSetParse struct {
	StatusList map[int][]*MsgSetItem // status:[{},{}]
	Set        []*pb.CronMsgSet
}
type MsgSetItem struct {
	MsgId         int   `json:"msg_id"`
	Status        int   `json:"status"`
	NotifyUserIds []int `json:"notify_user_ids"`
}

type MsgPushRequest struct {
	Status     int
	StatusDesc string
	Body       []byte
	Duration   float64
	RetryNum   int
	Args       map[string]any // 参数，如果默认存在将会被覆盖
}

// 文件信息
type File struct {
	Name string
	Byte []byte
}

package pb

// 用户列表
type JobStopRequest struct {
	RefId   int `json:"ref_id"`
	EntryId int `json:"entry_id"`
}
type JobStopReply struct{}

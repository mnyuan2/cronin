package pb

type DicGetsRequest struct {
	Types string `json:"types,omitempty"`
}

type DicGetsReply struct {
	Maps map[int]*DicGetsExtendList
}
type DicGetsExtendList struct {
	List []*DicGetExtendItem `json:"list"`
}
type DicGetExtendItem struct {
	// 键
	Id int32 `json:"id"`
	// 值
	Name string `json:"name"`
	// 其它数据，用于业务放关联操作
	Extend *DicExtendItem `json:"extend"`
	// 键2 (部分枚举采用的字符串)
	Key string `json:"key"`
}

type DicExtendItem struct{}

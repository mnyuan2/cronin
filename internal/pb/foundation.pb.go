package pb

type DicGetsRequest struct {
	Types string `json:"types" form:"types"`
	Env   string `json:"env" form:"env"`
}

type DicGetsReply struct {
	Maps map[int]*DicGetsList `json:"maps"`
}
type DicGetsList struct {
	List []*DicGetItem `json:"list"`
}
type DicGetItem struct {
	// 键
	Id int `json:"id"`
	// 值
	Name string `json:"name"`
	// 其它数据，用于业务放关联操作
	Extend *DicExtendItem `json:"extend"`
	// 键2 (部分枚举采用的字符串)
	Key string `json:"key"`
}

type DicExtendItem struct {
	Default   int    `json:"default"` // 默认 2.有效
	Driver    string `json:"driver"`  // 驱动·sql相关
	Remark    string `json:"remark"`
	Operation string `json:"operation,omitempty"`
}

// 系统信息
type SystemInfoRequest struct{}
type SystemInfoReply struct {
	Version     string `json:"version"`
	CmdName     string `json:"cmd_name"`
	Env         string `json:"env"`
	EnvName     string `json:"env_name"`
	CurrentDate string `json:"current_date"`
}

type ParseProtoRequest struct {
	Proto string `json:"proto"` // 文件内容
}
type ParseProtoReply struct {
	Actions []string `json:"actions"`
}

type ParseSpecRequest struct {
	Spec string `json:"spec"` // 时间表示
}
type ParseSpecReply struct {
	List []string `json:"list"`
}

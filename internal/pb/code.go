package pb

const (
	Success = "000000" // 状态 成功

	SysError         = "999999" // 系统错误
	NotFound         = "999998" // 访问不存在页面，确认URL地址是否正确
	IpDisabled       = "999997" // ip禁用，请联系管理员
	RequestOverclock = "999906" // 超过请求速率限制；请控制一下请求频率。
	UserNotExist     = "999907" // 身份不存在或权限不足
	UserError        = "999908" // 访问用户信息异常
	UserNotLogin     = "999909" // 用户未登录
	ParamNotFound    = "999910" // 参数未传递
	ParamError       = "999911" // 参数校验错误
	DataCoincide     = "999912" // 请求与目标数据一致
	OperationFailure = "999913" // 业务处理失败
	DataCreate       = "999801" // 数据新增失败
	DatalFind        = "999802" // 数据查询错误
	DataUpdate       = "999803" // 数据更新错误
)

var CodeList = map[string]string{
	Success:          "成功",
	SysError:         "系统错误！",
	NotFound:         "访问不存在页面，确认URL地址是否正确！",
	IpDisabled:       "ip禁用，请联系管理员！",
	RequestOverclock: "超过请求速率限制；请控制一下请求频率！",
	UserNotExist:     "身份不存在或权限不足！",
	UserError:        "访问用户信息异常！",
	UserNotLogin:     "用户未登录！",
	ParamNotFound:    "参数未传递！",
	ParamError:       "参数校验错误！",
	DataCoincide:     "请求与目标数据一致！",
	OperationFailure: "业务处理失败！",
	DataCreate:       "数据新增失败！",
	DatalFind:        "数据查询错误！",
	DataUpdate:       "数据更新错误！",
}

type Page struct {
	Size  int   `json:"size"`
	Page  int   `json:"page"`
	Total int64 `json:"total"`
}

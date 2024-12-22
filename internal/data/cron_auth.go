package data

const (
	AuthTypeOpen  = 1 // 无需登录,开放的
	AuthTypeLogin = 2 // 登录即可
	AuthTypeGrant = 3 // 授予权利
	//AuthTypeNotLogin = 4
)

// 权限项
type Permission struct {
	Id    int
	Path  string // 请求路径
	Title string // 标题
	Type  int    // 授权类型
	Group string // 分组名称
	Tag   string // 标签名称
}

// 权限
var authList = []*Permission{
	{
		Id:    1,
		Path:  "/index",
		Title: "首页·视图",
		Type:  AuthTypeOpen,
	},
	{
		Id:    2,
		Path:  "/login",
		Title: "登录·视图",
		Type:  AuthTypeOpen,
	},
	{Id: 11, Path: "/foundation/dic_gets", Title: "枚举选项", Group: "", Type: AuthTypeLogin},
	{Id: 12, Path: "/foundation/system_info", Title: "系统信息", Group: "基础", Type: AuthTypeLogin},
	{Id: 13, Path: "/foundation/parse_proto", Title: "解析proto", Group: "基础", Type: AuthTypeLogin},
	{Id: 14, Path: "/foundation/parse_spec", Title: "解析时间", Group: "基础", Type: AuthTypeLogin},

	{Id: 20, Path: "/config/", Title: "任务", Group: "", Type: AuthTypeGrant},
	{Id: 21, Path: "/config/list", Title: "列表", Group: "任务", Type: AuthTypeGrant, Tag: "config_list"},
	{Id: 22, Path: "/config/set", Title: "新增/编辑", Group: "任务", Type: AuthTypeGrant, Tag: "config_set"},
	{Id: 23, Path: "/config/change_status", Title: "提交/停用", Group: "任务", Type: AuthTypeGrant, Tag: "config_status"},
	{Id: 24, Path: "/config/change_status?auth_type=audit", Title: "审核", Group: "任务", Type: AuthTypeGrant, Tag: "config_audit"},
	{Id: 25, Path: "/config/detail", Title: "详情", Group: "任务", Type: AuthTypeGrant},
	{Id: 26, Path: "/config/run", Title: "执行一下", Group: "任务", Type: AuthTypeLogin},

	{Id: 30, Path: "/pipeline/", Title: "流水线", Group: "", Type: AuthTypeGrant},
	{Id: 31, Path: "/pipeline/list", Title: "列表", Group: "流水线", Type: AuthTypeGrant, Tag: "pipeline_list"},
	{Id: 32, Path: "/pipeline/set", Title: "新增/编辑", Group: "流水线", Type: AuthTypeGrant, Tag: "pipeline_set"},
	{Id: 33, Path: "/pipeline/change_status", Title: "提交/停用", Group: "流水线", Type: AuthTypeGrant, Tag: "pipeline_status"},
	{Id: 34, Path: "/pipeline/change_status?auth_type=audit", Title: "审核", Group: "流水线", Type: AuthTypeGrant, Tag: "pipeline_audit"},
	{Id: 35, Path: "/pipeline/detail", Title: "详情", Group: "流水线", Type: AuthTypeGrant},
	{Id: 36, Path: "/pipeline/run", Title: "执行一下", Group: "流水线", Type: AuthTypeLogin},

	// 150~159 接收
	{Id: 150, Path: "/receive/", Title: "接收", Group: "", Type: AuthTypeGrant},
	{Id: 151, Path: "/receive/list", Title: "列表", Group: "接收", Type: AuthTypeGrant, Tag: "receive_list"},
	{Id: 152, Path: "/receive/set", Title: "新增/编辑", Group: "接收", Type: AuthTypeGrant, Tag: "receive_set"},
	{Id: 153, Path: "/receive/change_status", Title: "提交/停用", Group: "接收", Type: AuthTypeGrant, Tag: "receive_status"},
	{Id: 154, Path: "/receive/change_status?auth_type=audit", Title: "审核", Group: "接收", Type: AuthTypeGrant, Tag: "receive_audit"},
	{Id: 155, Path: "/receive/detail", Title: "详情", Group: "接收", Type: AuthTypeGrant},
	{Id: 156, Path: "/receive/webhook/:key", Title: "开放钩子", Group: "接收", Type: AuthTypeOpen},

	//140~149 任务
	{Id: 140, Path: "/job/events", Title: "sse事件", Group: "注册任务", Type: AuthTypeLogin},
	{Id: 141, Path: "/job/stop", Title: "终止任务", Group: "注册任务", Type: AuthTypeLogin},
	{Id: 142, Path: "/job/list", Title: "已注册列表", Group: "注册任务", Type: AuthTypeLogin},

	{Id: 41, Path: "/work/table", Title: "工作表", Group: "我的", Type: AuthTypeLogin},
	{Id: 42, Path: "/work/task_del", Title: "任务删除", Group: "我的", Type: AuthTypeOpen},

	{Id: 50, Path: "/log/", Title: "日志", Group: "", Type: AuthTypeLogin},
	{Id: 51, Path: "/log/list", Title: "列表", Group: "日志", Type: AuthTypeLogin},
	{Id: 52, Path: "/log/traces", Title: "详情", Group: "日志", Type: AuthTypeLogin},
	{Id: 53, Path: "/log/del", Title: "删除", Group: "日志", Type: AuthTypeOpen},
	// 140~149
	{Id: 140, Path: "/change_log/list", Title: "列表", Group: "变更日志", Type: AuthTypeLogin},

	{Id: 60, Path: "/setting/", Title: "链接", Group: "", Type: AuthTypeGrant},
	{Id: 61, Path: "/setting/source_list", Title: "列表", Group: "链接", Type: AuthTypeGrant, Tag: "source_list"},
	{Id: 62, Path: "/setting/source_set", Title: "新增/编辑", Group: "链接", Type: AuthTypeGrant, Tag: "source_set"},
	{Id: 63, Path: "/setting/sql_source_change_status", Title: "状态变更", Group: "链接", Type: AuthTypeGrant, Tag: "source_status"},
	{Id: 64, Path: "/setting/source_ping", Title: "链接测试", Group: "链接", Type: AuthTypeLogin},

	{Id: 70, Path: "/setting/", Title: "环境", Group: "", Type: AuthTypeGrant},
	{Id: 71, Path: "/setting/env_list", Title: "列表", Group: "环境", Type: AuthTypeGrant, Tag: "env_list"},
	{Id: 72, Path: "/setting/env_set", Title: "新增/编辑", Group: "环境", Type: AuthTypeGrant, Tag: "env_set"},
	{Id: 73, Path: "/setting/env_set_content", Title: "-", Group: "环境", Type: AuthTypeLogin},
	{Id: 74, Path: "/setting/env_change_status", Title: "变更状态", Group: "环境", Type: AuthTypeGrant, Tag: "env_status"},
	{Id: 75, Path: "/setting/env_del", Title: "删除", Group: "环境", Type: AuthTypeGrant},

	{Id: 80, Path: "/setting/", Title: "消息", Group: "", Type: AuthTypeGrant},
	{Id: 81, Path: "/setting/message_list", Title: "列表", Group: "消息", Type: AuthTypeGrant, Tag: "message_list"},
	{Id: 82, Path: "/setting/message_set", Title: "新增/编辑", Group: "消息", Type: AuthTypeGrant, Tag: "message_set"},
	{Id: 83, Path: "/setting/message_run", Title: "执行一下", Group: "消息", Type: AuthTypeGrant},

	//120~139
	{Id: 120, Path: "/setting/", Title: "设置", Group: "", Type: AuthTypeGrant},
	{Id: 121, Path: "/setting/preference_set", Title: "偏好设置", Group: "设置", Type: AuthTypeGrant, Tag: "preference_set"},
	{Id: 122, Path: "/setting/preference_get", Title: "偏好查看", Group: "设置", Type: AuthTypeLogin},

	{Id: 90, Path: "/user/", Title: "用户", Group: "", Type: AuthTypeGrant},
	{Id: 91, Path: "/user/list", Title: "列表", Group: "用户", Type: AuthTypeGrant, Tag: "user_list"},
	{Id: 92, Path: "/user/set?auth_type=set", Title: "新增/编辑", Group: "用户", Type: AuthTypeGrant, Tag: "user_set"},
	{Id: 98, Path: "/user/set", Title: "编辑自己", Group: "用户", Type: AuthTypeLogin},
	{Id: 93, Path: "/user/change_password", Title: "修改密码", Group: "用户", Type: AuthTypeLogin},
	{Id: 94, Path: "/user/change_status", Title: "变更状态", Group: "用户", Type: AuthTypeLogin},
	{Id: 95, Path: "/user/change_account", Title: "设置账号", Group: "用户", Type: AuthTypeGrant, Tag: "user_account"},
	{Id: 96, Path: "/user/detail", Title: "详情", Group: "用户", Type: AuthTypeLogin},
	{Id: 97, Path: "/user/login", Title: "登录", Group: "用户", Type: AuthTypeOpen},

	{Id: 100, Path: "/role/", Title: "角色", Group: "", Type: AuthTypeGrant},
	{Id: 101, Path: "/role/list", Title: "列表", Group: "角色", Type: AuthTypeGrant, Tag: "role_list"},
	{Id: 102, Path: "/role/set", Title: "新增/编辑", Group: "角色", Type: AuthTypeGrant, Tag: "role_set"},
	{Id: 103, Path: "/role/auth_list", Title: "权限列表", Group: "角色", Type: AuthTypeLogin},
	{Id: 104, Path: "/role/auth_set", Title: "权限设置", Group: "角色", Type: AuthTypeGrant, Tag: "auth_set"},
	{Id: 105, Path: "/role/change_status", Title: "变更状态", Group: "角色", Type: AuthTypeGrant},
}
var authMap = map[string]*Permission{}

type AuthData struct {
}

func NewAuthData() *AuthData {
	return &AuthData{}
}

func (m AuthData) List() []*Permission {
	return authList
}

func (m AuthData) Map() map[string]*Permission {
	if len(authMap) != 0 {
		return authMap
	}
	for _, v := range authList {
		authMap[v.Path] = v
	}
	return authMap
}

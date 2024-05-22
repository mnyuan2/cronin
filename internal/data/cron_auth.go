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

	{Id: 20, Path: "/config/", Title: "任务", Group: "", Type: AuthTypeGrant},
	{Id: 21, Path: "/config/list", Title: "列表", Group: "任务", Type: AuthTypeGrant},
	{Id: 22, Path: "/config/set", Title: "设置", Group: "任务", Type: AuthTypeGrant},
	{Id: 23, Path: "/config/change_status", Title: "变更状态", Group: "任务", Type: AuthTypeGrant},
	{Id: 24, Path: "/config/get", Title: "详情", Group: "任务", Type: AuthTypeGrant},
	{Id: 25, Path: "/config/run", Title: "执行一下", Group: "任务", Type: AuthTypeLogin},
	{Id: 26, Path: "/config/register_list", Title: "已注册列表", Group: "任务", Type: AuthTypeLogin},

	{Id: 30, Path: "/pipeline/", Title: "流水线", Group: "", Type: AuthTypeGrant},
	{Id: 31, Path: "/pipeline/list", Title: "列表", Group: "流水线", Type: AuthTypeGrant},
	{Id: 32, Path: "/pipeline/set", Title: "设置", Group: "流水线", Type: AuthTypeGrant},
	{Id: 33, Path: "/pipeline/change_status", Title: "变更状态", Group: "流水线", Type: AuthTypeGrant},

	{Id: 41, Path: "/work/table", Title: "工作表", Group: "我的", Type: AuthTypeLogin},

	{Id: 50, Path: "/log/", Title: "日志", Group: "", Type: AuthTypeLogin},
	{Id: 51, Path: "/log/list", Title: "列表", Group: "日志", Type: AuthTypeLogin},
	{Id: 52, Path: "/log/traces", Title: "详情", Group: "日志", Type: AuthTypeLogin},
	{Id: 53, Path: "/log/del", Title: "删除", Group: "日志", Type: AuthTypeLogin},

	{Id: 60, Path: "/setting/", Title: "链接资源", Group: "", Type: AuthTypeGrant},
	{Id: 61, Path: "/setting/source_list", Title: "列表", Group: "链接资源", Type: AuthTypeGrant},
	{Id: 62, Path: "/setting/source_set", Title: "设置", Group: "链接资源", Type: AuthTypeGrant},
	{Id: 63, Path: "/setting/sql_source_change_status", Title: "状态变更", Group: "链接资源", Type: AuthTypeGrant},
	{Id: 64, Path: "/setting/source_ping", Title: "链接测试", Group: "链接资源", Type: AuthTypeLogin},

	{Id: 70, Path: "/setting/", Title: "环境", Group: "", Type: AuthTypeGrant},
	{Id: 71, Path: "/setting/env_list", Title: "列表", Group: "环境", Type: AuthTypeGrant},
	{Id: 72, Path: "/setting/env_set", Title: "设置", Group: "环境", Type: AuthTypeGrant},
	{Id: 73, Path: "/setting/env_set_content", Title: "-", Group: "环境", Type: AuthTypeLogin},
	{Id: 74, Path: "/setting/env_change_status", Title: "变更状态", Group: "环境", Type: AuthTypeGrant},
	{Id: 75, Path: "/setting/env_del", Title: "删除", Group: "环境", Type: AuthTypeGrant},

	{Id: 80, Path: "/setting/", Title: "消息", Group: "", Type: AuthTypeGrant},
	{Id: 81, Path: "/setting/message_list", Title: "列表", Group: "消息", Type: AuthTypeGrant},
	{Id: 82, Path: "/setting/message_set", Title: "设置", Group: "消息", Type: AuthTypeGrant},
	{Id: 83, Path: "/setting/message_run", Title: "执行一下", Group: "消息", Type: AuthTypeGrant},

	{Id: 90, Path: "/user/", Title: "用户", Group: "", Type: AuthTypeGrant},
	{Id: 91, Path: "/user/list", Title: "列表", Group: "用户", Type: AuthTypeGrant},
	{Id: 92, Path: "/user/set", Title: "设置", Group: "用户", Type: AuthTypeGrant},
	{Id: 93, Path: "/user/change_password", Title: "修改密码", Group: "用户", Type: AuthTypeLogin},
	{Id: 94, Path: "/user/change_status", Title: "变更状态", Group: "用户", Type: AuthTypeLogin},
	{Id: 95, Path: "/user/change_account", Title: "设置账号", Group: "用户", Type: AuthTypeGrant},
	{Id: 96, Path: "/user/detail", Title: "详情", Group: "用户", Type: AuthTypeLogin},
	{Id: 97, Path: "/user/login", Title: "登录", Group: "用户", Type: AuthTypeOpen},

	{Id: 100, Path: "/role/", Title: "角色", Group: "", Type: AuthTypeGrant},
	{Id: 101, Path: "/role/list", Title: "列表", Group: "角色", Type: AuthTypeGrant},
	{Id: 102, Path: "/role/set", Title: "设置", Group: "角色", Type: AuthTypeGrant},
	{Id: 103, Path: "/role/auth_list", Title: "权限列表", Group: "角色", Type: AuthTypeLogin},
	{Id: 104, Path: "/role/auth_set", Title: "权限设置", Group: "角色", Type: AuthTypeGrant},
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

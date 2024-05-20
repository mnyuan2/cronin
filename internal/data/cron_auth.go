package data

const (
	AuthTypeOpen  = 1 // 无需登录,开放的
	AuthTypeLogin = 2 // 登录即可
	AuthTypeGrant = 3 // 授予权利
	//AuthTypeNotLogin = 4
)

// 权限项
type Permission struct {
	Path  string // 请求路径
	Title string // 标题

	Type int // 授权类型
}

// 权限
var Permissions = map[string]*Permission{
	"/index": {
		Path:  "/index",
		Title: "首页·视图",
		Type:  AuthTypeOpen,
	},
	"/login": {
		Path:  "/login",
		Title: "登录·视图",
		Type:  AuthTypeOpen,
	},
	"/user/login": {
		Path:  "/user/login",
		Title: "登录",
		Type:  AuthTypeOpen,
	},
}

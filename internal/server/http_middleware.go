package server

import (
	"cron/internal/basic/auth"
	"cron/internal/basic/cache"
	"cron/internal/basic/errs"
	"cron/internal/data"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

// CORS跨域
func useCors() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Header("Access-Control-Allow-Origin", "*")
		ctx.Header("Access-Control-Allow-Methods", "GET,POST,OPTIONS,PUT,DELETE")
		ctx.Header("Access-Control-Allow-Headers", "X-Custom-Header,Content-Type,Authorization")
		ctx.Header("Access-Control-Expose-Headers", "Content-Disposition,Trace-Id") // 解决跨域自定义头获取
		// 放行所有的options请求
		if ctx.Request.Method == "OPTIONS" {
			ctx.AbortWithStatus(http.StatusOK)
		} else {
			ctx.Next()
		}
	}
}

// 授权中间件
// @param perms
func UseAuth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		path := ctx.FullPath()
		env := ctx.GetHeader("env")
		if auth_type := ctx.Request.URL.Query().Get("auth_type"); auth_type != "" { // 固定的参数校验头
			path += "?auth_type=" + auth_type
			ctx.Set("auth_type", auth_type) // 业务需要根据这个参数来判断权限
		}
		if path == "" || path == "/" {
			ctx.Next()
			return
		}
		node, ok := data.NewAuthData().Map()[path]
		if !ok {
			NewReply(ctx).SetReply(nil, errors.New("无权访问！")).RenderJson()
			ctx.Abort() // 终止
			return
			//node = &data.Permission{Type: data.AuthTypeLogin}
		}
		// 无需登录
		if node.Type == data.AuthTypeOpen {
			ctx.Next()
			return
		}
		user, err := auth.ParseJwtToken(ctx.Request.Header.Get("Authorization"))
		if err != nil {
			//authFailed(ctx.Writer, err.Error())
			NewReply(ctx).SetReply(nil, err).RenderJson()
			ctx.Abort() // 终止
			return
		}
		user.Env = env
		ctx.Set("user", user)
		// 仅登录
		if node.Type == data.AuthTypeLogin {
			ctx.Next()
			return
		}
		// 需要授权
		if node.Type == data.AuthTypeGrant {
			menu := cache.Get(fmt.Sprintf("user_menu_%v", user.UserId))
			if menu == nil {
				NewReply(ctx).SetReply(nil, errs.New(errors.New("请重新登录！"), errs.UserNotLogin)).RenderJson()
				ctx.Abort() // 终止
				return
			}
			if _, ok := menu.(map[int]int)[node.Id]; !ok {
				NewReply(ctx).SetReply(nil, errors.New("没有访问权限！")).RenderJson()
				ctx.Abort() // 终止
				return
			}

		} else {
			NewReply(ctx).SetReply(nil, errors.New("节点异常！")).RenderJson()
			ctx.Abort() // 终止
			return
		}

		ctx.Next()
	}
}

// 获得用户信息
func GetUser(ctx *gin.Context) (user *auth.UserToken, err error) {
	u, ok := ctx.Get("user")
	if !ok {
		return nil, errors.New("用户信息未找到！")
	}
	return u.(*auth.UserToken), nil
}

/*认证失败*/
func authFailed(w http.ResponseWriter, msg string) {
	w.Header().Set("WWW-Authenticate", `Basic realm="My METRICS"`)
	http.Error(w, msg, http.StatusUnauthorized)
}

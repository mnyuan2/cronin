package server

import (
	"cron/internal/basic/auth"
	"errors"
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
func UseAuth(NotUserRote map[string]string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		path := ctx.FullPath()
		env := ctx.GetHeader("env")
		if _, ok := NotUserRote[path]; !ok {
			user, err := auth.ParseToken(ctx)
			if err != nil {
				authFailed(ctx.Writer, err.Error())
				//rep.NewResult(ctx).SetErr(err).RenderJson()
				ctx.Abort() // 终止
				return
			}
			user.Env = env
			ctx.Set("user", user)
		} else {
			ctx.Set("user", &auth.UserToken{UserName: "无", Env: env})
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

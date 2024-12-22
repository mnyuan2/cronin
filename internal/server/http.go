package server

import (
	"cron/internal/basic/config"
	"cron/internal/basic/sse"
	"embed"
	"github.com/gin-gonic/gin"
	"html/template"
	"io/fs"
	"net/http"
)

// Init http 初始化
func InitHttp(Resource embed.FS, isBuildResource bool) *gin.Engine {
	r := gin.New()

	if isBuildResource {
		// 二进制版本,打包使用（优点，静态资源将被打包至二进制文件）
		s, e := fs.Sub(Resource, "web/static")
		if e != nil {
			panic("资源错误 " + e.Error())
		}
		c, e := fs.Sub(Resource, "web/components")
		if e != nil {
			panic("组件错误 " + e.Error())
		}
		r.StaticFS("/static", http.FS(s)).StaticFS("/components", http.FS(c))
		r.SetHTMLTemplate(template.Must(template.New("").Delims("[[", "]]").ParseFS(Resource, "web/*.html")))
	} else {
		r.Delims("[[", "]]")
		r.LoadHTMLGlob("web/*.html")
		r.Static("/static", "web/static")
		r.Static("/components", "web/components")
	}

	r.Use(gin.Recovery(), UseAuth())
	// api
	r.GET("/foundation/dic_gets", routerDicGets)
	r.GET("/foundation/system_info", routerSystemInfo)
	r.POST("/foundation/parse_proto", routerParseProto)
	r.POST("/foundation/parse_spec", routerParseSpec)

	r.GET("/config/list", httpList)
	r.GET("/config/detail", httpConfigDetail)
	r.POST("/config/set", httpSet)
	r.POST("/config/change_status", httpChangeStatus)
	r.POST("/config/run", httpRun)

	r.GET("/pipeline/list", routerPipelineList)
	r.GET("/pipeline/detail", httpPipelineDetail)
	r.POST("/pipeline/set", routerPipelineSet)
	r.POST("/pipeline/change_status", routerPipelineChangeStatus)
	r.POST("/pipeline/run", routerPipelineRun)

	r.POST("/receive/set", routerReceiveSet)
	r.GET("/receive/list", routerReceiveList)
	r.GET("/receive/detail", routerReceiveDetail)
	r.POST("/receive/change_status", routerReceiveChangeStatus)
	r.POST("/receive/webhook/:key", routerReceiveWebhook)

	r.GET("/job/events", func(ctx *gin.Context) {
		sse.Serve().ServeHTTP(ctx.Writer, ctx.Request)
	})
	r.GET("/job/list", httpRegister)
	r.POST("/job/stop", httpJobStop)

	r.GET("/work/table", routerWorkTable)
	r.POST("/work/task_del", routerWorkTaskDel)

	r.GET("/log/list", routerLogList)
	r.GET("/log/traces", routerLogTraces)
	r.POST("/log/del", routerLogDel)

	r.GET("/change_log/list", routerChangeLogList)

	r.GET("/setting/source_list", routerSqlList)
	r.POST("/setting/source_set", routerSqlSet)
	r.POST("/setting/sql_source_change_status", routerSqlChangeStatus)
	r.POST("/setting/source_ping", routerSqlPing)
	r.GET("/setting/env_list", routerEnvList)
	r.POST("/setting/env_set", routerEnvSet)
	r.POST("/setting/env_set_content", routerEnvSetContent)
	r.POST("/setting/env_change_status", routerEnvChangeStatus)
	r.POST("/setting/env_del", routerEnvDel)
	r.GET("/setting/message_list", routerMessageList)
	r.POST("/setting/message_set", routerMessageSet)
	r.POST("/setting/message_run", routerMessageRun)
	r.POST("/setting/preference_set", routerPreferenceSet)
	r.GET("/setting/preference_get", routerPreferenceGet)

	r.GET("/user/list", routerUserList)
	r.POST("/user/set", routerUserSet)
	r.POST("/user/change_password", routerUserChangePassword)
	r.POST("/user/change_status", routerUserChangeStatus)
	r.POST("/user/change_account", routerUserChangeAccount)
	r.GET("/user/detail", routerUserDetail)
	r.POST("/user/login", routerUserLogin)

	r.POST("/role/set", routerRoleSet)
	r.GET("/role/list", routerRoleList)
	r.GET("/role/auth_list", routerAuthList)
	r.POST("/role/auth_set", routerRoleAuthSet)
	r.POST("/role/change_status")

	// 视图
	r.GET("/", func(ctx *gin.Context) {
		ctx.Redirect(http.StatusMovedPermanently, "/index")
	})
	r.GET("/index", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "index.html", map[string]string{"version": config.Version})
	})
	r.GET("/login", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "login.html", map[string]string{"version": config.Version})
	})

	return r
}

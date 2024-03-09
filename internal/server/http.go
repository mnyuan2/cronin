package server

import (
	"cron/internal/basic/config"
	"embed"
	"github.com/gin-gonic/gin"
	"html/template"
	"io/fs"
	"net/http"
)

// Init http 初始化
func InitHttp(Resource embed.FS, isBuildResource bool) *gin.Engine {
	r := gin.Default()

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

	r.Use(UseAuth(nil))
	// api
	r.GET("/foundation/dic_gets", routerDicGets)
	r.GET("/foundation/system_info", routerSystemInfo)
	r.POST("/foundation/parse_proto", routerParseProto)

	r.GET("/config/list", httpList)
	r.POST("/config/set", httpSet)
	r.POST("/config/change_status", httpChangeStatus)
	r.GET("/config/get")
	r.POST("/config/run", httpRun)
	r.GET("/config/register_list", httpRegister)
	r.GET("/pipeline/list", routerPipelineList)
	r.POST("/pipeline/set", routerPipelineSet)
	r.POST("/pipeline/change_status", routerPipelineChangeStatus)

	r.GET("/log/list", routerLogList)
	r.GET("/log/traces", routerLogTraces)
	r.POST("/log/del", routerLogDel)

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

	r.GET("/user/list", routerUserList)
	r.POST("/user/set", routerUserSet)

	// 视图
	r.GET("/", func(ctx *gin.Context) {
		ctx.Redirect(http.StatusMovedPermanently, "/index")
	})
	r.GET("/index", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "index.html", map[string]string{"version": config.Version})
	})

	return r
}

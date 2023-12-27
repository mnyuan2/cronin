package server

import (
	"cron/internal/basic/config"
	"embed"
	"github.com/gin-gonic/gin"
	"html/template"
	"io/fs"
	"net/http"
	"runtime"
)

// Init http 初始化
func InitHttp(Resource embed.FS) *gin.Engine {

	r := gin.Default()

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

	//r.Delims("[[", "]]")
	//r.LoadHTMLGlob("web/*.html")
	//r.Static("/static", "web/static")
	//r.Static("/components", "web/components")

	r.Use(UseAuth(nil))

	r.GET("/foundation/dic_gets", routerDicGets)
	r.GET("/config/list", httpList)
	r.POST("/config/set", httpSet)
	r.POST("/config/change_status", httpChangeStatus)
	r.GET("/config/get")
	r.GET("/config/register_list", httpRegister)
	r.GET("/log/by_config", httpLogByConfig)
	r.POST("/log/del", httpLogDel)
	r.GET("/setting/sql_source_list", routerSqlList)
	r.POST("/setting/sql_source_set", routerSqlSet)
	r.POST("/setting/sql_source_change_status", routerSqlChangeStatus)
	r.POST("/setting/sql_source_ping", routerSqlPing)

	gv := r.Group("view")
	gv.GET("/cron/list", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "cron_list.html", map[string]string{"version": config.Version})
	})
	r.GET("/system/info", func(ctx *gin.Context) {
		cmd_name := "sh"
		if runtime.GOOS == "windows" {
			cmd_name = "cmd"
		}
		NewReply(ctx).SetSuccess(map[string]string{
			"cmd_name": cmd_name,
		}).RenderJson()
	})

	return r
}

package server

import (
	"cron/internal/basic/config"
	"cron/internal/biz"
	"cron/internal/pb"
	"embed"
	"github.com/gin-gonic/gin"
	"net/http"
	"runtime"
)

// Init http 初始化
func InitHttp(Resource embed.FS) *gin.Engine {

	r := gin.Default()

	// 二进制版本,打包使用（优点，静态资源将被打包至二进制文件）
	//s, e := fs.Sub(Resource, "web/static")
	//if e != nil {
	//	panic("资源错误 " + e.Error())
	//}
	//r.StaticFS("/static", http.FS(s))
	//r.SetHTMLTemplate(template.Must(template.New("").Delims("[[", "]]").ParseFS(Resource, "web/*.html")))

	r.Delims("[[", "]]")
	r.LoadHTMLGlob("web/*.html")
	r.Static("/static", "web/static")
	r.Static("/components", "web/components")

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

// 查看已注册任务
func httpRegister(ctx *gin.Context) {
	rep, err := biz.NewCronConfigService().RegisterList(ctx.Request.Context(), nil)
	NewReply(ctx).SetReply(rep, err).RenderJson()
}

// 任务列表
func httpList(ctx *gin.Context) {
	r := &pb.CronConfigListRequest{}
	if err := ctx.BindQuery(r); err != nil {
		NewReply(ctx).SetError(pb.ParamError, err.Error()).RenderJson()
		return
	}
	rep, err := biz.NewCronConfigService().List(ctx.Request.Context(), r)
	NewReply(ctx).SetReply(rep, err).RenderJson()
}

// 任务设置
func httpSet(ctx *gin.Context) {
	r := &pb.CronConfigSetRequest{}
	if err := ctx.BindJSON(r); err != nil {
		NewReply(ctx).SetError(pb.ParamError, err.Error()).RenderJson()
		return
	}
	rep, err := biz.NewCronConfigService().Set(ctx.Request.Context(), r)
	NewReply(ctx).SetReply(rep, err).RenderJson()
}

// 任务状态变更
func httpChangeStatus(ctx *gin.Context) {
	r := &pb.CronConfigSetRequest{}
	if err := ctx.BindJSON(r); err != nil {
		NewReply(ctx).SetError(pb.ParamError, err.Error()).RenderJson()
		return
	}
	rep, err := biz.NewCronConfigService().ChangeStatus(ctx.Request.Context(), r)
	NewReply(ctx).SetReply(rep, err).RenderJson()
}

// 任务状态变更
func httpLogByConfig(ctx *gin.Context) {
	r := &pb.CronLogByConfigRequest{}
	if err := ctx.BindQuery(r); err != nil {
		NewReply(ctx).SetError(pb.ParamError, err.Error()).RenderJson()
		return
	}
	rep, err := biz.NewCronLogService().ByConfig(ctx.Request.Context(), r)
	NewReply(ctx).SetReply(rep, err).RenderJson()
}

// 删除日志
func httpLogDel(ctx *gin.Context) {
	r := &pb.CronLogDelRequest{}
	if err := ctx.BindJSON(r); err != nil {
		NewReply(ctx).SetError(pb.ParamError, err.Error()).RenderJson()
		return
	}
	rep, err := biz.NewCronLogService().Del(ctx.Request.Context(), r)
	NewReply(ctx).SetReply(rep, err).RenderJson()
}

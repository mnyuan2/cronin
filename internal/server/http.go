package server

import (
	"cron/internal/biz"
	"cron/internal/pb"
	"github.com/gin-gonic/gin"
	"net/http"
)

// Init http 初始化
func InitHttp() *gin.Engine {

	r := gin.Default()
	r.Delims("[[", "]]")
	r.LoadHTMLGlob("web/*.html")
	r.Static("/static", "web/static")

	r.GET("/config/list", httpList)
	r.POST("/config/set", httpSet)
	r.POST("/config/edit", httpEdit)
	r.POST("/config/change_status", httpChangeStatus)
	r.GET("/config/get")
	r.POST("/config/del")

	gv := r.Group("view")
	gv.GET("/cron/list", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "cron_list.html", map[string]string{})
	})

	return r
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

// 任务编辑
func httpEdit(ctx *gin.Context) {
	r := &pb.CronConfigSetRequest{}
	if err := ctx.BindJSON(r); err != nil {
		NewReply(ctx).SetError(pb.ParamError, err.Error()).RenderJson()
		return
	}
	rep, err := biz.NewCronConfigService().Edit(ctx.Request.Context(), r)
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

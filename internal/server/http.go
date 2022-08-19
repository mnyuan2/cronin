package server

import (
	"cron/internal/biz"
	"cron/internal/pb"
	"github.com/gin-gonic/gin"
	"net/http"
)

// Init http 初始化
func Init() *gin.Engine {

	r := gin.Default()
	r.Delims("[[","]]")
	r.LoadHTMLGlob("web/*.html")
	r.Static("/static", "web/static")

	r.GET("/config/list", httpList)
	r.POST("/config/set")
	r.GET("/config/get")
	r.POST("/config/del")

	gv := r.Group("view")
	gv.GET("/cron/list", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "cron_list.html", nil)
	})



	return r
}

func httpList(ctx *gin.Context){
	r := &pb.CronConfigListRequest{}
	if err := ctx.BindQuery(r); err != nil{
		NewReply(ctx).SetError(pb.ParamError,err.Error()).RenderJson()
		return
	}
	rep, err := biz.NewCronConfigService().List(ctx.Request.Context(), r)
	NewReply(ctx).SetReply(rep, err).RenderJson()
}
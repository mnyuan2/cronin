package server

import (
	"cron/internal/biz"
	"cron/internal/pb"
	"github.com/gin-gonic/gin"
)

// 获取选项列表
func routerDicGets(ctx *gin.Context) {
	r := &pb.DicGetsRequest{}
	if err := ctx.BindQuery(r); err != nil {
		NewReply(ctx).SetError(pb.ParamError, err.Error()).RenderJson()
		return
	}
	rep, err := biz.NewDicService().Gets(ctx.Request.Context(), r)
	NewReply(ctx).SetReply(rep, err).RenderJson()
}

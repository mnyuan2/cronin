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
	user, err := GetUser(ctx)
	if err != nil {
		NewReply(ctx).SetError(pb.UserNotExist, err.Error()).RenderJson()
		return
	}
	rep, err := biz.NewDicService(ctx.Request.Context(), user).DicGets(r)
	NewReply(ctx).SetReply(rep, err).RenderJson()
}

func routerSystemInfo(ctx *gin.Context) {

	rep, err := biz.NewDicService(ctx.Request.Context(), nil).SystemInfo(nil)
	NewReply(ctx).SetReply(rep, err).RenderJson()
}

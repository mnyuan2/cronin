package server

import (
	"cron/internal/biz"
	"cron/internal/pb"
	"github.com/gin-gonic/gin"
)

// 工作表格
func routerWorkTable(ctx *gin.Context) {
	r := &pb.WorkTableRequest{}
	if err := ctx.BindQuery(r); err != nil {
		NewReply(ctx).SetError(pb.ParamError, err.Error()).RenderJson()
		return
	}
	user, err := GetUser(ctx)
	if err != nil {
		NewReply(ctx).SetError(pb.UserNotExist, err.Error()).RenderJson()
		return
	}
	rep, err := biz.NewWorkService(ctx.Request.Context(), user).Table(r)
	NewReply(ctx).SetReply(rep, err).RenderJson()
}

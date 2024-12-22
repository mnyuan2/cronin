package server

import (
	"cron/internal/biz"
	"cron/internal/pb"
	"github.com/gin-gonic/gin"
)

// 列表
func routerTagList(ctx *gin.Context) {
	//r := &pb.ReceiveListRequest{}
	//if err := ctx.BindQuery(r); err != nil {
	//	NewReply(ctx).SetError(pb.ParamError, err.Error()).RenderJson()
	//	return
	//}

	user, err := GetUser(ctx)
	if err != nil {
		NewReply(ctx).SetError(pb.UserNotExist, err.Error()).RenderJson()
		return
	}

	rep, err := biz.NewCronTagService(ctx.Request.Context(), user).List(nil)
	NewReply(ctx).SetReply(rep, err).RenderJson()
}

// 设置
func routerTagSet(ctx *gin.Context) {
	r := &pb.TagSetRequest{}
	if err := ctx.BindJSON(r); err != nil {
		NewReply(ctx).SetError(pb.ParamError, err.Error()).RenderJson()
		return
	}
	user, err := GetUser(ctx)
	if err != nil {
		NewReply(ctx).SetError(pb.UserNotExist, err.Error()).RenderJson()
		return
	}
	rep, err := biz.NewCronTagService(ctx.Request.Context(), user).Set(r)
	NewReply(ctx).SetReply(rep, err).RenderJson()
}

// 状态变更
func routerTagChangeStatus(ctx *gin.Context) {
	r := &pb.TagChangeStatusRequest{}
	if err := ctx.BindJSON(r); err != nil {
		NewReply(ctx).SetError(pb.ParamError, err.Error()).RenderJson()
		return
	}
	user, err := GetUser(ctx)
	if err != nil {
		NewReply(ctx).SetError(pb.UserNotExist, err.Error()).RenderJson()
		return
	}
	rep, err := biz.NewCronTagService(ctx.Request.Context(), user).ChangeStatus(r)
	NewReply(ctx).SetReply(rep, err).RenderJson()
}

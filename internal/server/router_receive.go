package server

import (
	"cron/internal/biz"
	"cron/internal/pb"
	"github.com/gin-gonic/gin"
)

// 列表
func routerReceiveList(ctx *gin.Context) {
	r := &pb.ReceiveListRequest{}
	if err := ctx.BindQuery(r); err != nil {
		NewReply(ctx).SetError(pb.ParamError, err.Error()).RenderJson()
		return
	}

	user, err := GetUser(ctx)
	if err != nil {
		NewReply(ctx).SetError(pb.UserNotExist, err.Error()).RenderJson()
		return
	}

	rep, err := biz.NewReceiveService(ctx.Request.Context(), user).List(r)
	NewReply(ctx).SetReply(rep, err).RenderJson()
}

// 设置
func routerReceiveSet(ctx *gin.Context) {
	r := &pb.ReceiveSetRequest{}
	if err := ctx.BindJSON(r); err != nil {
		NewReply(ctx).SetError(pb.ParamError, err.Error()).RenderJson()
		return
	}
	user, err := GetUser(ctx)
	if err != nil {
		NewReply(ctx).SetError(pb.UserNotExist, err.Error()).RenderJson()
		return
	}
	rep, err := biz.NewReceiveService(ctx.Request.Context(), user).Set(r)
	NewReply(ctx).SetReply(rep, err).RenderJson()
}

// 详情
func routerReceiveDetail(ctx *gin.Context) {
	r := &pb.ReceiveDetailRequest{}
	user, err := GetUser(ctx)
	if err != nil {
		NewReply(ctx).SetError(pb.UserNotExist, err.Error()).RenderJson()
		return
	}
	rep, err := biz.NewReceiveService(ctx.Request.Context(), user).Detail(r)
	NewReply(ctx).SetReply(rep, err).RenderJson()
}

// 状态变更
func routerReceiveChangeStatus(ctx *gin.Context) {
	r := &pb.ReceiveChangeStatusRequest{}
	if err := ctx.BindJSON(r); err != nil {
		NewReply(ctx).SetError(pb.ParamError, err.Error()).RenderJson()
		return
	}
	user, err := GetUser(ctx)
	if err != nil {
		NewReply(ctx).SetError(pb.UserNotExist, err.Error()).RenderJson()
		return
	}
	rep, err := biz.NewReceiveService(ctx.Request.Context(), user).ChangeStatus(r)
	NewReply(ctx).SetReply(rep, err).RenderJson()
}

// 接收钩子
func routerReceiveWebhook(ctx *gin.Context) {
	r := &pb.ReceiveWebhookRequest{}
	if err := ctx.BindJSON(r); err != nil {
		NewReply(ctx).SetError(pb.ParamError, err.Error()).RenderJson()
		return
	}
	rep, err := biz.NewReceiveService(ctx.Request.Context(), nil).Webhook(r)
	NewReply(ctx).SetReply(rep, err).RenderJson()
}

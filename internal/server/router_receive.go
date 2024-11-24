package server

import (
	"cron/internal/biz"
	"cron/internal/pb"
	"github.com/gin-gonic/gin"
	"io"
	"strconv"
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
	if err := ctx.BindQuery(r); err != nil {
		NewReply(ctx).SetError(pb.ParamError, err.Error()).RenderJson()
		return
	}
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
	var err error
	// 这里的接受结构体无法预定义，只能统一作为一个字符串接受
	r := &pb.ReceiveWebhookRequest{}
	r.Body, err = io.ReadAll(ctx.Request.Body)
	if err != nil {
		NewReply(ctx).SetError(pb.ParamError, err.Error()).RenderJson()
		return
	}
	key := ctx.Param("key")
	if key == "" {
		NewReply(ctx).SetError(pb.ParamError, "接受key未指定").RenderJson()
		return
	}

	r.Id, err = strconv.Atoi(key)
	if err != nil {
		r.Alias = key
	}

	rep, err := biz.NewReceiveService(ctx.Request.Context(), nil).Webhook(r)
	NewReply(ctx).SetReply(rep, err).RenderJson()
}

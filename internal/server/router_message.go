package server

import (
	"cron/internal/biz"
	"cron/internal/pb"
	"github.com/gin-gonic/gin"
)

// 设置 sql 连接源 列表
func routerMessageList(ctx *gin.Context) {
	r := &pb.SettingMessageListRequest{}
	if err := ctx.BindQuery(r); err != nil {
		NewReply(ctx).SetError(pb.ParamError, err.Error()).RenderJson()
		return
	}
	rep, err := biz.NewSettingMessageService(ctx.Request.Context()).List(r)
	NewReply(ctx).SetReply(rep, err).RenderJson()
}

// 任务设置
func routerMessageSet(ctx *gin.Context) {
	r := &pb.SettingMessageSetRequest{}
	if err := ctx.BindJSON(r); err != nil {
		NewReply(ctx).SetError(pb.ParamError, err.Error()).RenderJson()
		return
	}
	rep, err := biz.NewSettingMessageService(ctx.Request.Context()).Set(r)
	NewReply(ctx).SetReply(rep, err).RenderJson()
}

// 任务设置
func routerMessageRun(ctx *gin.Context) {
	r := &pb.SettingMessageSetRequest{}
	if err := ctx.BindJSON(r); err != nil {
		NewReply(ctx).SetError(pb.ParamError, err.Error()).RenderJson()
		return
	}
	rep, err := biz.NewSettingMessageService(ctx.Request.Context()).Run(r)
	NewReply(ctx).SetReply(rep, err).RenderJson()
}

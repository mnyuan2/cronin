package server

import (
	"cron/internal/biz"
	"cron/internal/pb"
	"github.com/gin-gonic/gin"
)

// 设置 sql 连接源 列表
func routerEnvList(ctx *gin.Context) {
	r := &pb.SettingEnvListRequest{}
	if err := ctx.BindQuery(r); err != nil {
		NewReply(ctx).SetError(pb.ParamError, err.Error()).RenderJson()
		return
	}
	rep, err := biz.NewSettingEnvService(ctx.Request.Context()).List(r)
	NewReply(ctx).SetReply(rep, err).RenderJson()
}

// 任务设置
func routerEnvSet(ctx *gin.Context) {
	r := &pb.SettingEnvSetRequest{}
	if err := ctx.BindJSON(r); err != nil {
		NewReply(ctx).SetError(pb.ParamError, err.Error()).RenderJson()
		return
	}
	rep, err := biz.NewSettingEnvService(ctx.Request.Context()).Set(r)
	NewReply(ctx).SetReply(rep, err).RenderJson()
}

// 任务设置
func routerEnvSetContent(ctx *gin.Context) {
	r := &pb.SettingEnvSetRequest{}
	if err := ctx.BindJSON(r); err != nil {
		NewReply(ctx).SetError(pb.ParamError, err.Error()).RenderJson()
		return
	}
	rep, err := biz.NewSettingEnvService(ctx.Request.Context()).SetContent(r)
	NewReply(ctx).SetReply(rep, err).RenderJson()
}

// 任务状态变更
func routerEnvChangeStatus(ctx *gin.Context) {
	r := &pb.SettingChangeStatusRequest{}
	if err := ctx.BindJSON(r); err != nil {
		NewReply(ctx).SetError(pb.ParamError, err.Error()).RenderJson()
		return
	}
	rep, err := biz.NewSettingEnvService(ctx.Request.Context()).ChangeStatus(r)
	NewReply(ctx).SetReply(rep, err).RenderJson()
}

// 任务状态变更
func routerEnvDel(ctx *gin.Context) {
	r := &pb.SettingEnvDelRequest{}
	if err := ctx.BindJSON(r); err != nil {
		NewReply(ctx).SetError(pb.ParamError, err.Error()).RenderJson()
		return
	}
	rep, err := biz.NewSettingEnvService(ctx.Request.Context()).Del(r)
	NewReply(ctx).SetReply(rep, err).RenderJson()
}

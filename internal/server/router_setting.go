package server

import (
	"cron/internal/biz"
	"cron/internal/pb"
	"github.com/gin-gonic/gin"
)

// 设置 sql 连接源 列表
func routerSqlList(ctx *gin.Context) {
	r := &pb.SettingSqlListRequest{}
	if err := ctx.BindQuery(r); err != nil {
		NewReply(ctx).SetError(pb.ParamError, err.Error()).RenderJson()
		return
	}
	rep, err := biz.NewSettingSqlService().List(ctx.Request.Context(), r)
	NewReply(ctx).SetReply(rep, err).RenderJson()
}

// 任务设置
func routerSqlSet(ctx *gin.Context) {
	r := &pb.SettingSqlSetRequest{}
	if err := ctx.BindJSON(r); err != nil {
		NewReply(ctx).SetError(pb.ParamError, err.Error()).RenderJson()
		return
	}
	rep, err := biz.NewSettingSqlService().Set(ctx.Request.Context(), r)
	NewReply(ctx).SetReply(rep, err).RenderJson()
}

// 任务状态变更
func routerSqlChangeStatus(ctx *gin.Context) {
	r := &pb.SettingChangeStatusRequest{}
	if err := ctx.BindJSON(r); err != nil {
		NewReply(ctx).SetError(pb.ParamError, err.Error()).RenderJson()
		return
	}
	rep, err := biz.NewSettingSqlService().ChangeStatus(ctx.Request.Context(), r)
	NewReply(ctx).SetReply(rep, err).RenderJson()
}

// sql连接ping
func routerSqlPing(ctx *gin.Context) {
	r := &pb.SettingSqlPingRequest{}
	if err := ctx.BindJSON(r); err != nil {
		NewReply(ctx).SetError(pb.ParamError, err.Error()).RenderJson()
		return
	}
	rep, err := biz.NewSettingSqlService().Ping(ctx.Request.Context(), r)
	NewReply(ctx).SetReply(rep, err).RenderJson()
}

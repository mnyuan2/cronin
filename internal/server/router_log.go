package server

import (
	"cron/internal/biz"
	"cron/internal/pb"
	"github.com/gin-gonic/gin"
)

// 日志列表
func routerLogList(ctx *gin.Context) {
	r := &pb.CronLogListRequest{}
	if err := ctx.BindQuery(r); err != nil {
		NewReply(ctx).SetError(pb.ParamError, err.Error()).RenderJson()
		return
	}
	user, err := GetUser(ctx)
	if err != nil {
		NewReply(ctx).SetError(pb.UserNotExist, err.Error()).RenderJson()
		return
	}
	rep, err := biz.NewCronLogService(ctx.Request.Context(), user).List(r)
	NewReply(ctx).SetReply(rep, err).RenderJson()
}

// 日志踪迹
func routerLogTraces(ctx *gin.Context) {
	r := &pb.CronLogTraceRequest{}
	if err := ctx.BindQuery(r); err != nil {
		NewReply(ctx).SetError(pb.ParamError, err.Error()).RenderJson()
		return
	}
	user, err := GetUser(ctx)
	if err != nil {
		NewReply(ctx).SetError(pb.UserNotExist, err.Error()).RenderJson()
		return
	}
	rep, err := biz.NewCronLogService(ctx.Request.Context(), user).Trace(r)
	NewReply(ctx).SetReply(rep, err).RenderJson()
}

// 删除日志
func routerLogDel(ctx *gin.Context) {
	r := &pb.CronLogDelRequest{}
	if err := ctx.BindJSON(r); err != nil {
		NewReply(ctx).SetError(pb.ParamError, err.Error()).RenderJson()
		return
	}
	//user, err := GetUser(ctx)
	//if err != nil {
	//	NewReply(ctx).SetError(pb.UserNotExist, err.Error()).RenderJson()
	//	return
	//}
	rep, err := biz.NewCronLogService(ctx.Request.Context(), nil).Del(r)
	NewReply(ctx).SetReply(rep, err).RenderJson()
}

// 变更日志列表
func routerChangeLogList(ctx *gin.Context) {
	r := &pb.CronChangeLogListRequest{}
	if err := ctx.BindQuery(r); err != nil {
		NewReply(ctx).SetError(pb.ParamError, err.Error()).RenderJson()
		return
	}
	user, err := GetUser(ctx)
	if err != nil {
		NewReply(ctx).SetError(pb.UserNotExist, err.Error()).RenderJson()
		return
	}
	rep, err := biz.NewCronChangeLogService(ctx.Request.Context(), user).List(r)
	NewReply(ctx).SetReply(rep, err).RenderJson()
}

package server

import (
	"cron/internal/biz"
	"cron/internal/pb"
	"github.com/gin-gonic/gin"
)

// 查看已注册任务
func httpRegister(ctx *gin.Context) {
	user, err := GetUser(ctx)
	if err != nil {
		NewReply(ctx).SetError(pb.UserNotExist, err.Error()).RenderJson()
		return
	}
	rep, err := biz.NewCronConfigService(ctx.Request.Context(), user).RegisterList(nil)
	NewReply(ctx).SetReply(rep, err).RenderJson()
}

// 查看已注册任务
func httpJobStop(ctx *gin.Context) {
	r := &pb.JobStopRequest{}
	if err := ctx.BindJSON(r); err != nil {
		NewReply(ctx).SetError(pb.ParamError, err.Error()).RenderJson()
		return
	}
	user, err := GetUser(ctx)
	if err != nil {
		NewReply(ctx).SetError(pb.UserNotExist, err.Error()).RenderJson()
		return
	}
	rep, err := biz.NewJobService(ctx.Request.Context(), user).Stop(r)
	NewReply(ctx).SetReply(rep, err).RenderJson()
}

// 执行任务踪迹
func httpJobTraces(ctx *gin.Context) {
	r := &pb.JobTracesRequest{}
	if err := ctx.BindQuery(r); err != nil {
		NewReply(ctx).SetError(pb.ParamError, err.Error()).RenderJson()
		return
	}
	user, err := GetUser(ctx)
	if err != nil {
		NewReply(ctx).SetError(pb.UserNotExist, err.Error()).RenderJson()
		return
	}
	rep, err := biz.NewJobService(ctx.Request.Context(), user).Traces(r)
	NewReply(ctx).SetReply(rep, err).RenderJson()
}

package server

import (
	"cron/internal/biz"
	"cron/internal/pb"
	"github.com/gin-gonic/gin"
)

// 任务状态变更
func httpLogByConfig(ctx *gin.Context) {
	r := &pb.CronLogByConfigRequest{}
	if err := ctx.BindQuery(r); err != nil {
		NewReply(ctx).SetError(pb.ParamError, err.Error()).RenderJson()
		return
	}
	user, err := GetUser(ctx)
	if err != nil {
		NewReply(ctx).SetError(pb.UserNotExist, err.Error()).RenderJson()
		return
	}
	rep, err := biz.NewCronLogService(ctx.Request.Context(), user).ByConfig(r)
	NewReply(ctx).SetReply(rep, err).RenderJson()
}

// 删除日志
func httpLogDel(ctx *gin.Context) {
	r := &pb.CronLogDelRequest{}
	if err := ctx.BindJSON(r); err != nil {
		NewReply(ctx).SetError(pb.ParamError, err.Error()).RenderJson()
		return
	}
	user, err := GetUser(ctx)
	if err != nil {
		NewReply(ctx).SetError(pb.UserNotExist, err.Error()).RenderJson()
		return
	}
	rep, err := biz.NewCronLogService(ctx.Request.Context(), user).Del(r)
	NewReply(ctx).SetReply(rep, err).RenderJson()
}

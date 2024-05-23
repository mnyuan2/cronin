package server

import (
	"cron/internal/biz"
	"cron/internal/models"
	"cron/internal/pb"
	"github.com/gin-gonic/gin"
)

// 流水线列表
func routerPipelineList(ctx *gin.Context) {
	r := &pb.CronPipelineListRequest{}
	if err := ctx.BindQuery(r); err != nil {
		NewReply(ctx).SetError(pb.ParamError, err.Error()).RenderJson()
		return
	}
	user, err := GetUser(ctx)
	if err != nil {
		NewReply(ctx).SetError(pb.UserNotExist, err.Error()).RenderJson()
		return
	}

	rep, err := biz.NewCronPipelineService(ctx.Request.Context(), user).List(r)
	NewReply(ctx).SetReply(rep, err).RenderJson()
}

// 流水线设置
func routerPipelineSet(ctx *gin.Context) {
	r := &pb.CronPipelineSetRequest{}
	if err := ctx.BindJSON(r); err != nil {
		NewReply(ctx).SetError(pb.ParamError, err.Error()).RenderJson()
		return
	}
	user, err := GetUser(ctx)
	if err != nil {
		NewReply(ctx).SetError(pb.UserNotExist, err.Error()).RenderJson()
		return
	}
	rep, err := biz.NewCronPipelineService(ctx.Request.Context(), user).Set(r)
	NewReply(ctx).SetReply(rep, err).RenderJson()
}

// 流水线状态变更
func routerPipelineChangeStatus(ctx *gin.Context) {
	r := &pb.CronPipelineChangeStatusRequest{}
	if err := ctx.BindJSON(r); err != nil {
		NewReply(ctx).SetError(pb.ParamError, err.Error()).RenderJson()
		return
	}
	user, err := GetUser(ctx)
	if err != nil {
		NewReply(ctx).SetError(pb.UserNotExist, err.Error()).RenderJson()
		return
	}
	if authType := ctx.GetString("auth_type"); authType == "audit" {
		if r.Status != models.ConfigStatusActive && r.Status != models.ConfigStatusReject { // 只能操作 通过、驳回
			NewReply(ctx).SetError(pb.ParamError, "权限与状态不匹配").RenderJson()
			return
		}
	} else {
		if r.Status == models.ConfigStatusActive || r.Status == models.ConfigStatusReject { // 不能操作 通过、驳回
			NewReply(ctx).SetError(pb.ParamError, "权限与状态不匹配").RenderJson()
			return
		}
	}

	rep, err := biz.NewCronPipelineService(ctx.Request.Context(), user).ChangeStatus(r)
	NewReply(ctx).SetReply(rep, err).RenderJson()
}

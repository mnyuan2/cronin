package server

import (
	"cron/internal/biz"
	"cron/internal/models"
	"cron/internal/pb"
	"github.com/gin-gonic/gin"
)

// 任务列表
func httpList(ctx *gin.Context) {
	r := &pb.CronConfigListRequest{Ids: []int{}, CreateUserIds: []int{}, HandleUserIds: []int{}}
	if err := ctx.BindQuery(r); err != nil {
		NewReply(ctx).SetError(pb.ParamError, err.Error()).RenderJson()
		return
	}
	user, err := GetUser(ctx)
	if err != nil {
		NewReply(ctx).SetError(pb.UserNotExist, err.Error()).RenderJson()
		return
	}

	rep, err := biz.NewCronConfigService(ctx.Request.Context(), user).List(r)
	NewReply(ctx).SetReply(rep, err).RenderJson()
}

// 匹配列表
func httpMatchList(ctx *gin.Context) {
	r := &pb.CronMatchListRequest{}
	if err := ctx.BindJSON(r); err != nil {
		NewReply(ctx).SetError(pb.ParamError, err.Error()).RenderJson()
		return
	}
	user, err := GetUser(ctx)
	if err != nil {
		NewReply(ctx).SetError(pb.UserNotExist, err.Error()).RenderJson()
		return
	}

	rep, err := biz.NewCronConfigService(ctx.Request.Context(), user).MatchList(r)
	NewReply(ctx).SetReply(rep, err).RenderJson()
}

// 任务详情
func httpConfigDetail(ctx *gin.Context) {
	r := &pb.CronConfigDetailRequest{}
	if err := ctx.BindQuery(r); err != nil {
		NewReply(ctx).SetError(pb.ParamError, err.Error()).RenderJson()
		return
	}
	user, err := GetUser(ctx)
	if err != nil {
		NewReply(ctx).SetError(pb.UserNotExist, err.Error()).RenderJson()
		return
	}

	rep, err := biz.NewCronConfigService(ctx.Request.Context(), user).Detail(r)
	NewReply(ctx).SetReply(rep, err).RenderJson()
}

// 任务设置
func httpSet(ctx *gin.Context) {
	r := &pb.CronConfigSetRequest{}
	if err := ctx.BindJSON(r); err != nil {
		NewReply(ctx).SetError(pb.ParamError, err.Error()).RenderJson()
		return
	}
	user, err := GetUser(ctx)
	if err != nil {
		NewReply(ctx).SetError(pb.UserNotExist, err.Error()).RenderJson()
		return
	}
	rep, err := biz.NewCronConfigService(ctx.Request.Context(), user).Set(r)
	NewReply(ctx).SetReply(rep, err).RenderJson()
}

// 任务状态变更
func httpChangeStatus(ctx *gin.Context) {
	r := &pb.CronConfigSetRequest{}
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

	rep, err := biz.NewCronConfigService(ctx.Request.Context(), user).ChangeStatus(r)
	NewReply(ctx).SetReply(rep, err).RenderJson()
}

// 任务 运行一次
func httpRun(ctx *gin.Context) {
	r := &pb.CronConfigRunRequest{}
	if err := ctx.BindJSON(r); err != nil {
		NewReply(ctx).SetError(pb.ParamError, err.Error()).RenderJson()
		return
	}
	user, err := GetUser(ctx)
	if err != nil {
		NewReply(ctx).SetError(pb.UserNotExist, err.Error()).RenderJson()
		return
	}
	rep, err := biz.NewCronConfigService(ctx.Request.Context(), user).Run(r)
	NewReply(ctx).SetReply(rep, err).RenderJson()
}

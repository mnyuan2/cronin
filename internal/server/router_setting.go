package server

import (
	"cron/internal/biz"
	"cron/internal/pb"
	"github.com/gin-gonic/gin"
)

// 设置 sql 连接源 列表
func routerSqlList(ctx *gin.Context) {
	r := &pb.SettingListRequest{}
	if err := ctx.BindQuery(r); err != nil {
		NewReply(ctx).SetError(pb.ParamError, err.Error()).RenderJson()
		return
	}
	user, err := GetUser(ctx)
	if err != nil {
		NewReply(ctx).SetError(pb.UserNotExist, err.Error()).RenderJson()
		return
	}
	rep, err := biz.NewSettingSqlService(ctx.Request.Context(), user).List(r)
	NewReply(ctx).SetReply(rep, err).RenderJson()
}

// 任务设置
func routerSqlSet(ctx *gin.Context) {
	r := &pb.SettingSqlSetRequest{}
	if err := ctx.BindJSON(r); err != nil {
		NewReply(ctx).SetError(pb.ParamError, err.Error()).RenderJson()
		return
	}
	user, err := GetUser(ctx)
	if err != nil {
		NewReply(ctx).SetError(pb.UserNotExist, err.Error()).RenderJson()
		return
	}
	rep, err := biz.NewSettingSqlService(ctx.Request.Context(), user).Set(r)
	NewReply(ctx).SetReply(rep, err).RenderJson()
}

// 任务状态变更
func routerSqlChangeStatus(ctx *gin.Context) {
	r := &pb.SettingChangeStatusRequest{}
	if err := ctx.BindJSON(r); err != nil {
		NewReply(ctx).SetError(pb.ParamError, err.Error()).RenderJson()
		return
	}
	user, err := GetUser(ctx)
	if err != nil {
		NewReply(ctx).SetError(pb.UserNotExist, err.Error()).RenderJson()
		return
	}
	rep, err := biz.NewSettingSqlService(ctx.Request.Context(), user).ChangeStatus(r)
	NewReply(ctx).SetReply(rep, err).RenderJson()
}

// sql连接ping
func routerSqlPing(ctx *gin.Context) {
	r := &pb.SettingSqlSetRequest{}
	if err := ctx.BindJSON(r); err != nil {
		NewReply(ctx).SetError(pb.ParamError, err.Error()).RenderJson()
		return
	}
	user, err := GetUser(ctx)
	if err != nil {
		NewReply(ctx).SetError(pb.UserNotExist, err.Error()).RenderJson()
		return
	}
	rep, err := biz.NewSettingSqlService(ctx.Request.Context(), user).Ping(r)
	NewReply(ctx).SetReply(rep, err).RenderJson()
}

// 偏好设置
func routerPreferenceSet(ctx *gin.Context) {
	r := &pb.SettingPreferenceSetRequest{}
	if err := ctx.BindJSON(r); err != nil {
		NewReply(ctx).SetError(pb.ParamError, err.Error()).RenderJson()
		return
	}

	rep, err := biz.NewSettingService(ctx.Request.Context()).PreferenceSet(r)
	NewReply(ctx).SetReply(rep, err).RenderJson()
}

// 偏好获取
func routerPreferenceGet(ctx *gin.Context) {
	r := &pb.SettingPreferenceGetRequest{}

	rep, err := biz.NewSettingService(ctx.Request.Context()).PreferenceGet(r)
	NewReply(ctx).SetReply(rep, err).RenderJson()
}

// 全局变量列表
func routerGlobalVariateList(ctx *gin.Context) {
	r := &pb.GlobalVariateListRequest{}

	rep, err := biz.NewSettingService(ctx.Request.Context()).GlobalVariateList(r)
	NewReply(ctx).SetReply(rep, err).RenderJson()
}

// 全局变量设置
func routerGlobalVariateSet(ctx *gin.Context) {
	r := &pb.GlobalVariateSetRequest{}
	if err := ctx.BindJSON(r); err != nil {
		NewReply(ctx).SetError(pb.ParamError, err.Error()).RenderJson()
		return
	}

	rep, err := biz.NewSettingService(ctx.Request.Context()).GlobalVariateSet(r)
	NewReply(ctx).SetReply(rep, err).RenderJson()
}

// 全局变量状态变更
func routerGlobalVariateChangeStatus(ctx *gin.Context) {
	r := &pb.GlobalVariateSetRequest{}
	if err := ctx.BindJSON(r); err != nil {
		NewReply(ctx).SetError(pb.ParamError, err.Error()).RenderJson()
		return
	}
	rep, err := biz.NewSettingService(ctx.Request.Context()).GlobalVariateChangeStatus(r)
	NewReply(ctx).SetReply(rep, err).RenderJson()
}

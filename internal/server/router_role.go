package server

import (
	"cron/internal/biz"
	"cron/internal/pb"
	"github.com/gin-gonic/gin"
)

// 用户列表
func routerRoleList(ctx *gin.Context) {
	r := &pb.RoleListRequest{}
	if err := ctx.BindQuery(r); err != nil {
		NewReply(ctx).SetError(pb.ParamError, err.Error()).RenderJson()
		return
	}
	rep, err := biz.NewRoleService(ctx.Request.Context()).List(r)
	NewReply(ctx).SetReply(rep, err).RenderJson()
}

// 用户设置
func routerRoleSet(ctx *gin.Context) {
	r := &pb.RoleSetRequest{}
	if err := ctx.BindJSON(r); err != nil {
		NewReply(ctx).SetError(pb.ParamError, err.Error()).RenderJson()
		return
	}
	rep, err := biz.NewRoleService(ctx.Request.Context()).Set(r)
	NewReply(ctx).SetReply(rep, err).RenderJson()
}

// 规则列表
func routerAuthList(ctx *gin.Context) {
	r := &pb.AuthListRequest{}

	rep, err := biz.NewRoleService(ctx.Request.Context()).AuthList(r)
	NewReply(ctx).SetReply(rep, err).RenderJson()
}

// 角色规则设置
func routerRoleAuthSet(ctx *gin.Context) {
	r := &pb.RoleAuthSetRequest{}
	if err := ctx.BindJSON(r); err != nil {
		NewReply(ctx).SetError(pb.ParamError, err.Error()).RenderJson()
		return
	}
	rep, err := biz.NewRoleService(ctx.Request.Context()).AuthSet(r)
	NewReply(ctx).SetReply(rep, err).RenderJson()
}

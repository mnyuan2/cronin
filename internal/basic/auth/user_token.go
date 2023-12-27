package auth

import (
	"cron/internal/basic/config"
	"errors"
	"github.com/gin-gonic/gin"
)

const (
	// http授权Header名称
	httpAuthHeader = "Authorization"
)

// 令牌用户信息
type UserToken struct {
	//Env int    `json:"env"`
	//UserId   int    `json:"user_id"`
	UserName string `json:"user_name"`
}

// 解析http令牌
func ParseHttpToken(ctx *gin.Context) (u *UserToken, err error) {
	conf := config.MainConf().User
	// 不校验账号
	if conf == nil || conf.AdminAccount == "" {
		return &UserToken{UserName: "无"}, nil
	}
	account, password, ok := ctx.Request.BasicAuth()
	if !ok {
		return nil, errors.New("401 Unauthorized!")
	}
	if account != conf.AdminAccount || password != conf.AdminPassword {
		return nil, errors.New("401 Password error!")
	}

	u = &UserToken{UserName: "管理员"}
	return u, nil
}

func GetUser(ctx *gin.Context) (user *UserToken, is bool) {
	u, ok := ctx.Get("user")
	if !ok {
		return nil, false
	}
	return u.(*UserToken), true
}

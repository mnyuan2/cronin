package auth

import (
	"cron/internal/basic/config"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"time"
)

const (
	// http授权Header名称
	httpAuthHeader = "Authorization"
)

// 令牌用户信息
type UserToken struct {
	jwt.StandardClaims
	//Env int    `json:"env"`
	UserId   int    `json:"user_id"`
	UserName string `json:"user_name"`
	Env      string `json:"env,omitempty"`
}

// 解析http令牌
func ParseToken(ctx *gin.Context) (u *UserToken, err error) {
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

// 生成令牌
func GenToken(id int, name string) (string, error) {
	unixTime := time.Now().Unix()
	u := &UserToken{
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  unixTime,
			NotBefore: unixTime,
			Issuer:    "cronin",
			Subject:   "cronin",
		},
		UserId:   id,
		UserName: name,
		Env:      "",
	}

	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, u)
	return tokenClaims.SignedString([]byte(config.MainConf().Crypto.Secret))
}

func GetUser(ctx *gin.Context) (user *UserToken, is bool) {
	u, ok := ctx.Get("user")
	if !ok {
		return nil, false
	}
	return u.(*UserToken), true
}

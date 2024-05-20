package auth

import (
	"cron/internal/basic/config"
	"cron/internal/basic/errs"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/segmentio/asm/ascii"
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

// 解析 http 令牌 v2
func ParseJwtToken(token string) (u *UserToken, err error) {
	if token == "" {
		return nil, errs.New(errors.New("用户未登录"), errs.UserNotLogin)
	}
	const prefix = "Bearer "
	preLen := len(prefix)
	if len(token) < preLen || !ascii.EqualFoldString(token[:preLen], prefix) {
		return nil, errors.New("令牌错误")
	}

	user := &UserToken{}
	jwtToken, err := jwt.ParseWithClaims(token[preLen:], user, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.MainConf().Crypto.Secret), nil
	})
	if err != nil {
		if err.Error() == "Token is expired" {
			return nil, errs.New(errors.New("令牌已过期"), "999802")
		}
		return nil, err
	}

	if jwtToken != nil && jwtToken.Valid {
		if claim, ok := jwtToken.Claims.(*UserToken); ok {
			return claim, nil
		}
	}
	return nil, errors.New("令牌解析失败")
}

// 生成令牌
func GenJwtToken(id int, name string) (string, error) {
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
	token, err := tokenClaims.SignedString([]byte(config.MainConf().Crypto.Secret))
	return "Bearer " + token, err
}

func GetUser(ctx *gin.Context) (user *UserToken, is bool) {
	u, ok := ctx.Get("user")
	if !ok {
		return nil, false
	}
	return u.(*UserToken), true
}

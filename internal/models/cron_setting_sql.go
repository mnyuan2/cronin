package models

import (
	"cron/internal/basic/config"
	"cron/internal/basic/conv"
)

const (
	SqlErrActionAbort    = 1 // 终止
	SqlErrActionProceed  = 2 // 继续
	SqlErrActionRollback = 3 // 事务回滚
)

var SqlErrActionMap = map[int]string{
	SqlErrActionAbort:    "终止任务",
	SqlErrActionProceed:  "跳过继续",
	SqlErrActionRollback: "事务回滚",
}

// sql驱动
const (
	SqlSourceMysql = "mysql"
)

// 加密
func SqlSourceEncrypt(data string) (string, error) {
	secret := config.MainConf().Crypto.Secret
	if len(secret) < 8 {
		return data, nil
	}
	str, err := conv.Des(secret, secret).Encrypt(data)
	if err != nil {
		return "", err
	}
	return "d." + str, nil
}

// 解密
func SqlSourceDecode(data string) (string, error) {
	if len(data) <= 2 || data[:2] != "d." {
		return data, nil
	}
	secret := config.MainConf().Crypto.Secret
	if len(secret) < 8 {
		return data, nil
	}
	return conv.Des(secret, secret).Decrypt(data[2:])
}

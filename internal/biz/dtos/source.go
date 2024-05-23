package dtos

import (
	"cron/internal/basic/errs"
	"cron/internal/models"
	"cron/internal/pb"
	jsoniter "github.com/json-iterator/go"
)

func ParseSource(data *models.CronSetting) (*pb.SettingSource, error) {
	s := &pb.SettingSource{}
	if er := jsoniter.UnmarshalFromString(data.Content, s); er != nil {
		return nil, errs.New(er, "连接配置解析异常")
	}
	return s, nil
}

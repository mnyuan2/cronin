package dtos

import (
	"cron/internal/basic/conv"
	"cron/internal/basic/errs"
	"cron/internal/pb"
	jsoniter "github.com/json-iterator/go"
)

func ParseParams(varFields []byte, after func(list map[string]any)) (map[string]any, errs.Errs) {
	varParams := map[string]any{}
	if len(varFields) > 5 {
		// 参数也可以通过模板初始化，以获得动态默认值
		str, err := conv.DefaultStringTemplate().Execute(varFields)
		if err != nil {
			return nil, errs.New(err, "变量模板错误")
		}

		temp := []*pb.KvItem{}
		if err := jsoniter.Unmarshal(str, &temp); err != nil {
			return nil, errs.New(err, "变量参数字段解析错误")
		}
		for _, item := range temp {
			if item.Key == "" {
				continue
			}
			varParams[item.Key] = item.Value
		}
	}
	after(varParams)
	return varParams, nil
}

// 解析
func ParseCommon(command []byte, params map[string]any) (*pb.CronConfigCommand, errs.Errs) {
	// 进行模板替换
	cmd, err := conv.DefaultStringTemplate().SetParam(params).Execute(command)
	if err != nil {
		return nil, errs.New(err, "模板错误")
	}

	commandParse := &pb.CronConfigCommand{}
	if err := jsoniter.Unmarshal(cmd, commandParse); err != nil {
		return nil, errs.New(err, "配置解析错误")
	}
	return commandParse, nil

}

package conv

import (
	"errors"
	"strconv"
)

type int64s struct{}

func Int64s() *int64s {
	return &int64s{}
}

// 将int64转换为字符串;
func (i *int64s) String(val int64) string {
	return strconv.FormatInt(val, 10)
}

// 将字符串转换为整数;
func (i *int64s) Parse(val string) (int64, error) {
	if val == "" {
		return 0, nil
	}

	result, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		return 0, err
	}
	return result, nil
}

func (i *int64s) ParseAny(val any) (int64, error) {
	switch val.(type) {
	case int32:
		return int64(val.(int32)), nil
	case int:
		return int64(val.(int)), nil
	case int64:
		return val.(int64), nil
	case string:
		return i.Parse(val.(string))
	case float64:
		return int64(val.(float64)), nil
	case float32:
		return int64(val.(float32)), nil
	}
	return 0, errors.New("未支持的输入类型")
}

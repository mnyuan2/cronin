package conv

import (
	"errors"
	"strconv"
	"strings"
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

func (i *int64s) Join(list any) (string, error) {
	//t := reflect.TypeOf(list)
	//if t.Kind() == reflect.Ptr {
	//	t = t.Elem()
	//}
	//if t.Elem().Kind() != reflect.Slice {
	//	return "", errors.New("必须为切片输入")
	//}

	s1 := []string{}
	switch list.(type) {
	case []int:
		for _, v := range list.([]int) {
			s1 = append(s1, strconv.Itoa(v))
		}
	case []int64:
		for _, v := range list.([]int64) {
			s1 = append(s1, strconv.Itoa(int(v)))
		}
	case []int32:
		for _, v := range list.([]int32) {
			s1 = append(s1, strconv.Itoa(int(v)))
		}
	default:
		return "", errors.New("不支持的输入内型")
	}

	return strings.Join(s1, ","), nil
}

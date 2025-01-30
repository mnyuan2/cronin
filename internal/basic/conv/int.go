package conv

import (
	"errors"
	"strconv"
	"strings"
)

type ints struct{}

func Ints() *ints {
	return &ints{}
}

// 转换为字符串
func (i *ints) String(val int) string {
	return strconv.Itoa(val)
}

// 将字符串转换为整数
func (i *ints) Parse(val string) (int, error) {
	if val == "" {
		return 0, nil
	}

	result, err := strconv.ParseInt(val, 10, 32)
	if err != nil {
		return 0, err
	}
	return int(result), nil
}

// 将任意类型转换为整数
func (i *ints) ParseAny(val interface{}) (int, error) {
	switch v := val.(type) {
	case int32:
		return int(v), nil
	case int:
		return v, nil
	case int64:
		return int(val.(int64)), nil
	case string:
		return strconv.Atoi(val.(string))
	case []uint8:
		return strconv.Atoi(string(val.([]uint8)))
	case float64:
		return int(val.(float64)), nil
	case float32:
		return int(val.(float32)), nil
	}
	return 0, errors.New("未支持的输入类型")
}

// 将逗号隔开的字符串（“1,2,3,4”）转换为切片
func (i *ints) Slice(val string) ([]int, error) {
	var out []int
	if val == "" {
		return out, nil
	}

	for _, v := range strings.Split(val, ",") {
		iv, err := i.Parse(v)
		if err != nil {
			return out, err
		}
		out = append(out, iv)
	}
	return out, nil
}

// 将逗号隔开的字符串（“1,2,3,4”）转换为[]interface{}
func (i *ints) ISlice(val string) ([]interface{}, error) {
	v, err := i.Slice(val)
	if err != nil {
		return nil, err
	}
	return i.Slice2I(v), nil
}

// 将切片转换为指定符号隔开的字符串(默认为逗号) “1,2,3,4”
func (i *ints) Join(s []int, split ...string) string {
	if len(s) <= 0 {
		return ""
	}

	join := ","
	// 取第一个
	if len(split) > 0 {
		join = split[0]
	}

	s1 := make([]string, 0, len(s))
	for _, v := range s {
		s1 = append(s1, i.String(v))
	}
	return strings.Join(s1, join)
}

// 将[]int对象转换为[]interface{}
func (i *ints) Slice2I(s []int) []interface{} {
	out := make([]interface{}, 0, len(s))
	for _, v := range s {
		out = append(out, v)
	}
	return out
}

// 元素是否包含
func (*ints) Contains(s []int, val int) (index int) {
	index = -1
	for i := 0; i < len(s); i++ {
		if s[i] == val {
			index = i
			return
		}
	}
	return
}

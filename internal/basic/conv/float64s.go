package conv

import (
	"errors"
	"fmt"
	"strconv"
)

type float64s struct {
}

func Float64s() *float64s {
	return &float64s{}
}

// 分到元 转换
func (f *float64s) FeeToYuan(fee float64) float64 {
	return fee / 100
}

// 元到分 转换
func (f *float64s) YuanToFee(yuan float64) float64 {
	return yuan * 100
}

// 将浮点型转换为字符串
// scale: 保留的位数,默认取2位
func (f *float64s) ToString(val float64, scales ...int) string {
	scale := 2
	if len(scales) > 0 {
		scale = scales[0]
	}

	return fmt.Sprintf("%."+strconv.Itoa(scale)+"f", val)
}

func (m *float64s) ParseAny(val any) (float64, error) {
	switch val.(type) {
	case int32:
		return float64(val.(int32)), nil
	case int:
		return float64(val.(int)), nil
	case int64:
		return float64(val.(int64)), nil
	case string:
		if val == "" {
			return 0, nil
		}
		result, err := strconv.ParseFloat(val.(string), 64)
		if err != nil {
			return 0, err
		}
		return result, nil
	case float64:
		return val.(float64), nil
	case float32:
		return float64(val.(float32)), nil
	}
	return 0, errors.New("未支持的输入类型")
}

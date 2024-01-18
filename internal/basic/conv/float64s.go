package conv

import (
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

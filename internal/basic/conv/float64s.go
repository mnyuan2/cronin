package conv

type float64s struct {

}

func Float64s()*float64s{
	return &float64s{}
}

// 分到元 转换
func (f *float64s)FeeToYuan(fee float64)float64{
	return fee / 100
}

// 元到分 转换
func (f *float64s)YuanToFee(yuan float64)float64{
	return yuan * 100
}
package conv

import (
	"fmt"
	"regexp"
	"strings"
	"unicode"
)

type Str struct {
	sep string
}

func NewStr() *Str {
	return &Str{
		sep: ",",
	}
}

// 是否包含数字
func (m *Str) IsNumber(s string) bool {
	re, _ := regexp.MatchString(`^[\+-]?\d+$`, s)
	return re
}

// 是否包含字母和数字
func (m *Str) IsLettersAndNumbers(str string) bool {
	return regexp.MustCompile(`^[A-Za-z0-9]+$`).MatchString(str)
}

// 检测字符串必须是同时包含字母和数字的组合
func (m *Str) ItIsLettersAndNumbers(str string) bool {
	if !m.IsLettersAndNumbers(str) {
		return false
	}

	letters := regexp.MustCompile(`[A-Za-z]+`).MatchString(str)
	numbers := regexp.MustCompile(`[0-9]+`).MatchString(str)

	if letters && numbers {
		return true
	}
	return false
}

// 是否包含中文
func (m *Str) IsChinese(str string) bool {
	for _, v := range str {
		if unicode.Is(unicode.Han, v) {
			return true
		}
	}
	return false
}

// 分切字符串
func (m *Str) Slice(val string, out []any) error {
	if val == "" {
		return nil
	}
	fmt.Println("输入", out)

	for _, v := range strings.Split(val, m.sep) {
		fmt.Println(v)
		//iv, err := i.Parse(v)
		//if err != nil {
		//	return out, err
		//}
		//out = append(out, iv)
	}
	return nil
}

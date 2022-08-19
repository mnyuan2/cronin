package conv

import (
	"regexp"
	"unicode"
)

type strings struct {
}

func Strings() *strings {
	return &strings{}
}

// 是否包含数字
func (m *strings) IsNumber(s string) bool {
	re, _ := regexp.MatchString(`^[\+-]?\d+$`, s)
	return re
}

// 是否包含字母和数字
func (m *strings) IsLettersAndNumbers(str string) bool {
	return regexp.MustCompile(`^[A-Za-z0-9]+$`).MatchString(str)
}

// 检测字符串必须是同时包含字母和数字的组合
func (m *strings) ItIsLettersAndNumbers(str string) bool {
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
func (m *strings) IsChinese(str string) bool {
	for _, v := range str {
		if unicode.Is(unicode.Han, v) {
			return true
		}
	}
	return false
}

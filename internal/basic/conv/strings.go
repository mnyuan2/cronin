package conv

import (
	"encoding/json"
	"errors"
	"reflect"
	"regexp"
	"strconv"
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

// 是否包含指定字符
func (m *Str) Contains(s, substr string) bool {
	return strings.Contains(s, substr)
}

// 分切字符串
func (m *Str) Slice(val string, out interface{}) (err error) {
	if val == "" {
		return nil
	}
	t := reflect.TypeOf(out)
	if t.Kind() != reflect.Ptr || t.Elem().Kind() != reflect.Slice {
		return errors.New("请指定切片指针out")
	}
	// 元素格式转换
	var itemVal func(v string) (reflect.Value, error)
	switch t.Elem().Elem().Kind() {
	case reflect.String:
		itemVal = func(v string) (reflect.Value, error) {
			return reflect.ValueOf(v), nil
		}
	case reflect.Int:
		itemVal = func(v string) (reflect.Value, error) {
			val, err := strconv.Atoi(v)
			return reflect.ValueOf(val), err
		}
	case reflect.Int32:
		itemVal = func(v string) (reflect.Value, error) {
			val, err := strconv.Atoi(v)
			return reflect.ValueOf(int32(val)), err
		}
	default:
		return errors.New("out 不支持的元素类型")
	}

	v := reflect.ValueOf(out).Elem()
	for _, item := range strings.Split(val, m.sep) {
		value, err := itemVal(item)
		if err != nil {
			return err
		}
		v = reflect.Append(v, value)
	}

	reflect.ValueOf(out).Elem().Set(v)
	return nil
}

// 将任意类型转换为字符串, 如果是普通类型, 直接转换, 复杂类型序列化
func (m *Str) ToString(value interface{}) string {
	var out string
	if value == nil {
		return out
	}

	switch value.(type) {
	case uint:
		it := value.(uint)
		out = strconv.Itoa(int(it))
	case int8:
		it := value.(int8)
		out = strconv.Itoa(int(it))
	case uint8:
		it := value.(uint8)
		out = strconv.Itoa(int(it))
	case int16:
		it := value.(int16)
		out = strconv.Itoa(int(it))
	case uint16:
		it := value.(uint16)
		out = strconv.Itoa(int(it))
	case int32:
		it := value.(int32)
		out = strconv.Itoa(int(it))
	case uint32:
		it := value.(uint32)
		out = strconv.Itoa(int(it))
	case int64:
		it := value.(int64)
		out = strconv.FormatInt(it, 10)
	case uint64:
		it := value.(uint64)
		out = strconv.FormatUint(it, 10)
	case int:
		it := value.(int)
		out = strconv.Itoa(it)
	case float32:
		ft := value.(float32)
		out = strconv.FormatFloat(float64(ft), 'f', -1, 64)
	case float64:
		ft := value.(float64)
		out = strconv.FormatFloat(ft, 'f', -1, 64)
	case string:
		out = value.(string)
	case []byte:
		out = string(value.([]byte))
	default:
		newValue, _ := json.Marshal(value)
		out = string(newValue)
	}

	return out
}

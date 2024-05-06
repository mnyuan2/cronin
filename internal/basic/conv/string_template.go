package conv

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gitlab.com/metakeule/fmtdate"
	"reflect"
	"text/template"
	"time"
)

type StringTemplate struct {
	funcs      template.FuncMap // 函数
	leftDelim  string           // 左界定符
	rightDelim string           // 右界定符
	params     any              // 参数
}

func NewStringTemplate() *StringTemplate {
	return &StringTemplate{
		funcs: map[string]any{},
	}
}

// 默认模板
func DefaultStringTemplate() *StringTemplate {
	return &StringTemplate{
		leftDelim:  "[[",
		rightDelim: "]]",
		funcs: map[string]any{
			// json 编码
			"jsonString": func(val any) any {
				v := reflect.ValueOf(val)
				switch v.Kind() {
				case reflect.Map, reflect.Slice:
					value, _ := json.Marshal(val)
					return string(value)
				default:
					return val
				}
			},
			// json 编码2次
			"jsonString2": func(val any) any {
				v := reflect.ValueOf(val)
				switch v.Kind() {
				case reflect.Map, reflect.Slice:
					value, _ := json.Marshal(val)
					value = bytes.ReplaceAll(value, []byte(`"`), []byte(`\"`))
					return string(value)
				case reflect.String:
					value := bytes.ReplaceAll([]byte(val.(string)), []byte(`"`), []byte(`\"`))
					return value
				default:
					return val
				}
			},
			// 格式化 时间/日期
			//  参数1：string format 格式表达式，默认 YYYY-MM-DD
			//  参数2：int64 timestamp 时间戳，默认 Unix 时间戳
			"date": func(param ...any) (date string, err error) {
				var timestamp *int64
				var format *string
				l := len(param)
				if l > 0 {
					temp := fmt.Sprintf("%v", param[0])
					format = &temp
				}
				if l > 1 && param[1] != nil {
					if temp, err := Int64s().ParseAny(param[1]); err != nil {
						return "", err
					} else {
						timestamp = &temp
					}
				}
				if format == nil {
					temp := "YYYY-MM-DD"
					format = &temp
				}
				if timestamp != nil {
					date = fmtdate.Format(*format, time.Unix(*timestamp, 0))
				} else {
					date = fmtdate.Format(*format, time.Now())
				}
				return date, err
			},
		},
	}
}

// AddFunc 添加处理函数
func (t *StringTemplate) AddFunc(name string, f any) *StringTemplate {
	t.funcs[name] = f
	return t
}

// SetParam 设置参数
func (t *StringTemplate) SetParam(params map[string]any) *StringTemplate {
	t.params = params
	return t
}

// SetDelim 设置边界符
func (t *StringTemplate) SetDelim(left, right string) *StringTemplate {
	t.leftDelim = left
	t.rightDelim = right
	return t
}

// Execute 模板执行
func (t *StringTemplate) Execute(text []byte) (newStr []byte, err error) {
	if text == nil {
		return nil, err
	}
	temp := template.New("tmpl")
	if len(t.funcs) > 0 {
		temp.Funcs(t.funcs)
	}

	if t.rightDelim != "" && t.leftDelim != "" {
		temp.Delims(t.leftDelim, t.rightDelim)
	} else if len(t.rightDelim+t.leftDelim) > 0 {
		return nil, fmt.Errorf("边界符为必填")
	}

	// 创建模板
	_tmpl, err := temp.Parse(string(text))
	if err != nil {
		return nil, fmt.Errorf("解析模板失败,%w", err)
	}

	// 应用模板到数据
	buf := bytes.NewBuffer([]byte{})
	if err = _tmpl.Execute(buf, t.params); err != nil {
		return nil, fmt.Errorf("模板执行失败,%w", err)
	}
	// 获取替换后的字符串
	//result := buf.String()
	return buf.Bytes(), nil
}

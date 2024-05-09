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
			// 兼容 null 参数，转换为 nil
			"null": func() any {
				return nil
			},
			// 获取时间
			//  参数1：string duration 持续时间字符串，示例 "300ms", "-1.5h" 或 "2h45m". 有效的时间单位是 "ns", "us" (or "µs"), "ms", "s", "m", "h".
			"time": func(param ...any) (ti time.Time, err error) { // 1.相对时间、2.时间戳、3.时间字符串；
				l := len(param)
				dur := time.Duration(0)
				if l > 0 && param[0] != nil && param[0] != "" {
					param1, ok := param[0].(string)
					if !ok {
						return time.Time{}, fmt.Errorf("time param 1 not string")
					}
					dur, err = time.ParseDuration(param1)
					if err != nil {
						return time.Time{}, fmt.Errorf("time param 1 error, %w", err)
					}
				}
				// 更多参数，后面根据情况扩展
				return time.Now().Add(dur), nil
			},
			// 格式化 时间/日期
			//  参数1：string format 格式表达式，默认 YYYY-MM-DD hh:mm:ss
			//  参数2：object time 时间对象，默认 当前时间
			"date": func(param ...any) (date string, err error) {
				var format *string
				l, t := len(param), time.Now()
				if l > 0 && param[0] != nil && param[0] != "" {
					temp, ok := param[0].(string)
					if !ok {
						return "", fmt.Errorf("date param 1 not string")
					}
					format = &temp
				}
				if l > 1 && param[1] != nil {
					if ti, ok := param[1].(time.Time); !ok {
						return "", fmt.Errorf("date param 2 not Time")
					} else {
						t = ti
					}
				}
				if format == nil {
					temp := "YYYY-MM-DD hh:mm:ss"
					format = &temp
				}
				date = fmtdate.Format(*format, t)
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

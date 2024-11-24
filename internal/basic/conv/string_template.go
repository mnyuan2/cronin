package conv

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Knetic/govaluate"
	jsoniter "github.com/json-iterator/go"
	"gitlab.com/metakeule/fmtdate"
	"net/url"
	"reflect"
	"regexp"
	"strconv"
	"strings"
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
			//  废弃 建议使用 json_encode
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
			//  后面废弃 使用 json_encode
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
			// 解码 json 字符串
			//  返回可能是切片、也可能是map
			"json_decode": func(str string) (any, error) {
				data := map[string]any{}
				err := jsoniter.UnmarshalFromString(str, &data)
				return data, err
			},
			// 编码为 json 字符串
			"json_encode": func(val any) (string, error) {
				v := reflect.ValueOf(val)
				switch v.Kind() {
				case reflect.Map, reflect.Slice:
					return jsoniter.MarshalToString(val)
				case reflect.String:
					value := strings.ReplaceAll(val.(string), `"`, `\"`)
					return value, nil
				default:
					return fmt.Sprintf("%v", val), nil
				}
			},
			// encodeURIComponent 遵循 RFC3986 url编码
			"rawurlencode": func(param string) string {
				str := url.QueryEscape(param)
				str = strings.ReplaceAll(str, "+", "%20")
				return str
			},
			// errorf 主动抛出错误
			//  format 错误值
			//  args 格式参数
			"errorf": func(format string, args ...any) (struct{}, error) {
				return struct{}{}, fmt.Errorf(format, args...)
			},
			// 兼容 null 参数，转换为 nil
			"null": func() any {
				return nil
			},
			"float64": func(val any) (float64, error) {
				return Float64s().ParseAny(val)
			},
			"string": func(any any) string {
				return fmt.Sprintf("%v", any)
			},
			// 合并 slice，返回新的 slice
			//  多个参数的类型定义必须一致
			"append_slice": func(slice ...any) (any, error) {
				values := reflect.ValueOf(slice)
				if values.Index(0).Elem().Kind() != reflect.Slice {
					return nil, fmt.Errorf("param 1 not slice")
				}

				newValues := reflect.MakeSlice(values.Index(0).Elem().Type(), 0, 0) // 创建一个空的原元素
				for i := 0; i < values.Len(); i++ {
					value := values.Index(0)
					if newValues.Kind() != values.Kind() {
						return nil, fmt.Errorf("param %v type inconsistency", i+1)
					}
					newValues = reflect.AppendSlice(newValues, value.Elem())
				}

				list := newValues.Interface()
				return list, nil
			},
			// 此方法目前处于试验阶段，应该找一种更通用的方法来声明元素。
			"make": func(t string) any {
				switch t {
				case "int":
					return int(0)
				case "[]map[string]any":
					return []map[string]any{}
				case "[]map[string]string":
					return []map[string]string{}
				default:
					return nil
				}
			},
			"append": func(slice any, elems ...any) (any, error) {
				values := reflect.ValueOf(slice)
				if values.Kind() != reflect.Slice {
					return nil, fmt.Errorf("method append param 1 not slice")
				}
				k := values.Type().Elem()
				for i, item := range elems {
					val := reflect.ValueOf(item)
					if k != val.Type() {
						return nil, fmt.Errorf("method append param %v type inconsistency(%s != %s)", i+2, k.String(), val.Type().String())
					}
					values = reflect.Append(values, val)
				}

				list := values.Interface()
				return list, nil
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
			// 字符串查最左找计算，并将结果替换原字符串
			//  参数1: string raw 原始字符串
			//  参数2: string regex 正则表达式，提取字符串中的数字，示例：`(\d+)(\D*$)`
			//  参数3: string 对数字进行计算的公式，示例 +1
			//  返回： 计算结果替换后的字符串
			"str_replace_calc": func(raw, regex, expr string) (str string, err error) {
				matches := regexp.MustCompile(regex).FindStringSubmatch(raw)
				if matches == nil {
					return raw, nil
				}
				exp, _ := govaluate.NewEvaluableExpression(matches[1] + expr)
				val, err := exp.Evaluate(nil)
				if err != nil {
					return "", fmt.Errorf("运算执行错误，%s", err.Error())
				}
				matches[1] = strconv.FormatFloat(val.(float64), 'f', -1, 64)
				str = raw[:len(raw)-len(matches[0])] + matches[1] + matches[2]
				return str, nil
			},
			// 字符串查找
			//  @param string raw 原始字符串
			//  @param string regex 正则匹配表达式
			//  @param string fields 输出结果keys, 多个key逗号相隔 与匹配结果顺序一致
			//  @return map[string]any
			"str_find": func(raw string, regex string) (out []string, err error) {
				matches := regexp.MustCompile(regex).FindStringSubmatch(raw)
				return matches, nil
			},
			// 字符串切割  str_slice_filter
			//  参数1: string str 原始字符串
			//  参数2: string sep 分隔符
			"str_split": func(str, sep string) []string {
				list := strings.Split(str, sep)
				return list
			},

			// 切片过滤
			//  参数1: []string data 原始切片数据
			//  参数3: string filter 过滤字符串，正则表示；符合条件的元素会被过滤掉
			"slice_filter": func(data any, filter string) (any, error) {
				// 目前仅支持字符串切片，后期根据需求扩展
				list, ok := data.([]string)
				if !ok {
					return data, nil
				}
				regex, err := regexp.Compile(filter)
				if err != nil {
					return nil, fmt.Errorf("过滤正则输入有误，%s", err.Error())
				}
				newData := []string{}
				for _, v := range list {
					if regex.MatchString(v) {
						continue
					}
					newData = append(newData, v)
				}
				return newData, nil
			},

			// slice 组合成 map
			//  @param []string list 原始切片数据
			//  @param string keys 字段集合，特殊字段："k"常规字段输入、""空值表示忽略元素、"k:v"分号后面为默认值，原元素不存在或为空会启用。
			"slice_combine": func(list []string, keys ...string) map[string]string {
				// 目前仅支持字符串切片，后期根据需求扩展
				l := len(list) - 1
				out := map[string]string{}
				for i, key := range keys {
					if key == "" {
						continue
					}
					kv := strings.Split(key, ":")
					val := ""
					if len(kv) > 1 {
						val = kv[1]
					}
					if i <= l {
						if v := list[i]; v != "" && v != "0" {
							val = v
						}
					}
					out[kv[0]] = val
				}
				return out
			},

			// map value 字符串按指定字符切割，并将结果生成新的map
			//  参数1: map data 原始map
			//  参数2: string sep 分隔符
			//  参数3: string keys 指定匹配字段，默认所有字段（注意匹配字段值必须为字符串）
			"map_split": func(data any, sep string, keys ...any) (any, error) {
				last := reflect.ValueOf(data)
				// 确保原始值是一个 map
				if last.Kind() != reflect.Map {
					return nil, fmt.Errorf("map_slice 参数 1 必须为map")
				}
				if len(keys) == 0 {
					for _, k := range last.MapKeys() {
						keys = append(keys, k.Interface())
					}
				}

				values := reflect.MakeSlice(reflect.SliceOf(last.Type()), 0, 0)
				for _, key := range keys {
					f, v := last.MapIndex(reflect.ValueOf(key)), ""
					if f.Type().Kind() == reflect.Interface {
						v = f.Elem().String()
					} else {
						v = f.String()
					}
					if ok := strings.Contains(v, sep); !ok {
						continue
					}
					list := strings.Split(v, sep)
					for _, item := range list {
						// 克隆 map， 并将指定字段替换。
						cloneValue := reflect.MakeMap(last.Type())
						for _, key := range last.MapKeys() {
							value := last.MapIndex(key)
							cloneValue.SetMapIndex(key, value)
						}
						cloneValue.SetMapIndex(reflect.ValueOf(key), reflect.ValueOf(item))
						values = reflect.Append(values, cloneValue)
						last = cloneValue
					}
				}
				if values.Len() == 0 {
					values = reflect.Append(values, last)
				}
				return values.Interface(), nil
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
		switch e := err.(type) {
		case template.ExecError:
			if e2 := errors.Unwrap(e.Unwrap()); e2 != nil {
				return nil, e2
			}
			return nil, e.Unwrap()
		default:
			return nil, e
		}
	}
	// 获取替换后的字符串
	//result := buf.String()
	return buf.Bytes(), nil
}

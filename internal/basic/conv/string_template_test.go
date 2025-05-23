package conv

import (
	"bytes"
	"cron/internal/pb"
	"encoding/json"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"gitlab.com/metakeule/fmtdate"
	"net/url"
	"reflect"
	"regexp"
	"strings"
	"testing"
	"text/template"
	"time"
)

func TestTemplate(t *testing.T) {
	var err error
	str := []byte(`{"http":{"method":"POST","url":"https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=xx","body":"{\n    \"msgtype\": \"text\",\n    \"text\": {\n        \"content\": \"时间：[[log.create_dt]]\\n任务 [[config.name]]执行[[log.status_name]]了 \\n耗时[[log.duration]]秒\\n响应：[[log.body]]\",\n        \"mentioned_mobile_list\": [[user.mobile]]\n    }\n}","header":[{"key":"","value":""}]}}`)

	// 提取模板变量
	// 重组临时变量，默认置空，有效的写入新值
	// 方案1 解析前监测双引号等关键、方案2让低层兼容
	args := map[string]string{
		"env":                  "测试环境",
		"config.name":          "xx任务",
		"config.protocol_name": "sql脚本",
		"log.status_name":      "成功",
		"log.status_desc":      "success",
		"log.body": strings.ReplaceAll(`Get "http://baidu.com": EOF
`, `"`, `\\\"`),
		"log.duration":  "3.2s",
		"log.create_dt": "2023-01-01 11:12:59",
		"user.username": "",
		"user.mobile":   "",
	}

	mobles := []string{"01987654321", "12345678910"}
	args["user.mobile"], err = jsoniter.MarshalToString(mobles)
	if err != nil {
		t.Fatal("数据转义错误", err.Error())
	}
	args["user.mobile"] = strings.ReplaceAll(args["user.mobile"], `"`, `\"`)

	username := []string{"大王", "二王"}
	args["user.username"], err = jsoniter.MarshalToString(username)
	if err != nil {
		t.Fatal("数据转义错误", err.Error())
	}
	args["user.username"] = strings.ReplaceAll(args["user.username"], `"`, `\"`)

	// 变量替换
	for k, v := range args {
		str = bytes.Replace(str, []byte("[["+k+"]]"), []byte(v), -1)
	}
	//fmt.Println(string(str))
	temp := &pb.SettingMessageTemplate{Http: &pb.CronHttp{Header: []*pb.KvItem{}}}
	if err := jsoniter.Unmarshal(str, temp); err != nil {
		t.Fatal(err, "解析错误")
	}
	fmt.Println(temp)
}

func TestTemplateV2(t *testing.T) {
	a := map[string]any{}
	a["a"] = ""
	fmt.Println(a)
	a["a"] = 0
	fmt.Println(a)
	a["a"] = map[string]string{"b": "BB"}
	fmt.Println(a)

	input := `切片: [[.name]] --> [[jsonString .name]] --> [[jsonString2 .name]]
数组：[[.c]] --> [[jsonString .c]] --> [[jsonString2 .c]] --> [[.c.cc]]
常量：age:[[.age]] | sex:[[.sex]] [[jsonString2 .sex]] [[jsonString .sex]] [[.b]]`
	paramStr := `{"sex": "男", "age": 180, "name": ["title2", "title1", 25], "c":{"cc":"CC"}}`

	// 定义一个 map 用于存储解析后的数据
	data := map[string]any{}
	if err := json.Unmarshal([]byte(paramStr), &data); err != nil {
		fmt.Println("解析JSON失败:", err)
		return
	}
	// 自定义模板函数
	f := template.FuncMap{
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
		// json 编码两次
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
	}
	// 创建模板
	_tmpl, err := template.New("tmpl").
		Funcs(f).
		Delims("[[", "]]").
		Parse(input)
	if err != nil {
		fmt.Println("解析模板失败:", err)
		return
	}

	// 应用模板到数据
	buf := bytes.NewBuffer([]byte{})
	err = _tmpl.Execute(buf, data)
	if err != nil {
		fmt.Println("模板执行失败:", err)
		return
	}
	// 获取替换后的字符串
	result := buf.Bytes()
	fmt.Printf(`%s\n`, result)
}

func TestParse(t *testing.T) {
	str := `{"cmd":{"host":{"id":-1,"ip":"","port":"","type":"","user":"","secret":""},"type":"bash","origin":"local","statement":{"git":{"ref":"","path":[""],"owner":"","link_id":0,"project":""},"type":"","local":"echo [[.a]] '\\n' [[.b]] [[string .b]] [[.b.b1]]","is_batch":0}},"git":{"events":[],"link_id":0},"rpc":{"addr":"","body":"","proto":"","action":"","header":[],"method":"GRPC","actions":[]},"sql":{"driver":"mysql","origin":"local","source":{"id":0,"port":"","title":"","database":"","hostname":"","password":"","username":""},"interval":0,"statement":[],"err_action":1,"err_action_name":""},"http":{"url":"","body":"","header":[{"key":"","value":""}],"method":"GET"},"jenkins":{"name":"","params":[{"key":"","value":""}],"source":{"id":0}}}`
	str = "[[date `YYYY-MM-DD+hh:mm:ss` (time `-10m`)]]"
	param := `{"a":"A","b":{"b1":"B1","B2":22},"c":3}`

	varParams := map[string]any{}
	if param != "" {
		if er := jsoniter.UnmarshalFromString(param, &varParams); er != nil {
			t.Fatal(er)
		}
	}

	b, e := DefaultStringTemplate().SetParam(varParams).Execute([]byte(str))
	if e != nil {
		t.Fatal(e)
	}
	fmt.Println(string(b))
}

// str_replace_calc 模板函数测试
func TestStrReplaceCalc(t *testing.T) {
	str := "[[if .tag_fixed]][[printf `%v%v` `release_` .tag_fixed]][[else if .tag_expr]][[str_replace_calc .raw_content `(\\d+)(\\D*$)` `+1`]][[end]]"
	p := map[string]any{
		"raw_content": "release_v3.5.87.2",
		//"tag_fixed":   "v3.2.1.0", // 固定值替换
		"tag_expr": "+1", // 运算规则替换
	}
	b, e := DefaultStringTemplate().SetParam(p).Execute([]byte(str))
	if e != nil {
		t.Fatal(e)
	}
	fmt.Println(str, "\n", p["raw_content"], "|", string(b), "|")
}

// 字符串切割并过滤不可见字符串元素
func TestStrSliceFilter(t *testing.T) {
	strs := []string{
		"",
		"   ",         // 只包含空格
		"\t\n",        // 只包含制表符和换行符
		"Hello Word!", // 包含可见字符
		"  Text  ",    // 也包含可见字符
	}
	for i, str := range strs {
		b, e := DefaultStringTemplate().SetParam(map[string]any{"data": str}).Execute([]byte("[[json_encode (slice_filter (str_split .data ` `) `^\\s*$`)]]"))
		if e != nil {
			t.Fatal(e)
		}

		fmt.Println(i, string(b))
	}
}

func TestStrFindMap(t *testing.T) {
	strs := []string{
		"https://gitee.com/mnyuan/cronin/pulls/15",
		"cronin/hotfix/user_3",
		"cronin/hotfix/user_3{serA,serB}",
	}
	for i, str := range strs {
		temp := DefaultStringTemplate().SetParam(map[string]any{"data": str})
		if i == 0 {
			b, e := temp.Execute([]byte("[[json_encode (str_find_map .data `https://gitee.com/(.+)/([^/]+)/pulls/(\\d+)` `owner,repo,number,type:pr`)]]"))
			if e != nil {
				t.Fatal(e)
			}
			fmt.Println(i, string(b))
		} else {
			b, e := temp.Execute([]byte("[[json_encode (str_find_map .data `([^/]+)(?:.*)(?:/|\\{(.*)\\})` `repo,service,type:jenkins`)]]"))
			if e != nil {
				t.Fatal(e)
			}
			fmt.Println(i, string(b))
		}
	}

	//resMap["type"] = "pr"
	//fmt.Println(resMap)
}

func TestTemplateJson(t *testing.T) {
	data := map[string]any{
		"a": map[string]string{"a1": "A1", "a2": "a2"},
		"b": []string{"B1", "B\"2"},
		"c": "C",
	}

	b, err := json.Marshal(data)
	if err != nil {
		t.Fatal("解析错误1", err)
	}
	fmt.Println(string(b))

	b2 := bytes.ReplaceAll(b, []byte(`"`), []byte(`\"`))
	fmt.Println(string(b2))

	b, err = json.Marshal(string(b))
	if err != nil {
		t.Fatal("解析错误2", err)
	}
	fmt.Println(string(b))

}

// 模板抛出错误
func TestTemplateErr(t *testing.T) {
	// 我可以将这个数据作为输入，通过模板语法确定某个变量的值，从而确定结果是否符合预期。
	tmpl := "[[if ne .a 0]]\n      [[- errorf `错误 %s %v` `数值` 5 -]]\n[[end]]"

	//buf := bytes.NewBuffer([]byte{})
	param := map[string]any{"a": 1}

	buf, err := DefaultStringTemplate().SetParam(param).Execute([]byte(tmpl))
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("模板输出", string(buf))
	// 可以根据这个输出，如果空白就是错误，仅接收 true 这个特定结果 为成功。
	fmt.Println(param)
}

// 模板方式 json 响应处理
func TestTemplateJsonResult(t *testing.T) {
	// 我可以将这个数据作为输入，通过模板语法确定某个变量的值，从而确定结果是否符合预期。
	str := `{"code":"000000","message":"成功","data":{"list":[],"page":{"size":0,"page":1,"total":0}}}`
	tmpl := "[[string 5]]"

	//buf := bytes.NewBuffer([]byte{})
	param := map[string]any{"result": str}

	buf, err := DefaultStringTemplate().SetParam(param).Execute([]byte(tmpl))
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("模板输出", string(buf))
	// 可以根据这个输出，如果空白就是错误，仅接收 true 这个特定结果 为成功。
	fmt.Println(param)

}

func TestFuncs(t *testing.T) {
	// 我希望获得30天前的时间，假设语法 {{}}
	const tmpl = "Now: {{ Now }} \n\r {{date `YYYY-MM-DD` (time `-720h`)}} \n\r {{rawurlencode (date `YYYY-MM-DD hh:mm:ss` (time `-23h`))}}"
	temp := template.Must(template.New("test").Funcs(template.FuncMap{
		"Now": func() time.Time { return time.Now() },
		"null": func() any {
			return nil
		},
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
			return time.Now().Add(dur), nil
		},
		"date": func(param ...any) (date string, err error) { // 参数：1.格式、2.时间对象
			var format *string
			l, t := len(param), time.Now()
			if l > 0 && param[0] != nil {
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
		"rawurlencode": func(param string) string {
			str := url.QueryEscape(param)
			str = strings.ReplaceAll(str, "+", "%20")
			return str
		},
	}).Parse(tmpl))

	buf := bytes.NewBuffer([]byte{})

	err := temp.Execute(buf, nil)
	if err != nil {
		panic(err)
	}
	fmt.Println(buf.String())
}

func TestName(t *testing.T) {
	str := "http://abc.com"
	//str := "Valid_Address123"
	// 定义匹配规则的正则表达式
	re := regexp.MustCompile(`^[a-zA-Z][\w-]{1,}[a-zA-Z0-9]$`)

	// 使用正则表达式进行匹配
	is := re.MatchString(str)
	fmt.Println("结果", is)

}

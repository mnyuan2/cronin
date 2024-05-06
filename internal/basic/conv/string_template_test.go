package conv

import (
	"bytes"
	"cron/internal/pb"
	"encoding/json"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"gitlab.com/metakeule/fmtdate"
	"reflect"
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
	param := `{"a":"A","b":{"b1":"B1","B2":22},"c":3}`

	varParams := map[string]any{}
	if param != "" {
		if er := jsoniter.UnmarshalFromString(param, varParams); er != nil {
			t.Fatal(er)
		}
	}

	b, e := DefaultStringTemplate().SetParam(varParams).Execute([]byte(str))
	if e != nil {
		t.Fatal(e)
	}
	fmt.Println(string(b))
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

func TestFuncs(t *testing.T) {
	const tmpl = `Now: {{ Now }} \n\r {{date}}`
	temp := template.Must(template.New("test").Funcs(template.FuncMap{
		"Now": func() time.Time { return time.Now() },
		"date": func(param ...any) (date string, err error) {
			var sec *int64
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
					sec = &temp
				}
			}
			if format == nil {
				temp := "YYYY-MM-DD hh:mm:ss"
				format = &temp
			}
			if sec != nil {
				date = fmtdate.Format(*format, time.Unix(*sec, 0))
			} else {
				date = fmtdate.Format(*format, time.Now())
			}
			return date, err
		},
	}).Parse(tmpl))

	buf := bytes.NewBuffer([]byte{})

	err := temp.Execute(buf, nil)
	if err != nil {
		panic(err)
	}
	fmt.Println(buf.String())
}

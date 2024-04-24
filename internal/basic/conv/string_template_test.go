package conv

import (
	"bytes"
	"encoding/json"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"reflect"
	"testing"
	"text/template"
)

func TestTemplateV2(t *testing.T) {
	a := map[string]any{}
	a["a"] = ""
	fmt.Println(a)
	a["a"] = 0
	fmt.Println(a)
	a["a"] = map[string]string{"b": "BB"}
	fmt.Println(a)

	input := `切片: [[.name]] --> [[string .name]] 
数组：[[.c]] --> [[string .c]] --> [[.c.cc]]
常量：age:[[.age]] | sex:[[.sex]]  [[.b]] [[b]]`
	paramStr := `{"sex": "男", "age": 180, "name": ["title2", "title1", 25], "c":{"cc":"CC"}}`

	// 定义一个 map 用于存储解析后的数据
	data := map[string]any{}
	if err := json.Unmarshal([]byte(paramStr), &data); err != nil {
		fmt.Println("解析JSON失败:", err)
		return
	}
	// 自定义模板函数
	f := template.FuncMap{
		// 任何数据转义为字符串
		"string": func(val any) string {
			v := reflect.ValueOf(val)
			switch v.Kind() {
			case reflect.Map, reflect.Slice:
				value, _ := json.Marshal(val)
				return string(value)
			default:
				return fmt.Sprintf("%v", val)
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
	result := buf.String()
	fmt.Println(result)
}

func TestParse(t *testing.T) {
	str := `{"cmd": {"host": {"id": -1, "ip": "", "port": "", "type": "", "user": "", "secret": ""}, "type": "bash", "origin": "local", "statement": {"git": {"ref": "", "path": [""], "owner": "", "link_id": 0, "project": ""}, "type": "", "local": "echo [[.a]] '\\n' [[.b]] [[string .b]] [[.b.b1]]", "is_batch": 0}}, "git": {"events": [], "link_id": 0}, "rpc": {"addr": "", "body": "", "proto": "", "action": "", "header": [], "method": "GRPC", "actions": []}, "sql": {"driver": "mysql", "origin": "local", "source": {"id": 0, "port": "", "title": "", "database": "", "hostname": "", "password": "", "username": ""}, "interval": 0, "statement": [], "err_action": 1, "err_action_name": ""}, "http": {"url": "", "body": "", "header": [{"key": "", "value": ""}], "method": "GET"}, "jenkins": {"name": "", "params": [{"key": "", "value": ""}], "source": {"id": 0}}}`
	param := ``

	varParams := map[string]any{}
	if er := jsoniter.UnmarshalFromString(param, varParams); er != nil {
		t.Fatal(er)
	}

	b, e := DefaultStringTemplate().SetParam(varParams).Execute([]byte(str))
	if e != nil {
		t.Fatal(e)
	}
	fmt.Println(string(b))
}

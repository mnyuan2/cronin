package util

import (
	"fmt"
	"regexp"
	"runtime/debug"
	"strings"
)

// 异常详情
func PanicInfo(recover any) string {
	if recover != nil {
		errStr := fmt.Sprintf("%v", recover)
		debugStack := errStr + "\n"
		for i, v := range strings.Split(string(debug.Stack()), "\n") {
			if i < 6 {
				continue
			}
			if i%2 != 0 {
				debugStack += "	-->> " + v + "\n"
			} else {
				debugStack += strings.Trim(v, " 	")
			}
		}
		return debugStack
	}
	return ""
}

// 解析sql类型名称
func ParseSqlTypeName(sql string) string {
	reg1 := regexp.MustCompile("^(?:[\\n\\t]*)(\\S*)(?:\\s*)")
	result := reg1.FindAllStringSubmatch(sql, 1)
	if len(result) > 0 {
		return strings.ToUpper(result[0][1])
	}
	return ""
}

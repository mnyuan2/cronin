package util

import (
	"fmt"
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

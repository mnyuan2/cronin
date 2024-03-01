package biz

import (
	"bytes"
	"fmt"
	"testing"
)

func TestSqlParse(t *testing.T) {
	sql := `select 1; 
select 2; select 3;`

	list := bytes.Split([]byte(sql), []byte(";"))
	for i, item := range list {
		s := bytes.TrimSpace(item)
		if s != nil {
			fmt.Println(i, s, string(s))
		}
	}
	fmt.Println(list)
}

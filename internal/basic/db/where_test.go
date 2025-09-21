package db

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"
)

func TestFindInSet(t *testing.T) {
	os.Chdir("/go/src/incron")
	w := NewWhere()

	_sql := "SELECT * FROM `cron_setting` %WHERE "
	w.FindInSet("env", []string{"test", "prod"})

	where, args := w.Build()
	_sql = strings.Replace(_sql, "%WHERE", "WHERE "+where, -1)

	list := []map[string]any{}
	err := New(context.Background()).Raw(_sql, args...).Scan(&list).Error
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(list)
}

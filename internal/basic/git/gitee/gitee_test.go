package gitee

import (
	"cron/internal/basic/git"
	"fmt"
	"testing"
)

func TestUrl(t *testing.T) {
	conf := &git.Config{AccessToken: "e6a28b06d79d492f9809069d5550b436"}
	api := NewApiV5(conf)

	res, err := api.ReposContents("mnyuan", "cronin", "work/mysql.sql", "master")
	if err != nil {
		t.Fatalf(err.Error())
	}
	fmt.Println(string(res))
}

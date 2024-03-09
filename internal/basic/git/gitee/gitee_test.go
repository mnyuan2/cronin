package gitee

import (
	"context"
	"cron/internal/basic/git"
	"fmt"
	"testing"
)

func TestUrl(t *testing.T) {
	conf := &git.Config{AccessToken: "e6a28b06d79d492f9809069d5550b436"}
	api := NewApiV5(conf)
	handler := NewHandler(context.Background())
	res, err := api.ReposContents(handler, "mnyuan", "cronin", "work/mysql.sql", "master")
	if err != nil {
		t.Fatalf(err.Error())
	}
	fmt.Println(handler)
	fmt.Println(string(res))
}

func TestUser(t *testing.T) {
	conf := &git.Config{AccessToken: "e6a28b06d79d492f9809069d5550b436"}
	api := NewApiV5(conf)
	handler := NewHandler(context.Background())
	res, err := api.User(handler)
	if err != nil {
		t.Fatalf(err.Error())
	}
	fmt.Println(handler)
	fmt.Println(string(res))
}

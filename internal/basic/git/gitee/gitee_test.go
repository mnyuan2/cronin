package gitee

import (
	"context"
	"cron/internal/basic/git"
	"fmt"
	"testing"
)

var conf = &git.Config{AccessToken: "e6a28b06d79d492f9809069d5550b436"}

func TestUrl(t *testing.T) {
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
	api := NewApiV5(conf)
	handler := NewHandler(context.Background())
	res, err := api.User(handler)
	if err != nil {
		t.Fatalf(err.Error())
	}
	fmt.Println(handler)
	fmt.Println(string(res))
}

func TestApiV5_PullsReview(t *testing.T) {
	api := NewApiV5(conf)
	handler := NewHandler(context.Background())

	res, err := api.PullsReview(handler, &PullsReviewRequest{
		BaseRequest: BaseRequest{
			Owner: "mnyuan",
			Repo:  "cronin",
		},
		Number: 9,
		Force:  false,
	})

	fmt.Println(handler.String())
	if err != nil {
		t.Fatalf(err.Error())
	}
	fmt.Println(string(res))
}

func TestApiV5_PullsTest(t *testing.T) {
	api := NewApiV5(conf)
	handler := NewHandler(context.Background())

	res, err := api.PullsTest(handler, &PullsTestRequest{
		BaseRequest: BaseRequest{
			Owner: "mnyuan",
			Repo:  "cronin",
		},
		Number: 9,
		Force:  false,
	})

	fmt.Println(handler.String())
	if err != nil {
		t.Fatalf(err.Error())
	}
	fmt.Println(string(res))
}

// 合并分支
func TestPullsMerge(t *testing.T) {
	api := NewApiV5(conf)
	handler := NewHandler(context.Background())

	res, err := api.PullsMerge(handler, &PullsMergeRequest{
		BaseRequest: BaseRequest{
			Owner: "mnyuan",
			Repo:  "cronin",
		},
		Number:            9,
		MergeMethod:       "merge",
		PruneSourceBranch: false,
		Title:             "A",
		Description:       "B",
	})

	fmt.Println(handler.String())
	if err != nil {
		t.Fatalf(err.Error())
	}
	fmt.Println(string(res))
}

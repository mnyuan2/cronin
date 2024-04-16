package gitee

import (
	"context"
	"cron/internal/basic/git"
	"fmt"
	"testing"
)

var conf = &git.Config{AccessToken: "f5c13d72c3f68dd6c92bb82641c8a7c9"}

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

func TestApiV5_PullsCreate(t *testing.T) {
	api := NewApiV5(conf)
	handler := NewHandler(context.Background())

	res, err := api.PullsCreate(handler, &PullsCreateRequest{
		BaseRequest: BaseRequest{
			Owner: "mnyuan",
			Repo:  "cronin",
		},
		Head:                  "master",
		Base:                  "test",
		Title:                 "test demo",
		Body:                  "pr body .",
		MilestoneNumber:       0,
		Labels:                "",
		Issue:                 "",
		Assignees:             "",
		Testers:               "",
		AssigneesNumber:       0,
		TestersNumber:         0,
		RefPullRequestNumbers: "",
		PruneSourceBranch:     false,
		CloseRelatedIssue:     false,
		Draft:                 false,
		Squash:                false,
	})

	fmt.Println(handler.String())
	if err != nil {
		t.Fatalf(err.Error())
	}
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

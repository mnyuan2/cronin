package github

import (
	"context"
	"cron/internal/basic/git"
	"encoding/base64"
	"fmt"
	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
	"log"
	"testing"
)

var token = "github_pat_***"

func TestGetFile(t *testing.T) {
	// 参数

	owner := "mnyuan2"  // 仓库所有者
	name := "cronin"    // 仓库名称
	path := "README.md" // 文件路径
	branch := "master"  // 分支

	// 基础
	src := oauth2.StaticTokenSource(
		&oauth2.Token{
			AccessToken: token,
			//AccessToken: os.Getenv(token),
		},
	)
	// 创建 GitHub 客户端
	httpClient := oauth2.NewClient(context.Background(), src)
	client := githubv4.NewClient(httpClient)

	// 查询
	// 发送 GraphQL 请求
	// 响应结果会再最内层
	q := &struct {
		Repository struct {
			Object struct {
				Blob struct {
					Text string `graphql:"text"`
				} `graphql:"...on Blob"`
			} `graphql:"object(expression: $path)"`
		} `graphql:"repository(owner: $owner, name: $name)"`
	}{}

	err := client.Query(context.Background(), q, map[string]interface{}{
		"owner": githubv4.String(owner),
		"name":  githubv4.String(name),
		"path":  githubv4.String(branch + ":" + path),
	})
	if err != nil {
		log.Fatalf("查询文件错误: %v", err)
	}
	// 检查并处理查询结果
	if q.Repository.Object.Blob.Text != "" {
		// 解码 Base64 编码的文件内容
		decodedContent, err := base64.StdEncoding.DecodeString(q.Repository.Object.Blob.Text)
		if err != nil {
			log.Fatalf("解码文件内容错误: %v", err)
		}
		// 输出文件内容
		fmt.Printf("文件内容:\n%s\n", string(decodedContent))
	} else {
		fmt.Println("未找到文件或内容为空.")
	}
}

// 获取历史提交
func TestApiV4_GetCommitHistory(t *testing.T) {
	api := NewApiV4(&git.Config{AccessToken: token})
	h := git.NewHandler(context.Background())

	owner := "mnyuan2" // 仓库所有者
	name := "cronin"   // 仓库名称
	branch := "master" // 分支

	res, err := api.CommitHistoryGet(h, owner, name, branch, 2)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(res)
}

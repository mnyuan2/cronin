package git

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
	"net/http"
	"strings"
	"time"
)

type GithubApiV4 struct {
	conf *Config
}

func NewGithubApiV4(c *Config) *GithubApiV4 {
	return &GithubApiV4{conf: c}
}

func (m *GithubApiV4) client(ctx context.Context) *githubv4.Client {
	// 基础
	src := oauth2.StaticTokenSource(
		&oauth2.Token{
			AccessToken: m.conf.GetAccessToken(),
		},
	)
	// 创建 GitHub 客户端
	httpClient := oauth2.NewClient(ctx, src)
	return githubv4.NewClient(httpClient)
}

// User 用户信息
func (m *GithubApiV4) User(h *Handler) (user *User, err error) {
	h.OnStartTime(time.Now())
	defer func() {
		h.OnEndTime(time.Now())
	}()
	cli := m.client(h.GetContext())

	// 查询
	q := &userGet{}
	err = cli.Query(context.Background(), q, nil)
	if err != nil {
		return nil, fmt.Errorf("user query error: %w", err)
	}

	return q.Viewer, nil
}

// 获取文件内容
func (m *GithubApiV4) FileGet(h *Handler, r *FileGetRequest) (res *FileGetResponse, err error) {
	h.OnStartTime(time.Now())
	defer func() {
		h.OnEndTime(time.Now())
	}()
	cli := m.client(h.GetContext())

	// 查询
	q := &fileGetBody{}
	err = cli.Query(context.Background(), q, map[string]interface{}{
		"owner": githubv4.String(r.Owner),
		"name":  githubv4.String(r.Repo),
		"path":  githubv4.String(r.Ref + ":" + r.Path),
	})
	if err != nil {
		return nil, fmt.Errorf("file query error: %w", err)
	}
	// 检查并处理查询结果
	res = &FileGetResponse{}
	if q.Repository.Object.Blob.Text != "" {
		res.Content = q.Repository.Object.Blob.Text
		h.OnResponseBody([]byte(res.Content))
	}
	return res, nil
}

// 更新文件
func (m *GithubApiV4) FileUpdate(h *Handler, r *FileUpdateRequest) (res *FileUpdateResponse, err error) {
	h.OnStartTime(time.Now())
	defer func() {
		h.OnEndTime(time.Now())
	}()
	cli := m.client(h.GetContext())

	// 获取最后一次 commit id
	res1, err := m.CommitHistoryGet(h, r.Owner, r.Repo, r.Branch, 1)
	if err != nil {
		return nil, err
	}

	req := &fileUpdateBody{}
	in := githubv4.CreateCommitOnBranchInput{ // 注意此处不能使用指针
		Branch: githubv4.CommittableBranch{
			RepositoryNameWithOwner: githubv4.NewString(githubv4.String(r.Owner + "/" + r.Repo)),
			BranchName:              githubv4.NewString(githubv4.String(r.Branch)),
		},
		Message: githubv4.CommitMessage{
			Headline: githubv4.String(r.Message),
			//Body:     githubv4.NewString(githubv4.String(r.Message)),
		},
		ExpectedHeadOid:  githubv4.GitObjectID(res1.LastOid), //
		ClientMutationID: githubv4.NewString(githubv4.String("cronin")),
		FileChanges: &githubv4.FileChanges{
			Additions: &[]githubv4.FileAddition{
				{
					Path:     githubv4.String(r.Path),
					Contents: githubv4.Base64String(base64.StdEncoding.EncodeToString([]byte(r.Content))),
				},
			},
		},
	}

	h.OnRequestBody([]byte(r.Content))
	err = cli.Mutate(h.GetContext(), req, in, nil)
	if err != nil {
		return nil, err
	}

	return &FileUpdateResponse{
		Commit: &Commit{
			CommitUrl: req.CreateCommitOnBranch.Commit.CommitUrl,
			Oid:       req.CreateCommitOnBranch.Commit.Oid,
		},
	}, nil
}

// 获取提交列表
func (m *GithubApiV4) CommitHistoryGet(h *Handler, owner, repo, branch string, limit int) (res *CommitHistoryGetResponse, err error) {
	h.OnStartTime(time.Now())
	defer func() {
		h.OnEndTime(time.Now())
	}()
	cli := m.client(h.GetContext())

	// 查询
	q := &getCommitHistory{}
	err = cli.Query(context.Background(), q, map[string]any{
		"owner":  githubv4.String(owner),
		"name":   githubv4.String(repo),
		"branch": githubv4.String(branch),
		"limit":  githubv4.Int(limit),
	})
	if err != nil {
		return nil, fmt.Errorf("commit history query error: %s", err.Error())
	}
	// 检查并处理查询结果
	res = &CommitHistoryGetResponse{
		LastOid: q.Repository.Ref.Target.Oid,
		Nodes:   q.Repository.Ref.Target.Commit.History.Nodes,
	}
	return res, nil
}

// 获取 pr 列表
func (m *GithubApiV4) Pulls(h *Handler, r *Pulls) (res *PullsResponse, err error) {
	h.OnStartTime(time.Now())
	defer func() {
		h.OnEndTime(time.Now())
	}()

	cli := m.client(h.GetContext())
	// 查询
	q := &getPulls{}
	v := map[string]any{
		"owner": githubv4.String(r.Owner),
		"name":  githubv4.String(r.Repo),
		"limit": githubv4.Int(r.PerPage),
	}
	if r.State != "" {
		v["states"] = []githubv4.PullRequestState{githubv4.PullRequestState(r.State)}
	} else {
		v["states"] = []githubv4.PullRequestState{}
	}

	err = cli.Query(context.Background(), q, v)
	if err != nil {
		return nil, fmt.Errorf("pulls query error: %s", err.Error())
	}
	// 检查并处理查询结果
	res = &PullsResponse{
		Total: q.Repository.PullRequests.TotalCount,
		List:  q.Repository.PullRequests.Nodes,
	}
	return res, nil
}

// 创建 pr
func (m *GithubApiV4) PullCreate(h *Handler, r *PullsCreateRequest) (res *Pull, err error) {
	h.OnStartTime(time.Now())
	defer func() {
		h.OnEndTime(time.Now())
	}()

	cli := m.client(h.GetContext())

	req := &createPull{}
	in := githubv4.CreatePullRequestInput{
		RepositoryID:     githubv4.ID(""), // 获取方式待定
		BaseRefName:      githubv4.String(r.Base),
		HeadRefName:      githubv4.String(r.Head),
		Title:            githubv4.String(r.Title),
		Body:             githubv4.NewString(githubv4.String(r.Body)),
		ClientMutationID: githubv4.NewString(githubv4.String("cronin")),
	}

	err = cli.Mutate(h.GetContext(), req, in, nil)

	//h.OnGeneral(req.Method, u.String(), resp.StatusCode)
	//h.OnRequestHeader(resp.Request.Header)
	//h.OnRequestBody(reqByte)
	//h.OnResponseHeader(resp.Header)
	//h.OnResponseBody(respByte)
	if err != nil {
		return nil, err
	}

	return req.CreatePullRequest.PullRequest, nil
}

// PullGet pr是否合并
func (m *GithubApiV4) PullGet(h *Handler, r *PullsMergeRequest) (res *Pull, err error) {
	h.OnStartTime(time.Now())
	defer func() {
		h.OnEndTime(time.Now())
	}()
	h.OnGeneral(http.MethodPost, "https://api.github.com/graphql", 0)

	cli := m.client(h.GetContext())
	// 查询
	q := &getPull{}
	v := map[string]any{
		"owner":  githubv4.String(r.Owner),
		"name":   githubv4.String(r.Repo),
		"number": githubv4.Int(r.Number),
	}
	err = cli.Query(context.Background(), q, v)
	if err != nil {
		return &Pull{
			Url: fmt.Sprintf("https://github.com/%s/%s/pull/%v", r.Owner, r.Repo, r.Number),
		}, fmt.Errorf("pull query error: %s", err.Error())
	}
	res = q.Repository.PullRequest
	res.State = strings.ToLower(res.State)
	res.Mergeable = strings.ToLower(res.Mergeable)
	return res, nil
}

// PullsIsMerge pr是否可合并
func (m *GithubApiV4) PullsIsMerge(handler *Handler, r *PullsMergeRequest) (err error) {
	res, err := m.PullGet(handler, r)
	if err != nil {
		return err
	}
	if res.State == "open" {
		if res.Mergeable == "mergeable" {
			return nil // ok
		}
		if res.Mergeable == "conflicting" {
			return errors.New("pr 存在冲突")
		}
	}
	if res.State == "merged" || res.Merged {
		return errors.New("pr 已合并")
	}
	if res.State == "closed" {
		return errors.New("pr 已关闭")
	}
	return errors.New("pr 错误")
}

// PullsMerge pr合并
func (m *GithubApiV4) PullMerge(h *Handler, r *PullsMergeRequest) (res *Pull, err error) {
	h.OnStartTime(time.Now())
	defer func() {
		h.OnEndTime(time.Now())
	}()

	// 查询pr 获取id
	pr, err := m.PullGet(h, r)
	if err != nil {
		return pr, err
	}
	if pr.Merged {
		return pr, errors.New("pr has been merged")
	}

	cli := m.client(h.GetContext())

	req := &pullMerge{}
	in := githubv4.MergePullRequestInput{
		PullRequestID:    githubv4.ID(pr.Id),
		ClientMutationID: githubv4.NewString(githubv4.String("cronin")),
	}
	if r.MergeMethod != "" {
		in.MergeMethod = newPullRequestMergeMethod(githubv4.PullRequestMergeMethod(strings.ToUpper(r.MergeMethod)))
	}
	if r.Title != "" {
		in.CommitHeadline = githubv4.NewString(githubv4.String(r.Title))
	}
	if r.Description != "" {
		in.CommitBody = githubv4.NewString(githubv4.String(r.Description))
	}

	err = cli.Mutate(h.GetContext(), req, in, nil)

	//h.OnGeneral(req.Method, u.String(), resp.StatusCode)
	//h.OnRequestHeader(resp.Request.Header)
	//h.OnRequestBody(reqByte)
	//h.OnResponseHeader(resp.Header)
	//h.OnResponseBody(respByte)
	if err != nil {
		return pr, err
	}

	return req.MergePullRequest.PullRequest, nil
}

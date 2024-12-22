package github

import (
	"context"
	"cron/internal/basic/git"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
	"time"
)

type ApiV4 struct {
	conf *git.Config
}

func NewApiV4(c *git.Config) *ApiV4 {
	return &ApiV4{conf: c}
}

func (m *ApiV4) client(ctx context.Context) *githubv4.Client {
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
func (m *ApiV4) User(h *git.Handler) (user *git.User, err error) {
	cli := m.client(h.GetContext())

	// 查询
	q := &userGet{}
	err = cli.Query(context.Background(), q, nil)
	if err != nil {
		return nil, fmt.Errorf("请求失败，%w", err)
	}

	return q.Viewer, nil
}

// 获取文件内容
func (m *ApiV4) FileGet(h *git.Handler, r *git.FileGetRequest) (res *git.FileGetResponse, err error) {
	cli := m.client(h.GetContext())

	// 查询
	q := &fileGetBody{}
	err = cli.Query(context.Background(), q, map[string]interface{}{
		"owner": githubv4.String(r.Owner),
		"name":  githubv4.String(r.Repo),
		"path":  githubv4.String(r.Ref + ":" + r.Path),
	})
	if err != nil {
		return nil, fmt.Errorf("请求失败，%w", err)
	}
	// 检查并处理查询结果
	res = &git.FileGetResponse{}
	if q.Repository.Object.Blob.Text != "" {
		decodedContent, err := base64.StdEncoding.DecodeString(q.Repository.Object.Blob.Text)
		if err != nil {
			return nil, fmt.Errorf("解码文件内容错误，%w", err)
		}
		res.Content = decodedContent
	}
	return res, nil
}

// 更新文件
func (m *ApiV4) FileUpdate(h *git.Handler, r *git.FileUpdateRequest) (res *git.FileUpdateResponse, err error) {
	cli := m.client(h.GetContext())

	in := &githubv4.CreateCommitOnBranchInput{
		Branch: githubv4.CommittableBranch{
			RepositoryNameWithOwner: githubv4.NewString(githubv4.String(r.Owner + "/" + r.Repo)),
			BranchName:              githubv4.NewString(githubv4.String(r.Branch)),
		},
		Message: githubv4.CommitMessage{
			Headline: githubv4.String(r.Message),
			//Body:     githubv4.NewString(githubv4.String(r.Message)),
		},
		ExpectedHeadOid:  "", //
		ClientMutationID: githubv4.NewString(githubv4.String("cronin")),
		FileChanges: &githubv4.FileChanges{
			Deletions: nil,
			Additions: &[]githubv4.FileAddition{
				{
					Path:     githubv4.String(r.Path),
					Contents: githubv4.Base64String(base64.StdEncoding.EncodeToString([]byte(r.Content))),
				},
			},
		},
	}

	// 获取最后一次 commit id
	res1, err := m.CommitHistoryGet(h, r.Owner, r.Repo, r.Branch, 1)
	if err != nil {
		return nil, fmt.Errorf("commit get error: %s", err.Error())
	}
	in.ExpectedHeadOid = githubv4.GitObjectID(res1.LastOid)

	req := &fileUpdateBody{}
	err = cli.Mutate(h.GetContext(), req, in, nil)
	if err != nil {
		return nil, err
	}

	return &git.FileUpdateResponse{
		Commit: &git.Commit{
			CommitUrl: req.CreateCommitOnBranch.Commit.CommitUrl,
			Oid:       req.CreateCommitOnBranch.Commit.Oid,
		},
	}, nil
}

// 获取提交列表
func (m *ApiV4) CommitHistoryGet(h *git.Handler, owner, repo, branch string, limit int) (res *CommitHistoryGetResponse, err error) {
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
		return nil, fmt.Errorf("get commit request error: %s", err.Error())
	}
	// 检查并处理查询结果
	res = &CommitHistoryGetResponse{
		LastOid: q.Repository.Ref.Target.Oid,
		Nodes:   q.Repository.Ref.Target.Commit.History.Nodes,
	}
	return res, nil
}

// 获取 pr 列表
func (m *ApiV4) Pulls(h *git.Handler, r *git.Pulls) (res *git.PullsResponse, err error) {
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
		return nil, fmt.Errorf("get pulls request error: %s", err.Error())
	}
	// 检查并处理查询结果
	res = &git.PullsResponse{
		Total: q.Repository.PullRequests.TotalCount,
		List:  q.Repository.PullRequests.Nodes,
	}
	return res, nil
}

// 创建 pr
func (m *ApiV4) PullCreate(h *git.Handler, r *git.PullsCreateRequest) (res *git.Pull, err error) {
	h.OnStartTime(time.Now())
	defer func() {
		h.OnEndTime(time.Now())
	}()

	cli := m.client(h.GetContext())

	in := &githubv4.CreatePullRequestInput{
		RepositoryID:     githubv4.ID(""), // 获取方式待定
		BaseRefName:      githubv4.String(r.Base),
		HeadRefName:      githubv4.String(r.Head),
		Title:            githubv4.String(r.Title),
		Body:             githubv4.NewString(githubv4.String(r.Body)),
		ClientMutationID: githubv4.NewString(githubv4.String("cronin")),
	}

	req := &createPull{}
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
func (m *ApiV4) PullGet(h *git.Handler, r *git.PullsMergeRequest) (res *git.Pull, err error) {
	h.OnStartTime(time.Now())
	defer func() {
		h.OnEndTime(time.Now())
	}()

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
		return nil, fmt.Errorf("get pull request error: %s", err.Error())
	}

	return q.Repository.PullRequest, nil
}

// PullsMerge pr合并
func (m *ApiV4) PullMerge(h *git.Handler, r *git.PullsMergeRequest) (res *git.Pull, err error) {
	h.OnStartTime(time.Now())
	defer func() {
		h.OnEndTime(time.Now())
	}()

	// 查询pr 获取id
	pr, err := m.PullGet(h, r)
	if err != nil {
		return nil, err
	}
	if pr.Merged {
		return nil, errors.New("pr 已合并")
	}

	cli := m.client(h.GetContext())
	in := &githubv4.MergePullRequestInput{
		PullRequestID:    githubv4.ID(pr.Id),
		ClientMutationID: githubv4.NewString(githubv4.String("cronin")),
		CommitHeadline:   githubv4.NewString(githubv4.String(r.Title)),
		CommitBody:       githubv4.NewString(githubv4.String(r.Description)),
		AuthorEmail:      nil,
	}
	*in.MergeMethod = githubv4.PullRequestMergeMethod(r.MergeMethod)

	req := &createPull{}
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

package git

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

const (
	apiV5BaseUrl = "https://gitee.com/api/v5"
)

type GiteeApiV5 struct {
	conf *Config
}

func NewGiteeApiV5(c *Config) *GiteeApiV5 {
	return &GiteeApiV5{conf: c}
}

// User 用户信息
func (m *GiteeApiV5) User(h *Handler) (res *User, err error) {
	h.OnStartTime(time.Now())
	defer func() {
		h.OnEndTime(time.Now())
	}()
	u, _ := url.Parse(fmt.Sprintf("%s/user", apiV5BaseUrl))
	params := url.Values{}
	params.Add("access_token", m.conf.GetAccessToken())
	u.RawQuery = params.Encode()

	resp, err := http.Get(u.String())
	h.OnGeneral(http.MethodGet, u.String(), resp.StatusCode)
	h.OnRequestHeader(resp.Request.Header)
	h.OnResponseHeader(resp.Header)

	if err != nil {
		return nil, fmt.Errorf("请求失败，%w", err)
	}
	defer resp.Body.Close()

	b, er := io.ReadAll(resp.Body)
	h.OnResponseBody(b)
	if er != nil {
		return nil, fmt.Errorf("响应获取失败，%w", err)
	}

	if resp.StatusCode != 200 { // {"message":"401 Unauthorized: Access token is expired"}
		out := map[string]any{}
		_ = jsoniter.Unmarshal(b, &out)
		if message, ok := out["message"]; ok {
			return nil, errors.New(message.(string))
		}
	}
	user := &giteeV5User{}
	if err := jsoniter.Unmarshal(b, &user); err != nil {
		return nil, fmt.Errorf("用户信息序列化失败，%w", err)
	}
	return &User{
		Id:        strconv.Itoa(user.Id),
		Login:     user.Login,
		Name:      user.Name,
		Bio:       user.Bio,
		AvatarUrl: user.AvatarUrl,
	}, nil
}

// FileGet 获取仓库具体路径下的文件内容
//
//	 https://gitee.com/api/v5/repos/{owner}/{repo}/contents(/{path})
//		@param string owner 仓库空间名称 仓库所属空间地址(企业、组织或个人的地址path)
//		@param string repo 项目名称 仓库路径(path)
//		@param string path 文件的路径
//		@param string ref 分支、tag或commit。默认: 仓库的默认分支(通常是master)
func (m *GiteeApiV5) FileGet(handler *Handler, r *FileGetRequest) (res *FileGetResponse, err error) {
	handler.OnStartTime(time.Now())
	defer func() {
		handler.OnEndTime(time.Now())
	}()
	u, _ := url.Parse(fmt.Sprintf("%s/repos/%s/%s/contents/%s", apiV5BaseUrl, r.Owner, r.Repo, url.PathEscape(r.Path)))
	params := url.Values{}
	if m.conf.AccessToken != "" {
		params.Add("access_token", m.conf.GetAccessToken())
	}
	if r.Ref != "" {
		params.Add("ref", r.Ref)
	}
	if len(params) > 0 {
		u.RawQuery = params.Encode()
	}
	resp, err := http.Get(u.String())
	handler.OnGeneral(http.MethodGet, u.String(), resp.StatusCode)
	handler.OnRequestHeader(resp.Request.Header)
	handler.OnResponseHeader(resp.Header)

	if err != nil {
		return nil, fmt.Errorf("请求失败，%w", err)
	}
	defer resp.Body.Close()

	b, er := io.ReadAll(resp.Body)
	handler.OnResponseBody(b)
	if er != nil {
		return nil, fmt.Errorf("响应获取失败，%w", err)
	}

	if resp.StatusCode != 200 { // {"message":"401 Unauthorized: Access token is expired"}
		out := &errorResponse{}
		_ = jsoniter.Unmarshal(b, &out)
		return nil, errors.New(out.Message)
	}
	out := &fileGetResponse{}
	_ = jsoniter.Unmarshal(b, &out)

	content, err := base64.StdEncoding.DecodeString(out.Content.Content)
	if err != nil {
		return nil, fmt.Errorf("文件内容解码错误，%w", err)
	}
	res = &FileGetResponse{
		Sha:     out.Sha,
		Content: string(content),
	}

	return res, nil
}

// FileUpdate 文件更新
//
//	https://gitee.com/api/v5/swagger#/putV5ReposOwnerRepoContentsPath
func (m *GiteeApiV5) FileUpdate(handler *Handler, r *FileUpdateRequest) (res *FileUpdateResponse, err error) {
	handler.OnStartTime(time.Now())
	defer func() {
		handler.OnEndTime(time.Now())
	}()
	u, _ := url.Parse(fmt.Sprintf("%s/repos/%s/%s/contents/%s", apiV5BaseUrl, r.Owner, r.Repo, url.PathEscape(r.Path)))

	request := map[string]any{
		"content": base64.StdEncoding.EncodeToString([]byte(r.Content)),
		"sha":     r.Sha,
		"message": r.Message,
	}
	if m.conf.GetAccessToken() != "" {
		request["access_token"] = m.conf.GetAccessToken()
	}
	if r.Branch != "" {
		request["branch"] = r.Branch
	}
	reqByte, _ := jsoniter.Marshal(request)

	req, err := http.NewRequest(http.MethodPut, u.String(), bytes.NewBuffer(reqByte))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json;charset=UTF-8")
	resp, err := http.DefaultClient.Do(req)

	handler.OnGeneral(req.Method, u.String(), resp.StatusCode)
	handler.OnRequestHeader(resp.Request.Header)
	handler.OnRequestBody(reqByte)
	handler.OnResponseHeader(resp.Header)

	if err != nil {
		return nil, fmt.Errorf("请求失败，%w", err)
	}
	defer resp.Body.Close()

	respByte, er := io.ReadAll(resp.Body)
	handler.OnResponseBody(respByte)
	if er != nil {
		return nil, fmt.Errorf("响应获取失败，%w", err)
	}
	res = &FileUpdateResponse{}
	_ = jsoniter.Unmarshal(respByte, res)
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(res.Message)
	}
	if res.Commit != nil {
		res.Commit.CommitUrl = fmt.Sprintf("https://gitee.com/%s/%s/commit/%s", r.Owner, r.Repo, res.Commit.Sha)
	} else {
		res.Commit = &Commit{}
	}
	return res, nil
}

// 获取 pr 列表
func (m *GiteeApiV5) Pulls(handler *Handler, r *Pulls) (res *PullsResponse, err error) {
	handler.OnStartTime(time.Now())
	defer func() {
		handler.OnEndTime(time.Now())
	}()
	u, _ := url.Parse(fmt.Sprintf("%s/repos/%s/%s/pulls", apiV5BaseUrl, r.Owner, r.Repo))
	params := url.Values{}
	if m.conf.GetAccessToken() != "" {
		params.Add("access_token", m.conf.GetAccessToken())
	}
	// 有很多参数需要补充，看一下怎么搞？
	b, _ := jsoniter.Marshal(r)
	p := map[string]any{}
	if err := jsoniter.Unmarshal(b, &p); err != nil {
		return nil, fmt.Errorf("请求数据序列化失败，%w", err)
	}
	for k, v := range p {
		if k == "owner" || k == "repo" {
			continue
		}
		switch val := v.(type) {
		case string:
			if val == "" {
				continue
			}
			params.Add(k, val)
		case int:
			if val == 0 {
				continue
			}
			params.Add(k, strconv.Itoa(val))
		case float64:
			if val == 0 {
				continue
			}
			value := strconv.FormatFloat(val, 'E', -1, 64)
			params.Add(k, value)
		default:
			return nil, fmt.Errorf("请求数据序类型异常，%v", v)
		}
	}
	if len(params) > 0 {
		u.RawQuery = params.Encode()
	}

	resp, err := http.Get(u.String())

	handler.OnGeneral(http.MethodGet, u.String(), resp.StatusCode)
	handler.OnRequestHeader(resp.Request.Header)
	handler.OnResponseHeader(resp.Header)

	if err != nil {
		return res, fmt.Errorf("请求失败，%w", err)
	}
	defer resp.Body.Close()

	respByte, er := io.ReadAll(resp.Body)
	handler.OnResponseBody(respByte)
	if er != nil {
		return res, fmt.Errorf("响应获取失败，%w", err)
	}

	if resp.StatusCode != http.StatusOK { //
		msg := &errorResponse{}
		_ = jsoniter.Unmarshal(respByte, msg)
		return res, errors.New(msg.Message)
	}
	res = &PullsResponse{}
	_ = jsoniter.Unmarshal(respByte, &res)

	return res, nil
}

// 创建 pr
//
//	https://gitee.com/api/v5/swagger#/postV5ReposOwnerRepoPulls
func (m *GiteeApiV5) PullCreate(handler *Handler, r *PullsCreateRequest) (res *Pull, err error) {
	handler.OnStartTime(time.Now())
	defer func() {
		handler.OnEndTime(time.Now())
	}()
	u, _ := url.Parse(fmt.Sprintf("%s/repos/%s/%s/pulls", apiV5BaseUrl, r.Owner, r.Repo))
	reqByte, _ := jsoniter.Marshal(map[string]any{
		"access_token": m.conf.GetAccessToken(),
		"head":         r.Head,
		"base":         r.Base,
		"title":        r.Title,
		"body":         r.Body,
	})

	req, err := http.NewRequest(http.MethodPost, u.String(), bytes.NewBuffer(reqByte))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json;charset=UTF-8")
	resp, err := http.DefaultClient.Do(req)

	handler.OnGeneral(req.Method, u.String(), resp.StatusCode)
	handler.OnRequestHeader(resp.Request.Header)
	handler.OnRequestBody(reqByte)
	handler.OnResponseHeader(resp.Header)

	if err != nil {
		return nil, fmt.Errorf("请求失败，%w", err)
	}
	defer resp.Body.Close()

	// 处理失败
	respByte, er := io.ReadAll(resp.Body)
	handler.OnResponseBody(respByte)
	if er != nil {
		return nil, fmt.Errorf("响应获取失败，%w", err)
	}

	if resp.StatusCode != http.StatusCreated {
		out := map[string]any{}
		_ = jsoniter.Unmarshal(respByte, &out)
		if message, ok := out["messages"]; ok {
			return nil, errors.New(message.([]any)[0].(string))
		} else if message, ok := out["message"]; ok {
			return nil, errors.New(message.(string))
		}
	}
	res = &Pull{}
	if err = jsoniter.Unmarshal(respByte, &res); err != nil {
		return nil, fmt.Errorf("响应解析失败，%w", err)
	}
	return res, nil
}

// pr 审查 确认
func (m *GiteeApiV5) PullsReview(handler *Handler, r *PullsReviewRequest) (res []byte, err error) {
	handler.OnStartTime(time.Now())
	defer func() {
		handler.OnEndTime(time.Now())
	}()
	u, _ := url.Parse(fmt.Sprintf("%s/repos/%s/%s/pulls/%v/review", apiV5BaseUrl, r.Owner, r.Repo, r.Number))

	data := map[string]any{
		"access_token": m.conf.GetAccessToken(),
		"force":        r.Force,
	}
	reqByte, _ := jsoniter.Marshal(data)

	req, err := http.NewRequest(http.MethodPost, u.String(), bytes.NewBuffer(reqByte))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json;charset=UTF-8")
	resp, err := http.DefaultClient.Do(req)

	handler.OnGeneral(req.Method, u.String(), resp.StatusCode)
	handler.OnRequestHeader(resp.Request.Header)
	handler.OnResponseHeader(resp.Header)

	if err != nil {
		return nil, fmt.Errorf("请求失败，%w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNoContent { // 成功
		return []byte{}, nil
	}
	// 处理失败
	respByte, er := io.ReadAll(resp.Body)
	handler.OnResponseBody(respByte)
	if er != nil {
		return nil, fmt.Errorf("响应获取失败，%w", err)
	}
	out := map[string]any{}
	_ = jsoniter.Unmarshal(respByte, &out)
	if message, ok := out["message"]; ok {
		return nil, errors.New(message.(string))
	}
	return respByte, nil
}

// pr 测试 确认
func (m *GiteeApiV5) PullsTest(handler *Handler, r *PullsTestRequest) (res []byte, err error) {
	handler.OnStartTime(time.Now())
	defer func() {
		handler.OnEndTime(time.Now())
	}()
	u, _ := url.Parse(fmt.Sprintf("%s/repos/%s/%s/pulls/%v/test", apiV5BaseUrl, r.Owner, r.Repo, r.Number))
	reqByte, _ := jsoniter.Marshal(map[string]any{
		"access_token": m.conf.GetAccessToken(),
		"force":        r.Force,
	})

	req, err := http.NewRequest(http.MethodPost, u.String(), bytes.NewBuffer(reqByte))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json;charset=UTF-8")
	resp, err := http.DefaultClient.Do(req)

	handler.OnGeneral(req.Method, u.String(), resp.StatusCode)
	handler.OnRequestHeader(resp.Request.Header)
	handler.OnRequestBody(reqByte)
	handler.OnResponseHeader(resp.Header)

	if err != nil {
		return nil, fmt.Errorf("请求失败，%w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNoContent { // 成功
		return []byte{}, nil
	}
	// 处理失败
	respByte, er := io.ReadAll(resp.Body)
	handler.OnResponseBody(respByte)
	if er != nil {
		return nil, fmt.Errorf("响应获取失败，%w", err)
	}
	out := map[string]any{}
	_ = jsoniter.Unmarshal(respByte, &out)
	if message, ok := out["message"]; ok {
		return nil, errors.New(message.(string))
	}
	return respByte, nil
}

// PullGet pr详情
//
//	https://gitee.com/api/v5/swagger#/getV5ReposOwnerRepoPullsNumber
func (m *GiteeApiV5) PullGet(handler *Handler, r *PullsMergeRequest) (res *Pull, err error) {
	handler.OnStartTime(time.Now())
	defer func() {
		handler.OnEndTime(time.Now())
	}()

	u, _ := url.Parse(fmt.Sprintf("%s/repos/%s/%s/pulls/%v", apiV5BaseUrl, r.Owner, r.Repo, r.Number))
	params := url.Values{}
	if m.conf.GetAccessToken() != "" {
		params.Add("access_token", m.conf.GetAccessToken())
	}
	if len(params) > 0 {
		u.RawQuery = params.Encode()
	}

	resp, err := http.Get(u.String())
	res = &Pull{
		Url: fmt.Sprintf("https://gitee.com/%s/%s/pulls/%v", r.Owner, r.Repo, r.Number),
	}

	handler.OnGeneral(http.MethodGet, u.String(), resp.StatusCode)
	handler.OnRequestHeader(resp.Request.Header)
	handler.OnResponseHeader(resp.Header)

	if err != nil {
		return res, fmt.Errorf("请求失败，%w", err)
	}
	defer resp.Body.Close()

	respByte, er := io.ReadAll(resp.Body)
	handler.OnResponseBody(respByte)
	if er != nil {
		return res, fmt.Errorf("响应获取失败，%w", err)
	}

	if resp.StatusCode != http.StatusOK {
		tmp := &errorResponse{}
		_ = jsoniter.Unmarshal(respByte, tmp)
		return res, errors.New(tmp.Error + tmp.Message)
	}
	body := &giteeV5Pull{}
	if err = jsoniter.Unmarshal(respByte, body); err != nil {
		return res, fmt.Errorf("响应解析失败，%w", err)
	}
	res.Id = strconv.Itoa(body.Id)
	res.Title = body.Title
	res.Number = body.Number
	res.State = body.State
	res.Merged = false
	res.Mergeable = "unknown"
	res.CreateAt = body.CreatedAt
	res.Url = body.HtmlUrl
	res.HeadRefName = body.Head.Ref
	res.BaseRefName = body.Base.Ref
	if body.MergedAt != nil {
		res.Merged = true
	}
	if body.Mergeable {
		res.Mergeable = "mergeable"
	}
	if !body.CanMergeCheck && res.State == "open" {
		res.Mergeable = "conflicting"
	}
	return res, nil
}

func (m *GiteeApiV5) PullsIsMerge(handler *Handler, r *PullsMergeRequest) (err error) {
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
	if res.State == "merged" {
		return errors.New("pr 已合并")
	}
	if res.State == "closed" {
		return errors.New("pr 已关闭")
	}
	return errors.New("pr 错误")
}

// PullMerge pr合并
//
//	https://gitee.com/api/v5/swagger#/putV5ReposOwnerRepoPullsNumberMerge
func (m *GiteeApiV5) PullMerge(handler *Handler, r *PullsMergeRequest) (res *Pull, err error) {
	handler.OnStartTime(time.Now())
	defer func() {
		handler.OnEndTime(time.Now())
	}()
	res = &Pull{
		Url: fmt.Sprintf("https://gitee.com/%s/%s/pulls/%v", r.Owner, r.Repo, r.Number),
	}

	u, _ := url.Parse(fmt.Sprintf("%s/repos/%s/%s/pulls/%v/merge", apiV5BaseUrl, r.Owner, r.Repo, r.Number))
	reqByte, _ := jsoniter.Marshal(map[string]any{
		"access_token": m.conf.GetAccessToken(),
		"merge_method": r.MergeMethod,
		"title":        r.Title,
		"description":  r.Description,
	})

	req, err := http.NewRequest(http.MethodPut, u.String(), bytes.NewBuffer(reqByte))
	if err != nil {
		return res, err
	}
	req.Header.Set("Content-Type", "application/json;charset=UTF-8")
	resp, err := http.DefaultClient.Do(req)

	handler.OnGeneral(req.Method, u.String(), resp.StatusCode)
	handler.OnRequestHeader(resp.Request.Header)
	handler.OnRequestBody(reqByte)
	handler.OnResponseHeader(resp.Header)

	if err != nil {
		return res, fmt.Errorf("请求失败，%w", err)
	}
	defer resp.Body.Close()

	respByte, er := io.ReadAll(resp.Body)
	handler.OnResponseBody(respByte)
	if er != nil {
		return res, fmt.Errorf("响应获取失败，%w", err)
	}

	if resp.StatusCode != http.StatusOK { // {"message":"此 Pull Request 未通过设置的审查"}、{"message":"此 Pull Request 未通过设置的测试"}
		tmp := &errorResponse{}
		_ = jsoniter.Unmarshal(respByte, tmp)
		return res, errors.New(tmp.Message)
	}

	if err = jsoniter.Unmarshal(respByte, res); err != nil {
		return nil, fmt.Errorf("响应解析失败，%w", err)
	}
	return res, nil
}

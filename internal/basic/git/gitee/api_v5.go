package gitee

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"io"
	"net/http"
	"net/url"
)

const (
	apiV5BaseUrl = "https://gitee.com/api/v5"
)

// 配置
type Config interface {
	GetAccessToken() string
}

type ApiV5 struct {
	conf Config
}

func NewApiV5(c Config) *ApiV5 {
	return &ApiV5{conf: c}
}

// ReposContents 获取仓库具体路径下的文件内容
//
//	 https://gitee.com/api/v5/repos/{owner}/{repo}/contents(/{path})
//		@param string owner 仓库空间名称 仓库所属空间地址(企业、组织或个人的地址path)
//		@param string repo 项目名称 仓库路径(path)
//		@param string path 文件的路径
//		@param string ref 分支、tag或commit。默认: 仓库的默认分支(通常是master)
func (m *ApiV5) ReposContents(handler *Handler, owner, repo, path, ref string) (res []byte, err error) {
	u, _ := url.Parse(fmt.Sprintf("%s/repos/%s/%s/contents/%s", apiV5BaseUrl, owner, repo, url.QueryEscape(path)))
	params := url.Values{}
	if m.conf.GetAccessToken() != "" {
		params.Add("access_token", m.conf.GetAccessToken())
	}
	if ref != "" {
		params.Add("ref", ref)
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

	out := map[string]any{}
	_ = jsoniter.Unmarshal(b, &out)
	if resp.StatusCode != 200 { // {"message":"401 Unauthorized: Access token is expired"}
		if message, ok := out["message"]; ok {
			return nil, errors.New(message.(string))
		}
	} else if content, oK := out["content"].(string); oK {
		return base64.StdEncoding.DecodeString(content)
	}
	return b, errors.New("请求异常")
}

// User 用户信息
func (m *ApiV5) User(handler *Handler) (res []byte, err error) {
	u, _ := url.Parse(fmt.Sprintf("%s/user", apiV5BaseUrl))
	params := url.Values{}
	params.Add("access_token", m.conf.GetAccessToken())
	u.RawQuery = params.Encode()

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
		out := map[string]any{}
		_ = jsoniter.Unmarshal(b, &out)
		if message, ok := out["message"]; ok {
			return nil, errors.New(message.(string))
		}
	}
	return b, nil
}

// 创建pr
//
//	https://gitee.com/api/v5/swagger#/postV5ReposOwnerRepoPulls
func (m *ApiV5) PullsCreate(handler *Handler, r *PullsCreateRequest) (res []byte, err error) {
	u, _ := url.Parse(fmt.Sprintf("%s/repos/%s/%s/pulls", apiV5BaseUrl, r.Owner, r.Repo))
	reqByte, _ := jsoniter.Marshal(map[string]any{
		"access_token": m.conf.GetAccessToken(),
		"head":         r.Head,
		"base":         r.Base,
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
		if message, ok := out["message"]; ok {
			return nil, errors.New(message.(string))
		}
	}
	return respByte, nil
}

// pr 审查 确认
func (m *ApiV5) PullsReview(handler *Handler, r *PullsReviewRequest) (res []byte, err error) {
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
func (m *ApiV5) PullsTest(handler *Handler, r *PullsTestRequest) (res []byte, err error) {
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

// PullsMerge pr合并
//
//	https://gitee.com/api/v5/swagger#/putV5ReposOwnerRepoPullsNumberMerge
func (m *ApiV5) PullsMerge(handler *Handler, r *PullsMergeRequest) (res []byte, err error) {
	u, _ := url.Parse(fmt.Sprintf("%s/repos/%s/%s/pulls/%v/merge", apiV5BaseUrl, r.Owner, r.Repo, r.Number))
	reqByte, _ := jsoniter.Marshal(map[string]any{
		"access_token": m.conf.GetAccessToken(),
		"merge_method": r.MergeMethod,
		"title":        r.Title,
		"description":  r.Description,
	})

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

	if resp.StatusCode != http.StatusOK { // {"message":"此 Pull Request 未通过设置的审查"}  {"message":"此 Pull Request 未通过设置的测试"}
		out := map[string]any{}
		_ = jsoniter.Unmarshal(respByte, &out)
		if message, ok := out["message"]; ok {
			return nil, errors.New(message.(string))
		}
	}
	return respByte, errors.New("请求异常")
}

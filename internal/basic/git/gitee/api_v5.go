package gitee

import (
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

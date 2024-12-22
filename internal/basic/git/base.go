package git

import (
	"cron/internal/basic/git/gitee"
	"cron/internal/basic/git/github"
)

type Config struct {
	Type        string `json:"type"`
	AccessToken string `json:"access_token"`
}

const (
	TypeGitee  = "gitee"
	TypeGithub = "github"
)

func (m *Config) GetAccessToken() string {
	return m.AccessToken
}

type Api interface {
	User(h *Handler) (user *User, err error)
	FileGet(h *Handler, r *FileGetRequest) (res *FileGetResponse, err error)
	FileUpdate(h *Handler, r *FileUpdateRequest) (res *FileUpdateResponse, err error)
	Pulls(h *Handler, r *Pulls) (res *PullsResponse, err error)
	PullCreate(h *Handler, r *PullsCreateRequest) (res *Pull, err error)
	PullGet(h *Handler, r *PullsMergeRequest) (res *Pull, err error)
	PullMerge(h *Handler, r *PullsMergeRequest) (res *Pull, err error)
}

func NewApi(conf Config) Api {
	switch conf.Type {
	case TypeGitee:
		return gitee.NewApiV5(&conf)
	case TypeGithub:
		return github.NewApiV4(&conf)
	default:
		return nil
	}
}

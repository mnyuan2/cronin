package git

type Config struct {
	Driver      string `json:"driver"`
	AccessToken string `json:"access_token"`
}

const (
	DriverGitee  = "gitee"
	DriverGithub = "github"
)

func (m *Config) GetAccessToken() string {
	return m.AccessToken
}

type Api interface {
	User(h *Handler) (user *User, err error)
	FileGet(h *Handler, r *FileGetRequest) (res *FileGetResponse, err error)
	PullsIsMerge(handler *Handler, r *PullsMergeRequest) (err error)
	FileUpdate(h *Handler, r *FileUpdateRequest) (res *FileUpdateResponse, err error)
	Pulls(h *Handler, r *Pulls) (res *PullsResponse, err error)
	PullCreate(h *Handler, r *PullsCreateRequest) (res *Pull, err error)
	PullGet(h *Handler, r *PullsMergeRequest) (res *Pull, err error)
	PullMerge(h *Handler, r *PullsMergeRequest) (res *Pull, err error)
}

func NewApi(conf Config) Api {
	switch conf.Driver {
	case DriverGitee:
		return NewGiteeApiV5(&conf)
	case DriverGithub:
		return NewGithubApiV4(&conf)
	default:
		return nil
	}
}

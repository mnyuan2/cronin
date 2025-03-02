package git

type errorResponse struct {
	Error   string `json:"error"` // 部分业务错误时会返回此字段
	Message string `json:"message"`
}

// 文件获取 响应
type fileGetResponse struct {
	Message string `json:"message"` // 错误描述
	Content
}
type Content struct {
	Type        string `json:"type"`
	Encoding    string `json:"encoding"`
	Size        int    `json:"size"`
	Name        string `json:"name"`
	Path        string `json:"path"`
	Content     string `json:"content"`
	Sha         string `json:"sha"`
	Url         string `json:"url"`
	HtmlUrl     string `json:"html_url"`
	DownloadUrl string `json:"download_url"`
	Links       struct {
		Self string `json:"self"`
		Html string `json:"html"`
	} `json:"_links"`
}

type giteeV5User struct {
	Id        int    `json:"id"`
	Login     string `json:"login"`
	Name      string `json:"name"`
	Bio       string `json:"bio"`
	AvatarUrl string `json:"avatar_url"`
}

// 拉取请求
type giteeV5Pull struct {
	Id            int     `json:"id"`
	Number        int     `json:"number"`
	State         string  `json:"state"`
	MergedAt      *string `json:"merged_at"`
	Mergeable     bool    `json:"mergeable"` // 是否可合并
	CanMergeCheck bool    `json:"can_merge_check"`
	Head          struct {
		Ref string `json:"ref"`
	} `json:"head"`
	Base struct {
		Ref string `json:"ref"`
	} `json:"base"`
	Title     string `json:"title"`
	Body      string `json:"body"`
	CreatedAt string `json:"created_at"`
	HtmlUrl   string `json:"html_url"`
}

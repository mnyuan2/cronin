package git

type BaseRequest struct {
	Owner string `json:"owner"` // 空间地址
	Repo  string `json:"repo"`  // 项目名称（仓库路径）
}

type User struct {
	Id        string `json:"id"`
	Login     string `json:"login"`
	Name      string `json:"name"`
	Bio       string `json:"bio"`
	AvatarUrl string `json:"avatar_url" graphql:"avatarUrl"`
}

// 文件获取
type FileGetRequest struct {
	BaseRequest
	Path string
	Ref  string
}

// 文件获取 响应
type FileGetResponse struct {
	Sha     string `json:"sha"`
	Content string `json:"content"`
}

// 文件更新
type FileUpdateRequest struct {
	BaseRequest
	Path    string // branch
	Content string // 文件内容, 要用 base64 编码
	Sha     string // 文件的 Blob SHA，可通过 [获取仓库具体路径下的内容] API 获取
	Message string // 提交(commit 描述)信息
	Branch  string // 分支名称。默认为仓库对默认分支
}

type FileUpdateResponse struct {
	Message string  `json:"message"` // 错误描述
	Commit  *Commit `json:"commit"`
}

type Commit struct {
	Sha           string `json:"sha" graphql:""` // gitee 使用
	CommitUrl     string `json:"commit_url"`
	Oid           string `json:"oid"`
	Message       string `json:"message"`
	CommittedDate string `json:"committed_date" graphql:"committedDate"`
	Author        struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	} `json:"author"`
}

// pr 列表请求
type Pulls struct {
	BaseRequest
	State   string `json:"state"`    // 可选。Pull Request 状态: open、closed、merged、all
	Head    string `json:"head"`     // 可选。Pull Request 提交的源分支。格式：branch 或者：username:branch
	Base    string `json:"base"`     // 可选。Pull Request 提交目标分支的名称。
	Page    int    `json:"page"`     // 当前的页码
	PerPage int    `json:"per_page"` // 每页的数量，最大为 100
}

// pr 列表响应
type PullsResponse struct {
	Total int     `json:"total"`
	List  []*Pull `json:"list"`
}
type Pull struct {
	Id          string `json:"id"`
	Title       string `json:"title"`
	Number      int    `json:"number"`
	State       string `json:"state"`
	Merged      bool   `json:"merged"`    // 是否已合并
	Mergeable   string `json:"mergeable"` // 是否可合并：conflicting·冲突、mergeable·可以合并、unknown·未知的
	CreateAt    string `json:"create_at" graphql:"createdAt"`
	Url         string `json:"url"`
	HeadRefName string `json:"head_ref_name" graphql:"headRefName"`
	BaseRefName string `json:"base_ref_name" graphql:"baseRefName"`
}

type PullsCreateRequest struct {
	BaseRequest
	// 必填。Pull Request 标题
	Title string
	// 必填。Pull Request 提交的源分支。格式：branch 或者：username:branch
	Head string
	// 必填。Pull Request 提交目标分支的名称
	Base string
	// 可选。Pull Request 内容
	Body string
	// 可选。里程碑序号(id)
	MilestoneNumber int
	// 用逗号分开的标签，名称要求长度在 2-20 之间且非特殊字符。如: bug,performance
	Labels string
	// 可选。Pull Request的标题和内容可以根据指定的Issue Id自动填充
	Issue string
	// 可选。审查人员username，可多个，半角逗号分隔，如：(username1,username2), 注意: 当仓库代码审查设置中已设置【指派审查人员】则此选项无效
	Assignees string
	// 可选。测试人员username，可多个，半角逗号分隔，如：(username1,username2), 注意: 当仓库代码审查设置中已设置【指派测试人员】则此选项无效
	Testers string
	// 可选。最少审查人数
	AssigneesNumber int
	// 可选。最少测试人数
	TestersNumber int
	// 可选。依赖的当前仓库下的PR编号，置空则清空依赖的PR。如：17,18,19
	RefPullRequestNumbers string
	// 可选。合并PR后是否删除源分支，默认false（不删除）
	PruneSourceBranch bool
	// 可选，合并后是否关闭关联的 Issue，默认根据仓库配置设置
	CloseRelatedIssue bool
	// 是否设置为草稿
	Draft bool
	// 接受 Pull Request 时使用扁平化（Squash）合并
	Squash bool
}

// pr 审查 确认
type PullsReviewRequest struct {
	BaseRequest
	// 第几个PR，即本仓库PR的序数
	Number int32
	// 是否强制测试通过（默认否），只对管理员生效
	Force bool
}

// pr 测试 确认
type PullsTestRequest struct {
	BaseRequest
	// 第几个PR，即本仓库PR的序数
	Number int32
	// 是否强制审查通过（默认否），只对管理员生效
	Force bool
}

// pr 合并
type PullsMergeRequest struct {
	BaseRequest
	// 第几个PR，即本仓库PR的序数
	Number int32
	// 可选。合并PR的方法，merge（合并所有提交）、squash（扁平化分支合并）和rebase（变基并合并）。默认为merge。
	MergeMethod string
	// 可选。合并PR后是否删除源分支，默认false（不删除）
	PruneSourceBranch bool
	// 可选。合并 commit 标题，默认为PR的标题
	Title string
	// 可选。合并 commit 描述，默认为 "Merge pull request !{pr_id} from {author}/{source_branch}"，与页面显示的默认一致。
	Description string
	State       string
}

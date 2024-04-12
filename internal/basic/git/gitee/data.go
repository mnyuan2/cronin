package gitee

type PullsCreateRequest struct {
	// owner 仓库所属空间地址(企业、组织或个人的地址path)
	Owner string
	// repo 仓库路径(path)
	Repo string
	// title 必填。Pull Request 标题
	Title string
	// head 必填。Pull Request 提交的源分支。格式：branch 或者：username:branch
	Head string
	// base 必填。Pull Request 提交目标分支的名称
	Base string
	// body 可选。Pull Request 内容
	Body string
	// milestone_number 可选。里程碑序号(id)
	MilestoneNumber int
	// labels 用逗号分开的标签，名称要求长度在 2-20 之间且非特殊字符。如: bug,performance
	Labels string
	// issue 可选。Pull Request的标题和内容可以根据指定的Issue Id自动填充
	Issue string
	// assignees 可选。审查人员username，可多个，半角逗号分隔，如：(username1,username2), 注意: 当仓库代码审查设置中已设置【指派审查人员】则此选项无效
	Assignees string
	// testers 可选。测试人员username，可多个，半角逗号分隔，如：(username1,username2), 注意: 当仓库代码审查设置中已设置【指派测试人员】则此选项无效
	Testers string
	// assignees_number 可选。最少审查人数
	AssigneesNumber int
	// testers_number 可选。最少测试人数
	TestersNumber int
	// ref_pull_request_numbers 可选。依赖的当前仓库下的PR编号，置空则清空依赖的PR。如：17,18,19
	RefPullRequestNumbers string
	// prune_source_branch 可选。合并PR后是否删除源分支，默认false（不删除）
	PruneSourceBranch bool
	// close_related_issue 可选，合并后是否关闭关联的 Issue，默认根据仓库配置设置
	CloseRelatedIssue bool
	// draft 是否设置为草稿
	Draft bool
	// squash 接受 Pull Request 时使用扁平化（Squash）合并
	Squash bool
}

type PullsMergeRequest struct {
	// number 第几个PR，即本仓库PR的序数
	Number int32
	// merge_method 可选。合并PR的方法，merge（合并所有提交）、squash（扁平化分支合并）和rebase（变基并合并）。默认为merge。
	MergeMethod string
	// prune_source_branch 可选。合并PR后是否删除源分支，默认false（不删除）
	PruneSourceBranch bool
	// title 可选。合并标题，默认为PR的标题
	Title string
	// description 可选。合并描述，默认为 "Merge pull request !{pr_id} from {author}/{source_branch}"，与页面显示的默认一致。
	Description string
}

package enum

const (
	/* pr系列 */
	// pr 列表 1
	GitEventPullsCreate = 2 // pr创建
	// pr 详情 3
	GitEventPullsDetail = 3 // pr详情
	// pr 更新 4
	// pr 操作日志 5
	// pr commit信息 6
	// pr commit文件列表 7
	// pr 是否合并 8
	GitEventPullsIsMerge = 8 // pr是否合并
	// pr 合并 9
	GitEventPullsMerge = 9 // pr合并
	// pr 审查确认 10
	GitEventPullsReview = 10
	// pr 测试确认 11
	GitEventPullsTest = 11
	// pr 指派审查人员 12
	// pr 取消审查人员 13
	// pr 重置审查状态 14
	// pr 指派测试人员 15
	// pr 取消测试人员 16
	// pr 重置测试状态 17
	// pr 获取关联issues 18
	// pr 获取所有评论 19
	// pr 提交评论 20
	// pr 获取所有标签 21
	// pr 创建标签 22
	// pr 替换标签 23
	// pr 删除标签 24
	// pr 获取指定评论 25
	// pr 编辑评论 26
	// pr 删除评论 27
	/* 仓库系列 */
	// 0101 获取所有分支 	getV5ReposOwnerRepoBranches
	// 0102 创建分支 	postV5ReposOwnerRepoBranches
	// 0103 获取单个分支 	getV5ReposOwnerRepoBranchesBranch
	// 0104 设置分支保护 	putV5ReposOwnerRepoBranchesBranchProtection
	// 0105 取消保护分支的设置 	deleteV5ReposOwnerRepoBranchesBranchProtection
	// 0106 更新保护分支规则 	putV5ReposOwnerRepoBranchesWildcardSetting
	// 0107 删除保护分支规则 	deleteV5ReposOwnerRepoBranchesWildcardSetting
	// 0108 新建保护分支规则 	putV5ReposOwnerRepoBranchesSettingNew
	// 0109 列出仓库所有的 tags 	getV5ReposOwnerRepoTags
	// 0110 创建一个仓库的 Tag 	postV5ReposOwnerRepoTags
	// 0111 获取仓库的 Commit 评论 	getV5ReposOwnerRepoComments
	// 0112 获取单个Commit的评论 	getV5ReposOwnerRepoCommitsRefComments
	// 0113 获取仓库的某条Commit评论 	getV5ReposOwnerRepoCommentsId
	// 0114 更新Commit评论 	patchV5ReposOwnerRepoCommentsId
	// 0115 删除Commit评论 	deleteV5ReposOwnerRepoCommentsId
	// 0116 创建Commit评论 	postV5ReposOwnerRepoCommitsShaComments
	// 0117 仓库的所有提交 	getV5ReposOwnerRepoCommits
	// 0118 提交多个文件变更 	postV5ReposOwnerRepoCommits
	// 0119 仓库的某个提交 	getV5ReposOwnerRepoCommitsSha
	// 0120 Commits 	getV5ReposOwnerRepoCompareBase
	// 0121 获取仓库已部署的公钥 	getV5ReposOwnerRepoKeys
	// 0122 为仓库添加公钥 	postV5ReposOwnerRepoKeys
	// 0123 获取仓库可部署的公钥 	getV5ReposOwnerRepoKeysAvailable
	// 0124 启用仓库公钥 	putV5ReposOwnerRepoKeysEnableId
	// 0125 停用仓库公钥 	deleteV5ReposOwnerRepoKeysEnableId
	// 0126 获取仓库的单个公钥 	getV5ReposOwnerRepoKeysId
	// 0127 删除一个仓库公钥 	deleteV5ReposOwnerRepoKeysId
	// 0128 获取仓库README 	getV5ReposOwnerRepoReadme
	// 0129 获取仓库具体路径下的内容 	getV5ReposOwnerRepoContents
	// 0130 新建文件 	postV5ReposOwnerRepoContentsPath

	// 0131 更新文件 	putV5ReposOwnerRepoContentsPath
	GitEventFileUpdate = 131
	// 0132 删除文件 	deleteV5ReposOwnerRepoContentsPath
	// 0133 Blame 	getV5ReposOwnerRepoBlamePath
	// 0134 获取 raw 文件（100MB 以内） 	getV5ReposOwnerRepoRawPath
	// 0135 下载仓库 zip 	getV5ReposOwnerRepoZipball
	// 0136 下载仓库 tar.gz 	getV5ReposOwnerRepoTarball
	// 0137 获取Pages信息 	getV5ReposOwnerRepoPages
	// 0138 上传设置 Pages SSL 证书和域名 	putV5ReposOwnerRepoPages
	// 0139 请求建立Pages 	postV5ReposOwnerRepoPagesBuilds
	// 0140 获取用户的某个仓库 	getV5ReposOwnerRepo
	// 0141 更新仓库设置 	patchV5ReposOwnerRepo
	// 0142 删除一个仓库 	deleteV5ReposOwnerRepo
	// 0143 修改代码审查设置 	putV5ReposOwnerRepoReviewer
	// 0144 获取仓库推送规则设置 	getV5ReposOwnerRepoPushConfig
	// 0145 修改仓库推送规则设置 	putV5ReposOwnerRepoPushConfig
	// 0146 获取仓库贡献者 	getV5ReposOwnerRepoContributors
	// 0147 清空一个仓库 	putV5ReposOwnerRepoClear
	// 0148 获取仓库的所有成员 	getV5ReposOwnerRepoCollaborators
	// 0149 判断用户是否为仓库成员 	getV5ReposOwnerRepoCollaboratorsUsername
	// 0150 添加仓库成员或更新仓库成员权限 	putV5ReposOwnerRepoCollaboratorsUsername
	// 0151 移除仓库成员 	deleteV5ReposOwnerRepoCollaboratorsUsername
	// 0152 查看仓库成员的权限 	getV5ReposOwnerRepoCollaboratorsUsernamePermission
	// 0153 查看仓库的Forks 	getV5ReposOwnerRepoForks
	// 0154 Fork一个仓库 	postV5ReposOwnerRepoForks
	// 0155 获取仓库的百度统计 key 	getV5ReposOwnerRepoBaiduStatisticKey
	// 0156 设置/更新仓库的百度统计 key 	postV5ReposOwnerRepoBaiduStatisticKey
	// 0157 删除仓库的百度统计 key 	deleteV5ReposOwnerRepoBaiduStatisticKey
	// 0158 获取最近30天的七日以内访问量 	postV5ReposOwnerRepoTrafficData
	// 0159 获取仓库的所有Releases 	getV5ReposOwnerRepoReleases
	// 0160 创建仓库Release 	postV5ReposOwnerRepoReleases
	// 0161 获取仓库的单个Releases 	getV5ReposOwnerRepoReleasesId
	// 0162 更新仓库Release 	patchV5ReposOwnerRepoReleasesId
	// 0163 删除仓库Release 	deleteV5ReposOwnerRepoReleasesId
	// 0164 获取仓库的最后更新的Release 	getV5ReposOwnerRepoReleasesLatest
	// 0165 根据Tag名称获取仓库的Release 	getV5ReposOwnerRepoReleasesTagsTag
	// 0166 获取仓库下的指定 Release 的所有附件 	getV5ReposOwnerRepoReleasesReleaseIdAttachFiles
	// 0167 上传附件到仓库指定 Release 	postV5ReposOwnerRepoReleasesReleaseIdAttachFiles
	// 0168 获取仓库下指定 Release 的单个附件 	getV5ReposOwnerRepoReleasesReleaseIdAttachFilesAttachFileId
	// 0169 删除仓库下指定 Release 的指定附件 	deleteV5ReposOwnerRepoReleasesReleaseIdAttachFilesAttachFileId
	// 0170 下载指定 Release 的单个附件 	getV5ReposOwnerRepoReleasesReleaseIdAttachFilesAttachFileIdDownload
	// 0171 开通Gitee Go 	postV5ReposOwnerRepoOpen
	// 0172 列出授权用户的所有仓库 	getV5UserRepos
	// 0173 创建一个仓库 	postV5UserRepos
	// 0174 获取某个用户的公开仓库 	getV5UsersUsernameRepos
	// 0175 获取一个组织的仓库 	getV5OrgsOrgRepos
	// 0176 创建组织仓库 	postV5OrgsOrgRepos
	// 0177 获取企业的所有仓库 	getV5EnterprisesEnterpriseRepos
	// 0178 创建企业仓库 	postV5EnterprisesEnterpriseRepos
)

var GitEventMap = map[int]string{
	GitEventPullsCreate:  "pr创建",
	GitEventPullsDetail:  "pr详情",
	GitEventPullsIsMerge: "pr校验是否合并",
	GitEventPullsMerge:   "pr合并",
	GitEventFileUpdate:   "文件更新",
}

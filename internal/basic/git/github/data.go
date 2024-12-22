package github

import "cron/internal/basic/git"

/*
* 观察用户

	query {
		viewer {
	        id
	        login
	        name
	        bio
	        avatarUrl
		}
	}

	{
	    "data": {
	        "viewer": {
	            "createdAt": "2017-05-01T08:09:44Z",
	            "id": "MDB**XNl**I4M**yGF9y",
	            "login": "mnyuan2",
	            "name": "慢鸟",
	            "bio": "php、go后端程序员。",
	            "avatarUrl": "https://avatars.githubusercontent.com/u/28252582?v=4"
	        }
	    }
	}
*/
type userGet struct {
	Viewer *git.User `graphql:"viewer"`
}

// 获取文件
type fileGetBody struct {
	Repository struct {
		Object struct {
			Blob struct {
				Text string `graphql:"text"`
			} `graphql:"...on Blob"`
		} `graphql:"object(expression: $path)"`
	} `graphql:"repository(owner: $owner, name: $name)"`
}

// 更新文件
type fileUpdateBody struct {
	CreateCommitOnBranch struct {
		Commit struct {
			CommitUrl string
			Oid       string
		}
	} `graphql:"createCommitOnBranch(input: $input)"`
}

/*
	获取提交历史

样例

	query ($owner: String!, $name: String!, $branch: String!, $limit:Int!) {
	    repository(owner: $owner, name: $name) {
	        ref(qualifiedName: $branch) {
	            target {
	                oid
	                ... on Commit {
	                    history(first: $limit) {
	                        nodes {
	                            oid
	                            message
	                            committedDate
	                            author {
	                                name
	                                email
	                            }
	                        }
	                    }
	                }
	            }
	        }
	    }
	}
*/
type getCommitHistory struct {
	Repository struct {
		Ref struct {
			Target struct {
				Oid    string
				Commit struct {
					History struct {
						Nodes []*git.Commit `graphql:"nodes"`
					} `graphql:"history(first: $limit)"`
				} `graphql:"... on Commit"`
			} `graphql:"target"`
		} `graphql:"ref(qualifiedName: $branch)"`
	} `graphql:"repository(owner: $owner, name: $name)"`
}

// 提交历史响应
type CommitHistoryGetResponse struct {
	LastOid string        `json:"last_oid"`
	Nodes   []*git.Commit `json:"nodes"`
}

/*
pr 创建

	mutation($input:CreatePullRequestInput!) {
	    createPullRequest(input: $input) {
	        pullRequest {
	            id
	            number
	            url
	        }
	    }
	}

	{
	    "input": {
	        "repositoryId": "R_kg**A9**rJ", // 获取仓库id 第一次创建pr失败信息中返回 next_global_id 取用后再次创建pr即可。
	        "baseRefName": "master",
	        "headRefName": "hotfix/001",
	        "title": "test demo",
	        "body": "pr body ."
	    }
	}

	{
	    "data": {
	        "createPullRequest": {
	            "pullRequest": {
	                "id": "PR_kwDOB4c8rc6F6nPm",
	                "number": 1,
	                "url": "https://github.com/mnyuan2/jquery/pull/1"
	            }
	        }
	    }
	}
*/
type createPull struct {
	CreatePullRequest struct {
		PullRequest *git.Pull `graphql:"pullRequest"`
	} `graphql:"createPullRequest(input: $input)"`
}

/*
pr 列表

	query ($owner: String!, $name: String!) {
	    repository(owner: $owner, name: $name) {
	        pullRequests(first: 10, states: [OPEN]) {
	            totalCount
	            nodes {
	                id
	                number
	                title
	                state
	                createdAt
	                url
	            }
	        }
	    }
	}

	{
	    "owner": "mnyuan2",
	    "name": "jquery"
	}

	{
	    "data": {
	        "repository": {
	            "pullRequests": {
	                "totalCount": 1,
	                "nodes": [
	                    {
	                        "id": "PR_kwDOB4c8rc6F6nPm",
	                        "number": 1,
	                        "title": "test demo",
	                        "state": "OPEN",
	                        "createdAt": "2024-12-20T13:43:31Z",
	                        "url": "https://github.com/mnyuan2/jquery/pull/1"
	                    }
	                ]
	            }
	        }
	    }
	}
*/
type getPulls struct {
	Repository struct {
		PullRequests struct {
			TotalCount int
			Nodes      []*git.Pull `graphql:"nodes"`
		} `graphql:"pullRequests(first: $limit, status: $status)"`
	} `graphql:"repository(owner: $owner, name: $name)"`
}

/*
验证pr是否合并

	query ($owner: String!, $name: String!, $number:Int!) {
	    repository(owner: $owner, name: $name) {
	        pullRequest(number: $number) {
	            id
	            state
	            merged
	            url
	        }
	    }
	}

	{
	    "owner": "mnyuan2",
	    "name": "jquery",
	    "number": 1
	}

	{
	    "data": {
	        "repository": {
	            "pullRequest": {
	                "id": "PR_kwDOB4c8rc6F6nPm",
	                "state": "OPEN", # 或 "OPEN"、"CLOSED"
	                "merged": false, # 表示是否已合并
	                "url": "https://github.com/mnyuan2/jquery/pull/1"
	            }
	        }
	    }
	}
*/
type getPull struct {
	Repository struct {
		PullRequest *git.Pull `graphql:"pullRequest(number: number)"`
	} `graphql:"repository(owner: $owner, name: $name)"`
}

/* pr 操作合并
注意：合并操作必须要用id，number 号可以查询出 id。
	所以很多更新操作，都是要先查询一次，才能正式执行写操作。
	思考一下这种联动操作，如何来规划链路的记录。
	传入记录器，内部就执行写入了，就可以不用把信息返回了。
	或者再外部分开执行。
		但这样封装的一致性就被破坏了。
mutation ($input: MergePullRequestInput!) {
    mergePullRequest(input: $input) {
        pullRequest {
            id
            number
            state
            merged
            url
        }
    }
}
{
    "input": {
        "pullRequestId": "PR_kwDOB4c8rc6F6nPm",
        "commitHeadline": "合并描述标题",
        "commitBody": "合并描述"
    }
}
{
    "data": {
        "mergePullRequest": {
            "pullRequest": {
                "id": "PR_kwDOB4c8rc6F6nPm",
                "number": 1,
                "state": "MERGED",
                "merged": true,
                "url": "https://github.com/mnyuan2/jquery/pull/1"
            }
        }
    }
}

*/

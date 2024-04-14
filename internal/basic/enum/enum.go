package enum

const (
	// pr 列表 1
	GitEventPullsCreate = 2 // pr创建
	// pr 详情 3
	// pr 更新 4
	// pr 操作日志 5
	// pr commit信息 6
	// pr commit文件列表 7
	// pr 是否合并 8
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
)

var GitEventMap = map[int]string{
	GitEventPullsCreate: "pr创建",
	GitEventPullsMerge:  "pr合并",
}

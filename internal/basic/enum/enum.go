package enum

const (
	GitEventPullsCreate = 1 // pr创建
	GitEventPullsMerge  = 2 // pr合并
)

var GitEventMap = map[int]string{
	GitEventPullsCreate: "pr创建",
	GitEventPullsMerge:  "pr合并",
}

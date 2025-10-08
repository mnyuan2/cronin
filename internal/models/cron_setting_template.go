package models

const (
	TemplateSceneConfigSearch = "config_search" // 模板场景·任务搜索
)

// 模板配置
type TemplateConfig struct {
	Temp string `json:"temp"` // 模板内容
	Hint string `json:"hint"` // 提示文本
}

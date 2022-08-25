package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"io/ioutil"
)

// 解析配置文件
// filePath.文件路径、data.载入到的结构体
func YamlParse(filePath string, data interface{}) error {

	b, err := ioutil.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("文件读取错误 %w", err)
	}
	if err = yaml.Unmarshal(b, data); err != nil {
		return fmt.Errorf("yaml 配置解析错误 %w", err)
	}
	return nil
}

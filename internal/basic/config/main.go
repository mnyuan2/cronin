package config

import (
	"sync"
)

type Main struct {
	Http   *HttpConf   `yaml:"http"`
	Task   *TaskConf   `yaml:"task"`
	Crypto *CryptoConf `yaml:"crypto"`
}

type HttpConf struct {
	Port string `yaml:"port"`
}
type TaskConf struct {
	LogRetention string `yaml:"log_retention"`
}
type CryptoConf struct {
	Secret string `yaml:"secret"`
}

var mainConf Main
var mainOnce sync.Once
var Version = "未定义" // 版本号

func MainConf() Main {
	mainOnce.Do(func() {
		mainConf = Main{}
		if err := YamlParse("configs/main.yaml", &mainConf); err != nil {
			panic(err)
		}
	})
	return mainConf
}

// Local 本机地址
func (m *HttpConf) Local() string {
	return "http://localhost:" + m.Port
}

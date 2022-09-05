package config

import (
	"sync"
)

type Main struct {
	Http *HttpConf `yaml:"http"`
	Task *TaskConf `yaml:"task"`
}

type HttpConf struct {
	Port string `yaml:"port"`
}
type TaskConf struct {
	LogRetention string `yaml:"log_retention"`
}

var mainConf Main
var mainOnce sync.Once

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

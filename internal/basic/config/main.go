package config

import (
	"sync"
)

type Main struct {
	Http   *HttpConf   `yaml:"http"`
	Task   *TaskConf   `yaml:"task"`
	Crypto *CryptoConf `yaml:"crypto"`
	User   *UserConf   `yaml:"user"` // 用户配置，配置后接口访问需要登录
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
type UserConf struct {
	AdminAccount  string `yaml:"admin_account"`
	AdminPassword string `yaml:"admin_password"`
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
		// 配置检测
		if mainConf.Crypto != nil {
			l := len(mainConf.Crypto.Secret)
			if l > 8 {
				panic("配置 crypto.secret 长度必须是8位字符串")
			}
		}
	})
	return mainConf
}

// Local 本机地址
func (m *HttpConf) Local() string {
	return "http://localhost:" + m.Port
}

package config

import (
	"sync"
)

type Main struct {
	UploadRootPath string             `yaml:"upload_root_path"`
	InnerService   map[string]Service `yaml:"inner_service"`
	Token          TokenConfig        `json:"token"`
}

type Service struct {
	BaseUrl string `yaml:"base_url"`
}

type TokenConfig struct {
	Secret string `json:"secret"`
	Expire int    `json:"expire"`
}

var mainConf Main
var mainOnce sync.Once

func MainConf() Main {
	mainOnce.Do(func() {
		mainConf = Main{}
		if err := YamlParse("./configs/main.yaml", &mainConf); err != nil {
			panic(err)
		}
	})
	return mainConf
}

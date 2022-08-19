package config

import "sync"

type EmailConfig struct {
	SenderCli SenderCli `yaml:"sender_cli"`
}

type SenderCli struct {
	UserName string `yaml:"user_name"`
	Addr     string `yaml:"addr"`
	Pass     string `yaml:"pass"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
}

var eConf EmailConfig
var eOnce sync.Once

func EmailConf() EmailConfig {
	eOnce.Do(func() {
		err := YamlParse("configs/email.yaml", &eConf)
		if err != nil {
			panic(err)
		}
	})
	return eConf
}

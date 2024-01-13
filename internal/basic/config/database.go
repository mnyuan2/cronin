package config

import (
	"sync"
)

type DataBaseConf struct {
	Driver string       `yaml:"driver"`
	Mysql  *MysqlSource `yaml:"mysql"`
}

type MysqlSource struct {
	Hostname string `json:"hostname"`
	Database string `json:"database"`
	Username string `json:"username"`
	Password string `json:"password"`
	Port     string `json:"port"`
	Debug    bool   `yaml:"debug"`
}

var dbConf DataBaseConf
var dbOnce sync.Once

func DbConf() DataBaseConf {
	dbOnce.Do(func() {
		err := YamlParse("configs/database.yaml", &dbConf)
		if err != nil {
			panic(err)
		}
	})
	return dbConf
}

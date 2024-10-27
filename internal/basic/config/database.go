package config

import (
	"sync"
)

type DataBaseConf struct {
	Driver string       `yaml:"driver"`
	Mysql  *MysqlSource `yaml:"mysql"`
	Sqlite *Sqlite      `json:"sqlite"`
}

type MysqlSource struct {
	Hostname string `json:"hostname"`
	Database string `json:"database"`
	Username string `json:"username"`
	Password string `json:"password"`
	Port     string `json:"port"`
	Debug    bool   `yaml:"debug"`
}

type Sqlite struct {
	Path  string `json:"path"`
	Debug bool   `yaml:"debug"`
}

var dbConf DataBaseConf
var dbOnce sync.Once

func DbConf() DataBaseConf {
	dbOnce.Do(func() {
		err := YamlParse("configs/database.yaml", &dbConf) // ../../../configs/database.yaml
		if err != nil {
			panic(err)
		}
	})
	return dbConf
}

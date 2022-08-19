package config

import (
	"sync"
)

type DataBaseConf struct {
	Driver string `yaml:"driver"`
	Source string	`yaml:"source"`
	Debug bool	`yaml:"debug"`
}

var dbConf map[string]DataBaseConf
var dbOnce sync.Once

func DbConf()map[string]DataBaseConf  {
	dbOnce.Do(func() {
		err := YamlParse("configs/database.yaml", &dbConf)
		if err != nil {
			panic(err)
		}
	})
	return dbConf
}

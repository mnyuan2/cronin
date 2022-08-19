package config

import "sync"

type elasticConf struct {
	Host string `yaml:"host"`
	Username string `yaml:"username"`
	Password string	`yaml:"password"`
}

var esConf elasticConf
var esOnce sync.Once


func ElasticConf()elasticConf  {
	esOnce.Do(func() {
		esConf = elasticConf{}
		YamlParse("configs/elastic.yaml", &esConf)
	})
	return esConf
}
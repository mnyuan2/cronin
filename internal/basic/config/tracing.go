package config

import "sync"

type tracingConf struct {
	CollectorUrl string	`yaml:"collector_url"`

}

var trOnce sync.Once
var trConf tracingConf

func TracingConf()tracingConf{
	trOnce.Do(func() {
		trConf = tracingConf{}
		if err := YamlParse("configs/tracing.yaml", &trConf); err != nil{
			panic(err)
		}
	})

	return trConf
}

package config

import "sync"

type httpConfs struct {
	Service map[string]HttpConfig
}

type HttpConfig struct {
	SrvName string
	Host string
	Port string
	Addr string
	File string	`description:"文件跟目录"`
}

var hOne sync.Once
var hConf httpConfs

func Http()*httpConfs{
	hOne.Do(func() {
		hConf = httpConfs{}
		if err := YamlParse("configs/http.yaml", &hConf); err != nil{
			panic("配置文件错误："+ err.Error())
		}
	})

	return &hConf
}

// 获得配置
// 如果没有，返回空配置；
func (m *httpConfs) GetConf(serverName string)HttpConfig  {
	conf, ok := m.Service[serverName]
	if !ok {
		panic("服务"+serverName+" 未配置！")
	}
	return conf
}

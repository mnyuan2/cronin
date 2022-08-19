package config

import "sync"

type redisConf struct {
	// redis地址
	Addr string `yaml:"addr"`
	// 库索引
	Db int `yaml:"db"`
	// 密码
	Password string `yaml:"password"`
	// 拨号超时时间，单位秒
	DialTimeout int `yaml:"dial_timeout"`
	// 程序池大小
	PoolSize int `yaml:"pool_size"`
	// 最小空闲数量
	MinIdleConns int `yaml:"min_idle_conns"`
	// 空闲回收时间(分钟)
	IdleTimeout int `yaml:"idle_timeout"`
}

var redisOnce sync.Once
var redisC redisConf

// redis
func Redis() redisConf {
	redisOnce.Do(func() {
		redisC = redisConf{}
		YamlParse("./configs/redis.yaml", &redisC)
	})
	return redisC
}

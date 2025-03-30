package cache

import "sync"

// 默认缓存
var def *Memory

func init() {
	// 目前采用内存缓存，后期根据配置而定
	def = &Memory{
		lock: sync.Map{},
	}
}

// Set 设置缓存
func Set(key string, data any) error {
	return def.Set(key, data)
}

// Get 获取缓存
func Get(key string) any {
	return def.Get(key)
}

func GetAll() map[string]any {
	return def.GetAll()
}

// 删除缓存
func Del(key string) {
	def.Del(key)
}

package cache

import "sync"

var list *memory

func init() {
	// 目前采用内存缓存，后期根据配置而定
	list = &memory{
		lock: sync.Map{},
	}
}

// 添加缓存
//
//	重复添加覆盖之前的值
func Add(key string, data any) error {
	//if _, ok := list.lock.Load(node); ok {
	//	return errors.New("请求正在执行中，请忽重复提交！")
	//}
	list.lock.Store(key, data)
	return nil
}

// 获取缓存
func Get(key string) any {
	val, ok := list.lock.Load(key)
	if !ok {
		return nil
	}
	return val
}

// 删除缓存
func Del(key string) {
	list.lock.Delete(key)
}

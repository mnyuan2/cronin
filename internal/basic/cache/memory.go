package cache

import "sync"

// 内存缓存方案
type Memory struct {
	lock sync.Map
}

// 添加缓存
//
//	重复添加覆盖之前的值
func (m *Memory) Set(key string, data any) error {
	//if _, ok := list.lock.Load(node); ok {
	//	return errors.New("请求正在执行中，请忽重复提交！")
	//}
	m.lock.Store(key, data)
	return nil
}

// 获取缓存
func (m *Memory) Get(key string) any {
	val, ok := m.lock.Load(key)
	if !ok {
		return nil
	}
	return val
}

func (m *Memory) GetAll() map[string]any {
	out := map[string]any{}
	m.lock.Range(func(key, value any) bool {
		out[key.(string)] = value
		return true
	})
	return out
}

// 删除缓存
func (m *Memory) Del(key string) {
	m.lock.Delete(key)
}

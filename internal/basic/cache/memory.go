package cache

import "sync"

// 内存缓存方案
type memory struct {
	lock sync.Map
}

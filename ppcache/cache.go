// Package ppcache
// @author    : MuXiang123
// @time      : 2022/7/27 18:38
package ppcache

import (
	"ppcache/lru"
	"sync"
)

//缓存操作实体，解决并发问题
type cache struct {
	mu         sync.Mutex //互斥锁
	lru        *lru.Cache //缓存数据结构
	cacheBytes int64      //缓存大小
}

//线程安全
func (c *cache) add(key string, value ByteView) {
	c.mu.Lock()         // 上锁
	defer c.mu.Unlock() //最后解锁
	//缓存为空进行初始化 延迟初始化 提高性能
	if c.lru == nil {
		c.lru = lru.New(c.cacheBytes, nil)
	}
	c.lru.Add(key, value)
}

//线程安全
func (c *cache) get(key string) (value ByteView, ok bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.lru == nil {
		return
	}
	if v, ok := c.lru.Get(key); ok {
		return v.(ByteView), ok
	}
	return
}

// Package singleflight
// @author    : MuXiang123
// @time      : 2022/8/1 18:27
// 解决缓存击穿的问题
package singleflight

import "sync"

// Call 表示正在进行中或者已经结束的请求
type Call struct {
	wg  sync.WaitGroup //避免重入
	val interface{}
	err error
}

// Group 主数据结构，管理不同可以的请求
type Group struct {
	mu sync.Mutex       //互斥锁，保护m不被并发读写
	m  map[string]*Call //key和请求的映射
}

// Do 对相同的key，无论Do被调用多少次，函数fn只会被调用一次
func (g *Group) Do(key string, fn func() (interface{}, error)) (interface{}, error) {
	g.mu.Lock()
	if g.m == nil {
		g.m = make(map[string]*Call)
	}
	if c, ok := g.m[key]; ok {
		g.mu.Unlock()
		//如果有请求正在进行则等待
		c.wg.Wait()
		//请求结束返回结果
		return c.val, c.err
	}
	c := new(Call)
	//计数器+1表示这个key已经被请求，请求前加锁
	c.wg.Add(1)
	//添加key和call的映射要保证线程安全
	//表明有key有对应的请求了
	g.m[key] = c
	g.mu.Unlock()

	c.val, c.err = fn()
	//请求结束 计数器-1，
	c.wg.Done()

	//调用结束后删除映射
	g.mu.Lock()
	delete(g.m, key)
	g.mu.Unlock()
	return c.val, c.err
}

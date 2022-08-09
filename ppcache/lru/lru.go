// Package lru
// @author    : MuXiang123
// @time      : 2022/7/26 21:32
package lru

import "container/list"

type Cache struct {
	cache     map[string]*list.Element      //字典 key是string value是双向链表中的指针
	list      *list.List                    //list底层是双向链表
	nBytes    int64                         //当前已经使用的内存
	maxBytes  int64                         //允许使用的最大内存
	OnEvicted func(Key string, value Value) //某个记录被移除时的回调函数
}

//双向链表节点的数据类型
type entry struct {
	//淘汰队首节点时，需要用key从字典中删除对应的映射
	key   string
	value Value
}

// Value 类型的接口
type Value interface {
	Length() int64 //返回所占用的内存大小
}

// New 实例化缓存方法
func New(maxBytes int64, onEvicted func(string, Value)) *Cache {
	return &Cache{
		cache:     make(map[string]*list.Element),
		list:      list.New(),
		maxBytes:  maxBytes,
		OnEvicted: onEvicted,
	}
}

// Get 通过key查找value
func (c *Cache) Get(key string) (value Value, ok bool) {
	//先从字典中找到对应的双向链表节点
	//将该节点移动到队尾
	if ele, ok := c.cache[key]; ok {
		//将链表中的节点ele移动到队尾，约定front为队尾，让队首成为最近最少访问
		c.list.MoveToFront(ele)
		//list存储的是任意类型，这里将链表节点转为对应的数据类型
		kv := ele.Value.(*entry)
		return kv.value, true
	}
	return
}

// RemoveOldest 缓存淘汰，删除最近最少访问的节点队首
func (c Cache) RemoveOldest() {
	ele := c.list.Back()
	if ele != nil {
		c.list.Remove(ele)
		//类型转换
		kv := ele.Value.(*entry)
		//从字典中删除c.cache节点的映射关系
		delete(c.cache, kv.key)
		//更新所用内存大小
		c.nBytes -= int64(len(kv.key)) + kv.value.Length()
		//调用回调函数
		if c.OnEvicted != nil {
			c.OnEvicted(kv.key, kv.value)
		}
	}
}

func (c Cache) Add(key string, value Value) {
	if ele, ok := c.cache[key]; ok {
		//如果键存在，更新对应节点的值，并且将该节点移动到队尾
		c.list.MoveToFront(ele)
		kv := ele.Value.(*entry)
		//更新已用内存的大小
		c.nBytes += value.Length() - kv.value.Length()
		kv.value = value
	} else {
		//不存在，新增key，在队尾添加新节点，并且在字典中添加key和节点的映射关系
		ele := c.list.PushFront(&entry{key, value})
		c.cache[key] = ele
		c.nBytes += int64(len(key)) + value.Length()
	}
	//超出内存大小，先淘汰缓存
	for c.maxBytes != 0 && c.maxBytes < c.nBytes {
		if c.Length() == 0 {
			panic("The first element exceeds the max bytes")
		}
		c.RemoveOldest()
	}
}

// Length 实现length方法，获取添加了多少条数据 方便测试
func (c *Cache) Length() int64 {
	return int64(c.list.Len())
}

// Clear 从缓存中清除所有存储
func (c *Cache) Clear() {
	if c.OnEvicted != nil {
		for _, e := range c.cache {
			kv := e.Value.(*entry)
			c.OnEvicted(kv.key, kv.value)
		}
	}
	c.list = nil
	c.cache = nil
}

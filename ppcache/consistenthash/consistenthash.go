// Package consistenthash
// @author    : MuXiang123
// @time      : 2022/7/30 22:06
package consistenthash

import (
	"hash/crc32"
	"sort"
	"strconv"
	"sync"
	"sync/atomic"
)

//Hash 将字节映射到无符号32位整形
type Hash func(data []byte) uint32

// Map 一致性哈希的主要数据结构，
type Map struct {
	sync.Mutex
	hash     Hash         //哈希函数
	replicas int          //虚拟节点倍数
	values   atomic.Value //原子的存取keys和hashMap
}

type values struct {
	keys    []int          // 哈希环
	hashMap map[int]string //虚拟节点和真实节点的映射表，key为虚拟节点，value为真实节点
}

// New 创建一个map实例 允许自定义虚拟节点倍数和hash函数
func New(replicas int, fn Hash) *Map {
	m := &Map{
		replicas: replicas,
		hash:     fn,
	}
	m.values.Store(&values{
		hashMap: make(map[int]string),
	})
	if m.hash == nil {
		//生成hashcode的默认算法
		m.hash = crc32.ChecksumIEEE
	}
	return m
}

// Add 传入0个或多个真实节点的名称
func (m *Map) Add(keys ...string) {
	m.Lock()
	defer m.Unlock()
	newValues := m.loadValues()
	for _, key := range keys {
		// 对每个 key(节点) 创建 m.replicas 个虚拟节点
		for i := 0; i < m.replicas; i++ {
			//生成虚拟节点的hash
			hash := int(m.hash([]byte(strconv.Itoa(i) + key)))
			//添加到环上
			newValues.keys = append(newValues.keys, hash)
			newValues.hashMap[hash] = key
		}
	}
	//对key排序
	sort.Ints(newValues.keys)
	m.values.Store(newValues)
}

// Get 传入缓存的key，返回真实节点
func (m *Map) Get(key string) string {
	values := m.loadValues()
	if len(values.keys) == 0 {
		return ""
	}
	//计算key的hash值
	hash := int(m.hash([]byte(key)))
	//顺时针找到第一个匹配虚拟节点的下标
	idx := sort.Search(len(values.keys), func(i int) bool {
		//返回最小索引的前提条件
		return values.keys[i] >= hash
	})
	//返回真实节点
	// 如果 idx == len(m.keys)，说明应选择 m.keys[0]，
	// 因为 m.keys 是一个环状结构，用取余数的方式来处理这种情况
	return values.hashMap[values.keys[idx%len(values.keys)]]
}

// Remove 用于删除keys和map上的节点及其虚拟节点
func (m *Map) Remove(key string) {
	m.Lock()
	defer m.Unlock()
	newValues := m.loadValues()

	for i := 0; i < m.replicas; i++ {
		hash := int(m.hash([]byte(strconv.Itoa(i) + key)))
		idx := sort.SearchInts(newValues.keys, hash)
		newValues.keys = append(newValues.keys[:idx], newValues.keys[idx+1:]...)
		delete(newValues.hashMap, hash)
	}
}

//从原子容器中加载values
func (m *Map) loadValues() *values {
	return m.values.Load().(*values)
}

func (m *Map) copyValues() *values {
	oldValues := m.loadValues()
	newValues := &values{
		keys:    make([]int, len(oldValues.keys)),
		hashMap: make(map[int]string),
	}
	copy(newValues.keys, oldValues.keys)
	for k, v := range oldValues.hashMap {
		newValues.hashMap[k] = v
	}
	return newValues
}

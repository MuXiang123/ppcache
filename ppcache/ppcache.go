// Package ppcache
// @author    : MuXiang123
// @time      : 2022/7/27 22:05
// 回调函数，当缓存不存在时，调用这个函数就可以得到源数据
package ppcache

import (
	"fmt"
	"log"
	pb "ppcache/ppcachepb"
	"ppcache/singleflight"
	"sync"
)

// Group 缓存的命名空间 负责与用户的交互，并且控制缓存值存储和获取的流程。
type Group struct {
	name      string              // 唯一名称
	getter    Getter              //缓存未命中时获取诗句的回调函数
	mainCache cache               //并发缓存实体
	peers     PeerPicker          //节点
	loader    *singleflight.Group //使用singleFilght 保证每个key只能获取一次
}

// Getter 通过key获取数据
type Getter interface {
	Get(key string) ([]byte, error)
}

// GetterFunc 实现getter方法
type GetterFunc func(key string) ([]byte, error)

// Get 实现getter接口的回调函数
func (f GetterFunc) Get(key string) ([]byte, error) {
	return f(key)
}

var (
	mu     sync.RWMutex
	groups = make(map[string]*Group)
)

func NewGroup(name string, cacheBytes int64, getter Getter) *Group {
	if getter == nil {
		panic("nil Getter")
	}
	mu.Lock()
	defer mu.Unlock()
	g := &Group{
		name:      name,
		getter:    getter,
		mainCache: cache{cacheBytes: cacheBytes},
		loader:    &singleflight.Group{},
	}
	groups[name] = g
	return g
}

// GetGroup 获取命名空间
func GetGroup(name string) *Group {
	mu.RLock()        //只读锁，当前函数不涉及写操作
	g := groups[name] //通过name获取实际的命名空间
	mu.RUnlock()
	return g
}

func (g *Group) Get(key string) (ByteView, error) {
	//key为空时返回空
	if key == "" {
		return ByteView{}, fmt.Errorf("key is require")
	}
	//从mainCache中查找缓存，如果存在就返回缓存值
	if v, ok := g.mainCache.get(key); ok {
		log.Println("[PPCache] hit")
		return v, nil
	}
	//不存在
	return g.load(key)
}

//load 缓存不存在时调用
func (g *Group) load(key string) (value ByteView, err error) {
	//无论并发多少次，每个key只能在同一时刻只能获取一次
	viewi, err := g.loader.Do(key, func() (interface{}, error) {
		if g.peers != nil {
			if peer, ok := g.peers.PickPeer(key); ok {
				if value, err = g.getFromPeer(peer, key); err == nil {
					return value, nil
				}
				log.Println("[GeeCache] Failed to get from peer", err)
			}
		}
		return g.getLocally(key)
	})
	if err == nil {
		return viewi.(ByteView), nil
	}
	return
}

// getLocally 调用用户的回调函数获取数据源，并且将数据院添加到缓存中
func (g *Group) getLocally(key string) (ByteView, error) {
	bytes, err := g.getter.Get(key)
	if err != nil {
		return ByteView{}, err
	}
	value := ByteView{b: cloneBytes(bytes)}
	g.populateCache(key, value)
	return value, nil
}

//往缓存填充key value
func (g *Group) populateCache(key string, value ByteView) {
	g.mainCache.add(key, value)
}

// RegisterPeers 注册节点
func (g *Group) RegisterPeers(peers PeerPicker) {
	if g.peers != nil {
		panic("RegisterPeerPicker called more than once")
	}
	g.peers = peers
}

//从节点中获取缓存
func (g *Group) getFromPeer(peer PeerGetter, key string) (ByteView, error) {
	//bytes, err := peer.Get(g.name, key)
	//改为使用protobuf进行通信
	req := &pb.Request{
		Group: g.name,
		Key:   key,
	}
	res := &pb.Response{}
	err := peer.Get(req, res)
	if err != nil {
		return ByteView{}, err
	}
	return ByteView{b: res.Value}, nil
}

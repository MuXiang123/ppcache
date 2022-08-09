![img.png](img.png)

# PPCache

## Introduction - 介绍

### 概要

基于Go语言实现的的分布式缓存。

### 功能

- 使用LRU算法，解决资源限制问题。 
- 使用互斥锁，实现单机并发功能，解决资源竞争问题。 
- 实现一致性哈希算法，解决远程节点的挑选问题。 
- 实现多节点间通过HTTP通信，解决节点间的通信问题。 
- 实现singlefight，解决缓存击穿问题。

## 环境

- GO 1.18版本

## 快速入门

### 使用demo运行

1. ``git clone https://gitee.com/MuXiang123/ppcache.git``
2. `./run.sh` 

### 具体使用

```
go build -o server
./server -port=8001 //开启多个端口
curl "http://localhost:9999/api?key=Tom" // 测试缓存
```

### 目录说明

```
CACHE
├─main.go  // 测试使用，其中模拟了数据库
└─ppcache  // 主要目录
	├─byteview.go // 并发读时的副本
	├─cache.go    // 缓存操作实体
	├─http.go    //  httpt通信
	├─peers.go    // 分布式节点
	├─ppcache.go  //缓存命名空间
    ├─consistenthash  // 一致性哈希 算法
    ├─lru             // lru缓存
    ├─ppcachepb       // protobuf通信
    └─singleflight    //singleflight防止缓存击穿

```



​	
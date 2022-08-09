// Package ppcache
// @author    : MuXiang123
// @time      : 2022/7/31 10:29
package ppcache

import pb "ppcache/ppcachepb"

// PeerPicker
// 每个节点都有自己特有的key
type PeerPicker interface {
	// PickPeer 根据传入key选择相应的节点
	PickPeer(key string) (peer PeerGetter, ok bool)
}

// PeerGetter 类似于http客户端的作用。从对应的group查找缓存值
type PeerGetter interface {
	//Get(group string, key string) ([]byte, error)
	//使用protobuf进行通信
	Get(in *pb.Request, out *pb.Response) error
}

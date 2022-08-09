// Package ppcache
// @author    : MuXiang123
// @time      : 2022/7/28 18:08
//http服务端
package ppcache

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"ppcache/consistenthash"
	pb "ppcache/ppcachepb"
	"strings"
	"sync"
)

const (
	defaultBasePath = "/_ppcache/" //默认基础路径
	defaultReplicas = 50           //默认虚拟节点倍数50
)

// HTTPPool HTTP通信的数据结构
type HTTPPool struct {
	//节点的url：https://example.net:8000
	self       string                 //主机/ip和端口号
	basePath   string                 //节点间通信地址的前缀
	mu         sync.Mutex             //互斥锁
	peers      *consistenthash.Map    //根据具体的key 选择节点
	httpGetter map[string]*httpGetter //映射远程节点和对应的httpGetter key : http://10.0.0.2:8008
}

// NewHTTPPool 初始化服务端数据
func NewHTTPPool(self string) *HTTPPool {
	return &HTTPPool{
		self:     self,
		basePath: defaultBasePath,
	}
}

//Log 带有服务器名称的信息
func (p *HTTPPool) Log(format string, v ...interface{}) {
	log.Printf("[Server %s] %s\n", p.self, fmt.Sprintf(format, v...))
}

//使用 proto.Marshal() 编码 HTTP 响应
func (p *HTTPPool) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//判断basePath路径是否存在
	if !strings.HasPrefix(r.URL.Path, p.basePath) {
		panic("HTTPPool serving unexpected path: " + r.URL.Path)
	}
	p.Log("%s, %s", r.Method, r.URL.Path)
	//判断是否存在groupName和key，约定访问路径为/<basepath>/<groupname>/<key>
	parts := strings.SplitN(r.URL.Path[len(p.basePath):], "/", 2)
	if len(parts) != 2 {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	groupName := parts[0]
	key := parts[1]

	//通过groupName得到group实例
	group := GetGroup(groupName)
	if group == nil {
		http.Error(w, "no such group: "+groupName, http.StatusNotFound)
		return
	}
	view, err := group.Get(key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// 将value 写入响应体中，
	body, err := proto.Marshal(&pb.Response{Value: view.ByteSlice()})
	//查缓存
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	//设置请求头
	w.Header().Set("Content-Type", "application/octet-stream")
	//放入response body
	w.Write(body)
}

//客户端
type httpGetter struct {
	baseURL string
}

// Get 从远程节点中获取缓存,使用proto.Unmarshal() 解码 HTTP 响应
func (h *httpGetter) Get(in *pb.Request, out *pb.Response) error {
	//打印访问远程节点的url
	u := fmt.Sprintf(
		"%v%v/%v",
		h.baseURL,
		url.QueryEscape(in.GetGroup()),
		url.QueryEscape(in.GetKey()),
	)
	//获取返回值
	res, err := http.Get(u)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("server return: %v", res.Status)
	}
	//转换为byte类型
	bytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("reading response body: %v", err)
	}
	if err = proto.Unmarshal(bytes, out); err != nil {
		return fmt.Errorf("decoding response body: %v", err)
	}
	return nil
}

var _ PeerGetter = (*httpGetter)(nil)

//Set 更新节点
func (p *HTTPPool) Set(peers ...string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	//实例化一致性哈希算法
	p.peers = consistenthash.New(defaultReplicas, nil)
	p.peers.Add(peers...)
	//为每个节点创建http客户端
	p.httpGetter = make(map[string]*httpGetter, len(peers))
	for _, peer := range peers {
		p.httpGetter[peer] = &httpGetter{
			baseURL: peer + p.basePath,
		}
	}
}

// PickPeer 根据key选择节点，返回节点对应的http客户端
func (p *HTTPPool) PickPeer(key string) (PeerGetter, bool) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if peer := p.peers.Get(key); peer != "" && peer != p.self {
		p.Log("pick peer %s", peer)
		return p.httpGetter[peer], true
	}
	return nil, false
}

var _ PeerPicker = (*HTTPPool)(nil)

// @author    : MuXiang123
// @time      : 2022/7/28 18:45
package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"ppcache"
)

//模拟数据库
var db = map[string]string{
	"Tom":  "630",
	"Jack": "589",
	"Sam":  "567",
}

func createGroup() *ppcache.Group {
	return ppcache.NewGroup("scores", 2<<10, ppcache.GetterFunc(
		func(key string) ([]byte, error) {
			log.Println("[slowDB] search key", key)
			if v, ok := db[key]; ok {
				return []byte(v), nil
			}
			return nil, fmt.Errorf("%s not exist", key)
		},
	))
}

//启动缓存服务器，创建HTTPPool 添加节点信息，注册到pp中，启动http服务
func startCacheServer(addr string, addrs []string, pp *ppcache.Group) {
	peers := ppcache.NewHTTPPool(addr)
	peers.Set(addrs...)
	pp.RegisterPeers(peers)
	log.Println("ppCache is running at", addr)
	log.Fatal(http.ListenAndServe(addr[7:], peers))
}

//启动一个api服务 和用户交互
func startAPIServer(apiAddr string, pp *ppcache.Group) {
	http.Handle("/api", http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			key := r.URL.Query().Get("key")
			view, err := pp.Get(key)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/octet-stream")
			w.Write(view.ByteSlice())
		}))
	log.Println("fontend server is running at", apiAddr)
	log.Fatal(http.ListenAndServe(apiAddr[7:], nil))
}

func main() {
	var port int
	var api bool
	flag.IntVar(&port, "port", 8001, "PPcache server port")
	flag.BoolVar(&api, "api", false, "Start a api server")
	flag.Parse()

	apiAddr := "http://localhost:9999"
	addrMap := map[int]string{
		8001: "http://localhost:8001",
		8002: "http://localhost:8002",
		8003: "http://localhost:8003",
	}

	var addrs []string
	for _, v := range addrMap {
		addrs = append(addrs, v)
	}

	pp := createGroup()
	if api {
		go startAPIServer(apiAddr, pp)
	}
	startCacheServer(addrMap[port], []string(addrs), pp)
}

//func main() {
//	//初始化group 并且写好查询数据库的方法
//	ppcache.NewGroup("scores", 2<<10, ppcache.GetterFunc(
//		func(key string) ([]byte, error) {
//			log.Println("[SlowDB] search key", key)
//			if v, ok := db[key]; ok {
//				return []byte(v), nil
//			}
//			return nil, fmt.Errorf("%s not exist", key)
//		}))
//
//	addr := "localhost:9999"
//	peers := ppcache.NewHTTPPool(addr)
//	log.Println("geecache is running at", addr)
//	log.Fatal(http.ListenAndServe(addr, peers))
//}

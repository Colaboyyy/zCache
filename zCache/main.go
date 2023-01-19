package main

import (
	"flag"
	"fmt"
	"log"
	"main/zCache"
	"net/http"
)

var db = map[string]string{
	"Tom":  "630",
	"Jack": "589",
	"Sam":  "567",
}

func createGroup() *zCache.Group {
	return zCache.NewGroup("scores", 2<<10, zCache.GetterFunc(
		func(key string) ([]byte, error) {
			log.Println("[[SlowDB] search key", key)
			if v, ok := db[key]; ok {
				return []byte(v), nil
			}
			return nil, fmt.Errorf("%s not exist", key)
		}))
}

// 启动缓存服务器：创建 HTTPPool，添加节点信息，注册到 http group 中，启动 HTTP 服务（共3个端口，8001/8002/8003），用户不感知。
func startCacheServer(addr string, addrs []string, zCacheGroup *zCache.Group) {
	peers := zCache.NewHTTPPool(addr)
	peers.Set(addrs...)
	zCacheGroup.RegisterPeers(peers)
	log.Println("zCache is running at:", addr)
	log.Fatal(http.ListenAndServe(addr[7:], peers))
}

// 启动一个 API 服务（端口 9999），与用户进行交互，用户感知。
func startAPIServer(apiAddr string, group *zCache.Group) {
	http.Handle("/api", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		key := r.URL.Query().Get("key")
		view, err := group.Get(key)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/octet-stream")
		w.Write(view.ByteSlice())
	}))
	log.Println("frontend server is running at:", apiAddr)
	log.Fatal(http.ListenAndServe(apiAddr[7:], nil))
}

// 命令行传入 port 和 api 2 个参数，用来在指定端口启动 HTTP 服务。选择api为1的作为服务端
func main() {
	var port int
	var api bool
	flag.IntVar(&port, "port", 8001, "zCache server port")
	flag.BoolVar(&api, "api", false, "Start a api Server")
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
	group := createGroup()
	if api {
		go startAPIServer(apiAddr, group)
	}
	startCacheServer(addrMap[port], []string(addrs), group)
}

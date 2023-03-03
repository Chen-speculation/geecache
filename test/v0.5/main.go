package main

import (
	"GeeCacheV0.5/gee_cache"
	"GeeCacheV0.5/utils"
	"flag"
	"fmt"
	"log"
	"net/http"
)

var db = map[string]string{
	"Tom":  "630",
	"Jack": "589",
	"Sam":  "567",
}

// 创建一个名为"scores"的缓存分组，缓存大小为 1024 字节
func createGroup() gee_cache.IGroup {
	return gee_cache.NewGroupByGetterFunc("scores", 1<<10, func(key string) ([]byte, error) {
		utils.Logger().Println("[SlowDB] search key", key)
		if v, ok := db[key]; ok {
			return []byte(v), nil
		}
		return nil, fmt.Errorf("%s not exist", key)
	})
}

// 创建一个gee的HTTP服务器，并将集群的addrs加进来组成一个HTTP的服务器分布式集群
func startCacheServer(addr string, addrs []string, gee gee_cache.IGroup) {
	peer := gee_cache.NewHTTPServer(addr)
	peer.Set(addrs...)
	gee.RegisterPeers(peer)
	utils.Logger().Println("gee_cache is running at", addr)
	if err := peer.Run(); err != nil {
		log.Fatal(err)
	}
}

func startAPIServer(apiAddr string, gee gee_cache.IGroup) {
	http.Handle("/api", http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			key := r.URL.Query().Get("key")
			view, err := gee.Get(key)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/octet-stream")
			w.Write(view.ByteSlice())
		}))
	utils.Logger().Println("HTTP server is running at", apiAddr)
	log.Fatal(http.ListenAndServe(apiAddr[7:], nil)) //监听端口
}

func main() {
	var port int
	var api bool
	flag.IntVar(&port, "port", 8001, "GeeCache server port")
	flag.BoolVar(&api, "api", false, "Start a api server?")
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

	gee := createGroup()
	if api { //这个算是一个多余的东西(与用户进行交互，用户感知。)
		go startAPIServer(apiAddr, gee)
	}
	startCacheServer(addrMap[port], addrs, gee)
}

package main

import (
	"GeeCacheV0.3/gee_cache"
	"GeeCacheV0.3/utils"
	"fmt"
)

var db = map[string]string{
	"Tom":  "630",
	"Jack": "589",
	"Sam":  "567",
}

func main() {
	gee_cache.NewGroupByGetterFunc("score", 1<<10, func(key string) ([]byte, error) {
		utils.Logger().Println("[SlowDB] search key", key)
		if v, ok := db[key]; ok {
			return []byte(v), nil
		}
		return nil, fmt.Errorf("%s not exist", key)
	})
	//HTTP服务
	peer := gee_cache.NewHTTPServer("localhost:9999")
	if err := peer.Run(); err != nil {
		panic(err)
	}
}

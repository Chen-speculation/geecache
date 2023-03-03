package main

import (
	"GeeCacheV0.1/lru"
	"fmt"
)

type String string

func (d String) Len() int {
	return len(d)
}

func main() {
	lru := lru.New(int64(10), nil)
	lru.Add("key1", String("1234"))
	lru.Add("key3", String("5555"))
	lru.Add("key2", String("6666"))
	if v, ok := lru.Get("key1"); ok {
		fmt.Println(v)
	}
	if val, ok := lru.Get("key2"); !ok {
		fmt.Println("can not found key2")
	} else {
		fmt.Println(val)
	}
}

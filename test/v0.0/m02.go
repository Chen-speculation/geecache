package main

import (
	"fmt"
	"sync"
	"time"
)

var m sync.Mutex
var set = make(map[int]bool, 0)

func printOnce(i, num int) {
	m.Lock()
	if _, exist := set[num]; !exist { //set存在共享读写，都得加锁
		fmt.Println(num)
	}
	set[num] = true
	m.Unlock()
	fmt.Println(i, "退出")
}

func main() {
	for i := 0; i < 10; i++ {
		go printOnce(i, 100)
	}
	time.Sleep(10 * time.Second)
}

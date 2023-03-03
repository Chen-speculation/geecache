package main

import "sync"

func main() {
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		go func(i int) {
			println("哈哈哈", i)
			wg.Done()
		}(i)
		wg.Add(1)
	}
	wg.Wait()
}

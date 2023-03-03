package main

import "fmt"

func main() {
	ok := t1()
	fmt.Println(ok)
}

func t1() (ok bool) {
	//上面的ok和这个ok不在用一个作用域
	if elem, ok := getOk(); ok {
		fmt.Println(elem)
	}
	return
}

func getOk() (int, bool) {
	return 1, true
}

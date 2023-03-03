package main

import "log"

func aaa(format string, v ...any) {
	log.Printf(format, v...)
}

func main() {
	aaa("%s %s", "111", "222")
}

package lru

type Value interface {
	Len() int //这个值的占多少字节
}

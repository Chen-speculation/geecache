package lru

import (
	"GeeCacheV0.7/utils"
	"container/list"
)

type ICache interface {
	Get(key string) (value Value, ok bool)
	Add(key string, value Value)
	Len() int
}

type cache struct {
	maxBytes   int64 //最大可以缓存多少字节
	curBytes   int64 //当前缓存的多少字节
	linkedList *list.List
	cacheMap   map[string]*list.Element
	onEvicted  func(key string, value Value) //可选的，在条目被清除时执行。evicted：驱逐
}

func New(maxBytes int64, onEvicted func(key string, value Value)) ICache {
	if maxBytes <= 0 {
		utils.Logger().Errorln("Your cache MaxBytes can not set with", maxBytes)
		panic("maxBytes should greater than 0")
	}
	utils.Logger().Println("create LRUCache successfully")
	return &cache{
		maxBytes:   maxBytes,
		linkedList: list.New(),
		cacheMap:   make(map[string]*list.Element),
		onEvicted:  onEvicted,
	}
}

// Get 获取元素
func (c *cache) Get(key string) (value Value, ok bool) {
	if elem, ok := c.cacheMap[key]; ok {
		c.linkedList.MoveToFront(elem)
		return elem.Value.(*entry).value, true
	}
	return
}

// Add 缓存添加
func (c *cache) Add(key string, value Value) {
	//判断是更新操作还是移除操作
	if elem, ok := c.cacheMap[key]; ok { //更新操作
		c.linkedList.MoveToFront(elem)                    //移动到头
		kv := elem.Value.(*entry)                         //类型断言转kv
		c.curBytes += int64(value.Len() - kv.value.Len()) //更新curBytes
		kv.value = value                                  //更新值
		utils.Logger().Infof("Update [key=%s,value=%v] successfully,current bytes is %d", key, value, c.curBytes)
	} else { //添加操作
		c.cacheMap[key] = c.linkedList.PushFront(newEntry(key, value)) //加入链表中、加入map中
		c.curBytes += int64(len(key) + value.Len())                    //更改curBytes值
		utils.Logger().Infof("Add [key=%s,value=%v] successfully,current bytes is %d", key, value, c.curBytes)
	}

	//清除缓存，知道maxBytes >= curBytes
	for c.maxBytes < c.curBytes {
		c.removeOldest()
	}

}

// Len 当前缓存个数
func (c *cache) Len() int {
	return c.linkedList.Len()
}

func (c *cache) removeOldest() {
	list := c.linkedList

	//获取最后一个元素
	revElem := list.Back()
	if revElem == nil { //无元素则之间返回
		return
	}

	//删除
	list.Remove(revElem)
	e := revElem.Value.(*entry)
	delete(c.cacheMap, e.key)
	c.curBytes -= int64(len(e.key) + e.value.Len())
	utils.Logger().Infof("remove [key=%s,value=%v] successfully,current bytes is %d", e.key, e.value, c.curBytes)

	//驱逐元素之后可以执行的操作
	if c.onEvicted != nil {
		utils.Logger().Printf("begin to call onEvicted![key=%s,value=%s]", e.key, e.value)
		c.onEvicted(e.key, e.value)
	}
}

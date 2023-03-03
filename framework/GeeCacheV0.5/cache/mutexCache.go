package cache

import (
	"GeeCacheV0.5/lru"
	"sync"
)

type IMutexCache interface {
	Add(key string, value ByteView)
	Get(key string) (value ByteView, ok bool)
}

// 由于使用了sync.Mutex，下面方法必须使用指针，不然mu失效
type mutexCache struct {
	mu         sync.Mutex
	lru        lru.ICache
	cacheBytes int64
}

func NewMutexCache(cacheBytes int64) IMutexCache {
	return &mutexCache{
		cacheBytes: cacheBytes,
	}
}

func (c *mutexCache) Add(key string, value ByteView) {
	//上锁、解锁
	c.mu.Lock()
	defer c.mu.Unlock()

	//判断lru是否初始化
	if c.lru == nil {
		//TODO 这里没有使用onEvicted，用户应该可以自定义才对
		c.lru = lru.New(c.cacheBytes, nil)
	}

	//添加缓存数据
	c.lru.Add(key, value)
}

func (c *mutexCache) Get(key string) (value ByteView, ok bool) {
	//上锁、解锁
	c.mu.Lock()
	defer c.mu.Unlock()

	//判断lru是否初始化
	if c.lru == nil {
		return
	}

	//获取缓存
	val, ok := c.lru.Get(key)
	if ok {
		value = val.(ByteView)
	}
	return
}

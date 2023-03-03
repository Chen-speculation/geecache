package cache

import (
	"GeeCacheV0.2/lru"
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
	cacheBytes int64 //这个缓存最大值
}

func NewMutexCache(cacheBytes int64) IMutexCache {
	return &mutexCache{
		lru:        lru.New(cacheBytes, nil),
		cacheBytes: cacheBytes,
	}
}

func (c *mutexCache) Add(key string, value ByteView) {
	//上锁、解锁
	c.mu.Lock()
	defer c.mu.Unlock()

	//添加缓存数据
	c.lru.Add(key, value)
}

func (c *mutexCache) Get(key string) (value ByteView, ok bool) {
	//上锁、解锁
	c.mu.Lock()
	defer c.mu.Unlock()

	//获取缓存
	val, ok := c.lru.Get(key)
	if ok {
		value = val.(ByteView)
	}
	return
}

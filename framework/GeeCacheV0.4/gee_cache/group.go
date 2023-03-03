package gee_cache

import (
	"GeeCacheV0.4/cache"
	"GeeCacheV0.4/utils"
	"errors"
	"sync"
)

type IGroup interface {
	Get(key string) (cache.ByteView, error)
}

// Group 命名空间、被加载的相关数据(例子：缓存学生成绩是一个Group)
type Group struct {
	name       string
	getter     Getter
	mutexCache cache.IMutexCache //因为我们使用的group都是指针的，所以这一块不会发生深拷贝而变化
}

var (
	mu     sync.RWMutex
	groups = make(map[string]IGroup)
)

// NewGroupByGetter 使用接口的
func NewGroupByGetter(name string, cacheBytes int64, getter Getter) (g IGroup) {
	if getter == nil {
		utils.Logger().Errorln("getter can not be nil")
		panic("getter is nil!")
	}
	mu.Lock()
	defer mu.Unlock()
	g = &Group{
		name:       name,
		getter:     getter,
		mutexCache: cache.NewMutexCache(cacheBytes),
	}
	groups[name] = g
	return
}

// NewGroupByGetterFunc 使用函数的
func NewGroupByGetterFunc(name string, cacheBytes int64, f GetterFunc) (g IGroup) {
	return NewGroupByGetter(name, cacheBytes, f)
}

func (g *Group) Get(key string) (cache.ByteView, error) {
	if len(key) == 0 {
		utils.Logger().Warningln("Your key len = 0")
		return cache.ByteView{}, errors.New("your key len = 0")
	}
	if val, ok := g.mutexCache.Get(key); ok {
		utils.Logger().Println("[GeeCache] Hit key is", key)
		return val, nil
	}
	return g.load(key)
}

func (g *Group) load(key string) (cache.ByteView, error) {
	return g.getLocally(key)
}

func (g *Group) getLocally(key string) (cache.ByteView, error) {
	bytes, err := g.getter.Get(key)
	if err != nil {
		return cache.ByteView{}, err
	}
	byteView := cache.NewByteView(bytes) //进行一次深拷贝，避免数据被恶意篡改
	g.populateCache(key, byteView)       //把数据加入缓存
	return byteView, nil
}

func (g *Group) populateCache(key string, view cache.ByteView) {
	g.mutexCache.Add(key, view)
}

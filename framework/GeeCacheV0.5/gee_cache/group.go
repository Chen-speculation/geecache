package gee_cache

import (
	"GeeCacheV0.5/cache"
	"GeeCacheV0.5/utils"
	"errors"
	"sync"
)

type IGroup interface {
	Get(key string) (cache.ByteView, error)
	RegisterPeers(peers PeerGetter)
}

// Group 命名空间、被加载的相关数据(例子：缓存学生成绩是一个Group)
type Group struct {
	name       string
	getter     Getter
	mainCache  cache.IMutexCache
	peerGetter PeerGetter //根据缓存的key可以获取主机的客户端
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
		name:      name,
		getter:    getter,
		mainCache: cache.NewMutexCache(cacheBytes),
	}
	groups[name] = g
	return
}

// NewGroupByGetterFunc 使用函数的
func NewGroupByGetterFunc(name string, cacheBytes int64, f GetterFunc) (g IGroup) {
	return NewGroupByGetter(name, cacheBytes, f)
}

// RegisterPeers registers a PeerGetter for choosing remote peer
func (g *Group) RegisterPeers(peers PeerGetter) {
	if g.peerGetter != nil {
		utils.Logger().Errorln("RegisterPeerPicker called more than once")
		panic("RegisterPeerPicker called more than once")
	}
	g.peerGetter = peers
}

func (g *Group) Get(key string) (cache.ByteView, error) {
	if len(key) == 0 {
		utils.Logger().Warningln("Your key len = 0")
		return cache.ByteView{}, errors.New("your key len = 0")
	}
	if val, ok := g.mainCache.Get(key); ok {
		utils.Logger().Println("[GeeCache] Hit key is", key)
		return val, nil
	}
	return g.load(key) //调用回调函数
}

func (g *Group) load(key string) (value cache.ByteView, err error) {
	if g.peerGetter != nil {
		if peerClient, ok := g.peerGetter.GetPeer(key); ok {
			if value, err = g.getFromPeer(peerClient, key); err == nil {
				return
			}
			//远程获取失败就继续从本地获取试一试
			utils.Logger().Println("[GeeCache] Failed to get from peer", err)
		}
	}
	return g.getByLocal(key)
}

func (g *Group) getByLocal(key string) (cache.ByteView, error) {
	bytes, err := g.getter.Get(key) //使用回调函数获取数据
	if err != nil {
		return cache.ByteView{}, err
	}
	byteView := cache.NewByteView(bytes) //进行一次深拷贝，避免数据被恶意篡改
	g.mainCache.Add(key, byteView)       //把数据加入缓存
	return byteView, nil
}

func (g *Group) getFromPeer(peer PeerClient, key string) (cache.ByteView, error) {
	bytes, err := peer.Get(g.name, key)
	if err != nil {
		return cache.ByteView{}, err
	}
	//TODO 从远程获取数据之后，没有存放在本地
	return cache.NewByteView(bytes), nil
}

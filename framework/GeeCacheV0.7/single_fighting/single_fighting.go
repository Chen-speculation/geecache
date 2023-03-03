package single_fighting

import "sync"

type IGroup interface {
	Do(key string, fn func() (any, error)) (any, error)
}

type call struct { //代表[正在进行中]或[已经结束]的请求。使用 sync.WaitGroup 锁避免重入。
	wg  sync.WaitGroup
	val any
	err error
}

type Group struct {
	mu sync.Mutex //保护下面的map
	m  map[string]*call
}

func NewGroup() IGroup {
	//还是开始的时候初始化一下爽一点
	return &Group{m: map[string]*call{}}
}

// Do 的作用就是，针对相同的 key，无论 Do 被调用多少次，函数 fn 都只会被调用一次，等待 fn 调用结束了，返回返回值或错误。
func (g *Group) Do(key string, fn func() (any, error)) (any, error) {
	g.mu.Lock()

	//这个函数正在调用中
	if c, ok := g.m[key]; ok {
		g.mu.Unlock() //读取完了，可以释放锁了
		c.wg.Wait()
		return c.val, c.err
	}

	//开始调用这个函数
	c := new(call)
	c.wg.Add(1)
	g.m[key] = c
	g.mu.Unlock() //对m操作完了，可以释放锁了

	c.val, c.err = fn()
	c.wg.Done()

	g.mu.Lock()
	delete(g.m, key)

	return c.val, c.err
}

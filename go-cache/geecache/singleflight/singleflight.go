package singleflight

import "sync"

/**
缓存雪崩avalanche: 缓存在同一时间同时失效,导致大量请求涌入到DB,压力陡增,引发雪崩
	一般是由于缓存服务宕机或者大量缓存同一时间过期
缓存击穿breakdown: 一个存在的key 在缓存过期的一刻,有大量的请求过来,请求打到DB
缓存穿透penetration: 大量请求查询不存在的key,请求都走DB
*/

// 正在进行中或者已经结束的请求
type call struct {
	wg  sync.WaitGroup
	val interface{}
	err error
}
type Group struct {
	mu sync.Mutex
	m  map[string]*call
}

// Do 相同的key 保证同一时间内,fn只会被调用一次
func (g *Group) Do(key string, fn func() (interface{}, error)) (interface{}, error) {
	g.mu.Lock()
	if g.m == nil {
		g.m = make(map[string]*call)
	}
	// 已经请求过或这正请求
	if c, ok := g.m[key]; ok {
		g.mu.Unlock()
		c.wg.Wait() // 等待请求完成
		return c.val, c.err
	}
	// 新请求
	c := new(call)
	c.wg.Add(1)
	g.m[key] = c
	g.mu.Unlock()

	c.val, c.err = fn() // 调用fn发起请求
	c.wg.Done()

	finish(g, key)
	return c.val, c.err
}

func finish(g *Group, key string) {
	g.mu.Lock()
	delete(g.m, key)
	g.mu.Unlock()
}

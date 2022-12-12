package geecache

import (
	"log"
	"sync"
)

// Getter define an interface to get data from datasource
// maybe file ,database,web and so on
type Getter interface {
	Get(key string) ([]byte, error)
}

// GetterFunc 函数类型实现接口, 接口型函数
type GetterFunc func(key string) ([]byte, error)

func (f GetterFunc) Get(key string) ([]byte, error) {
	return f(key)
}

type Group struct {
	name      string
	getter    Getter
	mainCache cache
}

var (
	mu     sync.RWMutex
	groups = make(map[string]*Group)
)

func (g *Group) Get(key string) (ByteView, error) {
	if v, ok := g.mainCache.get(key); ok {
		log.Printf("[GeeCache] hit key -> %s\n", key)
		return v, nil
	}
	// 未命中,需要从数据源加载数据
	log.Println("load from datasource...")
	return g.load(key)
}

func (g *Group) load(key string) (ByteView, error) {
	return g.getLocally(key)
}

func (g *Group) getLocally(key string) (ByteView, error) {
	bytes, err := g.getter.Get(key)
	if err != nil {
		return ByteView{}, err
	}
	value := ByteView{b: cloneBytes(bytes)}
	if len(bytes) != 0 {
		g.populateCache(key, value)
	}
	return value, nil
}

func (g *Group) populateCache(key string, value ByteView) {
	g.mainCache.add(key, value)
}
func NewGroup(name string, cacheBytes int64, getter Getter) *Group {
	if getter == nil {
		panic("nil Getter!")
	}
	mu.Lock()
	defer mu.Unlock()
	g := &Group{
		name:      name,
		getter:    getter,
		mainCache: cache{cacheBytes: cacheBytes},
	}
	groups[name] = g
	return g
}

// GetGroup 获取分组
func GetGroup(name string) *Group {
	mu.RLock()
	defer mu.RUnlock()
	return groups[name]
}

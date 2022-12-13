package geecache

import (
	"geecache/singleflight"
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
	picker    PeerPicker          // 分布式缓存节点选择器
	loader    *singleflight.Group // using this to make sure each key is only fetch once at the same time
}

var (
	mu     sync.RWMutex
	groups = make(map[string]*Group)
)

func (g *Group) RegisterPeerPicker(picker PeerPicker) {
	if g.picker != nil {
		panic("already register peers picker")
	}
	g.picker = picker
}
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
	view, err := g.loader.Do(key, func() (interface{}, error) {
		if g.picker != nil {
			if peer, ok := g.picker.pickPeer(key); ok {
				v, err := g.getFromPeer(peer, key)
				if err == nil {
					log.Println("get data form {}", peer)
					return v, nil
				}
				log.Println("fail to get from peer", err)
			}
		}
		// 从本地加载
		return g.getLocally(key)
	})
	if err == nil {
		return view.(ByteView), nil
	}
	return ByteView{}, nil
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

func (g *Group) getFromPeer(peer PeerGetter, key string) (ByteView, error) {
	bytes, err := peer.Get(g.name, key)
	if err != nil {
		return ByteView{}, err
	}
	return ByteView{b: bytes}, nil
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
		loader:    &singleflight.Group{},
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

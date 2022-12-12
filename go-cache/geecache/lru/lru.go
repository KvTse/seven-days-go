package lru

import (
	"container/list"
	"fmt"
)

type Cache struct {
	maxBytes  int64 //缓存最大字节数
	nBytes    int64
	cache     map[string]*list.Element      // list.Element the element of double linklist
	ll        *list.List                    // 双向链表
	OnEvicted func(key string, value Value) // the callback function of invalid cache event
}

// Value the cache value
type Value interface {
	// Len calculate how many bytes it takes
	Len() int
}

type entry struct {
	key   string
	value Value
}

// New create a new cache
// map to cache, and list to maintain request order
func New(maxBytes int64, onEvicted func(string, Value)) *Cache {
	return &Cache{
		maxBytes:  maxBytes,
		ll:        list.New(),
		cache:     make(map[string]*list.Element),
		OnEvicted: onEvicted,
	}
}

// Get ele from cache
func (c *Cache) Get(key string) (value Value, ok bool) {
	// 如果cache中取到元素
	if ele, ok := c.cache[key]; ok {
		c.ll.MoveToFront(ele) //移动到队尾 双向链表的首尾是相对的
		kv := ele.Value.(*entry)
		return kv.value, true
	}
	return
}

// Add cache ele
// put ele to map and let ele to the tail
func (c *Cache) Add(key string, value Value) {
	// 元素存在 返回元素
	if ele, ok := c.cache[key]; ok {
		c.ll.MoveToFront(ele)
		kv := ele.Value.(*entry)
		c.nBytes += int64(value.Len()) - int64(kv.value.Len())
		kv.value = value
	} else { // 元素不存在
		// 添加新节点
		ele := c.ll.PushFront(&entry{key, value})
		// 映射字典
		c.cache[key] = ele
		c.nBytes += int64(len(key)) + int64(value.Len())
	}
	// 如果内存不足了,移除队尾元素
	for c.maxBytes != 0 && c.maxBytes < c.nBytes {
		c.RemoveOldest()
	}
}

// RemoveOldest remove ele if above the container
func (c *Cache) RemoveOldest() {
	fmt.Printf("remove the oldest...")
	ele := c.ll.Back()
	if ele != nil {
		c.ll.Remove(ele)
		kv := ele.Value.(*entry)
		delete(c.cache, kv.key)
		c.nBytes -= int64(len(kv.key)) + int64(kv.value.Len())
		if c.OnEvicted != nil {
			c.OnEvicted(kv.key, kv.value)
		}
	}
}

// Len the number of elements of list
func (c *Cache) Len() int {
	return c.ll.Len()
}

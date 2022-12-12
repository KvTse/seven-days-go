package lru

import (
	"fmt"
	"testing"
)

type String string

func (s String) Len() int {
	return len(s)
}

func TestCache_Add(t *testing.T) {
	lruCache := New(int64(16), func(s string, value Value) {
		fmt.Printf("remove ele key %s\n", s)
	})
	lruCache.Add("zhangsan", String("123"))
	lruCache.Add("zhangsan", String("456"))

	fmt.Printf("the cache ele nBytes is %d\n", lruCache.nBytes)
}
func TestCache_Get(t *testing.T) {
	lruCache := New(int64(16), nil)
	lruCache.Add("zhangsan", String("123"))
	value, ok := lruCache.Get("zhangsan")
	if ok {
		fmt.Printf("hit cache %s\n", value.(String))
	}
}

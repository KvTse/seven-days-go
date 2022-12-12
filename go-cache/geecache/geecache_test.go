package geecache

import (
	"fmt"
	"log"
	"testing"
	"time"
)

var mockDb = map[string]string{
	"Tom":  "630",
	"Jack": "589",
	"Sam":  "567",
}

func TestGet(t *testing.T) {
	cacheGroup := NewGroup("cache1", 2<<10, GetterFunc(
		func(key string) ([]byte, error) {
			log.Println("load data from datasource...slow", key)
			time.Sleep(time.Duration(2) * time.Second)
			if value, ok := mockDb[key]; ok {
				return []byte(value), nil
			}
			return nil, fmt.Errorf("%s not exist", key)
		},
	))
	cacheGroup.Get("Tom")
	cacheGroup.Get("Tom")

}

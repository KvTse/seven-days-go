package main

import (
	"fmt"
	"geecache"
	"log"
	"net/http"
)

var mockDb = map[string]string{
	"Tom":  "630",
	"Jack": "589",
	"Sam":  "567",
}

func main() {
	geecache.NewGroup("scores", 2<<10, geecache.GetterFunc(
		func(key string) ([]byte, error) {
			log.Println("[SlowDB] search key", key)
			if v, ok := mockDb[key]; ok {
				return []byte(v), nil
			}
			return nil, fmt.Errorf("%s not exist", key)
		}))

	addr := "localhost:8099"
	peers := geecache.NewHttpPool(addr)
	log.Println("geecache is running at", addr)
	log.Fatal(http.ListenAndServe(addr, peers))
}

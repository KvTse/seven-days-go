package main

import (
	"flag"
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

func getOrCreateGroup(groupName string) *geecache.Group {
	group := geecache.GetGroup(groupName)
	if group == nil {
		geecache.NewGroup(groupName, 2<<10, geecache.GetterFunc(
			func(key string) ([]byte, error) {
				log.Println("[SlowDB] search key", key)
				if v, ok := mockDb[key]; ok {
					return []byte(v), nil
				}
				return nil, fmt.Errorf("%s not exist", key)
			}))
	}
	return geecache.GetGroup(groupName)
}
func startCacheServer(addr string, addrs []string, group *geecache.Group) {
	peers := geecache.NewHttpPool(addr)
	peers.Set(addrs...)
	group.RegisterPeerPicker(peers)
	log.Println("cache server is running at", addr)
	log.Fatal(http.ListenAndServe(addr[7:], peers))
}
func startApiServer(apiAddr string, group *geecache.Group) {
	http.Handle("/api", http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			key := r.URL.Query().Get("key")
			view, err := group.Get(key)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "text/html")
			_, _ = w.Write(view.ByteSlice())
		},
	))
	log.Println("frontend server is running at", apiAddr)
	log.Fatal(http.ListenAndServe(apiAddr[7:], nil))
}

func main() {
	var port int
	var api bool
	flag.IntVar(&port, "port", 8001, "Geecache server port")
	flag.BoolVar(&api, "api", true, "Start a api server?")
	flag.Parse()

	apiAddr := "http://localhost:8099"
	addrMap := map[int]string{
		8001: "http://localhost:8001",
		8002: "http://localhost:8002",
		8003: "http://localhost:8003",
	}

	var addrs []string
	for _, v := range addrMap {
		addrs = append(addrs, v)
	}

	group := getOrCreateGroup("scores")
	if api {
		go startApiServer(apiAddr, group)
	}
	startCacheServer(addrMap[port], []string(addrs), group)
}

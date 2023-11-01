package main

import (
	"cache/internal/data"
	"cache/internal/model/cachegroup"
	"cache/internal/model/pool"
	"flag"
	"log"
	"net/http"
)

// var db1 = map[string]string{
// 	"Tom":  "630",
// 	"Jack": "589",
// 	"Sam":  "567",
// }

func createGroup() *cachegroup.Group {
	return cachegroup.NewGroup("scores",
		data.Get(), 2<<10)
}

func startCacheServer(addr string, addrs []string, cache *cachegroup.Group) {
	peers := pool.NewHTTPPool(addr)
	peers.Set(addrs...)
	cache.RegisterPeer(peers)
	log.Println("running at", addr)
	log.Fatal(http.ListenAndServe(addr[7:], peers))
}
func StartAPIServer(apiAddr string, cache *cachegroup.Group) {
	http.Handle("/api", http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			key := r.URL.Query().Get("key")
			view, err := cache.Get(key)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Header().Set("Context-Type", "application/octet-stream")
			w.Write(view.ByteSlice())
		}))
	log.Println("fontend server is running at", apiAddr)
	log.Fatal(http.ListenAndServe(apiAddr[7:], nil))
}

func main() {
	var port int
	var api bool
	flag.IntVar(&port, "port", 8001, "server port")
	flag.BoolVar(&api, "api", false, "Start a api server?")
	flag.Parse()
	apiAddr := "http://localhost:9999"
	addrMap := map[int]string{
		8001: "http://localhost:8001",
		8002: "http://localhost:8002",
	}
	var addrs []string
	for _, v := range addrMap {
		addrs = append(addrs, v)
	}
	group := createGroup()
	if api {
		go StartAPIServer(apiAddr, group)
	}
	startCacheServer(addrMap[port], []string(addrs), group)
}

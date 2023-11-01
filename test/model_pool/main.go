package main

import (
	"cache/internal/interface/getter"
	"cache/internal/model/cachegroup"
	"cache/internal/model/pool"
	"fmt"
	"log"
	"net/http"
)

var db = map[string]string{
	"Tom":  "630",
	"Jack": "589",
	"Sam":  "567",
}

func main() {
	cachegroup.NewGroup("scores", getter.GetterFunc(
		func(key string) ([]byte, error) {
			log.Println("[SlowDB] search key", key)
			if v, ok := db[key]; ok {
				return []byte(v), nil
			}
			return nil, fmt.Errorf("%s no exist", key)
		}), 2<<10)
	addr := "localhost:9999"
	peers := pool.NewHTTPPool(addr)
	log.Println("running at", addr)
	log.Fatal(http.ListenAndServe(addr, peers))

}

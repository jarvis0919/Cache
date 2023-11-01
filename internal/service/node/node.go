package node

import (
	"bytes"
	"cache/internal/data"
	"cache/internal/model/cachegroup"
	"cache/internal/model/pool"
	clog "cache/pkg/log"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
)

type Node struct {
}

func (n *Node) Run() {
	n.validateArgs()
	port := flag.Int("port", 0, "指定节点运行端口")
	center := flag.String("center", "", "指定中心节点运行端口")
	flag.Parse()
	if *port < 1024 || *port > 49151 {
		clog.Info("[Node] 端口范围限制为 (1024-49151)")
		os.Exit(1)
	}
	// u := "http://127.0.0.1:3999/registry"
	resp, err := http.Post(*center+"registry", "application/json", bytes.NewBuffer([]byte(strconv.Itoa(*port))))
	if err != nil {
		clog.Panic("[Node] ", err)
	}
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		clog.Panic("[Node] ", err)
	}
	startCacheServer(string(data), createGroup())

}

func (n *Node) printUsage() {
	fmt.Println("用法:")
	fmt.Println("  端口 -port - 指定当前节点运行端口")
	fmt.Println("  端口 -center - 指定中心节点运行端口 例如: http://127.0.0.1:3999/")
	// fmt.Println("  端口 -cache  - 指定缓存大小 例如: http://127.0.0.1:3999/")
}

func (n *Node) validateArgs() {
	if len(os.Args) < 3 {
		n.printUsage()
		os.Exit(1)
	}
}
func createGroup() *cachegroup.Group {
	return cachegroup.NewGroup("scores",
		data.Get(), 2<<10)
}
func startCacheServer(url string, cache *cachegroup.Group) {
	peers := pool.NewHTTPPool(url)
	peers.Set(url)
	cache.RegisterPeer(peers)
	clog.Info("[Node] 节点运行成功:", url)
	if err := http.ListenAndServe(url[7:], peers); err != nil {
		clog.Panic(err)
	}
}

package registry

import (
	"cache/internal/interface/getter"
	"cache/internal/model/cachegroup"
	"cache/internal/model/pool"
	clog "cache/pkg/log"
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type Registry struct {
}

func (r *Registry) Run() {
	r.validateArgs()
	port := flag.Int("port", 0, "指定注册中心运行端口")
	flag.Parse()
	// fmt.Print(*port)
	if *port < 1024 || *port > 49151 {
		// fmt.Println("portFlag is out of range (1024-49151).")
		clog.Info("[center] 端口范围限制为 (1024-49151)")
		os.Exit(1)
	}
	group := createGroup()
	startAPIServer(*port, group)
}

func (r *Registry) printUsage() {
	fmt.Println("用法:")
	fmt.Println("  端口 -port - 指定注册中心运行端口")
}

func (r *Registry) validateArgs() {
	if len(os.Args) < 2 {
		r.printUsage()
		os.Exit(1)
	}
}
func createGroup() *cachegroup.Group {
	return cachegroup.NewGroup("scores",
		getter.GetterFunc(func(key string) ([]byte, error) {
			// clog.Warn("[center] 集群中无节点")
			return nil, errors.New("集群中无节点")
		}), 2<<10)
}
func startAPIServer(port int, cache *cachegroup.Group) {
	http.Handle("/get", http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			key := r.URL.Query().Get("key")
			clog.Info("[center] 来自", r.RemoteAddr, "的请求: |key:", key, "|")
			view, err := cache.Get(key)
			if err != nil {
				clog.Error("[center] ", err)
				w.Write([]byte("未查询到数据"))
				return
			}

			w.Header().Set("Context-Type", "application/octet-stream")
			if len(view.ByteSlice()) == 0 {
				clog.Info("[center] 未查询到数据")
				w.Write([]byte("未查询到数据"))
			} else {
				clog.Info("[center] 查询完毕发送 |", key, ":", string(view.ByteSlice()), "|")
				w.Write(view.ByteSlice())
			}
		}))
	http.Handle("/registry", http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			bytes, err := io.ReadAll(r.Body)
			if err != nil {
				clog.Error("[center] ", err)
			}
			var builder strings.Builder
			// 追加字符串
			builder.WriteString("http://")
			index := strings.Index(r.RemoteAddr, ":")
			builder.WriteString(r.RemoteAddr[:index+1])
			// fmt.Println(r.RemoteAddr[:index], "   ", string(bytes))
			builder.WriteString(string(bytes))
			// 获取最终拼接的字符串
			url := builder.String()
			clog.Info("[center] Node节点注册成功：", url)
			if cache.Peers == nil {
				peers := pool.NewHTTPPool("center")
				cache.RegisterPeer(peers)
			}
			cache.Peers.Set(url)
			w.Header().Set("Context-Type", "application/octet-stream")
			w.Write([]byte(url))
		}))
	clog.Info("[center] 中心节点部署在", port)
	if err := http.ListenAndServe("192.168.1.101:"+strconv.Itoa(port), nil); err != nil {
		clog.Panic("[center] ", err)
	}

}

package cachegroup

import (
	pb "cache/api"
	"cache/internal/interface/getter"
	"cache/internal/interface/peer"
	"cache/internal/model/singleflight"
	clog "cache/pkg/log"
	"fmt"
	"sync"
)

// 缓存组，添加getter回调函数
type Group struct {
	name      string
	getter    getter.Getter
	mainCache SafeCache
	Peers     peer.PeerPicker
	loader    *singleflight.ReqBuff
}

// 定义全局读写锁
// 定义组映射
var (
	mu     sync.RWMutex
	groups = make(map[string]*Group)
)

func (g *Group) RegisterPeer(peers peer.PeerPicker) {
	if g.Peers != nil {
		clog.Panic("[Group] 不可以再次注册 RegisterPeer")
	}
	g.Peers = peers
}

// 创建组
func NewGroup(name string, getter getter.Getter, cacheBytes int64) *Group {
	if getter == nil {
		clog.Panic("[Group] 空 Getter")
	}
	mu.Lock()
	defer mu.Unlock()
	g := &Group{
		name:      name,
		getter:    getter,
		mainCache: SafeCache{cacheBytes: cacheBytes},
		loader:    &singleflight.ReqBuff{},
	}
	groups[name] = g
	return g
}

// 获取组
func GetGroup(name string) *Group {
	mu.Lock()
	defer mu.Unlock()
	// g :=
	return groups[name]
}

// 获取指定key的组内缓存，如果缓存组没有数据，触发load方法获取缓存
func (g *Group) Get(key string) (ByteView, error) {
	if key == "" {
		return ByteView{}, fmt.Errorf("key is required")
	}
	if v, ok := g.mainCache.get(key); ok {
		clog.Info("[Group] 查询 |key: ", key, " Value: ", v, "|")
		// log.Println("[Cache] hit")
		return v, nil
	}
	return g.load(key)
}

// 调用转发
func (g *Group) load(key string) (value ByteView, err error) {
	view, err := g.loader.Do(key, func() (interface{}, error) {
		if g.Peers != nil {
		breakHere:
			if peer, ok, peerkey := g.Peers.PickPeer(key); ok {
				if value, err = g.getFromPeer(peer, key); err == nil {
					return value, nil
				}
				g.Peers.RemovePeer(peerkey)
				clog.Warn("[Group] 缓存节点失效", peerkey)
				goto breakHere
			}

		}
		return g.getLocally(key)
	})
	if err == nil {
		return view.(ByteView), nil
	}
	return
}
func (g *Group) getFromPeer(peer peer.PeerGetter, key string) (ByteView, error) {
	req := &pb.Request{
		Group: g.name,
		Key:   key,
	}
	res := &pb.Response{}
	err := peer.Get(req, res)
	if err != nil {
		return ByteView{}, err
	}
	return ByteView{b: res.Value}, nil
}

// 执行组的回调函数（传入的从数据获取数据的方法），并调用populateCache方法
func (g *Group) getLocally(key string) (value ByteView, err error) {
	bytes, err := g.getter.Get(key)
	if err != nil {
		return ByteView{}, err
	}
	value = ByteView{b: cloneBytes(bytes)}
	if string(bytes) != "数据库中相关无数据" {
		g.populateCache(key, value)
	}
	return value, nil
}

// 跟新缓存中数据
func (g *Group) populateCache(key string, value ByteView) {
	g.mainCache.add(key, value)
}

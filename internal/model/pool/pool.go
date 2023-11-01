package pool

import (
	pb "cache/api"
	"cache/internal/interface/peer"
	"cache/internal/model/cachegroup"
	"cache/internal/model/consistenthash"
	clog "cache/pkg/log"
	"fmt"
	"net/http"
	"strings"
	"sync"

	"google.golang.org/protobuf/proto"
)

const (
	defaultBasePath = "/cache/"
	defaultReplicas = 50
)

// handle
type HTTPPool struct {
	self        string
	basePath    string
	mu          sync.Mutex
	peers       *consistenthash.Map
	httpGetters map[string]*httpGetter
}

var _ peer.PeerPicker = (*HTTPPool)(nil)

// new方法
func NewHTTPPool(self string) *HTTPPool {
	return &HTTPPool{
		self:     self,
		basePath: defaultBasePath,
	}
}

// 服务处理方法
func (h *HTTPPool) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !strings.HasPrefix(r.URL.Path, h.basePath) {

		clog.Panic("[Server] HTTPPool serving unexpected path: " + r.URL.Path)
	}
	h.Log("收到请求%s %s", r.Method, r.URL.Path)
	//约定访问路径为/<basepath>/<groupname>/<key>，判断是否符合标准
	parts := strings.SplitN(r.URL.Path[len(h.basePath):], "/", 2)
	if len(parts) != 2 {
		http.Error(w, "[Server] bad request", http.StatusBadRequest)
		return
	}
	//获取路径中传参，并查找组
	groupName := parts[0]
	key := parts[1]
	group := cachegroup.GetGroup(groupName)
	if group == nil {
		http.Error(w, "[Server] no such group"+groupName, http.StatusNotFound)
		return
	}
	//获取组key对应数据
	view, err := group.Get(key)
	// h.Log(string(view.ByteSlice()), err)
	if err != nil {
		// w.Write([]byte("数据失败"))
		http.Error(w, "[Server] no such group"+groupName, http.StatusInternalServerError)
	}
	body, err := proto.Marshal(&pb.Response{Value: view.ByteSlice()})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Write(body)
}

func (h *HTTPPool) Set(peers ...string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if h.peers == nil {
		h.peers = consistenthash.New(defaultReplicas, nil)
	}
	h.peers.Add(peers...)
	if h.httpGetters == nil {
		h.httpGetters = make(map[string]*httpGetter)
	}
	for _, peer := range peers {
		// h.NodeList = append(h.NodeList, peer)
		h.httpGetters[peer] = &httpGetter{BaseURL: peer + h.basePath}
	}
}

func (h *HTTPPool) PickPeer(key string) (peer.PeerGetter, bool, string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if peer := h.peers.Get(key); peer != "" && peer != h.self {
		// h.Log("获取节点 %s", peer)
		return h.httpGetters[peer], true, peer
	} else {
		// h.Log("Pick peer %s", peer)
		return nil, false, peer
	}
}

func (h *HTTPPool) RemovePeer(peer string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.peers.Remove(peer)
	delete(h.httpGetters, peer)
	// return list
}

// 日志打印
func (h *HTTPPool) Log(format string, v ...interface{}) {
	clog.Info("[Server] ", h.self, fmt.Sprintf(format, v...))
}

package gocache

import (
	"fmt"
	"gocache/consistenthash"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"
)

const defaultBasePath = "/_gocache/"

type HTTPPool struct {
	self        string
	basePath    string
	mu          sync.Mutex
	peers       *consistenthash.Map    // 一致性hash 算法的key 用来选取节点
	httpGetters map[string]*httpGetter // key: eg :"http://10.0.0.2:8008" 远程节点对应的httpGetter
}

func NewHTTPPool(self string) *HTTPPool {

	return &HTTPPool{
		self:     self,
		basePath: defaultBasePath,
	}

}

func (h *HTTPPool) Log(format string, v ...interface{}) {

	log.Printf("[Server %s] %s", h.self, fmt.Sprintf(format, v...))

}
func (h *HTTPPool) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	if !strings.HasPrefix(r.URL.Path, h.basePath) {
		panic("HTTPPool serving unexpected path :" + r.URL.Path)

	}

	paths := strings.SplitN(r.URL.Path[len(defaultBasePath):], "/", 2)
	if len(paths) != 2 {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	groupName := paths[0]
	key := paths[1]

	g := GetGroup(groupName)
	if g == nil {
		http.Error(w, "no such group :"+g.name, http.StatusNotFound)
		return
	}

	v, err := g.Get(key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/octet-stream")
	w.Write(v.ByteSlice())
}

// Set  更新节点的list列表
func (h *HTTPPool) Set(peers ...string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.peers = consistenthash.New(defaultReplicas, nil)
	h.peers.Add(peers...)
	h.httpGetters = make(map[string]*httpGetter)
	for _, peer := range peers {
		h.httpGetters[peer] = &httpGetter{baseURL: peer + h.basePath}
	}

}
func (h *HTTPPool) PickPeer(key string) (PeerGetter, bool) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if peer := h.peers.Get(key); peer != "" && peer != h.self {
		h.Log("Pick peer %s", peer)
		return h.httpGetters[peer], true
	}
	return nil, false
}

type httpGetter struct {
	// 表示要访问的远程节点地址
	baseURL string
}

func (h *httpGetter) Get(group, key string) ([]byte, error) {
	// 对URL进行转义
	u := fmt.Sprintf("%v%v/%v",
		h.baseURL,
		url.QueryEscape(group),
		url.QueryEscape(key))
	res, err := http.Get(u)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server returned: %v", res.Status)
	}

	bytes, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return nil, err
	}

	return bytes, nil
}

var _ PeerGetter = (*httpGetter)(nil)

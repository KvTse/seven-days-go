package geecache

import (
	"fmt"
	"geecache/consistenthash"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
)

const (
	defaultBasePath = "/_cache/"
	defaultReplicas = 50
)

type HttpPool struct {
	self        string
	basePath    string
	mu          sync.Mutex
	peers       *consistenthash.Map
	httpGetters map[string]*httpGetter
}

func (p *HttpPool) Set(peers ...string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.peers = consistenthash.New(defaultReplicas, nil)
	p.peers.Add(peers...)
	// 初始化一个长度为len(peers)的map
	p.httpGetters = make(map[string]*httpGetter, len(peers))
	for _, peer := range peers {
		p.httpGetters[peer] = &httpGetter{baseURL: peer + p.basePath}
	}

}

// pick -> pickerGetter -> Get data
func (p *HttpPool) pickPeer(key string) (PeerGetter, bool) {
	p.mu.Lock()
	defer p.mu.Unlock()
	peer := p.peers.Get(key)
	if peer != "" && peer != p.self {
		return p.httpGetters[peer], true
	}
	return nil, false
}

func NewHttpPool(self string) *HttpPool {
	return &HttpPool{
		self:     self,
		basePath: defaultBasePath,
	}
}

// ServeHTTP 从缓存中获取数据提供给客户端
func (p *HttpPool) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !strings.HasPrefix(r.URL.Path, p.basePath) {
		panic("httpPool serving unexpected path " + r.URL.Path)
	}
	// the url path pattern such as /basePath/group/key
	parts := strings.SplitN(r.URL.Path[len(p.basePath):], "/", 2)
	if len(parts) != 2 {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	groupName := parts[0]
	key := parts[1]

	group := GetGroup(groupName)
	if group == nil {
		http.Error(w, "No Such Group :"+groupName, http.StatusNotFound)
		return
	}
	value, err := group.Get(key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write(value.ByteSlice())
}

type httpGetter struct {
	baseURL string
}

func (h httpGetter) Get(group string, key string) ([]byte, error) {
	u := fmt.Sprintf("%v%v/%v",
		h.baseURL,
		url.QueryEscape(group),
		url.QueryEscape(key))
	// 发送http请求连接cache server获取数据
	resp, err := http.Get(u)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server returned: %v", resp.Status)
	}
	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response body :%v", err)
	}
	return bytes, nil
}

// 利用强制类型转换来确保PeerGetter接口被httpGetter实现
// 未实现在编译期报错
var _ PeerGetter = (*httpGetter)(nil)
var _ PeerPicker = (*HttpPool)(nil)

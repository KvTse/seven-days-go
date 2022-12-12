package geecache

import (
	"net/http"
	"strings"
)

const defaultBasePath = "/_cache/"

type HttpPool struct {
	self     string
	basePath string
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

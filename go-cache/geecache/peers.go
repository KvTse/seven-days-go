package geecache

type PeerPicker interface {
	// pickPeer 获取节点
	pickPeer(key string) (peer PeerGetter, ok bool)
}

// PeerGetter 从节点上获取缓存数据接口
type PeerGetter interface {
	Get(group string, key string) ([]byte, error)
}

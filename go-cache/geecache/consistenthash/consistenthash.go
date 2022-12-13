package consistenthash

import (
	"hash/crc32"
	"sort"
	"strconv"
)

// Hash 函数类型hash,供用户自定义hash实现方法
type Hash func(data []byte) uint32

// Map 一致性hash算法的结构体
type Map struct {
	hash     Hash           // hash函数
	replicas int            // 虚拟节点倍数
	keys     []int          // hash环
	hashMap  map[int]string // 虚拟节点和真实节点的映射表
}

func New(replicas int, fn Hash) *Map {
	m := &Map{
		replicas: replicas,
		hash:     fn,
		hashMap:  make(map[int]string),
	}
	// 默认使用crc32.ChecksumIEEE算法
	if m.hash == nil {
		m.hash = crc32.ChecksumIEEE
	}
	return m
}

// Add to add one or more truly node to the hash ring
// the build [m.replicas] virtual nodes
func (m *Map) Add(keys ...string) {
	for _, key := range keys {
		for i := 0; i < m.replicas; i++ {
			virtualNodeName := strconv.Itoa(i) + key
			hash := int(m.hash([]byte(virtualNodeName)))
			m.keys = append(m.keys, hash)
			m.hashMap[hash] = key
		}
	}
	sort.Ints(m.keys)
}

// Get choose the node return the true node
func (m *Map) Get(key string) string {
	if len(m.keys) == 0 {
		return ""
	}
	hash := int(m.hash([]byte(key)))
	// 二分法查找到比key.hash大的第一个虚拟节点的index值
	idx := sort.Search(len(m.keys), func(i int) bool {
		return m.keys[i] >= hash
	})
	// 用取余的方式来处理环 防止下标越界
	ringIdx := m.keys[idx%len(m.keys)]

	return m.hashMap[ringIdx]
}

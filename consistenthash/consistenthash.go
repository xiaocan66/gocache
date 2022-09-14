package consistenthash

// 分布式节点的一致性hash算法

import (
	"hash/crc32"
	"sort"
	"strconv"
)

type Hash func([]byte) uint32

type Map struct {
	hash     Hash           // hash函数
	replicas int            // 复制的个数
	keys     []int          //sorted
	hashMap  map[int]string //存储复制的节点
}

func New(replicas int, fn Hash) *Map {

	m := &Map{
		replicas: replicas,
		hash:     fn,
		hashMap:  make(map[int]string),
	}
	if m.hash == nil {
		m.hash = crc32.ChecksumIEEE

	}
	return m
}

func (m *Map) Add(keys ...string) {
	for _, k := range keys {
		for i := 0; i < m.replicas; i++ {
			hash := int(m.hash([]byte(strconv.Itoa(i) + k))) // 添加编号用于区分不同的虚拟节点
			m.keys = append(m.keys, hash)
			m.hashMap[hash] = k // 将虚拟节点的hash值与真实节点建立对应关系
		}
	}
	sort.Ints(m.keys)
}

// Get 通过key的hash值来或者取对应的节点
func (m *Map) Get(key string) string {
	if len(m.keys) == 0 {
		return ""
	}
	hash := int(m.hash([]byte(key)))
	// 二分查找 大于等于 hash 的节点
	idx := sort.Search(len(m.keys), func(i int) bool {
		return m.keys[i] >= hash
	})
	return m.hashMap[m.keys[idx%len(m.keys)]]

}

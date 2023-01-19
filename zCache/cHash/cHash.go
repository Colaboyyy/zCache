package cHash

import (
	"hash/crc32"
	"sort"
	"strconv"
)

type Hash func(data []byte) uint32

// Map 致性哈希算法的主数据结构
type Map struct {
	hash     Hash  // 哈希函数
	replicas int   // 虚拟节点倍数
	keys     []int //哈希环的keys，已排序
	hashMap  map[int]string
}

// New 构造函数 允许自定义虚拟节点倍数和 Hash 函数。
func New(replicas int, fn Hash) *Map {
	m := &Map{
		replicas: replicas,
		hash:     fn,
		hashMap:  make(map[int]string),
	}
	if m.hash == nil {
		m.hash = crc32.ChecksumIEEE //默认哈希算法为 crc32.ChecksumIEEE 算法
	}
	return m
}

// Add 添加真实节点/机器
func (m *Map) Add(keys ...string) {
	for _, key := range keys {
		for i := 0; i < m.replicas; i++ {
			hash := int(m.hash([]byte(strconv.Itoa(i) + key))) //计算虚拟节点的哈希值，通过添加编号的方式区分不同虚拟节点
			m.keys = append(m.keys, hash)                      // 使用 append(m.keys, hash) 添加到环上
			m.hashMap[hash] = key                              //在 hashMap 中增加虚拟节点和真实节点的映射关系。
		}
	}
	sort.Ints(m.keys) // 环上哈希值排序
}

// Get 获取离key最近的节点
func (m *Map) Get(key string) string {
	if len(m.keys) == 0 {
		return ""
	}
	// 计算 key 的哈希值
	hash := int(m.hash([]byte(key)))
	// 二分搜索一个合适的倍数
	idx := sort.Search(len(m.keys), func(i int) bool {
		return m.keys[i] >= hash
	})
	return m.hashMap[m.keys[idx%len(m.keys)]] // 环状结构，所以用取余数的方式来处理这种情况。
}

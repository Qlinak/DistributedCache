package geeCache

import (
	"hash/crc32"
	"sort"
	"strconv"
)

type Hash func(bArr []byte) uint32

type Map struct {
	hash           Hash
	replicas       int
	keys           []int
	virtualNodeMap map[int]string
}

func NewMap(replicas int, fn Hash) *Map {
	m := &Map{
		replicas:       replicas,
		hash:           fn,
		keys:           make([]int, 0),
		virtualNodeMap: make(map[int]string),
	}
	if m.hash == nil {
		m.hash = crc32.ChecksumIEEE
	}
	return m
}

// Add - add node to the hash ring
func (m *Map) Add(keys ...string) {
	for _, value := range keys {
		for i := 0; i < m.replicas; i++ { // how many virtual nodes a real node has
			virtualNode := strconv.Itoa(i) + value   // which virtual node
			hash := int(m.hash([]byte(virtualNode))) // hash of the node
			m.keys = append(m.keys, hash)            // add to ring
			m.virtualNodeMap[hash] = value           // store mapping: virtual node - real node
		}
	}
	sort.Ints(m.keys)
}

// Get - get a node by key
func (m *Map) Get(key string) string {
	if len(key) == 0 {
		return ""
	}

	hash := int(m.hash([]byte(key)))
	// binary search, find the smallest idx that is greater than hash
	idx := sort.Search(len(m.keys), func(i int) bool {
		return m.keys[i] >= hash
	})

	return m.virtualNodeMap[m.keys[idx%len(m.keys)]]
}

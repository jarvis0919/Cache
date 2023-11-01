package model_consistenthash

import (
	"hash/crc32"
	"sort"
	"strconv"
	"sync"
	"testing"
)

type Hash func(data []byte) uint32

type Map struct {
	hash     Hash
	replicas int
	keys     []int
	hashMap  map[int]string
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
	for _, key := range keys {
		for i := 0; i < m.replicas; i++ {
			hash := int(m.hash([]byte(strconv.Itoa(i) + key)))
			m.keys = append(m.keys, hash)
			m.hashMap[hash] = key
		}
	}
	sort.Ints(m.keys)
}

func (m *Map) Get(key string) string {
	if len(m.keys) == 0 {
		return ""
	}
	hash := int(m.hash([]byte(key)))
	idx := sort.Search(len(m.keys), func(i int) bool {
		return m.keys[i] >= hash
	})

	return m.hashMap[m.keys[idx%len(m.keys)]]
}
func (m *Map) Remove(key string) {
	for i := 0; i < m.replicas; i++ {
		hash := int(m.hash([]byte(strconv.Itoa(i) + key)))
		idx := sort.SearchInts(m.keys, hash)
		if m.keys[idx] != hash {
			return
		}
		m.keys = append(m.keys[:idx], m.keys[idx+1:]...)
		delete(m.hashMap, hash)
	}
	// sort.Ints(m.keys)
}

func TestGo1(t *testing.T) {
	hash := New(3, func(key []byte) uint32 {
		i, _ := strconv.Atoi(string(key))
		return uint32(i)
	})
	// hash.Add("1", "4", "2")
	// testCases := map[string]string{
	// 	"1":  "2",
	// 	"2":  "2",
	// 	"3":  "4",
	// 	"4":  "2",
	// 	"5":  "2",
	// 	"6":  "2",
	// 	"7":  "4",
	// 	"8":  "2",
	// 	"9":  "2",
	// 	"10": "2",
	// 	"11": "4",
	// 	"12": "2",
	// 	"13": "2",
	// 	"14": "2",
	// 	"15": "4",
	// 	"16": "2",
	// }
	wg := sync.WaitGroup{}
	for k := 0; k < 10000; k++ {
		hash.Add(strconv.Itoa(k))
	}
	wg.Wait()
}

package iavl

import (
	"sync"
	"sync/atomic"

	xxhash "github.com/cespare/xxhash/v2"
)

type FixedMap struct {
	maxEntries uint64
	maps       [][]byte
	locks      []sync.Mutex

	// stats
	entries int64
	bytes   int64
	sets    int64
	gets    int64
	hits    int64
	misses  int64
}

func NewFixedMap(maxEntries int) *FixedMap {
	return &FixedMap{
		maxEntries: uint64(maxEntries),
		maps:       make([][]byte, maxEntries),
		locks:      make([]sync.Mutex, 1009),
	}
}

func (m *FixedMap) Set(key, value []byte) {
	kl, vl := len(key), len(value)
	if kl >= (1<<16) || vl >= (1<<16) {
		return
	}
	buf := make([]byte, 4+kl+vl)
	buf[0], buf[1] = byte(kl>>8), byte(kl)
	buf[2], buf[3] = byte(vl>>8), byte(vl)
	copy(buf[4:], key)
	copy(buf[4+kl:], value)

	h := xxhash.Sum64(key)
	lh := h % uint64(len(m.locks))
	h %= m.maxEntries
	n, sz := int64(1), int64(kl+vl)

	m.locks[lh].Lock()
	ob := m.maps[h]
	if ob != nil {
		okl := (uint(ob[0]) << 8) + uint(ob[1])
		ovl := (uint(ob[2]) << 8) + uint(ob[3])
		n = 0
		sz -= int64(okl + ovl)
	}
	m.maps[h] = buf
	m.locks[lh].Unlock()

	atomic.AddInt64(&m.entries, n)
	atomic.AddInt64(&m.bytes, sz)
}

func (m *FixedMap) Get(dst, key []byte) []byte {
	h := xxhash.Sum64(key)
	lh := h % uint64(len(m.locks))
	h %= m.maxEntries

	m.locks[lh].Lock()
	buf := m.maps[h]
	if buf == nil {
		m.locks[lh].Unlock()
		atomic.AddInt64(&m.misses, 1)
		return nil
	}
	kl := (uint(buf[0]) << 8) + uint(buf[1])
	vl := (uint(buf[2]) << 8) + uint(buf[3])
	if string(key) != string(buf[4:4+kl]) {
		m.locks[lh].Unlock()
		atomic.AddInt64(&m.misses, 1)
		return nil
	}
	rst := make([]byte, vl)
	copy(rst, buf[4+kl:])
	m.locks[lh].Unlock()

	atomic.AddInt64(&m.hits, 1)
	return rst
}

func (m *FixedMap) Has(key []byte) bool {
	h := xxhash.Sum64(key) % m.maxEntries
	return m.maps[h] != nil
}

func (m *FixedMap) Del(key []byte) {
	h := xxhash.Sum64(key)
	lh := h % uint64(len(m.locks))
	h %= m.maxEntries

	var n, sz int64
	m.locks[lh].Lock()
	ob := m.maps[h]
	if ob != nil {
		okl := (uint(ob[0]) << 8) + uint(ob[1])
		ovl := (uint(ob[2]) << 8) + uint(ob[3])
		n = -1
		sz -= int64(okl + ovl)
	}
	m.maps[h] = nil
	m.locks[lh].Unlock()

	atomic.AddInt64(&m.entries, n)
	atomic.AddInt64(&m.bytes, sz)
}

func (m *FixedMap) Stats() (hits, misses, entries, bytes uint64) {
	hits, misses, entries, bytes = uint64(m.hits), uint64(m.misses), uint64(m.entries), uint64(m.bytes)
	return
}

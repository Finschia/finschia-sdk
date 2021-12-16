package iavl

import (
	"github.com/VictoriaMetrics/fastcache"
	"github.com/coocood/freecache"
	"github.com/dgraph-io/ristretto"
)

type fastCache struct {
	*fastcache.Cache
}

func newFastCache(cacheSize int) *fastCache {
	return &fastCache{
		Cache: fastcache.New(cacheSize),
	}
}

func (c *fastCache) Stats() (hits, misses, entries, bytes uint64) {
	stats := fastcache.Stats{}
	c.UpdateStats(&stats)
	return stats.GetCalls - stats.Misses, stats.Misses, stats.EntriesCount, stats.BytesSize
}

type ristrettoCache struct {
	*ristretto.Cache
}

func newRistrettoCache(cacheSize int) *ristrettoCache {
	cache, err := ristretto.NewCache(&ristretto.Config{
		NumCounters: int64(cacheSize) * 10,
		MaxCost:     int64(cacheSize),
		BufferItems: 64,
		Metrics:     true,
	})
	if err != nil {
		panic(err)
	}
	return &ristrettoCache{
		Cache: cache,
	}
}

func (c *ristrettoCache) Set(key, value []byte) {
	c.Cache.Set(key, value, int64(len(value)))
}

func (c *ristrettoCache) Get(dst, key []byte) []byte {
	v, ok := c.Cache.Get(key)
	if ok {
		return v.([]byte)
	}
	return nil
}

func (c *ristrettoCache) Del(key []byte) {
	c.Cache.Del(key)
}

func (c *ristrettoCache) Has(key []byte) bool {
	_, ok := c.Cache.Get(key)
	return ok
}

func (c *ristrettoCache) Stats() (hits, misses, entries, bytes uint64) {
	m := c.Cache.Metrics
	hits = m.Hits()
	misses = m.Misses()
	entries = m.KeysAdded() - m.KeysEvicted()
	bytes = m.CostAdded() - m.CostEvicted()
	return
}

type freeCache struct {
	*freecache.Cache
}

func newFreeCache(cacheSize int) *freeCache {
	return &freeCache{
		Cache: freecache.NewCache(cacheSize),
	}
}

func (c *freeCache) Set(key, value []byte) {
	c.Cache.Set(key, value, 0)
}

func (c *freeCache) Get(dst, key []byte) []byte {
	v, err := c.Cache.Get(key)
	if err != nil {
		return nil
	}
	return v
}

func (c *freeCache) Has(key []byte) bool {
	v, err := c.Cache.Get(key)
	return err == nil && v != nil
}

func (c *freeCache) Del(key []byte) {
	c.Cache.Del(key)
}

func (c *freeCache) Stats() (hits, misses, entries, bytes uint64) {
	hits, misses, entries = uint64(c.HitCount()), uint64(c.MissCount()), uint64(c.EntryCount())
	return
}

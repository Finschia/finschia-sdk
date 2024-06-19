package cache

import (
	"sync"
	"time"

	"github.com/VictoriaMetrics/fastcache"

	"github.com/Finschia/finschia-sdk/store/cachekv"
	"github.com/Finschia/finschia-sdk/store/types"
)

const (
	DefaultCommitKVStoreCacheSize = 1024 * 1024 * 100 // 100 MB
)

var (
	_ types.CommitKVStore             = (*CommitKVStoreCache)(nil)
	_ types.MultiStorePersistentCache = (*CommitKVStoreCacheManager)(nil)
)

type (
	// CommitKVStoreCache implements an inter-block (persistent) cache that wraps a
	// CommitKVStore. Reads first hit the internal ARC (Adaptive Replacement Cache).
	// During a cache miss, the read is delegated to the underlying CommitKVStore
	// and cached. Deletes and writes always happen to both the cache and the
	// CommitKVStore in a write-through manner. Caching performed in the
	// CommitKVStore and below is completely irrelevant to this layer.
	CommitKVStoreCache struct {
		types.CommitKVStore
		cache   *fastcache.Cache
		prefix  []byte
		metrics *Metrics
	}

	// CommitKVStoreCacheManager maintains a mapping from a StoreKey to a
	// CommitKVStoreCache. Each CommitKVStore, per StoreKey, is meant to be used
	// in an inter-block (persistent) manner and typically provided by a
	// CommitMultiStore.
	CommitKVStoreCacheManager struct {
		mutex   sync.Mutex
		cache   *fastcache.Cache
		caches  map[string]types.CommitKVStore
		metrics *Metrics

		// All cache stores use the unique prefix that has one byte length
		// Contract: The number of all cache stores cannot exceed 127(max byte)
		prefixMap   map[string][]byte
		prefixOrder byte
	}
)

func NewCommitKVStoreCache(store types.CommitKVStore, prefix []byte, cache *fastcache.Cache,
	metrics *Metrics,
) *CommitKVStoreCache {
	return &CommitKVStoreCache{
		CommitKVStore: store,
		prefix:        prefix,
		cache:         cache,
		metrics:       metrics,
	}
}

func NewCommitKVStoreCacheManager(cacheSize int, provider MetricsProvider) *CommitKVStoreCacheManager {
	if cacheSize <= 0 {
		// This function was called because it intended to use the inter block cache, creating a cache of minimal size.
		cacheSize = DefaultCommitKVStoreCacheSize
	}
	cm := &CommitKVStoreCacheManager{
		cache:       fastcache.New(cacheSize),
		caches:      make(map[string]types.CommitKVStore),
		metrics:     provider(),
		prefixMap:   make(map[string][]byte),
		prefixOrder: 0,
	}
	startCacheMetricUpdator(cm.cache, cm.metrics)
	return cm
}

func startCacheMetricUpdator(cache *fastcache.Cache, metrics *Metrics) {
	// Execution time of `fastcache.UpdateStats()` can increase linearly as cache entries grows
	// So we update the metrics with a separate go route.
	go func() {
		for {
			stats := fastcache.Stats{}
			cache.UpdateStats(&stats)
			metrics.InterBlockCacheEntries.Set(float64(stats.EntriesCount))
			metrics.InterBlockCacheBytes.Set(float64(stats.BytesSize))
			time.Sleep(1 * time.Minute)
		}
	}()
}

// GetStoreCache returns a Cache from the CommitStoreCacheManager for a given
// StoreKey. If no Cache exists for the StoreKey, then one is created and set.
// The returned Cache is meant to be used in a persistent manner.
func (cmgr *CommitKVStoreCacheManager) GetStoreCache(key types.StoreKey, store types.CommitKVStore) types.CommitKVStore {
	if cmgr.caches[key.Name()] == nil {
		// After concurrent checkTx, delieverTx becomes to be possible, this should be protected by a mutex
		cmgr.mutex.Lock()
		if cmgr.caches[key.Name()] == nil { // recheck after acquiring lock
			cmgr.prefixMap[key.Name()] = []byte{cmgr.prefixOrder}
			cmgr.prefixOrder++
			if cmgr.prefixOrder <= 0 {
				panic("The number of cache stores exceed the maximum(127)")
			}
			cmgr.caches[key.Name()] = NewCommitKVStoreCache(store, cmgr.prefixMap[key.Name()], cmgr.cache, cmgr.metrics)
		}
		cmgr.mutex.Unlock()
	}

	return cmgr.caches[key.Name()]
}

// Unwrap returns the underlying CommitKVStore for a given StoreKey.
func (cmgr *CommitKVStoreCacheManager) Unwrap(key types.StoreKey) types.CommitKVStore {
	if ckv, ok := cmgr.caches[key.Name()]; ok {
		return ckv.(*CommitKVStoreCache).CommitKVStore
	}

	return nil
}

// Reset resets in the internal caches.
func (cmgr *CommitKVStoreCacheManager) Reset() {
	// Clear the map.
	// Please note that we are purposefully using the map clearing idiom.
	// See https://github.com/cosmos/cosmos-sdk/issues/6681.
	for key := range cmgr.caches {
		delete(cmgr.caches, key)
	}
}

// CacheWrap implements the CacheWrapper interface
func (ckv *CommitKVStoreCache) CacheWrap() types.CacheWrap {
	return cachekv.NewStore(ckv)
}

// Get retrieves a value by key. It will first look in the write-through cache.
// If the value doesn't exist in the write-through cache, the query is delegated
// to the underlying CommitKVStore.
func (ckv *CommitKVStoreCache) Get(key []byte) []byte {
	types.AssertValidKey(key)
	prefixedKey := append(ckv.prefix, key...)

	valueI := ckv.cache.Get(nil, prefixedKey)
	if valueI != nil {
		// cache hit
		ckv.metrics.InterBlockCacheHits.Add(1)
		return valueI
	}

	// cache miss; write to cache
	ckv.metrics.InterBlockCacheMisses.Add(1)
	value := ckv.CommitKVStore.Get(key)
	ckv.cache.Set(prefixedKey, value)
	return value
}

// Set inserts a key/value pair into both the write-through cache and the
// underlying CommitKVStore.
func (ckv *CommitKVStoreCache) Set(key, value []byte) {
	types.AssertValidKey(key)
	types.AssertValidValue(value)

	prefixedKey := append(ckv.prefix, key...)
	ckv.cache.Set(prefixedKey, value)
	ckv.CommitKVStore.Set(key, value)
}

// Delete removes a key/value pair from both the write-through cache and the
// underlying CommitKVStore.
func (ckv *CommitKVStoreCache) Delete(key []byte) {
	prefixedKey := append(ckv.prefix, key...)
	ckv.cache.Del(prefixedKey)
	ckv.CommitKVStore.Delete(key)
}

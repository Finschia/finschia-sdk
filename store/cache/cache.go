package cache

import (
	"github.com/VictoriaMetrics/fastcache"
	"github.com/line/lbm-sdk/v2/store/cachekv"
	"github.com/line/lbm-sdk/v2/store/types"
)

var (
	_ types.CommitKVStore             = (*CommitKVStoreCache)(nil)
	_ types.MultiStorePersistentCache = (*CommitKVStoreCacheManager)(nil)

	// DefaultCommitKVStoreCacheSize defines the persistent ARC cache size for a
	// CommitKVStoreCache.
	DefaultCommitKVStoreCacheSize uint = 1024*1024*2000 // 2000 MB
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
		cache *fastcache.Cache
		metrics *Metrics
	}

	// CommitKVStoreCacheManager maintains a mapping from a StoreKey to a
	// CommitKVStoreCache. Each CommitKVStore, per StoreKey, is meant to be used
	// in an inter-block (persistent) manner and typically provided by a
	// CommitMultiStore.
	CommitKVStoreCacheManager struct {
		cacheSize       uint
		caches          map[string]types.CommitKVStore
		metricsProvider func(storeName string) *Metrics
	}
)

func NewCommitKVStoreCache(store types.CommitKVStore, size uint, metrics *Metrics) *CommitKVStoreCache {
	cache := fastcache.New(int(size))

	return &CommitKVStoreCache{
		CommitKVStore: store,
		cache:         cache,
		metrics:       metrics,
	}
}

func NewCommitKVStoreCacheManager(size uint, metricsProvider MetricsProvider) *CommitKVStoreCacheManager {
	return &CommitKVStoreCacheManager{
		cacheSize:       size,
		caches:          make(map[string]types.CommitKVStore),
		metricsProvider: metricsProvider,
	}
}

// GetStoreCache returns a Cache from the CommitStoreCacheManager for a given
// StoreKey. If no Cache exists for the StoreKey, then one is created and set.
// The returned Cache is meant to be used in a persistent manner.
func (cmgr *CommitKVStoreCacheManager) GetStoreCache(key types.StoreKey, store types.CommitKVStore) types.CommitKVStore {
	if cmgr.caches[key.Name()] == nil {
		cmgr.caches[key.Name()] = NewCommitKVStoreCache(store, cmgr.cacheSize, cmgr.metricsProvider(key.Name()))
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

	valueI := ckv.cache.Get(nil, key)
	if valueI != nil {
		// cache hit
		ckv.metrics.InterBlockCacheHits.Add(1)
		return valueI
	}

	// cache miss; write to cache
	ckv.metrics.InterBlockCacheMisses.Add(1)
	value := ckv.CommitKVStore.Get(key)
	ckv.cache.Set(key, value)
	stats := fastcache.Stats{}
	ckv.cache.UpdateStats(&stats)
	ckv.metrics.InterBlockCacheEntries.Set(float64(stats.EntriesCount))
	ckv.metrics.InterBlockCacheBytes.Set(float64(stats.BytesSize))
	return value
}

// Set inserts a key/value pair into both the write-through cache and the
// underlying CommitKVStore.
func (ckv *CommitKVStoreCache) Set(key, value []byte) {
	types.AssertValidKey(key)
	types.AssertValidValue(value)

	ckv.cache.Set(key, value)
	ckv.CommitKVStore.Set(key, value)
}

// Delete removes a key/value pair from both the write-through cache and the
// underlying CommitKVStore.
func (ckv *CommitKVStoreCache) Delete(key []byte) {
	ckv.cache.Del(key)
	ckv.CommitKVStore.Delete(key)
}

package cache

import (
	"sync"

	"github.com/VictoriaMetrics/fastcache"
	"github.com/line/lbm-sdk/v2/store/cachekv"
	"github.com/line/lbm-sdk/v2/store/types"
	"github.com/line/lbm-sdk/v2/telemetry"
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
		cache  *fastcache.Cache
		prefix []byte
	}

	// CommitKVStoreCacheManager maintains a mapping from a StoreKey to a
	// CommitKVStoreCache. Each CommitKVStore, per StoreKey, is meant to be used
	// in an inter-block (persistent) manner and typically provided by a
	// CommitMultiStore.
	CommitKVStoreCacheManager struct {
		mutex  sync.Mutex
		cache  *fastcache.Cache
		caches map[string]types.CommitKVStore

		// All cache stores use the unique prefix that has one byte length
		// Contract: The number of all cache stores cannot exceed 127(max byte)
		prefixMap   map[string][]byte
		prefixOrder byte
	}
)

func NewCommitKVStoreCache(store types.CommitKVStore, prefix []byte, cache *fastcache.Cache) *CommitKVStoreCache {
	return &CommitKVStoreCache{
		CommitKVStore: store,
		prefix:        prefix,
		cache:         cache,
	}
}

func NewCommitKVStoreCacheManager(cacheSize int) *CommitKVStoreCacheManager {
	if cacheSize <= 0 {
		// This function was called because it intended to use the inter block cache, creating a cache of minimal size.
		cacheSize = DefaultCommitKVStoreCacheSize
	}
	return &CommitKVStoreCacheManager{
		cache:       fastcache.New(cacheSize),
		caches:      make(map[string]types.CommitKVStore),
		prefixMap:   make(map[string][]byte),
		prefixOrder: 0,
	}
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
			cmgr.caches[key.Name()] = NewCommitKVStoreCache(store, cmgr.prefixMap[key.Name()], cmgr.cache)
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
		telemetry.IncrCounter(1, "store", "inter-block-cache", "hits")
		return valueI
	}

	// cache miss; write to cache
	telemetry.IncrCounter(1, "store", "inter-block-cache", "misses")
	value := ckv.CommitKVStore.Get(key)
	ckv.cache.Set(prefixedKey, value)
	stats := fastcache.Stats{}
	ckv.cache.UpdateStats(&stats)
	telemetry.SetGauge(float32(stats.EntriesCount), "store", "inter-block-cache", "entries")
	telemetry.SetGauge(float32(stats.BytesSize), "store", "inter-block-cache", "bytes")
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

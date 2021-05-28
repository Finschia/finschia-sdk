package cache

import (
	"fmt"
	"reflect"
	"sync"
	"time"

	"github.com/VictoriaMetrics/fastcache"
	"github.com/dgraph-io/ristretto"
	"github.com/golang/protobuf/proto"
	"github.com/line/lbm-sdk/v2/codec"
	"github.com/line/lbm-sdk/v2/store/cachekv"
	"github.com/line/lbm-sdk/v2/store/types"
)

const (
	DefaultCommitKVStoreCacheSize = 1024 * 1024 * 100 // 100 MB
)

var (
	_ types.CommitKVObjectStore             = (*CommitKVStoreCache)(nil)
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
		cache   *ristretto.Cache
		prefix  []byte
		metrics *Metrics
	}

	// CommitKVStoreCacheManager maintains a mapping from a StoreKey to a
	// CommitKVStoreCache. Each CommitKVStore, per StoreKey, is meant to be used
	// in an inter-block (persistent) manner and typically provided by a
	// CommitMultiStore.
	CommitKVStoreCacheManager struct {
		mutex   sync.Mutex
		cache   *ristretto.Cache
		caches  map[string]types.CommitKVObjectStore
		metrics *Metrics

		// All cache stores use the unique prefix that has one byte length
		// Contract: The number of all cache stores cannot exceed 127(max byte)
		prefixMap   map[string][]byte
		prefixOrder byte
	}
)

func NewCommitKVStoreCache(store types.CommitKVStore, prefix []byte, cache *ristretto.Cache,
	metrics *Metrics) *CommitKVStoreCache {
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
	cache, err := ristretto.NewCache(&ristretto.Config{
		NumCounters: 1e7,               // number of keys to track frequency of (10M).
		MaxCost:     int64(cacheSize),
		BufferItems: 64,                // number of keys per Get buffer.
	})
	if err != nil {
		panic(fmt.Sprintf("Cannot create a ristretto cache: %s\n", err.Error()))
	}
	cm := &CommitKVStoreCacheManager{
		cache:       cache,
		caches:      make(map[string]types.CommitKVObjectStore),
		metrics:     provider(),
		prefixMap:   make(map[string][]byte),
		prefixOrder: 0,
	}
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
func (cmgr *CommitKVStoreCacheManager) GetStoreCache(key types.StoreKey,
	store types.CommitKVStore) types.CommitKVObjectStore {
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
func (ckv *CommitKVStoreCache) Get(key []byte, cdc codec.BinaryMarshaler, ptr interface{}) interface{} {
	types.AssertValidKey(key)
	prefixedKey := append(ckv.prefix, key...)

	val, exist := ckv.cache.Get(prefixedKey)
	if exist {
		ckv.metrics.InterBlockCacheHits.Add(1)
		return val
	}

	// cache miss; write to cache
	ckv.metrics.InterBlockCacheMisses.Add(1)
	value := ckv.CommitKVStore.Get(key)
	if err := cdc.UnmarshalInterface(value, ptr); err != nil {
		panic(fmt.Sprintf("Unable to unmarshal: %s\n", err.Error()))
	}
	ckv.cache.Set(prefixedKey, ptr, int64(reflect.TypeOf(ptr).Size()))
	return ptr
}

// Set inserts a key/value pair into both the write-through cache and the
// underlying CommitKVStore.
func (ckv *CommitKVStoreCache) Set(key []byte, cdc codec.BinaryMarshaler, obj proto.Message) {
	types.AssertValidKey(key)

	value, err := cdc.MarshalInterface(obj)
	if err != nil {
		panic(fmt.Sprintf("Unable to marshal: %s\n", err.Error()))
	}
	types.AssertValidValue(value)
	ckv.CommitKVStore.Set(key, value)

	prefixedKey := append(ckv.prefix, key...)
	ckv.cache.Set(prefixedKey, obj, int64(reflect.TypeOf(obj).Size()))
}

func (ckv *CommitKVStoreCache) SetObj(key []byte, cdc codec.BinaryMarshaler, obj proto.Message) {
	panic("This must not be called")
}

// Delete removes a key/value pair from both the write-through cache and the
// underlying CommitKVStore.
func (ckv *CommitKVStoreCache) Delete(key []byte) {
	prefixedKey := append(ckv.prefix, key...)
	ckv.cache.Del(prefixedKey)
	ckv.CommitKVStore.Delete(key)
}

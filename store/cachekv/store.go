package cachekv

import (
	"bytes"
	"io"
	"sort"
	"sync"
	"time"

	dbm "github.com/tendermint/tm-db"

	"github.com/line/lbm-sdk/internal/conv"
	"github.com/line/lbm-sdk/store/listenkv"
	"github.com/line/lbm-sdk/store/tracekv"
	"github.com/line/lbm-sdk/store/types"
	"github.com/line/lbm-sdk/telemetry"
	"github.com/line/lbm-sdk/types/kv"
)

// If value is nil but deleted is false, it means the parent doesn't have the
// key.  (No need to delete upon Write())
type cValue struct {
	value []byte
	dirty bool
}

// Store wraps an in-memory cache around an underlying types.KVStore.
// Set, Delete and Write for the same key must be called sequentially.
type Store struct {
	mtx           sync.RWMutex
	cache         sync.Map
	deleted       sync.Map
	unsortedCache sync.Map
	sortedCache   *dbm.MemDB // always ascending sorted
	parent        types.KVStore
}

var _ types.CacheKVStore = (*Store)(nil)

// NewStore creates a new Store object
func NewStore(parent types.KVStore) *Store {
	return &Store{
		cache:         sync.Map{},
		deleted:       sync.Map{},
		unsortedCache: sync.Map{},
		sortedCache:   dbm.NewMemDB(),
		parent:        parent,
	}
}

// GetStoreType implements Store.
func (store *Store) GetStoreType() types.StoreType {
	return store.parent.GetStoreType()
}

// Get implements types.KVStore.
func (store *Store) Get(key []byte) (value []byte) {
	types.AssertValidKey(key)
	store.mtx.RLock()
	defer store.mtx.RUnlock()
	cacheValue, ok := store.cache.Load(string(key))
	if ok {
		return cacheValue.(*cValue).value
	}

	value = store.parent.Get(key)
	store.setCacheValue(key, value, false, false)
	return value
}

// Set implements types.KVStore.
func (store *Store) Set(key []byte, value []byte) {
	types.AssertValidKey(key)
	types.AssertValidValue(value)

	store.mtx.Lock()
	defer store.mtx.Unlock()
	store.setCacheValue(key, value, false, true)
}

// Has implements types.KVStore.
func (store *Store) Has(key []byte) bool {
	value := store.Get(key)
	return value != nil
}

// Delete implements types.KVStore.
func (store *Store) Delete(key []byte) {
	defer telemetry.MeasureSince(time.Now(), "store", "cachekv", "delete")

	types.AssertValidKey(key)
	store.mtx.Lock()
	defer store.mtx.Unlock()
	store.setCacheValue(key, nil, true, true)
}

// Implements Cachetypes.KVStore.
func (store *Store) Write() {
	store.mtx.Lock()
	defer store.mtx.Unlock()
	defer telemetry.MeasureSince(time.Now(), "store", "cachekv", "write")

	// We need a copy of all of the keys.
	// Not the best, but probably not a bottleneck depending.
	keys := make([]string, 0)
	store.cache.Range(func(key, value interface{}) bool {
		if value.(*cValue).dirty {
			keys = append(keys, key.(string))
		}
		return true
	})

	sort.Strings(keys)

	// TODO: Consider allowing usage of Batch, which would allow the write to
	// at least happen atomically.
	for _, key := range keys {
		if store.isDeleted(key) {
			// We use []byte(key) instead of conv.UnsafeStrToBytes because we cannot
			// be sure if the underlying store might do a save with the byteslice or
			// not. Once we get confirmation that .Delete is guaranteed not to
			// save the byteslice, then we can assume only a read-only copy is sufficient.
			store.parent.Delete([]byte(key))
			continue
		}

		v, ok := store.cache.Load(key)
		cacheValue := v.(*cValue)
		if ok && cacheValue != nil {
			// It already exists in the parent, hence delete it.
			store.parent.Set([]byte(key), cacheValue.value)
		}
	}

	// Clear the cache
	store.cache = sync.Map{}
	store.deleted = sync.Map{}
	store.unsortedCache = sync.Map{}
	store.sortedCache = dbm.NewMemDB()
}

// CacheWrap implements CacheWrapper.
func (store *Store) CacheWrap() types.CacheWrap {
	return NewStore(store)
}

// CacheWrapWithTrace implements the CacheWrapper interface.
func (store *Store) CacheWrapWithTrace(w io.Writer, tc types.TraceContext) types.CacheWrap {
	return NewStore(tracekv.NewStore(store, w, tc))
}

// CacheWrapWithListeners implements the CacheWrapper interface.
func (store *Store) CacheWrapWithListeners(storeKey types.StoreKey, listeners []types.WriteListener) types.CacheWrap {
	return NewStore(listenkv.NewStore(store, storeKey, listeners))
}

//----------------------------------------
// Iteration

// Iterator implements types.KVStore.
func (store *Store) Iterator(start, end []byte) types.Iterator {
	return store.iterator(start, end, true)
}

// ReverseIterator implements types.KVStore.
func (store *Store) ReverseIterator(start, end []byte) types.Iterator {
	return store.iterator(start, end, false)
}

func (store *Store) iterator(start, end []byte, ascending bool) types.Iterator {
	store.mtx.Lock()
	defer store.mtx.Unlock()

	var parent, cache types.Iterator

	if ascending {
		parent = store.parent.Iterator(start, end)
	} else {
		parent = store.parent.ReverseIterator(start, end)
	}

	store.dirtyItems(start, end)
	cache = newMemIterator(start, end, store.sortedCache, &store.deleted, ascending)

	return newCacheMergeIterator(parent, cache, ascending)
}

// TODO(dudong2): need to bump up this func - (https://github.com/cosmos/cosmos-sdk/pull/10024)
// Constructs a slice of dirty items, to use w/ memIterator.
func (store *Store) dirtyItems(start, end []byte) {
	unsorted := make([]*kv.Pair, 0)
	// If the unsortedCache is too big, its costs too much to determine
	// whats in the subset we are concerned about.
	// If you are interleaving iterator calls with writes, this can easily become an
	// O(N^2) overhead.
	// Even without that, too many range checks eventually becomes more expensive
	// than just not having the cache.
	store.unsortedCache.Range(func(k, _ interface{}) bool {
		key := k.(string)
		if IsKeyInDomain(conv.UnsafeStrToBytes(key), start, end) {
			cacheValue, ok := store.cache.Load(key)
			if ok {
				unsorted = append(unsorted, &kv.Pair{Key: []byte(key), Value: cacheValue.(*cValue).value})
			}
		}
		return true
	})
	store.clearUnsortedCacheSubset(unsorted)
}

func (store *Store) clearUnsortedCacheSubset(unsorted []*kv.Pair) {
	for _, kv := range unsorted {
		store.unsortedCache.Delete(conv.UnsafeBytesToStr(kv.Key))
	}
	sort.Slice(unsorted, func(i, j int) bool {
		return bytes.Compare(unsorted[i].Key, unsorted[j].Key) < 0
	})

	for _, item := range unsorted {
		if item.Value == nil {
			// deleted element, tracked by store.deleted
			// setting arbitrary value
			store.sortedCache.Set(item.Key, []byte{})
			continue
		}
		err := store.sortedCache.Set(item.Key, item.Value)
		if err != nil {
			panic(err)
		}
	}
}

//----------------------------------------
// etc

// Only entrypoint to mutate store.cache.
func (store *Store) setCacheValue(key, value []byte, deleted bool, dirty bool) {
	keyStr := conv.UnsafeBytesToStr(key)
	store.cache.Store(keyStr, &cValue{
		value: value,
		dirty: dirty,
	})
	if deleted {
		store.deleted.Store(keyStr, struct{}{})
	} else {
		store.deleted.Delete(keyStr)
	}
	if dirty {
		store.unsortedCache.Store(string(key), struct{}{})
	}
}

func (store *Store) isDeleted(key string) bool {
	_, ok := store.deleted.Load(key)
	return ok
}

package iavl

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/line/lbm-sdk/store/types"
)

func testCache(t *testing.T, cache types.Cache) {
	key, val := []byte("abc"), []byte("defg")
	cache.Set(key, val)
	if _, ok := cache.(*ristrettoCache); ok {
		time.Sleep(100 * time.Millisecond)
	}
	require.True(t, cache.Has(key))
	require.Equal(t, val, cache.Get(nil, key))
	buf := make([]byte, 0, len(val))
	buf = cache.Get(buf, key)
	require.Equal(t, val, buf)
	cache.Del(key)
	hits, misses, entries, bytes := cache.Stats()
	_, _, _, _ = hits, misses, entries, bytes
}

func TestCache(t *testing.T) {
	cacheSize := 1 * 1024 * 1024
	testCache(t, NewFastCache(cacheSize))
	testCache(t, NewFreeCache(cacheSize))
	testCache(t, NewRistrettoCache(cacheSize))
}

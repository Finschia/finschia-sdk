package cache_test

import (
	"fmt"
	"testing"

	"github.com/line/lbm-sdk/v2/codec"
	store2 "github.com/line/lbm-sdk/v2/testutil/store"
	types2 "github.com/line/lbm-sdk/v2/x/auth/types"
	"github.com/stretchr/testify/require"

	"github.com/line/iavl/v2"
	"github.com/line/tm-db/v2/memdb"

	"github.com/line/lbm-sdk/v2/store/cache"
	iavlstore "github.com/line/lbm-sdk/v2/store/iavl"
	"github.com/line/lbm-sdk/v2/store/types"
)

func TestGetOrSetStoreCache(t *testing.T) {
	db := memdb.NewDB()
	mngr := cache.NewCommitKVStoreCacheManager(cache.DefaultCommitKVStoreCacheSize, cache.NopMetricsProvider())

	sKey := types.NewKVStoreKey("test")
	tree, err := iavl.NewMutableTree(db, 100)
	require.NoError(t, err)
	store := iavlstore.UnsafeNewStore(tree)
	store2 := mngr.GetStoreCache(sKey, store)

	require.NotNil(t, store2)
	require.Equal(t, store2, mngr.GetStoreCache(sKey, store))
}

func TestUnwrap(t *testing.T) {
	db := memdb.NewDB()
	mngr := cache.NewCommitKVStoreCacheManager(cache.DefaultCommitKVStoreCacheSize, cache.NopMetricsProvider())

	sKey := types.NewKVStoreKey("test")
	tree, err := iavl.NewMutableTree(db, 100)
	require.NoError(t, err)
	store := iavlstore.UnsafeNewStore(tree)
	_ = mngr.GetStoreCache(sKey, store)

	require.Equal(t, store, mngr.Unwrap(sKey))
	require.Nil(t, mngr.Unwrap(types.NewKVStoreKey("test2")))
}

func TestStoreCache(t *testing.T) {
	db := memdb.NewDB()
	mngr := cache.NewCommitKVStoreCacheManager(cache.DefaultCommitKVStoreCacheSize, cache.NopMetricsProvider())

	sKey := types.NewKVStoreKey("test")
	tree, err := iavl.NewMutableTree(db, 100)
	require.NoError(t, err)
	store := iavlstore.UnsafeNewStore(tree)
	kvStore := mngr.GetStoreCache(sKey, store)
	cdc := codec.NewProtoCodec(store2.CreateTestInterfaceRegistry())

	for i := 0; i < 10000; i++ {
		key := []byte(fmt.Sprintf("key_%d", i))
		value := store2.ValFmt(i)

		kvStore.Set(key, value, types2.GetAccountMarshalFunc(cdc))

		res := kvStore.Get(key, types2.GetAccountUnmarshalFunc(cdc))
		require.Equal(t, res, value)
		require.Equal(t, res, store.Get(key, types2.GetAccountUnmarshalFunc(cdc)))
	}

	// ristretto cache operates asynchronously.
	// Thus, when Set and Del are called consecutively, they sometimes get a value, so we separate the for loop.
	for i := 0; i < 10000; i++ {
		key := []byte(fmt.Sprintf("key_%d", i))
		kvStore.Delete(key)
	}

	for i := 0; i < 10000; i++ {
		key := []byte(fmt.Sprintf("key_%d", i))
		require.Nil(t, kvStore.Get(key, types2.GetAccountUnmarshalFunc(cdc)))
		require.Nil(t, store.Get(key, types2.GetAccountUnmarshalFunc(cdc)))
	}
}

package mock

import (
	"testing"

	"github.com/stretchr/testify/require"

	dbm "github.com/tendermint/tm-db"

	sdk "github.com/line/lbm-sdk/types"
)

func TestStore(t *testing.T) {
	db := dbm.NewMemDB()
	cms := NewCommitMultiStore()

	key := sdk.NewKVStoreKey("test")
	cms.MountStoreWithDB(key, sdk.StoreTypeIAVL, db)
	err := cms.LoadLatestVersion()
	require.Nil(t, err)

	store := cms.GetKVStore(key)
	require.NotNil(t, store)

	k := []byte("hello")
	v := []byte("world")
	require.False(t, store.Has(k))
	store.Set(k, v)
	require.True(t, store.Has(k))
	require.Equal(t, v, store.Get(k))
	store.Delete(k)
	require.False(t, store.Has(k))
	require.Panics(t, func() { store.Set([]byte(""), v) }, "setting an empty key should panic")
	require.Panics(t, func() { store.Set(nil, v) }, "setting a nil key should panic")
}

func TestMultiStore(t *testing.T) {
	store := multiStore{}
	require.Panics(t, func() { store.Snapshot(1, nil) }, "Snapshot should panic")
	require.Panics(t, func() { store.Restore(1, 1, nil) }, "Restore should panic")
}

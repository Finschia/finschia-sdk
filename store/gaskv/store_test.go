package gaskv_test

import (
	"fmt"
	"testing"

	"github.com/line/lbm-sdk/v2/codec"
	"github.com/line/lbm-sdk/v2/testutil/store"
	types2 "github.com/line/lbm-sdk/v2/x/auth/types"
	"github.com/line/tm-db/v2/memdb"

	"github.com/line/lbm-sdk/v2/store/dbadapter"
	"github.com/line/lbm-sdk/v2/store/gaskv"
	"github.com/line/lbm-sdk/v2/store/types"

	"github.com/stretchr/testify/require"
)

func bz(s string) []byte { return []byte(s) }

func keyFmt(i int) []byte { return bz(fmt.Sprintf("key%0.8d", i)) }

func TestGasKVStoreBasic(t *testing.T) {
	mem := dbadapter.Store{DB: memdb.NewDB()}
	meter := types.NewGasMeter(10000)
	st := gaskv.NewStore(mem, meter, types.KVGasConfig())
	cdc := codec.NewProtoCodec(store.CreateTestInterfaceRegistry())

	require.Equal(t, types.StoreTypeDB, st.GetStoreType())
	require.Panics(t, func() { st.CacheWrap() })
	require.Panics(t, func() { st.CacheWrapWithTrace(nil, nil) })

	require.Panics(t, func() { st.Set(nil, store.ValFmt(1), types2.GetAccountMarshalFunc(cdc)) }, "setting a nil key should panic")
	require.Panics(t, func() { st.Set([]byte(""), store.ValFmt(1), types2.GetAccountMarshalFunc(cdc)) }, "setting an empty key should panic")

	require.Empty(t, st.Get(keyFmt(1), types2.GetAccountUnmarshalFunc(cdc)), "Expected `key1` to be empty")
	st.Set(keyFmt(1), store.ValFmt(1), types2.GetAccountMarshalFunc(cdc))
	require.Equal(t, store.ValFmt(1), st.Get(keyFmt(1), types2.GetAccountUnmarshalFunc(cdc)))
	st.Delete(keyFmt(1))
	require.Empty(t, st.Get(keyFmt(1), types2.GetAccountUnmarshalFunc(cdc)), "Expected `key1` to be empty")
	require.Equal(t, meter.GasConsumed(), types.Gas(6264))
}

func TestGasKVStoreIterator(t *testing.T) {
	mem := dbadapter.Store{DB: memdb.NewDB()}
	meter := types.NewGasMeter(10000)
	st := gaskv.NewStore(mem, meter, types.KVGasConfig())
	cdc := codec.NewProtoCodec(store.CreateTestInterfaceRegistry())
	require.False(t, st.Has(keyFmt(1)))
	require.Empty(t, st.Get(keyFmt(1), types2.GetAccountUnmarshalFunc(cdc)), "Expected `key1` to be empty")
	require.Empty(t, st.Get(keyFmt(2), types2.GetAccountUnmarshalFunc(cdc)), "Expected `key2` to be empty")
	st.Set(keyFmt(1), store.ValFmt(1), types2.GetAccountMarshalFunc(cdc))
	require.True(t, st.Has(keyFmt(1)))
	st.Set(keyFmt(2), store.ValFmt(2), types2.GetAccountMarshalFunc(cdc))

	iterator := st.Iterator(nil, nil)
	require.NoError(t, iterator.Error())

	t.Cleanup(func() {
		if err := iterator.Close(); err != nil {
			t.Fatal(err)
		}
	})
	ka := iterator.Key()
	require.Equal(t, ka, keyFmt(1))
	va := iterator.ValueObject(types2.GetAccountUnmarshalFunc(cdc))
	require.Equal(t, va, store.ValFmt(1))
	iterator.Next()
	kb := iterator.Key()
	require.Equal(t, kb, keyFmt(2))
	vb := iterator.ValueObject(types2.GetAccountUnmarshalFunc(cdc))
	require.Equal(t, vb, store.ValFmt(2))
	iterator.Next()
	require.False(t, iterator.Valid())
	require.Panics(t, iterator.Next)
	require.NoError(t, iterator.Error())

	reverseIterator := st.ReverseIterator(nil, nil)
	t.Cleanup(func() {
		if err := reverseIterator.Close(); err != nil {
			t.Fatal(err)
		}
	})
	require.Equal(t, reverseIterator.Key(), keyFmt(2))
	reverseIterator.Next()
	require.Equal(t, reverseIterator.Key(), keyFmt(1))
	reverseIterator.Next()
	require.False(t, reverseIterator.Valid())
	require.Panics(t, reverseIterator.Next)

	require.Equal(t, types.Gas(8660), meter.GasConsumed())
}

func TestGasKVStoreOutOfGasSet(t *testing.T) {
	mem := dbadapter.Store{DB: memdb.NewDB()}
	meter := types.NewGasMeter(0)
	st := gaskv.NewStore(mem, meter, types.KVGasConfig())
	cdc := codec.NewProtoCodec(store.CreateTestInterfaceRegistry())
	require.Panics(t, func() { st.Set(keyFmt(1), store.ValFmt(1), types2.GetAccountMarshalFunc(cdc)) }, "Expected out-of-gas")
}

func TestGasKVStoreOutOfGasIterator(t *testing.T) {
	mem := dbadapter.Store{DB: memdb.NewDB()}
	meter := types.NewGasMeter(20000)
	st := gaskv.NewStore(mem, meter, types.KVGasConfig())
	cdc := codec.NewProtoCodec(store.CreateTestInterfaceRegistry())
	st.Set(keyFmt(1), store.ValFmt(1), types2.GetAccountMarshalFunc(cdc))
	iterator := st.Iterator(nil, nil)
	iterator.Next()
	require.Panics(t, func() { iterator.Value() }, "Expected out-of-gas")
}

package cachekv_test

import (
	"fmt"
	"sync"

	// "sync"
	"testing"

	"github.com/line/lbm-sdk/v2/codec"
	"github.com/line/lbm-sdk/v2/testutil/store"
	types3 "github.com/line/lbm-sdk/v2/x/auth/types"
	"github.com/line/ostracon/libs/rand"
	tmdb "github.com/line/tm-db/v2"
	"github.com/stretchr/testify/require"

	// ostrand "github.com/line/ostracon/libs/rand"
	// tmdb "github.com/line/tm-db/v2"
	"github.com/line/tm-db/v2/memdb"

	"github.com/line/lbm-sdk/v2/store/cachekv"
	"github.com/line/lbm-sdk/v2/store/dbadapter"
	"github.com/line/lbm-sdk/v2/store/types"
)

func newCacheKVStore() types.CacheKVStore {
	mem := dbadapter.Store{DB: memdb.NewDB()}
	return cachekv.NewStore(mem)
}

func keyFmt(i int) []byte { return bz(fmt.Sprintf("key%0.8d", i)) }

func TestCacheKVStore(t *testing.T) {
	mem := dbadapter.Store{DB: memdb.NewDB()}
	cdc := codec.NewProtoCodec(store.CreateTestInterfaceRegistry())
	st := cachekv.NewStore(mem)
	var value *types3.BaseAccount
	require.Empty(t, st.Get(keyFmt(1), types3.GetAccountUnmarshalFunc(cdc)), "Expected `key1` to be empty")

	// put something in mem and in cache
	value = store.ValFmt(1)
	mem.Set(keyFmt(1), value, types3.GetAccountMarshalFunc(cdc))
	st.Set(keyFmt(1), value, types3.GetAccountMarshalFunc(cdc))
	store.VerifyValue(t, store.ValFmt(1), mem.Get(keyFmt(1), types3.GetAccountUnmarshalFunc(cdc)))
	store.VerifyValue(t, store.ValFmt(1), st.Get(keyFmt(1), types3.GetAccountUnmarshalFunc(cdc)))

	// update it in cache, shoudn't change mem
	value = store.ValFmt(2)
	st.Set(keyFmt(1), value, types3.GetAccountMarshalFunc(cdc))
	store.VerifyValue(t, store.ValFmt(2), st.Get(keyFmt(1), types3.GetAccountUnmarshalFunc(cdc)))
	store.VerifyValue(t, store.ValFmt(1), mem.Get(keyFmt(1), types3.GetAccountUnmarshalFunc(cdc)))

	// write it. should change mem
	st.Write()
	store.VerifyValue(t, store.ValFmt(2), st.Get(keyFmt(1), types3.GetAccountUnmarshalFunc(cdc)))
	store.VerifyValue(t, store.ValFmt(2), mem.Get(keyFmt(1), types3.GetAccountUnmarshalFunc(cdc)))

	// more writes and checks
	st.Write()
	st.Write()
	store.VerifyValue(t, store.ValFmt(2), mem.Get(keyFmt(1), types3.GetAccountUnmarshalFunc(cdc)))
	store.VerifyValue(t, store.ValFmt(2), st.Get(keyFmt(1), types3.GetAccountUnmarshalFunc(cdc)))

	// make a new one, check it
	st = cachekv.NewStore(mem)
	store.VerifyValue(t, store.ValFmt(2), st.Get(keyFmt(1), types3.GetAccountUnmarshalFunc(cdc)))

	// make a new one and delete - should not be removed from mem
	st = cachekv.NewStore(mem)
	st.Delete(keyFmt(1))
	require.Empty(t, st.Get(keyFmt(1), types3.GetAccountUnmarshalFunc(cdc)))
	store.VerifyValue(t,  store.ValFmt(2), mem.Get(keyFmt(1), types3.GetAccountUnmarshalFunc(cdc)))

	// Write. should now be removed from both
	st.Write()
	require.Empty(t, st.Get(keyFmt(1), types3.GetAccountUnmarshalFunc(cdc)), "Expected `key1` to be empty")
	require.Empty(t, mem.Get(keyFmt(1), types3.GetAccountUnmarshalFunc(cdc)), "Expected `key1` to be empty")
}

func TestCacheKVStoreNoNilSet(t *testing.T) {
	mem := dbadapter.Store{DB: memdb.NewDB()}
	cdc := codec.NewProtoCodec(store.CreateTestInterfaceRegistry())
	st := cachekv.NewStore(mem)
	value := store.ValFmt(1)
	st.Set(keyFmt(1), value, types3.GetAccountMarshalFunc(cdc))
	require.Panics(t, func() { st.Set([]byte("key"), nil, types3.GetAccountMarshalFunc(cdc)) }, "setting a nil value should panic")
	require.Panics(t, func() { st.Set(nil, value, types3.GetAccountMarshalFunc(cdc)) }, "setting a nil key should panic")
	require.Panics(t, func() { st.Set([]byte(""), value, types3.GetAccountMarshalFunc(cdc)) }, "setting an empty key should panic")
}

func TestCacheKVStoreNested(t *testing.T) {
	mem := dbadapter.Store{DB: memdb.NewDB()}
	cdc := codec.NewProtoCodec(store.CreateTestInterfaceRegistry())
	st := cachekv.NewStore(mem)

	// set. check its there on st and not on mem.
	st.Set(keyFmt(1), store.ValFmt(1), types3.GetAccountMarshalFunc(cdc))
	require.Empty(t, mem.Get(keyFmt(1), types3.GetAccountUnmarshalFunc(cdc)))
	store.VerifyValue(t, store.ValFmt(1), st.Get(keyFmt(1), types3.GetAccountUnmarshalFunc(cdc)))

	// make a new from st and check
	st2 := cachekv.NewStore(st)
	store.VerifyValue(t, store.ValFmt(1), st2.Get(keyFmt(1), types3.GetAccountUnmarshalFunc(cdc)))

	// update the value on st2, check it only effects st2
	st2.Set(keyFmt(1), store.ValFmt(3), types3.GetAccountMarshalFunc(cdc))
	require.Empty(t, mem.Get(keyFmt(1), types3.GetAccountUnmarshalFunc(cdc)))
	store.VerifyValue(t, store.ValFmt(1), st.Get(keyFmt(1), types3.GetAccountUnmarshalFunc(cdc)))
	store.VerifyValue(t, store.ValFmt(3), st2.Get(keyFmt(1), types3.GetAccountUnmarshalFunc(cdc)))

	// st2 writes to its parent, st. doesnt effect mem
	st2.Write()
	require.Empty(t, mem.Get(keyFmt(1), types3.GetAccountUnmarshalFunc(cdc)))
	store.VerifyValue(t, store.ValFmt(3), st.Get(keyFmt(1), types3.GetAccountUnmarshalFunc(cdc)))

	// updates mem
	st.Write()
	store.VerifyValue(t, store.ValFmt(3), mem.Get(keyFmt(1), types3.GetAccountUnmarshalFunc(cdc)))
}

func TestCacheKVIteratorBounds(t *testing.T) {
	st := newCacheKVStore()
	cdc := codec.NewProtoCodec(store.CreateTestInterfaceRegistry())

	// set some items
	nItems := 5
	for i := 0; i < nItems; i++ {
		st.Set(keyFmt(i), store.ValFmt(i), types3.GetAccountMarshalFunc(cdc))
	}

	// iterate over all of them
	itr := st.Iterator(nil, nil)
	var i = 0
	for ; itr.Valid(); itr.Next() {
		k, v := itr.Key(), itr.ValueObject(types3.GetAccountUnmarshalFunc(cdc))
		require.Equal(t, keyFmt(i), k)
		require.Equal(t, store.ValFmt(i), v)
		i++
	}
	require.Equal(t, nItems, i)

	// iterate over none
	itr = st.Iterator(bz("money"), nil)
	i = 0
	for ; itr.Valid(); itr.Next() {
		i++
	}
	require.Equal(t, 0, i)

	// iterate over lower
	itr = st.Iterator(keyFmt(0), keyFmt(3))
	i = 0
	for ; itr.Valid(); itr.Next() {
		k, v := itr.Key(), itr.ValueObject(types3.GetAccountUnmarshalFunc(cdc))
		require.Equal(t, keyFmt(i), k)
		require.Equal(t, store.ValFmt(i), v)
		i++
	}
	require.Equal(t, 3, i)

	// iterate over upper
	itr = st.Iterator(keyFmt(2), keyFmt(4))
	i = 2
	for ; itr.Valid(); itr.Next() {
		k, v := itr.Key(), itr.ValueObject(types3.GetAccountUnmarshalFunc(cdc))
		require.Equal(t, keyFmt(i), k)
		require.Equal(t, store.ValFmt(i), v)
		i++
	}
	require.Equal(t, 4, i)
}

func TestCacheKVMergeIteratorBasics(t *testing.T) {
	st := newCacheKVStore()
	cdc := codec.NewProtoCodec(store.CreateTestInterfaceRegistry())

	// set and delete an item in the cache, iterator should be empty
	k, v := keyFmt(0), store.ValFmt(0)
	st.Set(k, v, types3.GetAccountMarshalFunc(cdc))
	st.Delete(k)
	assertIterateDomain(t, st, cdc, 0)

	// now set it and assert its there
	st.Set(k, v, types3.GetAccountMarshalFunc(cdc))
	assertIterateDomain(t, st, cdc, 1)

	// write it and assert its there
	st.Write()
	assertIterateDomain(t, st, cdc, 1)

	// remove it in cache and assert its not
	st.Delete(k)
	assertIterateDomain(t, st, cdc, 0)

	// write the delete and assert its not there
	st.Write()
	assertIterateDomain(t, st, cdc, 0)

	// add two keys and assert theyre there
	k1, v1 := keyFmt(1), store.ValFmt(1)
	st.Set(k, v, types3.GetAccountMarshalFunc(cdc))
	st.Set(k1, v1, types3.GetAccountMarshalFunc(cdc))
	assertIterateDomain(t, st, cdc, 2)

	// write it and assert theyre there
	st.Write()
	assertIterateDomain(t, st, cdc, 2)

	// remove one in cache and assert its not
	st.Delete(k1)
	assertIterateDomain(t, st, cdc, 1)

	// write the delete and assert its not there
	st.Write()
	assertIterateDomain(t, st, cdc, 1)

	// delete the other key in cache and asserts its empty
	st.Delete(k)
	assertIterateDomain(t, st, cdc, 0)
}

func TestCacheKVMergeIteratorDeleteLast(t *testing.T) {
	st := newCacheKVStore()
	cdc := codec.NewProtoCodec(store.CreateTestInterfaceRegistry())

	// set some items and write them
	nItems := 5
	for i := 0; i < nItems; i++ {
		st.Set(keyFmt(i), store.ValFmt(i), types3.GetAccountMarshalFunc(cdc))
	}
	st.Write()

	// set some more items and leave dirty
	for i := nItems; i < nItems*2; i++ {
		st.Set(keyFmt(i), store.ValFmt(i), types3.GetAccountMarshalFunc(cdc))
	}

	// iterate over all of them
	assertIterateDomain(t, st, cdc, nItems*2)

	// delete them all
	for i := 0; i < nItems*2; i++ {
		last := nItems*2 - 1 - i
		st.Delete(keyFmt(last))
		assertIterateDomain(t, st, cdc, last)
	}
}

func TestCacheKVMergeIteratorDeletes(t *testing.T) {
	st := newCacheKVStore()
	truth := memdb.NewDB()
	cdc := codec.NewProtoCodec(store.CreateTestInterfaceRegistry())

	// set some items and write them
	nItems := 10
	for i := 0; i < nItems; i++ {
		doOp(t, st, cdc, truth, opSet, i)
	}
	st.Write()

	// delete every other item, starting from 0
	for i := 0; i < nItems; i += 2 {
		doOp(t, st, cdc, truth, opDel, i)
		assertIterateDomainCompare(t, st, cdc, truth)
	}

	// reset
	st = newCacheKVStore()
	truth = memdb.NewDB()

	// set some items and write them
	for i := 0; i < nItems; i++ {
		doOp(t, st, cdc, truth, opSet, i)
	}
	st.Write()

	// delete every other item, starting from 1
	for i := 1; i < nItems; i += 2 {
		doOp(t, st, cdc, truth, opDel, i)
		assertIterateDomainCompare(t, st, cdc, truth)
	}
}

func TestCacheKVMergeIteratorChunks(t *testing.T) {
	st := newCacheKVStore()
	cdc := codec.NewProtoCodec(store.CreateTestInterfaceRegistry())

	// Use the truth to check values on the merge iterator
	truth := memdb.NewDB()

	// sets to the parent
	setRange(t, st, cdc, truth, 0, 20)
	setRange(t, st, cdc, truth, 40, 60)
	st.Write()

	// sets to the cache
	setRange(t, st, cdc, truth, 20, 40)
	setRange(t, st, cdc, truth, 60, 80)
	assertIterateDomainCheck(t, st, cdc, truth, []keyRange{{0, 80}})

	// remove some parents and some cache
	deleteRange(t, st, truth, 15, 25)
	assertIterateDomainCheck(t, st, cdc, truth, []keyRange{{0, 15}, {25, 80}})

	// remove some parents and some cache
	deleteRange(t, st, truth, 35, 45)
	assertIterateDomainCheck(t, st, cdc, truth, []keyRange{{0, 15}, {25, 35}, {45, 80}})

	// write, add more to the cache, and delete some cache
	st.Write()
	setRange(t, st, cdc, truth, 38, 42)
	deleteRange(t, st, truth, 40, 43)
	assertIterateDomainCheck(t, st, cdc, truth, []keyRange{{0, 15}, {25, 35}, {38, 40}, {45, 80}})
}

func TestCacheKVMergeIteratorRandom(t *testing.T) {
	st := newCacheKVStore()
	truth := memdb.NewDB()
	cdc := codec.NewProtoCodec(store.CreateTestInterfaceRegistry())

	start, end := 25, 975
	max := 1000
	setRange(t, st, cdc, truth, start, end)

	// do an op, test the iterator
	for i := 0; i < 2000; i++ {
		doRandomOp(t, st, cdc, truth, max)
		assertIterateDomainCompare(t, st, cdc, truth)
	}
}

// Set, Delete and Write for the same key must be called sequentially.
func TestCacheKVConcurrency(t *testing.T) {
	const NumOps = 2000

	st := newCacheKVStore()
	cdc := codec.NewProtoCodec(store.CreateTestInterfaceRegistry())
	wg := &sync.WaitGroup{}
	wg.Add(NumOps * 3)
	for i := 0; i < NumOps; i++ {
		i := i
		go func() {
			st.Set([]byte(fmt.Sprintf("key%d", i)), store.ValFmt(i), types3.GetAccountMarshalFunc(cdc))
			st.Write()
			wg.Done()
		}()
		go func() {
			st.Get([]byte(fmt.Sprintf("key%d", i)), types3.GetAccountUnmarshalFunc(cdc))
			wg.Done()
		}()
		go func() {
			iter := st.Iterator([]byte("key0"), []byte(fmt.Sprintf("key%d", NumOps)))
			for ; iter.Valid(); iter.Next() {
			}
			wg.Done()
		}()
	}
	wg.Wait()

	for i := 0; i < NumOps; i++ {
		require.Equal(t, store.ValFmt(i), st.Get([]byte(fmt.Sprintf("key%d", i)), types3.GetAccountUnmarshalFunc(cdc)))
	}
}

//-------------------------------------------------------------------------------------------
// do some random ops

const (
	opSet      = 0
	opSetRange = 1
	opDel      = 2
	opDelRange = 3
	opWrite    = 4

	totalOps = 5 // number of possible operations
)

func randInt(n int) int {
	return rand.NewRand().Int() % n
}

// useful for replaying a error case if we find one
func doOp(t *testing.T, st types.CacheKVStore, cdc codec.BinaryMarshaler, truth tmdb.DB, op int, args ...int) {
	switch op {
	case opSet:
		k := args[0]
		st.Set(keyFmt(k), store.ValFmt(k), types3.GetAccountMarshalFunc(cdc))
		v, err := cdc.MarshalInterface(store.ValFmt(k))
		require.NoError(t, err)
		err = truth.Set(keyFmt(k), v)
		require.NoError(t, err)
	case opSetRange:
		start := args[0]
		end := args[1]
		setRange(t, st, cdc, truth, start, end)
	case opDel:
		k := args[0]
		st.Delete(keyFmt(k))
		err := truth.Delete(keyFmt(k))
		require.NoError(t, err)
	case opDelRange:
		start := args[0]
		end := args[1]
		deleteRange(t, st, truth, start, end)
	case opWrite:
		st.Write()
	}
}

func doRandomOp(t *testing.T, st types.CacheKVStore, cdc codec.BinaryMarshaler, truth tmdb.DB, maxKey int) {
	r := randInt(totalOps)
	switch r {
	case opSet:
		k := randInt(maxKey)
		st.Set(keyFmt(k), store.ValFmt(k), types3.GetAccountMarshalFunc(cdc))
		v, err := cdc.MarshalInterface(store.ValFmt(k))
		require.NoError(t, err)
		err = truth.Set(keyFmt(k), v)
		require.NoError(t, err)
	case opSetRange:
		start := randInt(maxKey - 2)
		end := randInt(maxKey-start) + start
		setRange(t, st, cdc, truth, start, end)
	case opDel:
		k := randInt(maxKey)
		st.Delete(keyFmt(k))
		err := truth.Delete(keyFmt(k))
		require.NoError(t, err)
	case opDelRange:
		start := randInt(maxKey - 2)
		end := randInt(maxKey-start) + start
		deleteRange(t, st, truth, start, end)
	case opWrite:
		st.Write()
	}
}

//-------------------------------------------------------------------------------------------

// iterate over whole domain
func assertIterateDomain(t *testing.T, st types.KVStore, cdc codec.BinaryMarshaler, expectedN int) {
	itr := st.Iterator(nil, nil)
	var i = 0
	for ; itr.Valid(); itr.Next() {
		k, v := itr.Key(), itr.ValueObject(types3.GetAccountUnmarshalFunc(cdc))
		require.Equal(t, keyFmt(i), k)
		require.Equal(t, store.ValFmt(i), v)
		i++
	}
	require.Equal(t, expectedN, i)
}

func assertIterateDomainCheck(t *testing.T, st types.KVStore, cdc codec.BinaryMarshaler, mem tmdb.DB, r []keyRange) {
	// iterate over each and check they match the other
	itr := st.Iterator(nil, nil)
	dbIter, err := mem.Iterator(nil, nil) // ground truth
	require.NoError(t, err)
	itr2 := dbadapter.NewDbAdapterIterator(dbIter)

	krc := newKeyRangeCounter(r)
	i := 0
	for ; krc.valid(); krc.next() {
		require.True(t, itr.Valid())
		require.True(t, itr2.Valid())

		// check the key/val matches the ground truth
		k, v := itr.Key(), itr.ValueObject(types3.GetAccountUnmarshalFunc(cdc))
		k2, v2 := itr2.Key(), itr2.ValueObject(types3.GetAccountUnmarshalFunc(cdc))
		require.Equal(t, k, k2)
		require.Equal(t, v, v2)

		// check they match the counter
		require.Equal(t, k, keyFmt(krc.key()))

		itr.Next()
		itr2.Next()
		i++
	}

	require.False(t, itr.Valid())
	require.False(t, itr2.Valid())
}

func assertIterateDomainCompare(t *testing.T, st types.KVStore, cdc codec.BinaryMarshaler, mem tmdb.DB) {
	// iterate over each and check they match the other
	itr := st.Iterator(nil, nil)
	dbIter, err := mem.Iterator(nil, nil) // ground truth
	require.NoError(t, err)
	itr2 := dbadapter.NewDbAdapterIterator(dbIter)
	checkIterators(t, cdc, itr, itr2)
	checkIterators(t, cdc, itr2, itr)
}

func checkIterators(t *testing.T, cdc codec.BinaryMarshaler, itr, itr2 types.Iterator) {
	for ; itr.Valid(); itr.Next() {
		require.True(t, itr2.Valid())
		k, v := itr.Key(), itr.ValueObject(types3.GetAccountUnmarshalFunc(cdc))
		k2, v2 := itr2.Key(), itr2.ValueObject(types3.GetAccountUnmarshalFunc(cdc))
		require.Equal(t, k, k2)
		require.Equal(t, v, v2)
		itr2.Next()
	}
	require.False(t, itr.Valid())
	require.False(t, itr2.Valid())
}

//--------------------------------------------------------

func setRange(t *testing.T, st types.KVStore, cdc codec.BinaryMarshaler, mem tmdb.DB, start, end int) {
	for i := start; i < end; i++ {
		st.Set(keyFmt(i), store.ValFmt(i), types3.GetAccountMarshalFunc(cdc))
		v, err := cdc.MarshalInterface(store.ValFmt(i))
		require.NoError(t, err)
		err = mem.Set(keyFmt(i), v)
		require.NoError(t, err)
	}
}

func deleteRange(t *testing.T, st types.KVStore, mem tmdb.DB, start, end int) {
	for i := start; i < end; i++ {
		st.Delete(keyFmt(i))
		err := mem.Delete(keyFmt(i))
		require.NoError(t, err)
	}
}

//--------------------------------------------------------

type keyRange struct {
	start int
	end   int
}

func (kr keyRange) len() int {
	return kr.end - kr.start
}

func newKeyRangeCounter(kr []keyRange) *keyRangeCounter {
	return &keyRangeCounter{keyRanges: kr}
}

// we can iterate over this and make sure our real iterators have all the right keys
type keyRangeCounter struct {
	rangeIdx  int
	idx       int
	keyRanges []keyRange
}

func (krc *keyRangeCounter) valid() bool {
	maxRangeIdx := len(krc.keyRanges) - 1
	maxRange := krc.keyRanges[maxRangeIdx]

	// if we're not in the max range, we're valid
	if krc.rangeIdx <= maxRangeIdx &&
		krc.idx < maxRange.len() {
		return true
	}

	return false
}

func (krc *keyRangeCounter) next() {
	thisKeyRange := krc.keyRanges[krc.rangeIdx]
	if krc.idx == thisKeyRange.len()-1 {
		krc.rangeIdx++
		krc.idx = 0
	} else {
		krc.idx++
	}
}

func (krc *keyRangeCounter) key() int {
	thisKeyRange := krc.keyRanges[krc.rangeIdx]
	return thisKeyRange.start + krc.idx
}

//--------------------------------------------------------

func bz(s string) []byte { return []byte(s) }

func BenchmarkCacheKVStoreGetNoKeyFound(b *testing.B) {
	st := newCacheKVStore()
	b.ResetTimer()
	// assumes b.N < 2**24
	for i := 0; i < b.N; i++ {
		st.Get([]byte{byte((i & 0xFF0000) >> 16), byte((i & 0xFF00) >> 8), byte(i & 0xFF)}, types.GetBytesUnmarshalFunc())
	}
}

func BenchmarkCacheKVStoreGetKeyFound(b *testing.B) {
	st := newCacheKVStore()
	cdc := codec.NewProtoCodec(store.CreateTestInterfaceRegistry())
	value := store.ValFmt(1)
	for i := 0; i < b.N; i++ {
		arr := []byte{byte((i & 0xFF0000) >> 16), byte((i & 0xFF00) >> 8), byte(i & 0xFF)}
		st.Set(arr, value, types3.GetAccountMarshalFunc(cdc))
	}
	b.ResetTimer()
	// assumes b.N < 2**24
	for i := 0; i < b.N; i++ {
		st.Get([]byte{byte((i & 0xFF0000) >> 16), byte((i & 0xFF00) >> 8), byte(i & 0xFF)}, types3.GetAccountUnmarshalFunc(cdc))
	}
}

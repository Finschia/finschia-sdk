package dbadapter

import (
	"io"

	tmdb "github.com/line/tm-db/v2"

	"github.com/line/lbm-sdk/v2/store/cachekv"
	"github.com/line/lbm-sdk/v2/store/tracekv"
	"github.com/line/lbm-sdk/v2/store/types"
)

// Wrapper type for tmdb.Db with implementation of KVStore
type Store struct {
	tmdb.DB
}

// Get wraps the underlying DB's Get method panicing on error.
func (dsa Store) Get(key []byte, unmarshal func(value []byte) interface{}) interface{} {
	v, err := dsa.DB.Get(key)
	if err != nil {
		panic(err)
	}

	if v == nil {
		return nil
	}
	return unmarshal(v)
}

// Has wraps the underlying DB's Has method panicing on error.
func (dsa Store) Has(key []byte) bool {
	ok, err := dsa.DB.Has(key)
	if err != nil {
		panic(err)
	}

	return ok
}

// Set wraps the underlying DB's Set method panicing on error.
func (dsa Store) Set(key []byte, obj interface{}, marshal func(obj interface{}) []byte) {
	types.AssertValidKey(key)
	types.AssertValidObjectValue(obj)
	value := marshal(obj)
	if err := dsa.DB.Set(key, value); err != nil {
		panic(err)
	}
}

// Delete wraps the underlying DB's Delete method panicing on error.
func (dsa Store) Delete(key []byte) {
	if err := dsa.DB.Delete(key); err != nil {
		panic(err)
	}
}

type dbAdapterIterator struct {
	dbIterator tmdb.Iterator
}

var _ types.Iterator = (*dbAdapterIterator)(nil)

func NewDbAdapterIterator(iterator tmdb.Iterator) types.Iterator {
	return &dbAdapterIterator{
		dbIterator: iterator,
	}
}

func (iter *dbAdapterIterator) Valid() bool {
	return iter.dbIterator.Valid()
}

func (iter *dbAdapterIterator) Next() {
	iter.dbIterator.Next()
}

func (iter *dbAdapterIterator) Key() (key []byte) {
	return iter.dbIterator.Key()
}

func (iter *dbAdapterIterator) Value() (value []byte) {
	return iter.dbIterator.Value()
}

func (iter *dbAdapterIterator) IsValueNil() bool {
	return iter.dbIterator.Value() == nil
}

func (iter *dbAdapterIterator) ValueObject(unmarshal func(value []byte) interface{}) interface{} {
	v := iter.Value()
	if v == nil {
		return nil
	}
	return unmarshal(v)
}

func (iter *dbAdapterIterator) Error() error {
	return iter.dbIterator.Error()
}

func (iter *dbAdapterIterator) Close() error {
	return iter.dbIterator.Close()
}

// Iterator wraps the underlying DB's Iterator method panicing on error.
func (dsa Store) Iterator(start, end []byte) types.Iterator {
	iter, err := dsa.DB.Iterator(start, end)
	if err != nil {
		panic(err)
	}

	return NewDbAdapterIterator(iter)
}

// ReverseIterator wraps the underlying DB's ReverseIterator method panicing on error.
func (dsa Store) ReverseIterator(start, end []byte) types.Iterator {
	iter, err := dsa.DB.ReverseIterator(start, end)
	if err != nil {
		panic(err)
	}

	return NewDbAdapterIterator(iter)
}

// GetStoreType returns the type of the store.
func (Store) GetStoreType() types.StoreType {
	return types.StoreTypeDB
}

// CacheWrap branches the underlying store.
func (dsa Store) CacheWrap() types.CacheWrap {
	return cachekv.NewStore(dsa)
}

// CacheWrapWithTrace implements KVStore.
func (dsa Store) CacheWrapWithTrace(w io.Writer, tc types.TraceContext) types.CacheWrap {
	return cachekv.NewStore(tracekv.NewStore(dsa, w, tc))
}

// tmdb.DB implements KVStore so we can CacheKVStore it.
var _ types.KVStore = Store{}

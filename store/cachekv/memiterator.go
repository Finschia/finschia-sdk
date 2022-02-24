package cachekv

import (
	"bytes"
	"sync"

	"github.com/line/tm-db/v2/memdb"

	"github.com/line/lbm-sdk/store/types"
)

// Iterates over iterKVCache items.
// if key is nil, means it was deleted.
// Implements Iterator.
type memIterator struct {
	types.Iterator

	deleted map[string]struct{}
}

func IsKeyInDomain(key, start, end []byte) bool {
	if bytes.Compare(key, start) < 0 {
		return false
	}
	if end != nil && bytes.Compare(end, key) <= 0 {
		return false
	}
	return true
}

func newMemIterator(start, end []byte, items *memdb.MemDB, deleted *sync.Map, ascending bool) *memIterator {
	var iter types.Iterator
	var err error

	if ascending {
		iter, err = items.Iterator(start, end)
	} else {
		iter, err = items.ReverseIterator(start, end)
	}

	if err != nil {
		panic(err)
	}

	newDeleted := make(map[string]struct{})
	deleted.Range(func(key, value interface{}) bool {
		newDeleted[key.(string)] = value.(struct{})
		return true
	})

	return &memIterator{
		Iterator: iter,

		deleted: newDeleted,
	}
}

func (mi *memIterator) Value() []byte {
	key := mi.Iterator.Key()
	if _, ok := mi.deleted[string(key)]; ok {
		return nil
	}
	return mi.Iterator.Value()
}

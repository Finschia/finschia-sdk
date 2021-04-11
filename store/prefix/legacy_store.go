package prefix

import (
	"bytes"

	"github.com/line/lbm-sdk/v2/store/types"
	tmdb "github.com/tendermint/tm-db"
)


// LegacyStore is an implementation to pass the store
// to the wasmvm for tendermint/tm-db dependency.
// line/tm-db/v2 has a modification to the Iterator interface, occured a type mismatch.
// Do not use it for any purpose other than passing it to wasmvm in x/wasm.
type LegacyStore struct {
	Store
}

func NewLegacyStore(parent types.KVStore, prefix []byte) LegacyStore {
	return LegacyStore{
		Store{
			parent: parent,
			prefix: prefix,	
		},
	}
}

func (s LegacyStore) Iterator(start, end []byte) tmdb.Iterator {
	newstart := cloneAppend(s.prefix, start)

	var newend []byte
	if end == nil {
		newend = cpIncr(s.prefix)
	} else {
		newend = cloneAppend(s.prefix, end)
	}

	iter := s.parent.Iterator(newstart, newend)

	return newLegacyPrefixIterator(s.prefix, start, end, iter)
}

func (s LegacyStore) ReverseIterator(start, end []byte) tmdb.Iterator {
	newstart := cloneAppend(s.prefix, start)

	var newend []byte
	if end == nil {
		newend = cpIncr(s.prefix)
	} else {
		newend = cloneAppend(s.prefix, end)
	}

	iter := s.parent.ReverseIterator(newstart, newend)

	return newLegacyPrefixIterator(s.prefix, start, end, iter)
}

func (s LegacyStore) SdkIterator(start, end []byte) types.Iterator {
	newstart := cloneAppend(s.prefix, start)

	var newend []byte
	if end == nil {
		newend = cpIncr(s.prefix)
	} else {
		newend = cloneAppend(s.prefix, end)
	}

	iter := s.parent.Iterator(newstart, newend)

	return newPrefixIterator(s.prefix, start, end, iter)
}

type legacyPrefixIterator struct {
	prefixIterator
}


func newLegacyPrefixIterator(prefix, start, end []byte, parent types.Iterator) *legacyPrefixIterator {
	return &legacyPrefixIterator{
		prefixIterator{
			prefix: prefix,
			start:  start,
			end:    end,
			iter:   parent,
			valid:  parent.Valid() && bytes.HasPrefix(parent.Key(), prefix),	
		},
	}
}


// Implements Iterator
func (pi *legacyPrefixIterator) Close() {
	pi.iter.Close()
}

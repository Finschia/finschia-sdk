package types

import (
	"bytes"
)

// Iterator over all the keys with a certain prefix in ascending order
func KVStorePrefixIterator(kvs KVStore, prefix []byte) Iterator {
	return kvs.Iterator(prefix, PrefixEndBytes(prefix))
}

// Iterator over all the keys with a certain prefix in descending order.
func KVStoreReversePrefixIterator(kvs KVStore, prefix []byte) Iterator {
	return kvs.ReverseIterator(prefix, PrefixEndBytes(prefix))
}

// DiffKVStores compares two KVstores and returns all the key/value pairs
// that differ from one another. It also skips value comparison for a set of provided prefixes.
func DiffKVStores(a KVStore, b KVStore, prefixesToSkip [][]byte, marshal func(key []byte) func(obj interface{}) []byte,
	unmarshal func(key []byte) func(value []byte) interface{}) (kvAs, kvBs []KOPair) {
	iterA := a.Iterator(nil, nil)

	defer iterA.Close()

	iterB := b.Iterator(nil, nil)

	defer iterB.Close()

	for {
		if !iterA.Valid() && !iterB.Valid() {
			return kvAs, kvBs
		}

		var kvA, kvB KOPair
		if iterA.Valid() {
			kvA = KOPair{Key: iterA.Key(), Value: iterA.ValueObject(unmarshal(iterA.Key()))}

			iterA.Next()
		}

		if iterB.Valid() {
			kvB = KOPair{Key: iterB.Key(), Value: iterB.ValueObject(unmarshal(iterB.Key()))}

			iterB.Next()
		}

		compareValue := true

		for _, prefix := range prefixesToSkip {
			// Skip value comparison if we matched a prefix
			if bytes.HasPrefix(kvA.Key, prefix) || bytes.HasPrefix(kvB.Key, prefix) {
				compareValue = false
				break
			}
		}

		valA := marshal(kvA.Key)(kvA.Value)
		valB := marshal(kvB.Key)(kvB.Value)
		if compareValue && (!bytes.Equal(kvA.Key, kvB.Key) || !bytes.Equal(valA, valB)) {
			kvAs = append(kvAs, kvA)
			kvBs = append(kvBs, kvB)
		}
	}
}

// PrefixEndBytes returns the []byte that would end a
// range query for all []byte with a certain prefix
// Deals with last byte of prefix being FF without overflowing
func PrefixEndBytes(prefix []byte) []byte {
	if len(prefix) == 0 {
		return nil
	}

	end := make([]byte, len(prefix))
	copy(end, prefix)

	for {
		if end[len(end)-1] != byte(255) {
			end[len(end)-1]++
			break
		}

		end = end[:len(end)-1]

		if len(end) == 0 {
			end = nil
			break
		}
	}

	return end
}

// InclusiveEndBytes returns the []byte that would end a
// range query such that the input would be included
func InclusiveEndBytes(inclusiveBytes []byte) []byte {
	return append(inclusiveBytes, byte(0x00))
}

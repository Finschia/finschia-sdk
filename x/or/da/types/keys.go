package types

import (
	"encoding/binary"
)

const (
	// ModuleName defines the module name
	ModuleName = "da"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey defines the module's message routing key
	RouterKey = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_da"

	ParamsKey = byte(0x00)

	CCStoreKey         = byte(0x10)
	CCMetadataKey      = byte(0x11)
	CCBatchIndexPrefix = byte(0x12)

	SCCStoreKey         = byte(0x20)
	SCCMetadataKey      = byte(0x21)
	SCCBatchIndexPrefix = byte(0x22)
)

func GetCCBatchIndexKey(i uint64) []byte {
	return genPrefixIndexKey(i, []byte{CCBatchIndexPrefix})
}

func GetSCCBatchIndexKey(i uint64) []byte {
	return genPrefixIndexKey(i, []byte{SCCBatchIndexPrefix})
}

func genPrefixIndexKey(i uint64, prefix []byte) []byte {
	l := len(prefix)
	k := make([]byte, 8+l)
	copy(k[:l], prefix)
	binary.BigEndian.PutUint64(k[l:], i)
	return k
}

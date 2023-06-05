package types

import (
	"encoding/binary"
)

const (
	// ModuleName defines the module name
	ModuleName = "orda"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey defines the module's message routing key
	RouterKey = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_orda"

	ParamsKey = byte(0x00)

	CCStoreKey         = byte(0x10)
	CCMetadataKey      = byte(0x11)
	CCBatchIndexPrefix = byte(0x12)

	SCCStoreKey         = byte(0x20)
	SCCMetadataKey      = byte(0x21)
	SCCBatchIndexPrefix = byte(0x22)
)

func GetCCBatchIndexKey(rollupName string, i uint64) []byte {
	if i < 1 {
		panic("batch index must be positive")
	}

	return genPrefixIndexKey(append([]byte(rollupName), []byte{CCBatchIndexPrefix}...), i)
}

func GetSCCBatchIndexKey(rollupName string, i uint64) []byte {
	if i < 1 {
		panic("batch index must be positive")
	}

	return genPrefixIndexKey(append([]byte(rollupName), []byte{SCCBatchIndexPrefix}...), i)
}

func genPrefixIndexKey(prefix []byte, i uint64) []byte {
	l := len(prefix)
	k := make([]byte, 8+l)
	copy(k[:l], prefix)
	binary.BigEndian.PutUint64(k[l:], i)
	return k
}

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

	CCStateStoreKey         = byte(0x10)
	QueueTxStateStoreKey    = byte(0x11)
	CCBatchIndexPrefix      = byte(0x12)
	CCQueueTxPrefix         = byte(0x13)
	CCL2HeightToBatchPrefix = byte(0x14)

	SCCStateStoreKey    = byte(0x20)
	SCCMetadataKey      = byte(0x21)
	SCCBatchIndexPrefix = byte(0x22)
)

func GetCCBatchIndexKey(rollupName string, i uint64) []byte {
	if i < 1 {
		panic("batch index must be positive")
	}

	return genPrefixIndexKey(GenRollupPrefix(rollupName, CCBatchIndexPrefix), i)
}

func GetCCQueueTxKey(rollupName string, i uint64) []byte {
	if i < 1 {
		panic("queue tx index must be positive")
	}

	return genPrefixIndexKey(GenRollupPrefix(rollupName, CCQueueTxPrefix), i)
}

func GetSCCBatchIndexKey(rollupName string, i uint64) []byte {
	if i < 1 {
		panic("batch index must be positive")
	}

	return genPrefixIndexKey(GenRollupPrefix(rollupName, SCCBatchIndexPrefix), i)
}

func GetCCStateStoreKey(rollupName string) []byte {
	return GenRollupPrefix(rollupName, CCStateStoreKey)
}

func GetQueueTxStateStoreKey(rollupName string) []byte {
	return GenRollupPrefix(rollupName, QueueTxStateStoreKey)
}

func GetSCCStateStoreKey(rollupName string) []byte {
	return GenRollupPrefix(rollupName, SCCStateStoreKey)
}

func GetCCL2HeightToBatchKey(rollupName string, l2height uint64) []byte {
	return genPrefixIndexKey(GenRollupPrefix(rollupName, CCL2HeightToBatchPrefix), l2height)
}

func genPrefixIndexKey(prefix []byte, i uint64) []byte {
	l := len(prefix)
	k := make([]byte, 8+l)
	copy(k[:l], prefix)
	binary.BigEndian.PutUint64(k[l:], i)
	return k
}

func GenRollupPrefix(rollupName string, prefix byte) []byte {
	key := make([]byte, 1+1+len(rollupName))
	key[0] = prefix
	key[1] = byte(len(rollupName))
	copy(key[2:], []byte(rollupName))
	return key
}

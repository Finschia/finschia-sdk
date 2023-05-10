package types

import "fmt"

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

	CTCStoreKey         = byte(0x10)
	CTCMetadataKey      = byte(0x11)
	CTCBatchIndexPrefix = byte(0x12)

	SCCStoreKey         = byte(0x20)
	SCCMetadataKey      = byte(0x21)
	SCCBatchIndexPrefix = byte(0x22)
)

func BatchIdxKey(idx uint64) []byte {
	return []byte(fmt.Sprintf("BI:%d", idx))
}

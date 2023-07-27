package types

import "encoding/binary"

var _ binary.ByteOrder

const (
	// ModuleName defines the module name
	ModuleName = "rollup"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey is the message route for slashing
	RouterKey = ModuleName

	// QuerierRoute defines the module's query routing key
	QuerierRoute = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_rollup"
)

var (
	RollupKeyPrefix = []byte{0x11}

	SequencerKeyPrefix          = []byte{0x21}
	SequencersByRollupKeyPrefix = []byte{0x22}

	DepositKeyPrefix = []byte{0x31}
)

func RollupKey(rollupName string) []byte {
	return append(RollupKeyPrefix, []byte(rollupName)...)
}

func SequencerKey(sequencerAddress string) []byte {
	return append(SequencerKeyPrefix, []byte(sequencerAddress)...)
}

func DepositKey(rollupName, sequencerAddress string) []byte {
	return append(DepositKeyPrefix, []byte(rollupName+"/"+sequencerAddress)...)
}

func SequencersByRollupKey(rollupName string) []byte {
	return append(SequencersByRollupKeyPrefix, []byte(rollupName)...)
}

package types

import "cosmossdk.io/collections"

const (
	// ModuleName defines the module name
	ModuleName = "fswap"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey is the message route for slashing
	RouterKey = ModuleName

	// QuerierRoute defines the module's query routing key
	QuerierRoute = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_fswap"
)

var (
	paramsKeyValue                 = collections.NewPrefix(0)
	swappedKeyValue                = collections.NewPrefix(1)
	swappableNewCoinAmountKeyValue = collections.NewPrefix(2)
)

func ParamsKey() []byte {
	return paramsKeyValue.Bytes()
}

func SwappedKey() []byte {
	return swappedKeyValue.Bytes()
}

func SwappableNewCoinAmountKey() []byte {
	return swappableNewCoinAmountKeyValue.Bytes()
}

func KeyPrefix(p string) []byte {
	return []byte(p)
}

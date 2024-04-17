package keeper

import "cosmossdk.io/collections"

var (
	paramsKeyValue  = collections.NewPrefix(0)
	swappedKeyValue = collections.NewPrefix(1)
	swapCapKeyValue = collections.NewPrefix(2)
)

func swappedKey() []byte {
	return swappedKeyValue.Bytes()
}

func swapCapKey() []byte {
	return swapCapKeyValue.Bytes()
}

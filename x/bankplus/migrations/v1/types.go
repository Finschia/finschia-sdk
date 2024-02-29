package v1

import sdk "github.com/cosmos/cosmos-sdk/types"

// Keys for bankplus store but this prefix must not be overlap with bank key prefix.
var InactiveAddrsKeyPrefix = []byte{0xa0}

// InactiveAddrKey key of a specific inactiveAddr from store
func InactiveAddrKey(addr sdk.AccAddress) []byte {
	return append(InactiveAddrsKeyPrefix, addr.Bytes()...)
}

package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	ModuleName = "iam"

	StoreKey  = ModuleName
	RouterKey = ModuleName
)

var (
	AddressStoreKeyPrefix = []byte{0x01}
)

func AddressStoreKey(addr sdk.AccAddress) []byte {
	return append(AddressStoreKeyPrefix, addr.Bytes()...)
}

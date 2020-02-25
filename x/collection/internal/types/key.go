package types

import (
	"github.com/cosmos/cosmos-sdk/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	ModuleName = "collection"

	StoreKey  = ModuleName
	RouterKey = ModuleName
)

var (
	AccountKeyPrefix            = []byte{0x00}
	CollectionKeyPrefix         = []byte{0x01}
	SupplyKeyPrefix             = []byte{0x02}
	TokenTypeKeyPrefix          = []byte{0x03}
	TokenChildToParentKeyPrefix = []byte{0x04}
	TokenParentToChildKeyPrefix = []byte{0x05}
	CollectionApprovedKeyPrefix = []byte{0x06}
)

func AccountKey(symbol string, acc sdk.AccAddress) []byte {
	return append(append(AccountKeyPrefix, []byte(symbol)...), acc...)
}

func SupplyKey(symbol string) []byte {
	return append(SupplyKeyPrefix, []byte(symbol)...)
}

func CollectionKey(symbol string) []byte {
	return append(CollectionKeyPrefix, []byte(symbol)...)
}

func TokenTypeKey(symbol, tokenType string) []byte {
	return append(append(TokenTypeKeyPrefix, []byte(symbol)...), []byte(tokenType)...)
}

func TokenChildToParentKey(symbol, tokenID string) []byte {
	return append(append(TokenChildToParentKeyPrefix, []byte(symbol)...), []byte(tokenID)...)
}

func TokenParentToChildSubKey(symbol, parent string) []byte {
	return append(append(TokenParentToChildKeyPrefix, []byte(symbol)...), []byte(parent)...)
}

func TokenParentToChildKey(symbol, parent, child string) []byte {
	return append(append(append(TokenParentToChildKeyPrefix, []byte(symbol)...), []byte(parent)...), []byte(child)...)
}

func CollectionApprovedKey(symbol string, proxy types.AccAddress, approver types.AccAddress) []byte {
	return append(append(append(CollectionApprovedKeyPrefix, []byte(symbol)...), proxy.Bytes()...), approver.Bytes()...)
}

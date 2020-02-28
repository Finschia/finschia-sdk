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
	TokenKeyPrefix              = []byte{0x03}
	TokenTypeKeyPrefix          = []byte{0x04}
	TokenChildToParentKeyPrefix = []byte{0x05}
	TokenParentToChildKeyPrefix = []byte{0x06}
	CollectionApprovedKeyPrefix = []byte{0x07}
)

func AccountKey(contractID string, acc sdk.AccAddress) []byte {
	return append(append(AccountKeyPrefix, []byte(contractID)...), acc...)
}

func SupplyKey(contractID string) []byte {
	return append(SupplyKeyPrefix, []byte(contractID)...)
}

func CollectionKey(contractID string) []byte {
	return append(CollectionKeyPrefix, []byte(contractID)...)
}

func TokenKey(contractID, tokenID string) []byte {
	return append(append(TokenKeyPrefix, []byte(contractID)...), []byte(tokenID)...)
}

func TokenTypeKey(contractID, tokenType string) []byte {
	return append(append(TokenTypeKeyPrefix, []byte(contractID)...), []byte(tokenType)...)
}

func TokenChildToParentKey(contractID, tokenID string) []byte {
	return append(append(TokenChildToParentKeyPrefix, []byte(contractID)...), []byte(tokenID)...)
}

func TokenParentToChildSubKey(contractID, parent string) []byte {
	return append(append(TokenParentToChildKeyPrefix, []byte(contractID)...), []byte(parent)...)
}

func TokenParentToChildKey(contractID, parent, child string) []byte {
	return append(append(append(TokenParentToChildKeyPrefix, []byte(contractID)...), []byte(parent)...), []byte(child)...)
}

func CollectionApprovedKey(contractID string, proxy types.AccAddress, approver types.AccAddress) []byte {
	return append(append(append(CollectionApprovedKeyPrefix, []byte(contractID)...), proxy.Bytes()...), approver.Bytes()...)
}

package types

import (
	sdk "github.com/line/lbm-sdk/types"
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
	NextTokenTypeFTKeyPrefix    = []byte{0x08}
	NextTokenTypeNFTKeyPrefix   = []byte{0x09}
	NextTokenIDNFTKeyPrefix     = []byte{0x0a}
	PermKeyPrefix               = []byte{0x0b}
	AccountOwnNFTKeyPrefix      = []byte{0x0c}
	TokenTypeMintCountPrefix    = []byte{0x0d}
	TokenTypeBurnCountPrefix    = []byte{0x0e}
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

func CollectionApprovedKey(contractID string, proxy sdk.AccAddress, approver sdk.AccAddress) []byte {
	return append(CollectionApproversKey(contractID, proxy), approver.Bytes()...)
}

func CollectionApproversKey(contractID string, proxy sdk.AccAddress) []byte {
	return append(append(CollectionApprovedKeyPrefix, []byte(contractID)...), proxy.Bytes()...)
}

func NextTokenTypeFTKey(contractID string) []byte {
	return append(NextTokenTypeFTKeyPrefix, []byte(contractID)...)
}
func NextTokenTypeNFTKey(contractID string) []byte {
	return append(NextTokenTypeNFTKeyPrefix, []byte(contractID)...)
}
func NextTokenIDNFTKey(contractID, tokenType string) []byte {
	return append(append(NextTokenIDNFTKeyPrefix, []byte(contractID)...), []byte(tokenType)...)
}

func PermKey(contractID string, addr sdk.AccAddress) []byte {
	return append(append(PermKeyPrefix, []byte(contractID)...), addr...)
}

func AccountOwnNFTKey(contractID string, owner sdk.AccAddress, tokenID string) []byte {
	return append(append(append(AccountOwnNFTKeyPrefix, []byte(contractID)...), owner.Bytes()...), []byte(tokenID)...)
}

func TokenTypeMintCount(contractID, tokenType string) []byte {
	return append(append(TokenTypeMintCountPrefix, []byte(contractID)...), []byte(tokenType)...)
}

func TokenTypeBurnCount(contractID, tokenType string) []byte {
	return append(append(TokenTypeBurnCountPrefix, []byte(contractID)...), []byte(tokenType)...)
}

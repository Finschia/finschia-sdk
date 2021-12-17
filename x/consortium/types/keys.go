package types

import (
	sdk "github.com/line/lbm-sdk/types"
)

const (
	// ModuleName is the name of the module
	ModuleName = "consortium"

	// StoreKey is the store key string for consortium
	StoreKey = ModuleName

	// RouterKey is the message route for consortium
	RouterKey = ModuleName

	// QuerierKey is used to handle abci_query requests
	QuerierKey = ModuleName
)

// Keys for consortium store
// Items are stored with the following key: values
//
// - 0x00: Params
//
// - 0x01<valAddress_Bytes>: bool
var (
	ParamsKey              = []byte{0x00}
	ValidatorAuthKeyPrefix = []byte{0x01}
)

// ValidatorAuthKey key for a specific validator from the store
func ValidatorAuthKey(valAddr sdk.ValAddress) []byte {
	return append(ValidatorAuthKeyPrefix, valAddr.Bytes()...)
}

// SplitValidatorAuthKey splits the validator auth key and returns validator
func SplitValidatorAuthKey(key []byte) sdk.ValAddress {
	return sdk.ValAddress(key[1:]) // remove prefix
}

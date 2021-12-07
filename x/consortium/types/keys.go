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
// - 0x00: bool
//
// - 0x01<valAddress_Bytes>: bool
//
// - 0x02<valAddress_Bytes>: bool
var (
	EnabledKeyPrefix = []byte{0x00}

	AllowedValidatorKeyPrefix = []byte{0x01}
	DeniedValidatorKeyPrefix  = []byte{0x02}
)

// AllowedValidatorKey key for a specific validator from the store
func AllowedValidatorKey(valAddr sdk.ValAddress) []byte {
	return append(AllowedValidatorKeyPrefix, valAddr.Bytes()...)
}

// DeniedValidatorKey key for a specific validator from the store
func DeniedValidatorKey(valAddr sdk.ValAddress) []byte {
	return append(DeniedValidatorKeyPrefix, valAddr.Bytes()...)
}

// SplitAllowedValidatorKey split the proposal key and returns the proposal id
func SplitValidatorKey(key []byte) (valAddr sdk.ValAddress) {
	// if len(key[1:]) != 8 {
	// 	panic(fmt.Sprintf("unexpected key length (%d ≠ 8)", len(key[1:])))
	// }

	return sdk.ValAddress(key[1:])
}

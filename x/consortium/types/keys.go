package types

import (
	"fmt"
	"strings"

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
//
// - 0x02<valAddress_Bytes>: bool
var (
	ParamsKey = []byte{0x00}
	ValidatorAuthKeyPrefix = []byte{0x01}

	PendingRejectedDelegationKeyPrefix = []byte{0x10}

	AddressDelimiter = ","
)

// ValidatorAuthKey key for a specific validator from the store
func ValidatorAuthKey(valAddr sdk.ValAddress) []byte {
	return append(ValidatorAuthKeyPrefix, valAddr.Bytes()...)
}

// PendingRejectedDelegationKey key for a specific delegator-validator pair from the store
func PendingRejectedDelegationKey(delAddr sdk.AccAddress, valAddr sdk.ValAddress) []byte {
	return append(append(append(PendingRejectedDelegationKeyPrefix, valAddr.Bytes()...), AddressDelimiter...), delAddr.Bytes()...)
}

// SplitValidatorAuthKey splits the validator auth key and returns validator
func SplitValidatorAuthKey(key []byte) sdk.ValAddress {
	return sdk.ValAddress(key[1:]) // remove prefix
}

// SplitPendingRejectedDelegationKey splits the pending rejected delegation key and returns the delegator-validator pair
func SplitPendingRejectedDelegationKey(key []byte) (sdk.AccAddress, sdk.ValAddress) {
	addrsStr := key[1:] // remove prefix
	addrs := strings.Split(string(addrsStr), AddressDelimiter)
	if len(addrs) != 2 {
		panic(fmt.Sprintf("%s does not contain two addresses", addrsStr))
	}
	return sdk.AccAddress(addrs[1]), sdk.ValAddress(addrs[0])
}

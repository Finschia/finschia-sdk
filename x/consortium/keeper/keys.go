package keeper

import (
	sdk "github.com/line/lbm-sdk/types"
)

// Keys for consortium store
// Items are stored with the following key: values
//
// - 0x00: Params
//
// - 0x01<valAddress_Bytes>: bool
var (
	paramsKey              = []byte{0x00}
	validatorAuthKeyPrefix = []byte{0x01}
)

// validatorAuthKey key for a specific validator from the store
func validatorAuthKey(valAddr sdk.ValAddress) []byte {
	key := make([]byte, len(validatorAuthKeyPrefix)+len(valAddr))
	copy(key, validatorAuthKeyPrefix)
	copy(key[len(validatorAuthKeyPrefix):], valAddr)
	return key
}

// splitValidatorAuthKey splits the validator auth key and returns validator
func splitValidatorAuthKey(key []byte) sdk.ValAddress {
	return sdk.ValAddress(key[1:]) // remove prefix
}

package keeper

import (
	sdk "github.com/line/lbm-sdk/types"
)

var (
	balanceKeyPrefix     = []byte{0x00}
	contractKeyPrefix    = []byte{0x01}
	classKeyPrefix       = []byte{0x02}
	nextClassIDKeyPrefix = []byte{0x03}
)

func balanceKey(contractID string, address sdk.AccAddress, tokenID string) []byte {
	prefix := balanceKeyPrefixByAddress(contractID, address)
	key := make([]byte, len(prefix)+len(tokenID))

	copy(key, prefix)
	copy(key[len(prefix):], tokenID)

	return key
}

func balanceKeyPrefixByAddress(contractID string, address sdk.AccAddress) []byte {
	prefix := balanceKeyPrefixByContractID(contractID)
	key := make([]byte, len(prefix)+1+len(address))

	begin := 0
	copy(key, prefix)

	begin += len(prefix)
	key[begin] = byte(len(address))

	begin++
	copy(key[begin:], address)

	return key
}

func balanceKeyPrefixByContractID(contractID string) []byte {
	key := make([]byte, len(balanceKeyPrefix)+1+len(contractID))

	begin := 0
	copy(key, balanceKeyPrefix)

	begin += len(balanceKeyPrefix)
	key[begin] = byte(len(contractID))

	begin++
	copy(key[begin:], contractID)

	return key
}

// func splitBalanceKey(key []byte) (contractID string, address sdk.AccAddress, tokenID string) {
// 	begin := len(balanceKeyPrefix) + 1
// 	end := begin + int(key[begin-1])
// 	contractID = string(key[begin:end])

// 	begin = end + 1
// 	end = begin + int(key[begin-1])
// 	address = sdk.AccAddress(key[begin:end])

// 	begin = end
// 	tokenID = string(key[begin:])

// 	return
// }

func classKey(contractID string, classID string) []byte {
	prefix := classKeyPrefixByContractID(contractID)
	key := make([]byte, len(prefix)+len(classID))

	copy(key, prefix)
	copy(key[len(prefix):], classID)

	return key
}

func classKeyPrefixByContractID(contractID string) []byte {
	key := make([]byte, len(classKeyPrefix)+1+len(contractID))

	begin := 0
	copy(key, classKeyPrefix)

	begin += len(classKeyPrefix)
	key[begin] = byte(len(contractID))

	begin++
	copy(key[begin:], contractID)

	return key
}

func contractKey(contractID string) []byte {
	key := make([]byte, len(contractKeyPrefix)+len(contractID))

	copy(key, contractKeyPrefix)
	copy(key[len(contractKeyPrefix):], contractID)

	return key
}

func nextClassIDKey(contractID string) []byte {
	key := make([]byte, len(nextClassIDKeyPrefix)+len(contractID))

	copy(key, nextClassIDKeyPrefix)
	copy(key[len(nextClassIDKeyPrefix):], contractID)

	return key
}

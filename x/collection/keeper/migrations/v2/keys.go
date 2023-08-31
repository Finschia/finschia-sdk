package v2

import (
	sdk "github.com/Finschia/finschia-sdk/types"
)

var (
	contractKeyPrefix    = []byte{0x10}
	classKeyPrefix       = []byte{0x11}
	nextClassIDKeyPrefix = []byte{0x12}

	balanceKeyPrefix = []byte{0x20}

	supplyKeyPrefix = []byte{0x40}
	mintedKeyPrefix = []byte{0x41}
	burntKeyPrefix  = []byte{0x42}
)

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

func splitBalanceKey(key []byte) (contractID string, address sdk.AccAddress, tokenID string) {
	begin := len(balanceKeyPrefix) + 1
	end := begin + int(key[begin-1])
	contractID = string(key[begin:end])

	begin = end + 1
	end = begin + int(key[begin-1])
	address = sdk.AccAddress(key[begin:end])

	begin = end
	tokenID = string(key[begin:])

	return
}

func statisticKey(keyPrefix []byte, contractID string, classID string) []byte {
	prefix := statisticKeyPrefixByContractID(keyPrefix, contractID)
	key := make([]byte, len(prefix)+len(classID))

	copy(key, prefix)
	copy(key[len(prefix):], classID)

	return key
}

func statisticKeyPrefixByContractID(keyPrefix []byte, contractID string) []byte {
	key := make([]byte, len(keyPrefix)+1+len(contractID))

	begin := 0
	copy(key, keyPrefix)

	begin += len(keyPrefix)
	key[begin] = byte(len(contractID))

	begin++
	copy(key[begin:], contractID)

	return key
}

func splitStatisticKey(keyPrefix, key []byte) (contractID string, classID string) {
	begin := len(keyPrefix) + 1
	end := begin + int(key[begin-1])
	contractID = string(key[begin:end])

	begin = end
	classID = string(key[begin:])

	return
}

func nextClassIDKey(contractID string) []byte {
	key := make([]byte, len(nextClassIDKeyPrefix)+len(contractID))

	copy(key, nextClassIDKeyPrefix)
	copy(key[len(nextClassIDKeyPrefix):], contractID)

	return key
}

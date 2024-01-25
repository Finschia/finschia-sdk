package v2

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	contractKeyPrefix    = []byte{0x10}
	nextClassIDKeyPrefix = []byte{0x12}

	balanceKeyPrefix = []byte{0x20}

	SupplyKeyPrefix = []byte{0x40}
	MintedKeyPrefix = []byte{0x41}
	BurntKeyPrefix  = []byte{0x42}
)

func ContractKey(contractID string) []byte {
	key := make([]byte, len(contractKeyPrefix)+len(contractID))

	copy(key, contractKeyPrefix)
	copy(key[len(contractKeyPrefix):], contractID)

	return key
}

func BalanceKey(contractID string, address sdk.AccAddress, tokenID string) []byte {
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

func StatisticKey(keyPrefix []byte, contractID, classID string) []byte {
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

func NextClassIDKey(contractID string) []byte {
	key := make([]byte, len(nextClassIDKeyPrefix)+len(contractID))

	copy(key, nextClassIDKeyPrefix)
	copy(key[len(nextClassIDKeyPrefix):], contractID)

	return key
}

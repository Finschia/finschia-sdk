package keeper

import (
	sdk "github.com/Finschia/finschia-sdk/types"
	"github.com/Finschia/finschia-sdk/x/token"
)

var (
	balanceKeyPrefix       = []byte{0x00}
	classKeyPrefix         = []byte{0x01}
	grantKeyPrefix         = []byte{0x02}
	authorizationKeyPrefix = []byte{0x03}

	// statistics keys
	supplyKeyPrefix = []byte{0x04}
	mintKeyPrefix   = []byte{0x05}
	burnKeyPrefix   = []byte{0x06}
)

func classKey(id string) []byte {
	key := make([]byte, len(classKeyPrefix)+len(id))
	copy(key, classKeyPrefix)
	copy(key[len(classKeyPrefix):], id)
	return key
}

func balanceKey(contractID string, address sdk.AccAddress) []byte {
	prefix := balanceKeyPrefixByContractID(contractID)
	key := make([]byte, len(prefix)+len(address))

	copy(key, prefix)
	copy(key[len(prefix):], address)

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

func splitBalanceKey(key []byte) (contractID string, address sdk.AccAddress) {
	begin := len(balanceKeyPrefix) + 1
	end := begin + int(key[begin-1])
	contractID = string(key[begin:end])

	begin = end
	address = key[begin:]

	return
}

func statisticsKey(keyPrefix []byte, contractID string) []byte {
	key := make([]byte, len(keyPrefix)+len(contractID))
	copy(key, keyPrefix)
	copy(key[len(keyPrefix):], contractID)
	return key
}

// func supplyKey(contractID string) []byte {
// 	return statisticsKey(supplyKeyPrefix, contractID)
// }

// func mintKey(contractID string) []byte {
// 	return statisticsKey(mintKeyPrefix, contractID)
// }

// func burnKey(contractID string) []byte {
// 	return statisticsKey(burnKeyPrefix, contractID)
// }

func splitStatisticsKey(key, keyPrefix []byte) (contractID string) {
	return string(key[len(keyPrefix):])
}

// func splitSupplyKey(key []byte) (contractID string) {
// 	return splitStatisticsKey(key, supplyKeyPrefix)
// }

// func splitMintKey(key []byte) (contractID string) {
// 	return splitStatisticsKey(key, mintKeyPrefix)
// }

// func splitBurnKey(key []byte) (contractID string) {
// 	return splitStatisticsKey(key, burnKeyPrefix)
// }

func grantKey(contractID string, grantee sdk.AccAddress, permission token.Permission) []byte {
	prefix := grantKeyPrefixByGrantee(contractID, grantee)
	key := make([]byte, len(prefix)+1)

	copy(key, prefix)
	key[len(prefix)] = byte(permission)

	return key
}

func grantKeyPrefixByGrantee(contractID string, grantee sdk.AccAddress) []byte {
	prefix := grantKeyPrefixByContractID(contractID)
	key := make([]byte, len(prefix)+1+len(grantee))

	begin := 0
	copy(key, prefix)

	begin += len(prefix)
	key[begin] = byte(len(grantee))

	begin++
	copy(key[begin:], grantee)

	return key
}

func grantKeyPrefixByContractID(contractID string) []byte {
	key := make([]byte, len(grantKeyPrefix)+1+len(contractID))

	begin := 0
	copy(key, grantKeyPrefix)

	begin += len(grantKeyPrefix)
	key[begin] = byte(len(contractID))

	begin++
	copy(key[begin:], contractID)

	return key
}

func splitGrantKey(key []byte) (contractID string, grantee sdk.AccAddress, permission token.Permission) {
	begin := len(grantKeyPrefix) + 1
	end := begin + int(key[begin-1])
	contractID = string(key[begin:end])

	begin = end + 1
	end = begin + int(key[begin-1])
	grantee = key[begin:end]

	begin = end
	permission = token.Permission(key[begin])

	return
}

func authorizationKey(contractID string, operator, holder sdk.AccAddress) []byte {
	prefix := authorizationKeyPrefixByOperator(contractID, operator)
	key := make([]byte, len(prefix)+len(holder))

	copy(key, prefix)
	copy(key[len(prefix):], holder)

	return key
}

func authorizationKeyPrefixByOperator(contractID string, operator sdk.AccAddress) []byte {
	prefix := authorizationKeyPrefixByContractID(contractID)
	key := make([]byte, len(prefix)+1+len(operator))

	begin := 0
	copy(key, prefix)

	begin += len(prefix)
	key[begin] = byte(len(operator))

	begin++
	copy(key[begin:], operator)

	return key
}

func authorizationKeyPrefixByContractID(contractID string) []byte {
	key := make([]byte, len(authorizationKeyPrefix)+1+len(contractID))

	begin := 0
	copy(key, authorizationKeyPrefix)

	begin += len(authorizationKeyPrefix)
	key[begin] = byte(len(contractID))

	begin++
	copy(key[begin:], contractID)

	return key
}

func splitAuthorizationKey(key []byte) (contractID string, operator, holder sdk.AccAddress) {
	begin := len(authorizationKeyPrefix) + 1
	end := begin + int(key[begin-1])
	contractID = string(key[begin:end])

	begin = end + 1
	end = begin + int(key[begin-1])
	operator = key[begin:end]

	begin = end
	holder = key[begin:]

	return
}

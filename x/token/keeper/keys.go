package keeper

import (
	sdk "github.com/line/lbm-sdk/types"
	"github.com/line/lbm-sdk/x/token"
)

var (
	classKeyPrefix   = []byte{0x01}
	balanceKeyPrefix = []byte{0x02}
	grantKeyPrefix   = []byte{0x03}
	authorizationKeyPrefix = []byte{0x04}

	// statistics keys
	supplyKeyPrefix = []byte{0x05}
	mintKeyPrefix   = []byte{0x06}
	burnKeyPrefix   = []byte{0x07}
)

func classKey(id string) []byte {
	key := make([]byte, len(classKeyPrefix)+len(id))
	copy(key, classKeyPrefix)
	copy(key[len(classKeyPrefix):], id)
	return key
}

func balanceKey(classID string, address sdk.AccAddress) []byte {
	prefix := balanceKeyPrefixByContractID(classID)
	key := make([]byte, len(prefix)+len(address))

	copy(key, prefix)
	copy(key[len(prefix):], address)

	return key
}

func balanceKeyPrefixByContractID(classID string) []byte {
	key := make([]byte, len(balanceKeyPrefix)+1+len(classID))

	begin := 0
	copy(key, balanceKeyPrefix)

	begin += len(balanceKeyPrefix)
	key[begin] = byte(len(classID))

	begin++
	copy(key[begin:], classID)

	return key
}

func splitBalanceKey(key []byte) (classID string, address sdk.AccAddress) {
	begin := len(balanceKeyPrefix) + 1
	end := begin + int(key[begin-1])
	classID = string(key[begin:end])

	begin = end
	address = sdk.AccAddress(key[begin:])

	return
}

func statisticsKey(keyPrefix []byte, classID string) []byte {
	key := make([]byte, len(keyPrefix)+len(classID))
	copy(key, keyPrefix)
	copy(key[len(keyPrefix):], classID)
	return key
}

// func supplyKey(classID string) []byte {
// 	return statisticsKey(supplyKeyPrefix, classID)
// }

// func mintKey(classID string) []byte {
// 	return statisticsKey(mintKeyPrefix, classID)
// }

// func burnKey(classID string) []byte {
// 	return statisticsKey(burnKeyPrefix, classID)
// }

func splitStatisticsKey(key, keyPrefix []byte) (classID string) {
	return string(key[len(keyPrefix):])
}

// func splitSupplyKey(key []byte) (classID string) {
// 	return splitStatisticsKey(key, supplyKeyPrefix)
// }

// func splitMintKey(key []byte) (classID string) {
// 	return splitStatisticsKey(key, mintKeyPrefix)
// }

// func splitBurnKey(key []byte) (classID string) {
// 	return splitStatisticsKey(key, burnKeyPrefix)
// }

func grantKey(classID string, grantee sdk.AccAddress, permission token.Permission) []byte {
	prefix := grantKeyPrefixByGrantee(classID, grantee)
	key := make([]byte, len(prefix)+1)

	copy(key, prefix)
	key[len(prefix)] = byte(permission)

	return key
}

func grantKeyPrefixByGrantee(classID string, grantee sdk.AccAddress) []byte {
	prefix := grantKeyPrefixByContractID(classID)
	key := make([]byte, len(prefix)+1+len(grantee))

	begin := 0
	copy(key, prefix)

	begin += len(prefix)
	key[begin] = byte(len(grantee))

	begin++
	copy(key[begin:], grantee)

	return key
}

func grantKeyPrefixByContractID(classID string) []byte {
	key := make([]byte, len(grantKeyPrefix)+1+len(classID))

	begin := 0
	copy(key, grantKeyPrefix)

	begin += len(grantKeyPrefix)
	key[begin] = byte(len(classID))

	begin++
	copy(key[begin:], classID)

	return key
}

func splitGrantKey(key []byte) (classID string, grantee sdk.AccAddress, permission token.Permission) {
	begin := len(grantKeyPrefix) + 1
	end := begin + int(key[begin-1])
	classID = string(key[begin:end])

	begin = end + 1
	end = begin + int(key[begin-1])
	grantee = sdk.AccAddress(key[begin:end])

	begin = end
	permission = token.Permission(key[begin])

	return
}

func authorizationKey(classID string, proxy, approver sdk.AccAddress) []byte {
	prefix := authorizationKeyPrefixByProxy(classID, proxy)
	key := make([]byte, len(prefix)+len(approver))

	copy(key, prefix)
	copy(key[len(prefix):], approver)

	return key
}

func authorizationKeyPrefixByProxy(classID string, proxy sdk.AccAddress) []byte {
	prefix := authorizationKeyPrefixByContractID(classID)
	key := make([]byte, len(prefix)+1+len(proxy))

	begin := 0
	copy(key, prefix)

	begin += len(prefix)
	key[begin] = byte(len(proxy))

	begin++
	copy(key[begin:], proxy)

	return key
}

func authorizationKeyPrefixByContractID(classID string) []byte {
	key := make([]byte, len(authorizationKeyPrefix)+1+len(classID))

	begin := 0
	copy(key, authorizationKeyPrefix)

	begin += len(authorizationKeyPrefix)
	key[begin] = byte(len(classID))

	begin++
	copy(key[begin:], classID)

	return key
}

func splitAuthorizationKey(key []byte) (classID string, proxy, approver sdk.AccAddress) {
	begin := len(authorizationKeyPrefix) + 1
	end := begin + int(key[begin-1])
	classID = string(key[begin:end])

	begin = end + 1
	end = begin + int(key[begin-1])
	proxy = sdk.AccAddress(key[begin:end])

	begin = end
	approver = sdk.AccAddress(key[begin:])

	return
}

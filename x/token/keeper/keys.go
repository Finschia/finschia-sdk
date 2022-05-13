package keeper

import (
	sdk "github.com/line/lbm-sdk/types"
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

func balanceKey(addr sdk.AccAddress, classID string) []byte {
	key := make([]byte, len(balanceKeyPrefix)+1+len(addr)+len(classID))

	begin := 0
	copy(key, balanceKeyPrefix)

	begin += len(balanceKeyPrefix)
	key[begin] = byte(len(addr))

	begin++
	copy(key[begin:], addr)

	begin += len(addr)
	copy(key[begin:], classID)

	return key
}

func splitBalanceKey(key []byte) (addr sdk.AccAddress, classID string) {
	begin := len(balanceKeyPrefix) + 1
	end := begin + int(key[begin-1])
	addr = sdk.AccAddress(key[begin:end])

	begin = end
	classID = string(key[begin:])

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

func grantKey(grantee sdk.AccAddress, classID, permission string) []byte {
	key := make([]byte, len(grantKeyPrefix)+1+len(grantee)+1+len(classID)+len(permission))

	begin := 0
	copy(key, grantKeyPrefix)

	begin += len(grantKeyPrefix)
	key[begin] = byte(len(grantee))

	begin++
	copy(key[begin:], grantee)

	begin += len(grantee)
	key[begin] = byte(len(classID))

	begin++
	copy(key[begin:], classID)

	begin += len(classID)
	copy(key[begin:], permission)

	return key
}

func splitGrantKey(key []byte) (grantee sdk.AccAddress, classID, permission string) {
	begin := len(grantKeyPrefix) + 1
	end := begin + int(key[begin-1])
	grantee = sdk.AccAddress(key[begin:end])

	begin = end + 1
	end = begin + int(key[begin-1])
	classID = string(key[begin:end])

	begin = end
	permission = string(key[begin:])

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
	key := make([]byte, len(authorizationKeyPrefix)+1+len(classID)+1+len(proxy))

	begin := 0
	copy(key, authorizationKeyPrefix)

	begin += len(authorizationKeyPrefix)
	key[begin] = byte(len(classID))

	begin++
	copy(key[begin:], classID)

	begin += len(classID)
	key[begin] = byte(len(proxy))

	begin++
	copy(key[begin:], proxy)

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

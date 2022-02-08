package keeper

import (
	sdk "github.com/line/lbm-sdk/types"
)

var (
	classKeyPrefix = []byte{0x01}
	balanceKeyPrefix = []byte{0x02}
	grantKeyPrefix = []byte{0x03}
	proxyKeyPrefix = []byte{0x04}

	// statistics keys
	supplyKeyPrefix = []byte{0x05}
	mintKeyPrefix = []byte{0x06}
	burnKeyPrefix = []byte{0x07}
)

func classKey(id string) []byte {
	key := make([]byte, len(classKeyPrefix)+len(id))
	copy(key, classKeyPrefix)
	copy(key[len(classKeyPrefix):], id)
	return key
}

func balanceKey(addr sdk.AccAddress, classId string) []byte {
	key := make([]byte, len(balanceKeyPrefix)+1+len(addr)+len(classId))

	begin := 0
	copy(key, balanceKeyPrefix)

	begin += len(balanceKeyPrefix)
	key[begin] = byte(len(addr))

	begin += 1
	copy(key[begin:], addr)

	begin += len(addr)
	copy(key[begin:], classId)

	return key
}

func splitBalanceKey(key []byte) (addr sdk.AccAddress, classId string) {
	begin := len(balanceKeyPrefix) + 1
	end := begin + int(key[begin - 1])
	addr = sdk.AccAddress(key[begin:end])

	begin = end
	classId = string(key[begin:])

	return
}

func statisticsKey(keyPrefix []byte, classId string) []byte {
	key := make([]byte, len(keyPrefix)+len(classId))
	copy(key, keyPrefix)
	copy(key[len(keyPrefix):], classId)
	return key
}

func supplyKey(classId string) []byte {
	return statisticsKey(supplyKeyPrefix, classId)
}

func mintKey(classId string) []byte {
	return statisticsKey(mintKeyPrefix, classId)
}

func burnKey(classId string) []byte {
	return statisticsKey(burnKeyPrefix, classId)
}

func splitStatisticsKey(key, keyPrefix []byte) (classId string) {
	return string(key[len(keyPrefix):])
}

func splitSupplyKey(key []byte) (classId string) {
	return splitStatisticsKey(key, supplyKeyPrefix)
}

func splitMintKey(key []byte) (classId string) {
	return splitStatisticsKey(key, mintKeyPrefix)
}

func splitBurnKey(key []byte) (classId string) {
	return splitStatisticsKey(key, burnKeyPrefix)
}

func grantKey(grantee sdk.AccAddress, classId, action string) []byte {
	key := make([]byte, len(grantKeyPrefix)+1+len(grantee)+1+len(classId)+len(action))

	begin := 0
	copy(key, proxyKeyPrefix)

	begin += len(proxyKeyPrefix)
	key[begin] = byte(len(grantee))

	begin += 1
	copy(key[begin:], grantee)

	begin += len(grantee)
	key[begin] = byte(len(classId))

	begin += 1
	copy(key[begin:], classId)

	begin += len(classId)
	copy(key[begin:], action)

	return key
}

func splitGrantKey(key []byte) (grantee sdk.AccAddress, classId, action string) {
	begin := len(grantKeyPrefix) + 1
	end := begin + int(key[begin - 1])
	grantee = sdk.AccAddress(key[begin:end])

	begin = end + 1
	end = begin + int(key[begin - 1])
	classId = string(key[begin:end])

	begin = end
	action = string(key[begin:])

	return
}

func proxyKey(classId string, proxy, approver sdk.AccAddress) []byte {
	key := make([]byte, len(proxyKeyPrefix)+1+len(classId)+1+len(proxy)+len(approver))

	begin := 0
	copy(key, grantKeyPrefix)

	begin += len(grantKeyPrefix)
	key[begin] = byte(len(classId))

	begin += 1
	copy(key[begin:], classId)

	begin += len(classId)
	key[begin] = byte(len(proxy))

	begin += 1
	copy(key[begin:], proxy)

	begin += len(proxy)
	copy(key[begin:], approver)

	return key
}

func splitProxyKey(key []byte) (classId string, proxy, approver sdk.AccAddress) {
	begin := len(proxyKeyPrefix) + 1
	end := begin + int(key[begin - 1])
	classId = string(key[begin:end])

	begin = end + 1
	end = begin + int(key[begin - 1])
	proxy = sdk.AccAddress(key[begin:end])

	begin = end
	approver = sdk.AccAddress(key[begin:])

	return
}

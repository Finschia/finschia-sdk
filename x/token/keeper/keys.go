package keeper

import (
	sdk "github.com/line/lbm-sdk/types"
)

var (
	classKeyPrefix = []byte{0x01}
	balanceKeyPrefix = []byte{0x02}
	supplyKeyPrefix = []byte{0x03}
	grantKeyPrefix = []byte{0x04}
	proxyKeyPrefix = []byte{0x05}
)

func classKey(id string) []byte {
	key := make([]byte, len(classKeyPrefix)+len(id))
	copy(key, classKeyPrefix)
	copy(key[len(classKeyPrefix):], id)
	return key
}

func balanceKey(addr sdk.AccAddress, classId string) []byte {
	key := make([]byte, len(balanceKeyPrefix)+len(addr)+len(classId))
	copy(key, balanceKeyPrefix)
	copy(key[len(balanceKeyPrefix):], addr)
	copy(key[len(balanceKeyPrefix)+len(addr):], classId)
	return key
}

func supplyKey(supplyType, classId string) []byte {
	key := make([]byte, len(supplyKeyPrefix)+len(supplyType)+len(classId))
	copy(key, supplyKeyPrefix)
	copy(key[len(supplyKeyPrefix):], supplyType)
	copy(key[len(supplyKeyPrefix)+len(supplyType):], classId)
	return key
}

func grantKey(grantee sdk.AccAddress, classId, action string) []byte {
	key := make([]byte, len(grantKeyPrefix)+len(grantee)+len(classId)+len(action))
	copy(key, proxyKeyPrefix)
	copy(key[len(proxyKeyPrefix):], grantee)
	copy(key[len(proxyKeyPrefix)+len(grantee):], classId)
	copy(key[len(proxyKeyPrefix)+len(grantee)+len(classId):], action)
	return key
}

func proxyKey(classId string, proxy, approver sdk.AccAddress) []byte {
	key := make([]byte, len(proxyKeyPrefix)+len(classId)+len(proxy)+len(approver))
	copy(key, grantKeyPrefix)
	copy(key[len(grantKeyPrefix):], classId)
	copy(key[len(grantKeyPrefix)+len(classId):], proxy)
	copy(key[len(grantKeyPrefix)+len(classId)+len(proxy):], approver)
	return key
}

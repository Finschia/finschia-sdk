package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/Finschia/finschia-sdk/x/collection"
)

var (
	paramsKey = []byte{0x00}

	contractKeyPrefix    = []byte{0x10}
	classKeyPrefix       = []byte{0x11}
	nextClassIDKeyPrefix = []byte{0x12}
	nextTokenIDKeyPrefix = []byte{0x13}

	balanceKeyPrefix = []byte{0x20}
	ownerKeyPrefix   = []byte{0x21}
	nftKeyPrefix     = []byte{0x22}
	parentKeyPrefix  = []byte{0x23}
	childKeyPrefix   = []byte{0x24}

	authorizationKeyPrefix = []byte{0x30}
	grantKeyPrefix         = []byte{0x31}

	supplyKeyPrefix = []byte{0x40}
	mintedKeyPrefix = []byte{0x41}
	burntKeyPrefix  = []byte{0x42}

	legacyTokenKeyPrefix     = []byte{0xf0}
	legacyTokenTypeKeyPrefix = []byte{0xf1}
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

// ----------------------------------------------------------------------------
// owner
func ownerKey(contractID, tokenID string) []byte {
	prefix := ownerKeyPrefixByContractID(contractID)
	key := make([]byte, len(prefix)+len(tokenID))

	copy(key, prefix)
	copy(key[len(prefix):], tokenID)

	return key
}

func ownerKeyPrefixByContractID(contractID string) []byte {
	key := make([]byte, len(ownerKeyPrefix)+1+len(contractID))

	begin := 0
	copy(key, ownerKeyPrefix)

	begin += len(ownerKeyPrefix)
	key[begin] = byte(len(contractID))

	begin++
	copy(key[begin:], contractID)

	return key
}

// ----------------------------------------------------------------------------
// nft
func nftKey(contractID, tokenID string) []byte {
	prefix := nftKeyPrefixByContractID(contractID)
	key := make([]byte, len(prefix)+len(tokenID))

	copy(key, prefix)
	copy(key[len(prefix):], tokenID)

	return key
}

func nftKeyPrefixByContractID(contractID string) []byte {
	key := make([]byte, len(nftKeyPrefix)+1+len(contractID))

	begin := 0
	copy(key, nftKeyPrefix)

	begin += len(nftKeyPrefix)
	key[begin] = byte(len(contractID))

	begin++
	copy(key[begin:], contractID)

	return key
}

func splitNFTKey(key []byte) (contractID, tokenID string) {
	begin := len(nftKeyPrefix) + 1
	end := begin + int(key[begin-1])
	contractID = string(key[begin:end])

	begin = end
	tokenID = string(key[begin:])

	return
}

// ----------------------------------------------------------------------------
// parent
func parentKey(contractID, tokenID string) []byte {
	prefix := parentKeyPrefixByContractID(contractID)
	key := make([]byte, len(prefix)+len(tokenID))

	copy(key, prefix)
	copy(key[len(prefix):], tokenID)

	return key
}

func parentKeyPrefixByContractID(contractID string) []byte {
	key := make([]byte, len(parentKeyPrefix)+1+len(contractID))

	begin := 0
	copy(key, parentKeyPrefix)

	begin += len(parentKeyPrefix)
	key[begin] = byte(len(contractID))

	begin++
	copy(key[begin:], contractID)

	return key
}

func splitParentKey(key []byte) (contractID, tokenID string) {
	begin := len(parentKeyPrefix) + 1
	end := begin + int(key[begin-1])
	contractID = string(key[begin:end])

	begin = end
	tokenID = string(key[begin:])

	return
}

// ----------------------------------------------------------------------------
// child
func childKey(contractID, tokenID, childID string) []byte {
	prefix := childKeyPrefixByTokenID(contractID, tokenID)
	key := make([]byte, len(prefix)+len(childID))

	copy(key, prefix)
	copy(key[len(prefix):], childID)

	return key
}

func childKeyPrefixByTokenID(contractID, tokenID string) []byte {
	prefix := childKeyPrefixByContractID(contractID)
	key := make([]byte, len(prefix)+1+len(tokenID))

	begin := 0
	copy(key, prefix)

	begin += len(prefix)
	key[begin] = byte(len(tokenID))

	begin++
	copy(key[begin:], tokenID)

	return key
}

func childKeyPrefixByContractID(contractID string) []byte {
	key := make([]byte, len(childKeyPrefix)+1+len(contractID))

	begin := 0
	copy(key, childKeyPrefix)

	begin += len(childKeyPrefix)
	key[begin] = byte(len(contractID))

	begin++
	copy(key[begin:], contractID)

	return key
}

func splitChildKey(key []byte) (contractID, tokenID, childID string) {
	begin := len(childKeyPrefix) + 1
	end := begin + int(key[begin-1])
	contractID = string(key[begin:end])

	begin = end + 1
	end = begin + int(key[begin-1])
	tokenID = string(key[begin:end])

	begin = end
	childID = string(key[begin:])

	return
}

// ----------------------------------------------------------------------------
func contractKey(contractID string) []byte {
	key := make([]byte, len(contractKeyPrefix)+len(contractID))

	copy(key, contractKeyPrefix)
	copy(key[len(contractKeyPrefix):], contractID)

	return key
}

func classKey(contractID, classID string) []byte {
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

func nextTokenIDKey(contractID, classID string) []byte {
	prefix := nextTokenIDKeyPrefixByContractID(contractID)
	key := make([]byte, len(prefix)+len(classID))

	copy(key, prefix)
	copy(key[len(prefix):], classID)

	return key
}

func nextTokenIDKeyPrefixByContractID(contractID string) []byte {
	key := make([]byte, len(nextTokenIDKeyPrefix)+1+len(contractID))

	begin := 0
	copy(key, nextTokenIDKeyPrefix)

	begin += len(nextTokenIDKeyPrefix)
	key[begin] = byte(len(contractID))

	begin++
	copy(key[begin:], contractID)

	return key
}

func splitNextTokenIDKey(key []byte) (contractID, classID string) {
	begin := len(nextTokenIDKeyPrefix) + 1
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

// ----------------------------------------------------------------------------
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
	operator = sdk.AccAddress(key[begin:end])

	begin = end
	holder = sdk.AccAddress(key[begin:])

	return
}

// ----------------------------------------------------------------------------
func grantKey(contractID string, grantee sdk.AccAddress, permission collection.Permission) []byte {
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

func splitGrantKey(key []byte) (contractID string, grantee sdk.AccAddress, permission collection.Permission) {
	begin := len(grantKeyPrefix) + 1
	end := begin + int(key[begin-1])
	contractID = string(key[begin:end])

	begin = end + 1
	end = begin + int(key[begin-1])
	grantee = sdk.AccAddress(key[begin:end])

	begin = end
	permission = collection.Permission(key[begin])

	return
}

// ----------------------------------------------------------------------------
// statistics
func statisticKey(keyPrefix []byte, contractID, classID string) []byte {
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

func splitStatisticKey(keyPrefix, key []byte) (contractID, classID string) {
	begin := len(keyPrefix) + 1
	end := begin + int(key[begin-1])
	contractID = string(key[begin:end])

	begin = end
	classID = string(key[begin:])

	return
}

// ----------------------------------------------------------------------------
// legacy keys
func legacyTokenKey(contractID, tokenID string) []byte {
	prefix := legacyTokenKeyPrefixByContractID(contractID)
	key := make([]byte, len(prefix)+len(tokenID))

	copy(key, prefix)
	copy(key[len(prefix):], tokenID)

	return key
}

func legacyTokenKeyPrefixByContractID(contractID string) []byte {
	key := make([]byte, len(legacyTokenKeyPrefix)+1+len(contractID))

	begin := 0
	copy(key, legacyTokenKeyPrefix)

	begin += len(legacyTokenKeyPrefix)
	key[begin] = byte(len(contractID))

	begin++
	copy(key[begin:], contractID)

	return key
}

func legacyTokenTypeKey(contractID, tokenType string) []byte {
	prefix := legacyTokenTypeKeyPrefixByContractID(contractID)
	key := make([]byte, len(prefix)+len(tokenType))

	copy(key, prefix)
	copy(key[len(prefix):], tokenType)

	return key
}

func legacyTokenTypeKeyPrefixByContractID(contractID string) []byte {
	key := make([]byte, len(legacyTokenTypeKeyPrefix)+1+len(contractID))

	begin := 0
	copy(key, legacyTokenTypeKeyPrefix)

	begin += len(legacyTokenTypeKeyPrefix)
	key[begin] = byte(len(contractID))

	begin++
	copy(key[begin:], contractID)

	return key
}

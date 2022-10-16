package keeper

import (
	"encoding/binary"
	"time"

	sdk "github.com/line/lbm-sdk/types"
)

// Keys for foundation store
// Items are stored with the following key: values
var (
	paramsKey         = []byte{0x00}
	foundationInfoKey = []byte{0x01}

	memberKeyPrefix          = []byte{0x10}
	previousProposalIDKey    = []byte{0x11}
	proposalKeyPrefix        = []byte{0x12}
	proposalByVPEndKeyPrefix = []byte{0x13}
	voteKeyPrefix            = []byte{0x14}

	grantKeyPrefix = []byte{0x20}

	poolKey = []byte{0x30}

	govMintKey = []byte{0x40}
)

// Uint64FromBytes converts a byte array to uint64
func Uint64FromBytes(bz []byte) uint64 {
	return binary.BigEndian.Uint64(bz)
}

// Uint64ToBytes converts a number in uint64 to a byte array
func Uint64ToBytes(number uint64) []byte {
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, number)
	return bz
}

// memberKey key for a specific member from the store
func memberKey(address sdk.AccAddress) []byte {
	key := make([]byte, len(memberKeyPrefix)+len(address))
	copy(key, memberKeyPrefix)
	copy(key[len(memberKeyPrefix):], address)
	return key
}

// proposalKey key for a specific proposal from the store
func proposalKey(id uint64) []byte {
	idBz := Uint64ToBytes(id)

	key := make([]byte, len(proposalKeyPrefix)+len(idBz))
	copy(key, proposalKeyPrefix)
	copy(key[len(proposalKeyPrefix):], idBz)
	return key
}

func voteKey(proposalID uint64, voter sdk.AccAddress) []byte {
	idBz := Uint64ToBytes(proposalID)
	key := make([]byte, len(voteKeyPrefix)+len(idBz)+len(voter))

	begin := 0
	copy(key[begin:], voteKeyPrefix)

	begin += len(voteKeyPrefix)
	copy(key[begin:], idBz)

	begin += len(idBz)
	copy(key[begin:], voter)

	return key
}

func proposalByVPEndKey(id uint64, end time.Time) []byte {
	idBz := Uint64ToBytes(id)
	endBz := sdk.FormatTimeBytes(end)
	key := make([]byte, len(proposalByVPEndKeyPrefix)+1+len(idBz)+len(endBz))

	begin := 0
	copy(key[begin:], proposalByVPEndKeyPrefix)

	begin += len(proposalByVPEndKeyPrefix)
	key[begin] = byte(len(idBz))

	begin++
	copy(key[begin:], idBz)

	begin += len(idBz)
	copy(key[begin:], endBz)

	return key
}

// func splitProposalByVPEndKey(key []byte) (proposalID uint64, vpEnd time.Time) {
// 	begin := len(proposalByVPEndKeyPrefix) + 1
// 	end := begin + int(key[begin-1]) // uint64
// 	proposalID = Uint64FromBytes(key[begin:end])

// 	begin = end
// 	vpEnd, err := sdk.ParseTimeBytes(key[begin:])
// 	if err != nil {
// 		panic(err)
// 	}

// 	return
// }

func grantKey(grantee sdk.AccAddress, url string) []byte {
	prefix := grantKeyPrefixByGrantee(grantee)
	key := make([]byte, len(prefix)+len(url))

	copy(key, prefix)
	copy(key[len(prefix):], url)

	return key
}

func grantKeyPrefixByGrantee(grantee sdk.AccAddress) []byte {
	key := make([]byte, len(grantKeyPrefix)+1+len(grantee))

	begin := 0
	copy(key[begin:], grantKeyPrefix)

	begin += len(grantKeyPrefix)
	key[begin] = byte(len(grantee))

	begin++
	copy(key[begin:], grantee)

	return key
}

func splitGrantKey(key []byte) (grantee sdk.AccAddress, url string) {
	begin := len(grantKeyPrefix) + 1
	end := begin + int(key[begin-1])
	grantee = key[begin:end]

	begin = end
	url = string(key[begin:])

	return
}

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

	treasuryKey = []byte{0x10}

	validatorAuthKeyPrefix = []byte{0x20}

	memberKeyPrefix          = []byte{0x30}
	previousProposalIdKey    = []byte{0x31}
	proposalKeyPrefix        = []byte{0x32}
	proposalByVPEndKeyPrefix = []byte{0x33}
	voteKeyPrefix            = []byte{0x34}
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

// validatorAuthKey key for a specific validator from the store
func validatorAuthKey(valAddr sdk.ValAddress) []byte {
	key := make([]byte, len(validatorAuthKeyPrefix)+len(valAddr))
	copy(key, validatorAuthKeyPrefix)
	copy(key[len(validatorAuthKeyPrefix):], valAddr)
	return key
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

func voteKey(proposalId uint64, voter sdk.AccAddress) []byte {
	idBz := Uint64ToBytes(proposalId)
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

	begin += 1
	copy(key[begin:], idBz)

	begin += len(idBz)
	copy(key[begin:], endBz)

	return key
}

func splitProposalByVPEndKey(key []byte) (proposalId uint64, vpEnd time.Time) {
	begin := len(proposalByVPEndKeyPrefix) + 1
	end := begin + int(key[begin-1]) // uint64
	proposalId = Uint64FromBytes(key[begin:end])

	begin = end
	vpEnd, err := sdk.ParseTimeBytes(key[begin:])
	if err != nil {
		panic(err)
	}

	return
}

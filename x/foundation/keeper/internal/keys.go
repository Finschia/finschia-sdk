package internal

import (
	"encoding/binary"
	"time"

	sdk "github.com/Finschia/finschia-sdk/types"
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

	censorshipKeyPrefix = []byte{0x20}
	grantKeyPrefix      = []byte{0x21}

	poolKey = []byte{0x30}

	// deprecatedGovMintKey Deprecated. Don't use it again.
	deprecatedGovMintKey = []byte{0x40}
	_                    = deprecatedGovMintKey
)

// must be constant
var lenTime = len(sdk.FormatTimeBytes(time.Now()))

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
	prefix := memberKeyPrefix
	key := make([]byte, len(prefix)+len(address))

	copy(key, prefix)
	copy(key[len(prefix):], address)

	return key
}

// proposalKey key for a specific proposal from the store
func proposalKey(id uint64) []byte {
	prefix := proposalKeyPrefix
	idBz := Uint64ToBytes(id)
	key := make([]byte, len(prefix)+len(idBz))

	copy(key, prefix)
	copy(key[len(prefix):], idBz)

	return key
}

func voteKey(proposalID uint64, voter sdk.AccAddress) []byte {
	prefix := voteKeyPrefix
	idBz := Uint64ToBytes(proposalID)
	key := make([]byte, len(prefix)+len(idBz)+len(voter))

	begin := 0
	copy(key[begin:], prefix)

	begin += len(prefix)
	copy(key[begin:], idBz)

	begin += len(idBz)
	copy(key[begin:], voter)

	return key
}

func proposalByVPEndKey(vpEnd time.Time, id uint64) []byte {
	prefix := proposalByVPEndKeyPrefix
	vpEndBz := sdk.FormatTimeBytes(vpEnd)
	idBz := Uint64ToBytes(id)
	key := make([]byte, len(prefix)+lenTime+len(idBz))

	begin := 0
	copy(key[begin:], prefix)

	begin += len(prefix)
	copy(key[begin:], vpEndBz)

	begin += len(vpEndBz)
	copy(key[begin:], idBz)

	return key
}

func splitProposalByVPEndKey(key []byte) (vpEnd time.Time, id uint64) {
	prefix := proposalByVPEndKeyPrefix
	begin := len(prefix)
	end := begin + lenTime
	vpEnd, err := sdk.ParseTimeBytes(key[begin:end])
	if err != nil {
		panic(err)
	}

	begin = end
	id = Uint64FromBytes(key[begin:])

	return
}

// memberKey key for a specific member from the store
func censorshipKey(url string) []byte {
	prefix := censorshipKeyPrefix
	key := make([]byte, len(prefix)+len(url))

	copy(key, prefix)
	copy(key[len(prefix):], url)

	return key
}

func grantKey(grantee sdk.AccAddress, url string) []byte {
	prefix := grantKeyPrefixByGrantee(grantee)
	key := make([]byte, len(prefix)+len(url))

	copy(key, prefix)
	copy(key[len(prefix):], url)

	return key
}

func grantKeyPrefixByGrantee(grantee sdk.AccAddress) []byte {
	prefix := grantKeyPrefix
	key := make([]byte, len(prefix)+1+len(grantee))

	begin := 0
	copy(key[begin:], prefix)

	begin += len(prefix)
	key[begin] = byte(len(grantee))

	begin++
	copy(key[begin:], grantee)

	return key
}

func splitGrantKey(key []byte) (grantee sdk.AccAddress, url string) {
	prefix := grantKeyPrefix

	begin := len(prefix) + 1
	end := begin + int(key[begin-1])
	grantee = key[begin:end]

	begin = end
	url = string(key[begin:])

	return
}

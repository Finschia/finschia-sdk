package types

import (
	"encoding/base64"
	"encoding/binary"
	"errors"
	"math/big"
	"strings"

	"github.com/Finschia/finschia-sdk/types"
	"github.com/tendermint/crypto/blake2b"
)

// AccAddressFromAddressSeed create an AccAddress from addressSeed string and iss string
// AccAddress = blake2b_256(iss_L, iss, addressSeed)
func AccAddressFromAddressSeed(addrSeed, issBase64 string) (types.AccAddress, error) {
	if len(strings.TrimSpace(addrSeed)) == 0 {
		return types.AccAddress{}, errors.New("empty address seed string is not allowed")
	}

	// convert addrSeed string to big endian bytes
	addrSeedBigInt, ok := new(big.Int).SetString(addrSeed, 10)
	if !ok {
		return types.AccAddress{}, errors.New("invalid address seed")
	}
	addrSeedBytes := addrSeedBigInt.Bytes()
	issBytes, err := base64.StdEncoding.DecodeString(issBase64)
	if err != nil {
		return types.AccAddress{}, err
	}

	// convert the issBytes length to big endian 2 bytes
	issL := make([]byte, 2)
	binary.BigEndian.PutUint16(issL, uint16(len(issBytes)))

	// hash by blake2b
	hasher, err := blake2b.New256(nil)
	if err != nil {
		return types.AccAddress{}, err
	}

	hasher.Write(issL)
	hasher.Write(issBytes)
	hasher.Write(addrSeedBytes)

	addrBytes := hasher.Sum(nil)

	return addrBytes, nil
}

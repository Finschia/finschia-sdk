package types

import (
	"encoding/binary"
	"errors"
	"github.com/Finschia/finschia-sdk/types"
	"github.com/tendermint/crypto/blake2b"
	"strings"
)

// AccAddressFromAddressSeed create an AccAddress from addressSeed string and iss string
// AccAddress = blake2b_256(iss_L, iss, addressSeed)
func AccAddressFromAddressSeed(addrSeed, iss string) (types.AccAddress, error) {
	if len(strings.TrimSpace(addrSeed)) == 0 {
		return types.AccAddress{}, errors.New("empty address seed string is not allowed")
	}

	addrSeedBytes := []byte(addrSeed)

	if iss == "accounts.google.com" {
		iss = "https://accounts.google.com"
	}
	issBytes := []byte(iss)

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

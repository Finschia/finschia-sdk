package keeper

import (
	"github.com/Finschia/finschia-sdk/types/address"
	sdkerrors "github.com/Finschia/finschia-sdk/types/errors"
)

var (
	swapPrefix       = []byte{0x01}
	swapStatsKey     = []byte{0x02}
	swappedKeyPrefix = []byte{0x03}
)

// swapKey key(prefix + fromDenom + toDenom)
func swapKey(fromDenom, toDenom string) []byte {
	denoms := combineDenoms(fromDenom, toDenom)
	return append(swapPrefix, denoms...)
}

// swappedKey key(prefix + (lengthPrefixed+)fromDenom + (lengthPrefixed+)toDenom)
func swappedKey(fromDenom, toDenom string) []byte {
	denoms := combineDenoms(fromDenom, toDenom)
	return append(swappedKeyPrefix, denoms...)
}

func combineDenoms(fromDenom, toDenom string) []byte {
	lengthPrefixedFromDenom, err := address.LengthPrefix([]byte(fromDenom))
	if err != nil {
		panic(sdkerrors.ErrInvalidRequest.Wrapf("fromDenom length should be max %d bytes, got %d", address.MaxAddrLen, len(fromDenom)))
	}
	lengthPrefixedToDenom, err := address.LengthPrefix([]byte(toDenom))
	if err != nil {
		panic(sdkerrors.ErrInvalidRequest.Wrapf("toDenom length should be max %d bytes, got %d", address.MaxAddrLen, len(toDenom)))
	}
	return append(lengthPrefixedFromDenom, lengthPrefixedToDenom...)
}

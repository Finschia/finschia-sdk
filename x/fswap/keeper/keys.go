package keeper

var (
	swapPrefix       = []byte{0x01}
	swapStatsKey     = []byte{0x02}
	swappedKeyPrefix = []byte{0x03}
)

// swapKey key(prefix + (lengthPrefixed+)fromDenom + (lengthPrefixed+)toDenom)
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
	lengthPrefixedFromDenom := lengthPrefix([]byte(fromDenom))
	lengthPrefixedToDenom := lengthPrefix([]byte(toDenom))
	return append(lengthPrefixedFromDenom, lengthPrefixedToDenom...)
}

// lengthPrefix prefixes the address bytes with its length, this is used
// for example for variable-length components in store keys.
func lengthPrefix(bz []byte) []byte {
	bzLen := len(bz)
	if bzLen == 0 {
		return bz
	}

	return append([]byte{byte(bzLen)}, bz...)
}

package keeper

var (
	swapPrefix       = []byte{0x01}
	swapStatsKey     = []byte{0x02}
	swappedKeyPrefix = []byte{0x03}
)

// swapKey key(prefix + fromDenom + toDenom)
func swapKey(fromDenom, toDenom string) []byte {
	key := append(swapPrefix, fromDenom...)
	return append(key, toDenom...)
}

// swappedKey key(prefix + fromDenom + toDenom)
func swappedKey(fromDenom, toDenom string) []byte {
	key := append(swappedKeyPrefix, fromDenom...)
	return append(key, toDenom...)
}

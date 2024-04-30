package keeper

var (
	swapInitPrefix   = []byte{0x01}
	swappedKeyPrefix = []byte{0x02}
)

// swapInitKey key(prefix + toDenom)
func swapInitKey(toDenom string) []byte {
	return append(swapInitPrefix, toDenom...)
}

// swappedKey key(prefix + toDenom)
func swappedKey(toDenom string) []byte {
	return append(swappedKeyPrefix, toDenom...)
}

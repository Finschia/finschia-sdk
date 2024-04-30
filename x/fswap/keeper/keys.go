package keeper

var (
	fswapInitPrefix         = []byte{0x01}
	swappedKeyPrefix        = []byte{0x02}
	allowFswapInitOnceValue = []byte{0x03}
)

// fswapInitKey key(prefix + toDenom)
func fswapInitKey(toDenom string) []byte {
	return append(fswapInitPrefix, toDenom...)
}

// swappedKey key(prefix + toDenom)
func swappedKey(toDenom string) []byte {
	return append(swappedKeyPrefix, toDenom...)
}

func allowFswapInitOnceKey() []byte {
	return allowFswapInitOnceValue
}

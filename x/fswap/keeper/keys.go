package keeper

var (
	fswapInitPrefix  = []byte{0x01}
	swappedKeyPrefix = []byte{0x02}
)

// fswapInitKey key(prefix + toDenom)
func fswapInitKey(toDenom string) []byte {
	return append(fswapInitPrefix, toDenom...)
}

func swappedKey(toDenom string) []byte {
	return append(swappedKeyPrefix, toDenom...)
}

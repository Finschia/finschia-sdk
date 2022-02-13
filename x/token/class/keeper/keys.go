package keeper

var (
	nonceKey    = []byte{0x01}
	idKeyPrefix = []byte{0x02}
)

func idKey(id string) []byte {
	key := make([]byte, len(idKeyPrefix)+len(id))
	copy(key, idKeyPrefix)
	copy(key[len(idKeyPrefix):], id)
	return key
}

func splitIDKey(key []byte) (id string) {
	return string(key[len(idKeyPrefix):])
}

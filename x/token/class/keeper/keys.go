package keeper

var (
	nextIdKey = []byte{0x01}
	idKeyPrefix = []byte{0x02}
)

func idKey(id string) []byte {
	key := make([]byte, len(idKeyPrefix)+len(id))
	copy(key, idKeyPrefix)
	copy(key[len(idKeyPrefix):], id)
	return key
}

package v3

const ClassStoreKey = "class"

var (
	nonceKey    = []byte{0x01}
	idKeyPrefix = []byte{0x02}
)

func splitIDKey(key []byte) (id string) {
	return string(key[len(idKeyPrefix):])
}

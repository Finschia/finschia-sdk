package keeper

var (
	paramsKeyPrefix = []byte{0x00}
)

func paramsKey() []byte {
	return paramsKeyPrefix
}

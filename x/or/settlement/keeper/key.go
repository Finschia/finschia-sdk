package keeper

var (
	paramsKeyPrefix    = []byte{0x00}
	challengeKeyPrefix = []byte{0x01}
)

func paramsKey() []byte {
	return paramsKeyPrefix
}

func challengeKey(challengID string) []byte {
	return append(challengeKeyPrefix, challengID...)
}

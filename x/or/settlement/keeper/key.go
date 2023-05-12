package keeper

import (
	"encoding/binary"
)

var (
	paramsKeyPrefix    = []byte{0x00}
	challengeKeyPrefix = []byte{0x01}
)

func paramsKey() []byte {
	return paramsKeyPrefix
}

func challengeKey(challengID int64) []byte {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, uint64(challengID))
	return append(challengeKeyPrefix, b...)
}

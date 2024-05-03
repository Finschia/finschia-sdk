package types

import "encoding/binary"

const (
	// ModuleName is the module name constant used in many places
	ModuleName = "fbridge"

	// StoreKey is the store key string for distribution
	StoreKey = ModuleName
)

var (
	KeyParams              = []byte{0x01} // key for fbridge module params
	KeyNextSeqSend         = []byte{0x02} // key for the next bridge send sequence
	KeySeqToBlocknumPrefix = []byte{0x03} // key prefix for the sequence to block number mapping
)

func SeqToBlocknumKey(seq uint64) []byte {
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, seq)
	return append(KeySeqToBlocknumPrefix, bz...)
}

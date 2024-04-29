package types

const (
	// ModuleName is the module name constant used in many places
	ModuleName = "fbridge"

	// StoreKey is the store key string for distribution
	StoreKey = ModuleName
)

var (
	KeyParams      = []byte{0x01} // key for fbridge module params
	KeyNextSeqSend = []byte{0x02} // key for the next bridge send sequence
)

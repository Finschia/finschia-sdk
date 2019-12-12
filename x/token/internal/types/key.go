package types

const (
	ModuleName = "token"

	StoreKey  = ModuleName
	RouterKey = ModuleName
)

var (
	TokenSymbolKeyPrefix = []byte{0x00}
	CollectionKeyPrefix  = []byte{0x01}
)

func TokenSymbolKey(symbol string) []byte {
	return append(TokenSymbolKeyPrefix, []byte(symbol)...)
}

func CollectionKey(symbol string) []byte {
	return append(CollectionKeyPrefix, []byte(symbol)...)
}

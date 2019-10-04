package types

const (
	ModuleName = "token"

	StoreKey  = ModuleName
	RouterKey = ModuleName
)

var (
	TokenSymbolKeyPrefix = []byte{0x00}
)

func TokenSymbolKey(symbol string) []byte {
	return append(TokenSymbolKeyPrefix, []byte(symbol)...)
}

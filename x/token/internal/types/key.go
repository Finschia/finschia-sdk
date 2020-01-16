package types

const (
	ModuleName = "token"

	StoreKey  = ModuleName
	RouterKey = ModuleName
)

var (
	TokenDenomKeyPrefix = []byte{0x00}
	CollectionKeyPrefix = []byte{0x01}
)

func TokenDenomKey(denom string) []byte {
	return append(TokenDenomKeyPrefix, []byte(denom)...)
}

func CollectionKey(denom string) []byte {
	return append(CollectionKeyPrefix, []byte(denom)...)
}

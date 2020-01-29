package types

const (
	ModuleName = "token"

	StoreKey  = ModuleName
	RouterKey = ModuleName
)

var (
	TokenDenomKeyPrefix            = []byte{0x00}
	CollectionKeyPrefix            = []byte{0x01}
	TokenChildToParentKeyPrefix    = []byte{0x02}
	TokenParentToChildKeyPrefix    = []byte{0x03}
	TokenParentToChildSubKeyPrefix = []byte{0x04}
)

func TokenDenomKey(denom string) []byte {
	return append(TokenDenomKeyPrefix, []byte(denom)...)
}

func CollectionKey(denom string) []byte {
	return append(CollectionKeyPrefix, []byte(denom)...)
}

func TokenChildToParentKey(token *BaseIDNFT) []byte {
	return append(TokenChildToParentKeyPrefix, []byte(token.GetDenom())...)
}

func TokenParentToChildKey(token *BaseIDNFT) []byte {
	return append(TokenParentToChildKeyPrefix, []byte(token.GetDenom())...)
}

func TokenParentToChildSubKey(token *BaseIDNFT) []byte {
	return append(TokenParentToChildSubKeyPrefix, []byte(token.GetDenom())...)
}

func ParentToChildSubKeyToToken(prefix []byte, key []byte) (tokenDenom string) {
	return string(key[len(prefix)+1:])
}

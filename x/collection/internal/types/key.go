package types

import "github.com/cosmos/cosmos-sdk/types"

const (
	ModuleName = "collection"

	StoreKey  = ModuleName
	RouterKey = ModuleName
)

var (
	CollectionKeyPrefix            = []byte{0x01}
	TokenTypeKeyPrefix             = []byte{0x02}
	TokenChildToParentKeyPrefix    = []byte{0x03}
	TokenParentToChildKeyPrefix    = []byte{0x04}
	TokenParentToChildSubKeyPrefix = []byte{0x05}
	CollectionApprovedKeyPrefix    = []byte{0x06}
)

func CollectionKey(denom string) []byte {
	return append(CollectionKeyPrefix, []byte(denom)...)
}

func TokenTypeKey(symbol, tokenType string) []byte {
	key := append(TokenTypeKeyPrefix, []byte(symbol)...)
	return append(key, []byte(tokenType)...)
}

func TokenChildToParentKey(token Token) []byte {
	return append(TokenChildToParentKeyPrefix, []byte(token.GetDenom())...)
}

func TokenParentToChildKey(token Token) []byte {
	return append(TokenParentToChildKeyPrefix, []byte(token.GetDenom())...)
}

func TokenParentToChildSubKey(token Token) []byte {
	return append(TokenParentToChildSubKeyPrefix, []byte(token.GetDenom())...)
}

func ParentToChildSubKeyToToken(prefix []byte, key []byte) (tokenDenom string) {
	return string(key[len(prefix)+1:])
}

func CollectionApprovedKey(proxy types.AccAddress, approver types.AccAddress, symbol string) []byte {
	return append(append(append(CollectionApprovedKeyPrefix, proxy.Bytes()...), approver.Bytes()...), symbol...)
}

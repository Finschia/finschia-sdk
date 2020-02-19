package types

import "github.com/cosmos/cosmos-sdk/types"

const (
	ModuleName = "token"

	StoreKey  = ModuleName
	RouterKey = ModuleName
)

var (
	TokenSymbolKeyPrefix = []byte{0x00}
	BlacklistKeyPrefix   = []byte{0x01}
)

func BlacklistKey(addr types.AccAddress, action string) []byte {
	key := append(BlacklistKeyPrefix, addr...)
	key = append(key, []byte(":"+action)...)
	return key
}

func TokenSymbolKey(symbol string) []byte {
	return append(TokenSymbolKeyPrefix, []byte(symbol)...)
}

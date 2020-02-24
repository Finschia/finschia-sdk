package types

import sdk "github.com/cosmos/cosmos-sdk/types"

const (
	ModuleName = "token"

	StoreKey  = ModuleName
	RouterKey = ModuleName
)

var (
	TokenSymbolKeyPrefix = []byte{0x00}
	BlacklistKeyPrefix   = []byte{0x01}
	AccountKeyPrefix     = []byte{0x02}
	SupplyKeyPrefix      = []byte{0x03}
)

func BlacklistKey(addr sdk.AccAddress, action string) []byte {
	key := append(BlacklistKeyPrefix, addr...)
	key = append(key, []byte(":"+action)...)
	return key
}

func TokenSymbolKey(symbol string) []byte {
	return append(TokenSymbolKeyPrefix, []byte(symbol)...)
}

func AccountKey(symbol string, acc sdk.AccAddress) []byte {
	return append(append(AccountKeyPrefix, []byte(symbol)...), acc...)
}

func SupplyKey(symbol string) []byte {
	return append(SupplyKeyPrefix, []byte(symbol)...)
}

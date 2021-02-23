package types

import sdk "github.com/cosmos/cosmos-sdk/types"

const (
	// module name
	ModuleName = "coin"
	StoreKey   = ModuleName

	ActionTransferTo = "transferTo"
)

var (
	BlacklistKeyPrefix = []byte{0x00}
)

func BlacklistKey(addr sdk.AccAddress, action string) []byte {
	key := append(BlacklistKeyPrefix, addr...)
	key = append(key, []byte(":"+action)...)
	return key
}

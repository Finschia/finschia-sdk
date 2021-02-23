package types

import sdk "github.com/line/lbm-sdk/types"

const (
	ModuleName = "token"

	StoreKey  = ModuleName
	RouterKey = ModuleName
)

var (
	TokenKeyPrefix         = []byte{0x00}
	BlacklistKeyPrefix     = []byte{0x01}
	AccountKeyPrefix       = []byte{0x02}
	SupplyKeyPrefix        = []byte{0x03}
	PermKeyPrefix          = []byte{0x04}
	TokenApprovedKeyPrefix = []byte{0x05}
)

func BlacklistKey(addr sdk.AccAddress, action string) []byte {
	key := append(BlacklistKeyPrefix, addr...)
	key = append(key, []byte(":"+action)...)
	return key
}

func TokenKey(contractID string) []byte {
	return append(TokenKeyPrefix, []byte(contractID)...)
}

func SupplyKey(contractID string) []byte {
	return append(SupplyKeyPrefix, []byte(contractID)...)
}

func AccountKey(contractID string, addr sdk.AccAddress) []byte {
	return append(append(AccountKeyPrefix, []byte(contractID)...), addr...)
}

func PermKey(contractID string, addr sdk.AccAddress) []byte {
	return append(append(PermKeyPrefix, []byte(contractID)...), addr...)
}

func TokenApprovedKey(contractID string, proxy sdk.AccAddress, approver sdk.AccAddress) []byte {
	return append(append(append(TokenApprovedKeyPrefix, []byte(contractID)...), proxy.Bytes()...), approver.Bytes()...)
}

func TokenApproversKey(contractID string, proxy sdk.AccAddress) []byte {
	return append(append(TokenApprovedKeyPrefix, []byte(contractID)...), proxy.Bytes()...)
}

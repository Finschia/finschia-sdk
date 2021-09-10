package types

import (
	"fmt"
	"strings"

	sdk "github.com/line/lbm-sdk/types"
)

const (
	// ModuleName defines the module name
	ModuleName = "bank"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey defines the module's message routing key
	RouterKey = ModuleName

	// QuerierRoute defines the module's query routing key
	QuerierRoute = ModuleName
)

// KVStore keys
var (
	BalancesPrefix      = []byte("balances")
	SupplyKey           = []byte{0x00}
	DenomMetadataPrefix = []byte{0x1}

	// Contract: Address must not contain this character
	AddressDenomDelimiter = ","
)

// DenomMetadataKey returns the denomination metadata key.
func DenomMetadataKey(denom string) []byte {
	d := []byte(denom)
	return append(DenomMetadataPrefix, d...)
}

// AddressFromBalancesStore returns an account address from a balances prefix
// store. The key must not contain the perfix BalancesPrefix as the prefix store
// iterator discards the actual prefix.
func AddressFromBalancesStore(key []byte) sdk.AccAddress {
	addr := string(key)
	if !strings.Contains(addr, sdk.Bech32MainPrefix) {
		panic(fmt.Sprintf("unexpected account address key; key does not start with (%s): %s",
			sdk.Bech32MainPrefix, addr))
	}
	index := strings.Index(addr, AddressDenomDelimiter)
	if index <= 0 {
		panic(fmt.Sprintf("AddressBalance store key does not contain the delimiter(,): %s", addr))
	}
	return sdk.AccAddress(addr[:index])
}

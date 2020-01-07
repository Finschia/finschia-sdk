package types

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"regexp"
)

const (
	reSymbolString            = `[a-z][a-z0-9]{2,15}`
	reSymbolStringReserved    = `[a-z][a-z0-9]{2,4}`
	reSymbolStringUserDefined = `[a-z][a-z0-9]{5,7}`
	reSymbolStringTokenID     = `[a-z0-9]{8}`
)

var (
	reSymbol                = regexp.MustCompile(fmt.Sprintf(`^%s$`, reSymbolString))
	reSymbolReserved        = regexp.MustCompile(fmt.Sprintf(`^%s$`, reSymbolStringReserved))
	reSymbolUserDefined     = regexp.MustCompile(fmt.Sprintf(`^%s$`, reSymbolStringUserDefined))
	reSymbolTokenID         = regexp.MustCompile(fmt.Sprintf(`^%s$`, reSymbolStringTokenID))
	reSymbolCollectionToken = regexp.MustCompile(fmt.Sprintf(`^%s%s$`, reSymbolStringUserDefined, reSymbolStringTokenID))
)

const (
	AccAddrSuffixLen = 3
)

func ValidateReg(symbol string, reg *regexp.Regexp) error {
	if !reg.MatchString(symbol) {
		return fmt.Errorf("symbol [%s] mismatched to [%s]", symbol, reg.String())
	}
	return nil
}

func ValidateSymbol(symbol string) error         { return ValidateReg(symbol, reSymbol) }
func ValidateSymbolReserved(symbol string) error { return ValidateReg(symbol, reSymbolReserved) }
func ValidateSymbolCollectionToken(symbol string) error {
	return ValidateReg(symbol, reSymbolCollectionToken)
}
func ValidateSymbolUserDefined(symbol string) error { return ValidateReg(symbol, reSymbolUserDefined) }
func ValidateSymbolTokenID(symbol string) error     { return ValidateReg(symbol, reSymbolTokenID) }

func SymbolCollectionToken(collection, tokenID string) string {
	return fmt.Sprintf("%s%s", collection, tokenID)
}

func AccAddrSuffix(addr sdk.AccAddress) string {
	bech32Addr := addr.String()
	return bech32Addr[len(bech32Addr)-AccAddrSuffixLen:]
}

package types

import (
	"fmt"
	"regexp"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	reSymbolString            = `[a-z][a-z0-9]{2,15}`
	reSymbolStringReserved    = `[a-z][a-z0-9]{2,4}`
	reSymbolStringUserDefined = `[a-z][a-z0-9]{5,7}`
	/* #nosec */
	reSymbolStringTokenID = `[a-z0-9]{8}`
	/* #nosec */
	reSymbolStringBaseTokenID = `[a-z0-9]{4}`
)

var (
	reSymbol                = regexp.MustCompile(fmt.Sprintf(`^%s$`, reSymbolString))
	reSymbolReserved        = regexp.MustCompile(fmt.Sprintf(`^%s$`, reSymbolStringReserved))
	reSymbolUserDefined     = regexp.MustCompile(fmt.Sprintf(`^%s$`, reSymbolStringUserDefined))
	reSymbolTokenID         = regexp.MustCompile(fmt.Sprintf(`^%s$`, reSymbolStringTokenID))
	reSymbolBaseTokenID     = regexp.MustCompile(fmt.Sprintf(`^%s$`, reSymbolStringBaseTokenID))
	reSymbolCollectionToken = regexp.MustCompile(fmt.Sprintf(`^%s%s$`, reSymbolStringUserDefined, reSymbolStringTokenID))
)

const (
	AccAddrSuffixLen = 3
	TokenIDLen       = 8
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
func ValidateSymbolBaseTokenID(symbol string) error { return ValidateReg(symbol, reSymbolBaseTokenID) }

func SymbolCollectionToken(collection, tokenID string) string {
	return fmt.Sprintf("%s%s", collection, tokenID)
}

func AccAddrSuffix(addr sdk.AccAddress) string {
	bech32Addr := addr.String()
	return bech32Addr[len(bech32Addr)-AccAddrSuffixLen:]
}

func ParseDenom(denom string) (string, string, string) {
	var tokenID string

	if ValidateSymbol(denom) != nil {
		return "", "", ""
	}

	if ValidateSymbolReserved(denom) == nil {
		return denom, "", ""
	}

	if ValidateSymbolCollectionToken(denom) == nil {
		tokenID = denom[len(denom)-8:]
		denom = denom[:len(denom)-8]
	}

	if ValidateSymbolUserDefined(denom) == nil {
		aas := denom[len(denom)-3:]
		ticker := denom[:len(denom)-3]
		return ticker, aas, tokenID
	}
	return "", "", ""
}

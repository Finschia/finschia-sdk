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
	reTokenIDString = `[a-z0-9]{16}`
	/* #nosec */
	reTokenTypeNFTString = `[a-z1-9][a-z0-9]{7}`
	/* #nosec */
	reTokenIndexString = `[a-z0-9]{8}`
)

var (
	reSymbol                = regexp.MustCompile(fmt.Sprintf(`^%s$`, reSymbolString))
	reSymbolReserved        = regexp.MustCompile(fmt.Sprintf(`^%s$`, reSymbolStringReserved))
	reSymbolUserDefined     = regexp.MustCompile(fmt.Sprintf(`^%s$`, reSymbolStringUserDefined))
	reTokenID               = regexp.MustCompile(fmt.Sprintf(`^%s$`, reTokenIDString))
	reTokenTypeNFT          = regexp.MustCompile(fmt.Sprintf(`^%s$`, reTokenTypeNFTString))
	reTokenIndex            = regexp.MustCompile(fmt.Sprintf(`^%s$`, reTokenIndexString))
	reSymbolCollectionToken = regexp.MustCompile(fmt.Sprintf(`^%s%s$`, reSymbolStringUserDefined, reTokenIDString))
)

const (
	AccAddrSuffixLen = 3
	TokenTypeLen     = 8
	TokenIDLen       = 16
)

func ValidateReg(symbol string, reg *regexp.Regexp) error {
	if !reg.MatchString(symbol) {
		return fmt.Errorf("symbol [%s] mismatched to [%s]", symbol, reg.String())
	}
	return nil
}

func ValidateSymbol(symbol string) error            { return ValidateReg(symbol, reSymbol) }
func ValidateSymbolReserved(symbol string) error    { return ValidateReg(symbol, reSymbolReserved) }
func ValidateSymbolUserDefined(symbol string) error { return ValidateReg(symbol, reSymbolUserDefined) }
func ValidateTokenID(symbol string) error           { return ValidateReg(symbol, reTokenID) }
func ValidateTokenTypeNFT(symbol string) error      { return ValidateReg(symbol, reTokenTypeNFT) }
func ValidateTokenIndex(index string) error         { return ValidateReg(index, reTokenIndex) }
func ValidateCollectionToken(symbol string) error   { return ValidateReg(symbol, reSymbolCollectionToken) }
func AccAddrSuffix(addr sdk.AccAddress) string {
	bech32Addr := addr.String()
	return bech32Addr[len(bech32Addr)-AccAddrSuffixLen:]
}

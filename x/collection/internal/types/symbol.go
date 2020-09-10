package types

import (
	"fmt"
	"regexp"
)

const (
	/* #nosec */
	reTokenIDString = `[a-f0-9]{16}`
	/* #nosec */
	reTokenTypeString = `[a-f0-9]{8}`
	/* #nosec */
	reTokenTypeFTString = `0[a-f0-9]{7}`
	/* #nosec */
	reTokenTypeNFTString = `[a-f1-9][a-f0-9]{7}`
	/* #nosec */
	reTokenIndexString = `[a-f0-9]{8}`
)

var (
	reTokenID      = regexp.MustCompile(fmt.Sprintf(`^%s$`, reTokenIDString))
	reTokenType    = regexp.MustCompile(fmt.Sprintf(`^%s$`, reTokenTypeString))
	reTokenTypeFT  = regexp.MustCompile(fmt.Sprintf(`^%s$`, reTokenTypeFTString))
	reTokenTypeNFT = regexp.MustCompile(fmt.Sprintf(`^%s$`, reTokenTypeNFTString))
	reTokenIndex   = regexp.MustCompile(fmt.Sprintf(`^%s$`, reTokenIndexString))
)

func ValidateReg(symbol string, reg *regexp.Regexp) error {
	if !reg.MatchString(symbol) {
		return fmt.Errorf("symbol [%s] mismatched to [%s]", symbol, reg.String())
	}
	return nil
}

func ValidateTokenID(tokenID string) error        { return ValidateReg(tokenID, reTokenID) }
func ValidateTokenType(tokenType string) error    { return ValidateReg(tokenType, reTokenType) }
func ValidateTokenTypeFT(tokenType string) error  { return ValidateReg(tokenType, reTokenTypeFT) }
func ValidateTokenTypeNFT(tokenType string) error { return ValidateReg(tokenType, reTokenTypeNFT) }
func ValidateTokenIndex(index string) error       { return ValidateReg(index, reTokenIndex) }

package types

import (
	"fmt"
	"regexp"
)

const (
	/* #nosec */
	reUserTokenSymbolString = `[A-Z][A-Z0-9]{1,4}`
)

var (
	reUserTokenSymbol = regexp.MustCompile(fmt.Sprintf(`^%s$`, reUserTokenSymbolString))
)

func ValidateReg(symbol string, reg *regexp.Regexp) error {
	if !reg.MatchString(symbol) {
		return fmt.Errorf("symbol [%s] mismatched to [%s]", symbol, reg.String())
	}
	return nil
}

func ValidateTokenSymbol(symbol string) error { return ValidateReg(symbol, reUserTokenSymbol) }

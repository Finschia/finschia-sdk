package types

import (
	"fmt"
	"regexp"
)

const (
	/* #nosec */
	reSymbolStringReserved = `[a-z][a-z0-9]{2,4}`
)

var (
	reSymbolReserved = regexp.MustCompile(fmt.Sprintf(`^%s$`, reSymbolStringReserved))
)

func ValidateReg(symbol string, reg *regexp.Regexp) error {
	if !reg.MatchString(symbol) {
		return fmt.Errorf("symbol [%s] mismatched to [%s]", symbol, reg.String())
	}
	return nil
}

func ValidateSymbolReserved(symbol string) error { return ValidateReg(symbol, reSymbolReserved) }

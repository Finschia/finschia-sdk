package token

import (
	"fmt"
	"regexp"
	"unicode/utf8"

	sdk "github.com/line/lbm-sdk/types"
	sdkerrors "github.com/line/lbm-sdk/types/errors"
)

const (
	maxName     = 20
	maxImageURI = 1000
	maxMeta     = 1000
)

var (
	reSymbolString = `[A-Z][A-Z0-9]{1,4}`
	reSymbol       = regexp.MustCompile(fmt.Sprintf(`^%s$`, reSymbolString))
)

func stringInSize(str string, size int) bool {
	return utf8.RuneCountInString(str) <= size
}

func validateName(name string) error {
	if len(name) == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "Name cannot be empty")
	} else if !stringInSize(name, maxName) {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "Name cannot be longer than %d", maxName)
	}
	return nil
}

func validateSymbol(symbol string) error {
	if !reSymbol.MatchString(symbol) {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "Invalid symbol: %s, valid expression is: %s", symbol, reSymbolString)
	}
	return nil
}

func validateImageURI(uri string) error {
	if !stringInSize(uri, maxImageURI) {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "ImageUri cannot be longer than %d", maxImageURI)
	}
	return nil
}

func validateMeta(meta string) error {
	if !stringInSize(meta, maxMeta) {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "Meta cannot be longer than %d", maxMeta)
	}
	return nil
}

func validateDecimals(decimals int32) error {
	if decimals < 0 || decimals > 18 {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "Invalid decimals: %d", decimals)
	}
	return nil
}

func validateAmount(amount sdk.Int) error {
	if !amount.IsPositive() {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "Amount must be positive: %s", amount)
	}
	return nil
}

func validateAction(action string) error {
	actions := []string{
		ActionMint,
		ActionBurn,
		ActionModify,
	}
	for _, a := range actions {
		if action == a {
			return nil
		}
	}
	return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "Invalid action: %s", action)
}

func validateChange(change Pair) error {
	validators := map[string]func(string) error{
		AttributeKeyName:     validateName,
		AttributeKeyImageURI: validateImageURI,
		AttributeKeyMeta:     validateMeta,
	}
	validator, ok := validators[change.Key]
	if !ok {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "Invalid field: %s", change.Key)
	}
	return validator(change.Value)
}

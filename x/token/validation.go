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
		return sdkerrors.ErrInvalidRequest.Wrap("Name cannot be empty")
	} else if !stringInSize(name, maxName) {
		return sdkerrors.ErrInvalidRequest.Wrapf("Name cannot be longer than %d", maxName)
	}
	return nil
}

func validateSymbol(symbol string) error {
	if !reSymbol.MatchString(symbol) {
		return sdkerrors.ErrInvalidRequest.Wrapf("Invalid symbol: %s, valid expression is: %s", symbol, reSymbolString)
	}
	return nil
}

func validateImageURI(uri string) error {
	if !stringInSize(uri, maxImageURI) {
		return sdkerrors.ErrInvalidRequest.Wrapf("ImageUri cannot be longer than %d", maxImageURI)
	}
	return nil
}

func validateMeta(meta string) error {
	if !stringInSize(meta, maxMeta) {
		return sdkerrors.ErrInvalidRequest.Wrapf("Meta cannot be longer than %d", maxMeta)
	}
	return nil
}

func validateDecimals(decimals int32) error {
	if decimals < 0 || decimals > 18 {
		return sdkerrors.ErrInvalidRequest.Wrapf("Invalid decimals: %d", decimals)
	}
	return nil
}

func validateAmount(amount sdk.Int) error {
	if !amount.IsPositive() {
		return sdkerrors.ErrInvalidRequest.Wrapf("Amount must be positive: %s", amount)
	}
	return nil
}

func validatePermission(permission string) error {
	if value := Permission_value[permission]; value == 0 {
		return sdkerrors.ErrInvalidRequest.Wrapf("Invalid permission: %s", permission)
	}
	return nil
}

func validateChange(change Pair) error {
	validators := map[AttributeKey]func(string) error{
		AttributeKey_Name:     validateName,
		AttributeKey_ImageURI: validateImageURI,
		AttributeKey_Meta:     validateMeta,
	}

	validator, ok := validators[AttributeKey(AttributeKey_value[change.Field])]
	if !ok {
		return sdkerrors.ErrInvalidRequest.Wrapf("Invalid field: %s", change.Field)
	}
	return validator(change.Value)
}

package token

import (
	"fmt"
	"regexp"
	"unicode/utf8"

	sdk "github.com/line/lbm-sdk/types"
	sdkerrors "github.com/line/lbm-sdk/types/errors"
	"github.com/line/lbm-sdk/x/token/class"
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
		return sdkerrors.ErrInvalidRequest.Wrap("name cannot be empty")
	} else if !stringInSize(name, maxName) {
		return sdkerrors.ErrInvalidRequest.Wrapf("name cannot be longer than %d", maxName)
	}
	return nil
}

func validateSymbol(symbol string) error {
	if !reSymbol.MatchString(symbol) {
		return sdkerrors.ErrInvalidRequest.Wrapf("invalid symbol: %s, valid expression is: %s", symbol, reSymbolString)
	}
	return nil
}

func validateImageURI(uri string) error {
	if !stringInSize(uri, maxImageURI) {
		return sdkerrors.ErrInvalidRequest.Wrapf("image_uri cannot be longer than %d", maxImageURI)
	}
	return nil
}

func validateMeta(meta string) error {
	if !stringInSize(meta, maxMeta) {
		return sdkerrors.ErrInvalidRequest.Wrapf("meta cannot be longer than %d", maxMeta)
	}
	return nil
}

func validateDecimals(decimals int32) error {
	if decimals < 0 || decimals > 18 {
		return sdkerrors.ErrInvalidRequest.Wrapf("invalid decimals: %d", decimals)
	}
	return nil
}

func validateAmount(amount sdk.Int) error {
	if !amount.IsPositive() {
		return sdkerrors.ErrInvalidRequest.Wrapf("amount must be positive: %s", amount)
	}
	return nil
}

func validateLegacyPermission(permission string) error {
	return ValidatePermission(Permission(LegacyPermissionFromString(permission)))
}

func ValidatePermission(permission Permission) error {
	if p := Permission_value[Permission_name[int32(permission)]]; p == 0 {
		return sdkerrors.ErrInvalidRequest.Wrapf("invalid permission: %s", permission)
	}
	return nil
}

func validateChange(change Pair) error {
	validators := map[AttributeKey]func(string) error{
		AttributeKeyName:     validateName,
		AttributeKeyImageURI: validateImageURI,
		AttributeKeyMeta:     validateMeta,
	}

	validator, ok := validators[AttributeKeyFromString(change.Field)]
	if !ok {
		return sdkerrors.ErrInvalidRequest.Wrapf("invalid field: %s", change.Field)
	}
	return validator(change.Value)
}

func ValidateContractID(id string) error {
	return class.ValidateID(id)
}

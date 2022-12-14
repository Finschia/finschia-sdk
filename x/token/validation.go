package token

import (
	"fmt"
	"regexp"
	"unicode/utf8"

	sdk "github.com/line/lbm-sdk/types"
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
		return ErrInvalidName.Wrap("cannot be empty")
	} else if !stringInSize(name, maxName) {
		return ErrInvalidName.Wrapf("cannot be longer than %d", maxName)
	}
	return nil
}

func validateSymbol(symbol string) error {
	if !reSymbol.MatchString(symbol) {
		return ErrInvalidSymbol.Wrapf("got; %s, valid expression is; %s", symbol, reSymbolString)
	}
	return nil
}

func validateImageURI(uri string) error {
	if !stringInSize(uri, maxImageURI) {
		return ErrInvalidImageURI.Wrapf("cannot be longer than %d", maxImageURI)
	}
	return nil
}

func validateMeta(meta string) error {
	if !stringInSize(meta, maxMeta) {
		return ErrInvalidMeta.Wrapf("cannot be longer than %d", maxMeta)
	}
	return nil
}

func validateDecimals(decimals int32) error {
	if decimals < 0 || decimals > 18 {
		return ErrInvalidDecimals.Wrapf("must be >=0 and <=18, got; %d", decimals)
	}
	return nil
}

func validateAmount(amount sdk.Int) error {
	if !amount.IsPositive() {
		return ErrInvalidAmount.Wrapf("must be positive: %s", amount)
	}
	return nil
}

func validateLegacyPermission(permission string) error {
	return ValidatePermission(Permission(LegacyPermissionFromString(permission)))
}

func ValidatePermission(permission Permission) error {
	if p := Permission_value[Permission_name[int32(permission)]]; p == 0 {
		return ErrInvalidPermission.Wrap(permission.String())
	}
	return nil
}

func validateChange(change Pair) error {
	validators := map[string]func(string) error{
		AttributeKeyName.String():     validateName,
		AttributeKeyImageURI.String(): validateImageURI,
		AttributeKeyMeta.String():     validateMeta,
	}

	validator, ok := validators[change.Field]
	if !ok {
		return ErrInvalidChanges.Wrapf("invalid key: %s", change.Field)
	}
	return validator(change.Value)
}

func ValidateContractID(id string) error {
	if err := class.ValidateID(id); err != nil {
		return ErrInvalidContractID.Wrap(id)
	}

	return nil
}

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
		return ErrEmpty.Wrap("name")
	} else if !stringInSize(name, maxName) {
		return ErrMaxLimit.Wrapf("length of name must be <=%d", maxName)
	}
	return nil
}

func validateSymbol(symbol string) error {
	if !reSymbol.MatchString(symbol) {
		return ErrInvalid.Wrapf("valid expression for symbol is %s", reSymbolString)
	}
	return nil
}

func validateImageURI(uri string) error {
	if !stringInSize(uri, maxImageURI) {
		return ErrMaxLimit.Wrapf("length of image uri must be <=%d", maxImageURI)
	}
	return nil
}

func validateMeta(meta string) error {
	if !stringInSize(meta, maxMeta) {
		return ErrMaxLimit.Wrapf("length of meta must be <=%d", maxMeta)
	}
	return nil
}

func validateDecimals(decimals int32) error {
	min, max := int32(0), int32(18)
	if decimals < min || decimals > max {
		return ErrInvalid.Wrapf("decimals must be >=%d and <=%d", min, max)
	}
	return nil
}

func validateAmount(amount sdk.Int) error {
	if !amount.IsPositive() {
		return ErrInvalid.Wrapf("amount must be positive; %s", amount)
	}
	return nil
}

func validateLegacyPermission(permission string) error {
	p := Permission(LegacyPermissionFromString(permission))
	if err := ValidatePermission(p); err != nil {
		return ErrInvalid.Wrapf("permission; %s", permission)
	}
	return nil
}

func ValidatePermission(permission Permission) error {
	if p := Permission_value[Permission_name[int32(permission)]]; p == 0 {
		return ErrInvalid.Wrapf("permission; %s", permission)
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
		return ErrInvalid.Wrapf("key; %s", change.Field)
	}
	return validator(change.Value)
}

func ValidateContractID(id string) error {
	if err := class.ValidateID(id); err != nil {
		return ErrInvalid.Wrapf("contract id; %s", id)
	}

	return nil
}

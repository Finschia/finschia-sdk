package token

import (
	"fmt"
	"regexp"
	"unicode/utf8"

	sdk "github.com/Finschia/finschia-sdk/types"
	sdkerrors "github.com/Finschia/finschia-sdk/types/errors"
	"github.com/Finschia/finschia-sdk/x/token/class"
)

const (
	maxName = 20
	maxURI  = 1000
	maxMeta = 1000
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
		return ErrInvalidTokenName.Wrap("name cannot be empty")
	} else if !stringInSize(name, maxName) {
		return ErrInvalidNameLength.Wrapf("name cannot be longer than %d", maxName)
	}
	return nil
}

func validateSymbol(symbol string) error {
	if !reSymbol.MatchString(symbol) {
		return ErrInvalidTokenSymbol.Wrapf("invalid symbol: %s, valid expression is: %s", symbol, reSymbolString)
	}
	return nil
}

func validateURI(uri string) error {
	if !stringInSize(uri, maxURI) {
		return ErrInvalidImageURILength.Wrapf("uri cannot be longer than %d", maxURI)
	}
	return nil
}

func validateMeta(meta string) error {
	if !stringInSize(meta, maxMeta) {
		return ErrInvalidMetaLength.Wrapf("meta cannot be longer than %d", maxMeta)
	}
	return nil
}

func validateDecimals(decimals int32) error {
	if decimals < 0 || decimals > 18 {
		return ErrInvalidTokenDecimals.Wrapf("invalid decimals: %d", decimals)
	}
	return nil
}

func validateAmount(amount sdk.Int) error {
	if !amount.IsPositive() {
		return ErrInvalidAmount.Wrapf("amount must be positive: %s", amount)
	}
	return nil
}

func validateLegacyPermission(permission string) error {
	return ValidatePermission(Permission(LegacyPermissionFromString(permission)))
}

func ValidatePermission(permission Permission) error {
	if p := Permission_value[Permission_name[int32(permission)]]; p == 0 {
		return sdkerrors.ErrInvalidPermission.Wrap(permission.String())
	}
	return nil
}

func validateChange(change Attribute) error {
	validators := map[string]func(string) error{
		AttributeKeyName.String():     validateName,
		AttributeKeyImageURI.String(): validateURI,
		AttributeKeyMeta.String():     validateMeta,
		AttributeKeyURI.String():      validateURI,
	}

	validator, ok := validators[change.Key]
	if !ok {
		return ErrInvalidChangesField.Wrapf("invalid field of key: %s", change.Key)
	}
	return validator(change.Value)
}

func ValidateContractID(id string) error {
	return class.ValidateID(id)
}

func canonicalKey(key string) string {
	convert := map[string]string{
		AttributeKeyImageURI.String(): AttributeKeyURI.String(),
	}
	if converted, ok := convert[key]; ok {
		return converted
	}
	return key
}

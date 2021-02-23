package types

import (
	"unicode/utf8"

	sdkerrors "github.com/line/lbm-sdk/types/errors"
)

const (
	MaxImageURILength    = 1000
	MaxTokenNameLength   = 20
	MaxTokenMetaLength   = 1000
	MaxChangeFieldsCount = 100
)

var (
	TokenModifiableFields = ModifiableFields{
		AttributeKeyName:     true,
		AttributeKeyMeta:     true,
		AttributeKeyImageURI: true,
	}
)

type ModifiableFields map[string]bool

func ValidateName(name string) bool {
	return utf8.RuneCountInString(name) <= MaxTokenNameLength
}

func ValidateMeta(meta string) bool {
	return utf8.RuneCountInString(meta) <= MaxTokenMetaLength
}

func ValidateImageURI(imageURI string) bool {
	return utf8.RuneCountInString(imageURI) <= MaxImageURILength
}

type ChangesValidator struct {
	modifiableFields ModifiableFields
	handlers         map[string]func(value string) error
}

func NewChangesValidator() *ChangesValidator {
	hs := make(map[string]func(value string) error)
	hs[AttributeKeyName] = func(value string) error {
		if !ValidateName(value) {
			return sdkerrors.Wrapf(ErrInvalidNameLength, "[%s] should be shorter than [%d] UTF-8 characters, current length: [%d]", value, MaxTokenNameLength, utf8.RuneCountInString(value))
		}
		return nil
	}
	hs[AttributeKeyImageURI] = func(value string) error {
		if !ValidateImageURI(value) {
			return sdkerrors.Wrapf(ErrInvalidImageURILength, "[%s] should be shorter than [%d] UTF-8 characters, current length: [%d]", value, MaxImageURILength, utf8.RuneCountInString(value))
		}
		return nil
	}
	hs[AttributeKeyMeta] = func(value string) error {
		if !ValidateMeta(value) {
			return sdkerrors.Wrapf(ErrInvalidMetaLength, "[%s] should be shorter than [%d] UTF-8 characters, current length: [%d]", value, MaxTokenMetaLength, utf8.RuneCountInString(value))
		}
		return nil
	}
	return &ChangesValidator{
		modifiableFields: TokenModifiableFields,
		handlers:         hs,
	}
}

func (c *ChangesValidator) Validate(changes Changes) error {
	if len(changes) == 0 {
		return ErrEmptyChanges
	}

	if len(changes) > MaxChangeFieldsCount {
		return sdkerrors.Wrapf(ErrInvalidChangesFieldCount, "You can not change fields more than [%d] at once, current count: [%d]", MaxChangeFieldsCount, len(changes))
	}

	checkedFields := map[string]bool{}
	for _, change := range changes {
		if !c.modifiableFields[change.Field] {
			return sdkerrors.Wrapf(ErrInvalidChangesField, "Field: %s", change.Field)
		}
		if checkedFields[change.Field] {
			return sdkerrors.Wrapf(ErrDuplicateChangesField, "Field: %s", change.Field)
		}

		validateHandler, ok := c.handlers[change.Field]
		if !ok {
			return sdkerrors.Wrapf(ErrInvalidChangesField, "Field: %s", change.Field)
		}

		if err := validateHandler(change.Value); err != nil {
			return err
		}
		checkedFields[change.Field] = true
	}
	return nil
}

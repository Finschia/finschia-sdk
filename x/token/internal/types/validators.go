package types

import (
	"unicode/utf8"

	sdk "github.com/cosmos/cosmos-sdk/types"
	linktype "github.com/line/link/types"
)

const (
	MaxImageURILength    = 1000
	MaxTokenNameLength   = 1000
	MaxTokenSymbolLength = 5
	MaxChangeFieldsCount = 100
)

var (
	TokenModifiableFields = ModifiableFields{
		AttributeKeyName:     true,
		AttributeKeyTokenURI: true,
	}
)

type ModifiableFields map[string]bool

func ValidateName(name string) bool {
	return utf8.RuneCountInString(name) < MaxTokenNameLength
}

func ValidateImageURI(tokenURI string) bool {
	return utf8.RuneCountInString(tokenURI) < MaxImageURILength
}

type ChangesValidator struct {
	modifiableFields ModifiableFields
	handlers         map[string]func(value string) sdk.Error
}

func NewChangesValidator() *ChangesValidator {
	hs := make(map[string]func(value string) sdk.Error)
	hs[AttributeKeyName] = func(value string) sdk.Error {
		if !ValidateName(value) {
			return ErrInvalidNameLength(DefaultCodespace, value)
		}
		return nil
	}
	hs[AttributeKeyTokenURI] = func(value string) sdk.Error {
		if !ValidateImageURI(value) {
			return ErrInvalidImageURILength(DefaultCodespace, value)
		}
		return nil
	}
	return &ChangesValidator{
		modifiableFields: TokenModifiableFields,
		handlers:         hs,
	}
}

func (c *ChangesValidator) Validate(changes linktype.Changes) sdk.Error {
	if len(changes) == 0 {
		return ErrEmptyChanges(DefaultCodespace)
	}

	if len(changes) > MaxChangeFieldsCount {
		return ErrInvalidChangesFieldCount(DefaultCodespace, len(changes))
	}

	checkedFields := map[string]bool{}
	for _, change := range changes {
		if !c.modifiableFields[change.Field] {
			return ErrInvalidChangesField(DefaultCodespace, change.Field)
		}
		if checkedFields[change.Field] {
			return ErrDuplicateChangesField(DefaultCodespace, change.Field)
		}

		validateHandler, ok := c.handlers[change.Field]
		if !ok {
			return ErrInvalidChangesField(DefaultCodespace, change.Field)
		}

		if err := validateHandler(change.Value); err != nil {
			return err
		}
		checkedFields[change.Field] = true
	}
	return nil
}

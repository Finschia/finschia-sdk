package types

import (
	"unicode/utf8"

	sdk "github.com/cosmos/cosmos-sdk/types"
	linktype "github.com/line/link/types"
)

const (
	MaxBaseImgURILength  = 1000
	MaxTokenNameLength   = 1000
	MaxChangeFieldsCount = 100
)

var (
	CollectionModifiableFields = ModifiableFields{
		AttributeKeyName:       true,
		AttributeKeyBaseImgURI: true,
	}
	TokenTypeModifiableFields = ModifiableFields{
		AttributeKeyName: true,
	}
	TokenModifiableFields = ModifiableFields{
		AttributeKeyName: true,
	}
)

type ModifiableFields map[string]bool

func ValidateName(name string) bool {
	return utf8.RuneCountInString(name) < MaxTokenNameLength
}

func ValidateBaseImgURI(baseImgURI string) bool {
	return utf8.RuneCountInString(baseImgURI) < MaxBaseImgURILength
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
	hs[AttributeKeyBaseImgURI] = func(value string) sdk.Error {
		if !ValidateBaseImgURI(value) {
			return ErrInvalidBaseImgURILength(DefaultCodespace, value)
		}
		return nil
	}
	return &ChangesValidator{
		handlers: hs,
	}
}

func (c *ChangesValidator) Validate(changes linktype.Changes) sdk.Error {
	if len(changes) == 0 {
		return ErrEmptyChanges(DefaultCodespace)
	}

	if len(changes) > MaxChangeFieldsCount {
		return ErrInvalidChangesFieldCount(DefaultCodespace, len(changes))
	}

	validator := NewChangesValidator()
	for _, change := range changes {
		if !c.modifiableFields[change.Field] {
			return ErrInvalidChangesField(DefaultCodespace, change.Field)
		}

		validateHandler, ok := validator.handlers[change.Field]
		if !ok {
			return ErrInvalidChangesField(DefaultCodespace, change.Field)
		}

		if err := validateHandler(change.Value); err != nil {
			return err
		}
	}
	return nil
}

func (c *ChangesValidator) SetMode(tokenType, tokenIndex string) sdk.Error {
	if tokenType != "" {
		if tokenIndex == "" {
			c.forTokenType()
		} else {
			c.forToken()
		}
	} else {
		if tokenIndex == "" {
			c.forCollection()
		} else {
			return ErrTokenIndexWithoutType(DefaultCodespace)
		}
	}
	return nil
}

func (c *ChangesValidator) forCollection() {
	c.modifiableFields = CollectionModifiableFields
}

func (c *ChangesValidator) forTokenType() {
	c.modifiableFields = TokenTypeModifiableFields
}

func (c *ChangesValidator) forToken() {
	c.modifiableFields = TokenModifiableFields
}

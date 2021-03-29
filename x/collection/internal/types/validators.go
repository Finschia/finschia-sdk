package types

import (
	"unicode/utf8"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	MaxBaseImgURILength  = 1000
	MaxTokenNameLength   = 20
	MaxChangeFieldsCount = 100
	MaxTokenMetaLength   = 1000
)

var (
	CollectionModifiableFields = ModifiableFields{
		AttributeKeyName:       true,
		AttributeKeyBaseImgURI: true,
		AttributeKeyMeta:       true,
	}
	TokenTypeModifiableFields = ModifiableFields{
		AttributeKeyName: true,
		AttributeKeyMeta: true,
	}
	TokenModifiableFields = ModifiableFields{
		AttributeKeyName: true,
		AttributeKeyMeta: true,
	}
)

type ModifiableFields map[string]bool

func ValidateName(name string) bool {
	return utf8.RuneCountInString(name) <= MaxTokenNameLength
}

func ValidateBaseImgURI(baseImgURI string) bool {
	return utf8.RuneCountInString(baseImgURI) <= MaxBaseImgURILength
}
func ValidateMeta(meta string) bool {
	return utf8.RuneCountInString(meta) <= MaxTokenMetaLength
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
	hs[AttributeKeyBaseImgURI] = func(value string) error {
		if !ValidateBaseImgURI(value) {
			return sdkerrors.Wrapf(ErrInvalidBaseImgURILength, "[%s] should be shorter than [%d] UTF-8 characters, current length: [%d]", value, MaxBaseImgURILength, utf8.RuneCountInString(value))
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
		handlers: hs,
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

func (c *ChangesValidator) SetMode(tokenType, tokenIndex string) error {
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
			return ErrTokenIndexWithoutType
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

package types

import (
	"strings"
	"testing"
	"unicode/utf8"

	sdkerrors "github.com/line/lbm-sdk/types/errors"
	"github.com/stretchr/testify/require"
)

var length1001String = strings.Repeat("Eng글자日本語はスゲ", 91) // 11 * 91 = 1001

func TestValidateName(t *testing.T) {
	t.Log("Given valid name")
	{
		var length20String = strings.Repeat("Eng글자日本語はス", 2) // 10 * 2 = 20
		require.True(t, ValidateName(length20String))
	}
	t.Log("Given invalid name")
	{
		var length21String = strings.Repeat("Eng글자日本", 3) // 7 * 3 = 21
		require.False(t, ValidateName(length21String))
	}
}

func TestValidateBaseImgURI(t *testing.T) {
	t.Log("Given valid base_img_uri")
	{
		var length990String = strings.Repeat("Eng글자日本語はスゲ", 90) // 11 * 90 = 990
		require.True(t, ValidateBaseImgURI(length990String))
	}
	t.Log("Given invalid base_img_uri")
	{
		require.False(t, ValidateBaseImgURI(length1001String))
	}
}

func TestValidateChangesForCollection(t *testing.T) {
	// Given ChangesValidator for collection
	validator := NewChangesValidator()
	err := validator.SetMode("", "")
	require.NoError(t, err)

	t.Log("Test with valid changes")
	{
		changes := NewChangesWithMap(map[string]string{
			"name":         "new_name",
			"base_img_uri": "new_base_uri",
		})

		require.Nil(t, validator.Validate(changes))
	}
	t.Log("Test with empty changes")
	{
		changes := Changes{}
		require.EqualError(t, validator.Validate(changes), ErrEmptyChanges.Error())
	}
	t.Log("Test with base_img_uri too long")
	{
		length1001String := strings.Repeat("Eng글자日本語はスゲ", 91) // 11 * 91 = 1001
		changes := NewChangesWithMap(map[string]string{
			"name":         "new_name",
			"base_img_uri": length1001String,
		})

		require.EqualError(
			t,
			validator.Validate(changes),
			sdkerrors.Wrapf(ErrInvalidBaseImgURILength, "[%s] should be shorter than [%d] UTF-8 characters, current length: [%d]", length1001String, MaxBaseImgURILength, utf8.RuneCountInString(length1001String)).Error(),
		)
	}
	t.Log("Test with invalid changes field")
	{
		// Given changes with invalid fields
		changes := NewChanges(
			NewChange("invalid_field", "value"),
		)

		// Then error is occurred
		require.EqualError(t, validator.Validate(changes), sdkerrors.Wrapf(ErrInvalidChangesField, "Field: invalid_field").Error())
	}
	t.Log("Test with changes more than max")
	{
		// Given changes more than max
		changeList := make([]Change, MaxChangeFieldsCount+1)
		changes := Changes(changeList)

		// Then error is occurred
		require.EqualError(t, validator.Validate(changes), sdkerrors.Wrapf(ErrInvalidChangesFieldCount, "You can not change fields more than [%d] at once, current count: [%d]", MaxChangeFieldsCount, len(changeList)).Error())
	}
	t.Log("Test with duplicate fields")
	{
		// Given changes with duplicate fields
		changes := NewChanges(
			NewChange("name", "value"),
			NewChange("name", "value2"),
		)

		// Then error is occurred
		require.EqualError(t, validator.Validate(changes), sdkerrors.Wrapf(ErrDuplicateChangesField, "Field: name").Error())
	}
}

func TestValidateChangesForTokenType(t *testing.T) {
	// Given ChangesValidator for token type
	validator := NewChangesValidator()
	err := validator.SetMode(defaultTokenType, "")
	require.NoError(t, err)

	t.Log("Test with valid changes")
	{
		changes := NewChangesWithMap(map[string]string{
			"name": "new_name",
		})

		require.Nil(t, validator.Validate(changes))
	}
	t.Log("Test with base_img_uri")
	{
		changes := NewChangesWithMap(map[string]string{
			"name":         "new_name",
			"base_img_uri": "new_base_uri",
		})

		require.EqualError(t, validator.Validate(changes), sdkerrors.Wrap(ErrInvalidChangesField, "Field: base_img_uri").Error())
	}
}

func TestValidateChangesForToken(t *testing.T) {
	// Given ChangesValidator for token
	validator := NewChangesValidator()
	err := validator.SetMode(defaultTokenType, defaultTokenIndex)
	require.NoError(t, err)

	t.Log("Test with valid changes")
	{
		changes := NewChangesWithMap(map[string]string{
			"name": "new_name",
		})

		require.Nil(t, validator.Validate(changes))
	}
	t.Log("Test with base_img_uri")
	{
		changes := NewChangesWithMap(map[string]string{
			"name":         "new_name",
			"base_img_uri": "new_base_uri",
		})

		require.EqualError(t, validator.Validate(changes), sdkerrors.Wrap(ErrInvalidChangesField, "Field: base_img_uri").Error())
	}
}

func TestValidateChangesForTokenWithoutType(t *testing.T) {
	validator := NewChangesValidator()
	require.EqualError(t, validator.SetMode("", defaultTokenIndex), ErrTokenIndexWithoutType.Error())
}

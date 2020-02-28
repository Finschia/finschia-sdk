package types

import (
	"strings"
	"testing"

	linktype "github.com/line/link/types"
	"github.com/stretchr/testify/require"
)

var length990String = strings.Repeat("Eng글자日本語はスゲ", 90)  // 11 * 90 = 990
var length1001String = strings.Repeat("Eng글자日本語はスゲ", 91) // 11 * 91 = 1001

func TestValidateName(t *testing.T) {
	t.Log("Given valid name")
	{
		require.True(t, ValidateName(length990String))
	}
	t.Log("Given invalid name")
	{
		require.False(t, ValidateName(length1001String))
	}
}

func TestValidateBaseImgURI(t *testing.T) {
	t.Log("Given valid base_img_uri")
	{
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
		changes := linktype.NewChangesWithMap(map[string]string{
			"name":         "new_name",
			"base_img_uri": "new_base_uri",
		})

		require.Nil(t, validator.Validate(changes))
	}
	t.Log("Test with empty changes")
	{
		changes := linktype.Changes{}
		require.EqualError(t, validator.Validate(changes), ErrEmptyChanges(DefaultCodespace).Error())
	}
	t.Log("Test with base_img_uri too long")
	{
		length1001String := strings.Repeat("Eng글자日本語はスゲ", 91) // 11 * 91 = 1001
		changes := linktype.NewChangesWithMap(map[string]string{
			"name":         "new_name",
			"base_img_uri": length1001String,
		})

		require.EqualError(
			t,
			validator.Validate(changes),
			ErrInvalidBaseImgURILength(DefaultCodespace, length1001String).Error(),
		)
	}
	t.Log("Test with invalid changes field")
	{
		changes := linktype.NewChanges(
			linktype.NewChange("invalid_field", "value"),
		)
		require.EqualError(t, validator.Validate(changes), ErrInvalidChangesField(DefaultCodespace,
			"invalid_field").Error())
	}
	t.Log("Test with changes more than max")
	{
		// Given changes more than max
		changeList := make([]linktype.Change, MaxChangeFieldsCount+1)
		changes := linktype.Changes(changeList)

		require.EqualError(t, validator.Validate(changes), ErrInvalidChangesFieldCount(DefaultCodespace,
			len(changeList)).Error())
	}
}

func TestValidateChangesForTokenType(t *testing.T) {
	// Given ChangesValidator for token type
	validator := NewChangesValidator()
	err := validator.SetMode(defaultTokenType, "")
	require.NoError(t, err)

	t.Log("Test with valid changes")
	{
		changes := linktype.NewChangesWithMap(map[string]string{
			"name": "new_name",
		})

		require.Nil(t, validator.Validate(changes))
	}
	t.Log("Test with base_img_uri")
	{
		changes := linktype.NewChangesWithMap(map[string]string{
			"name":         "new_name",
			"base_img_uri": "new_base_uri",
		})

		require.EqualError(t, validator.Validate(changes), ErrInvalidChangesField(DefaultCodespace, "base_img_uri").Error())
	}
}

func TestValidateChangesForToken(t *testing.T) {
	// Given ChangesValidator for token
	validator := NewChangesValidator()
	err := validator.SetMode(defaultTokenType, defaultTokenIndex)
	require.NoError(t, err)

	t.Log("Test with valid changes")
	{
		changes := linktype.NewChangesWithMap(map[string]string{
			"name": "new_name",
		})

		require.Nil(t, validator.Validate(changes))
	}
	t.Log("Test with base_img_uri")
	{
		changes := linktype.NewChangesWithMap(map[string]string{
			"name":         "new_name",
			"base_img_uri": "new_base_uri",
		})

		require.EqualError(t, validator.Validate(changes), ErrInvalidChangesField(DefaultCodespace, "base_img_uri").Error())
	}
}

func TestValidateChangesForTokenWithoutType(t *testing.T) {
	validator := NewChangesValidator()
	require.EqualError(t, validator.SetMode("", defaultTokenIndex), ErrTokenIndexWithoutType(DefaultCodespace).Error())
}

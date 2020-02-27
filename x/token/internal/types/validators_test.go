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

func TestValidateTokenURI(t *testing.T) {
	t.Log("Given valid token_uri")
	{
		require.True(t, ValidateTokenURI(length990String))
	}
	t.Log("Given invalid token_uri")
	{
		require.False(t, ValidateTokenURI(length1001String))
	}
}

func TestValidateChanges(t *testing.T) {
	validator := NewChangesValidator()
	t.Log("Test with valid changes")
	{
		changes := linktype.NewChangesWithMap(map[string]string{
			"name":      "new_name",
			"token_uri": "new_torken_uri",
		})

		require.Nil(t, validator.Validate(changes))
	}
	t.Log("Test with empty changes")
	{
		changes := linktype.Changes{}
		require.EqualError(t, validator.Validate(changes), ErrEmptyChanges(DefaultCodespace).Error())
	}
	t.Log("Test with token_uri too long")
	{
		length1001String := strings.Repeat("Eng글자日本語はスゲ", 91) // 11 * 91 = 1001
		changes := linktype.NewChangesWithMap(map[string]string{
			"name":      "new_name",
			"token_uri": length1001String,
		})

		require.EqualError(
			t,
			validator.Validate(changes),
			ErrInvalidTokenURILength(DefaultCodespace, length1001String).Error(),
		)
	}
	t.Log("Test with invalid changes field")
	{
		changes := linktype.NewChanges(
			linktype.NewChange("invalid_field", "value"),
		)
		require.EqualError(t, validator.Validate(changes), ErrInvalidChangesField(DefaultCodespace, "invalid_field").Error())
	}
	t.Log("Test with changes more than max")
	{
		// Given changes more than max
		changeList := make([]linktype.Change, MaxChangeFieldsCount+1)
		changes := linktype.Changes(changeList)

		require.EqualError(t, validator.Validate(changes), ErrInvalidChangesFieldCount(DefaultCodespace, len(changeList)).Error())
	}
}

package types

import (
	"strings"
	"testing"

	linktype "github.com/line/link/types"
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

func TestValidateTokenURI(t *testing.T) {
	t.Log("Given valid base_img_uri")
	{
		var length990String = strings.Repeat("Eng글자日本語はスゲ", 90) // 11 * 90 = 990
		require.True(t, ValidateImageURI(length990String))
	}
	t.Log("Given invalid token_uri")
	{
		require.False(t, ValidateImageURI(length1001String))
	}
}

func TestValidateChanges(t *testing.T) {
	validator := NewChangesValidator()
	t.Log("Test with valid changes")
	{
		changes := linktype.NewChangesWithMap(map[string]string{
			"name":    "new_name",
			"img_uri": "new_img_uri",
		})

		require.Nil(t, validator.Validate(changes))
	}
	t.Log("Test with empty changes")
	{
		changes := linktype.Changes{}
		require.EqualError(t, validator.Validate(changes), ErrEmptyChanges(DefaultCodespace).Error())
	}
	t.Log("Test with img_uri too long")
	{
		length1001String := strings.Repeat("Eng글자日本語はスゲ", 91) // 11 * 91 = 1001
		changes := linktype.NewChangesWithMap(map[string]string{
			"name":    "new_name",
			"img_uri": length1001String,
		})

		require.EqualError(
			t,
			validator.Validate(changes),
			ErrInvalidImageURILength(DefaultCodespace, length1001String).Error(),
		)
	}
	t.Log("Test with invalid changes field")
	{
		// Given changes with invalid fields
		changes := linktype.NewChanges(
			linktype.NewChange("invalid_field", "value"),
		)

		// Then error is occurred
		require.EqualError(t, validator.Validate(changes), ErrInvalidChangesField(DefaultCodespace, "invalid_field").Error())
	}
	t.Log("Test with changes more than max")
	{
		// Given changes more than max
		changeList := make([]linktype.Change, MaxChangeFieldsCount+1)
		changes := linktype.Changes(changeList)

		// Then error is occurred
		require.EqualError(t, validator.Validate(changes), ErrInvalidChangesFieldCount(DefaultCodespace, len(changeList)).Error())
	}
	t.Log("Test with duplicate fields")
	{
		// Given changes with duplicate fields
		changes := linktype.NewChanges(
			linktype.NewChange("name", "value"),
			linktype.NewChange("name", "value2"),
		)

		// Then error is occurred
		require.EqualError(t, validator.Validate(changes), ErrDuplicateChangesField(DefaultCodespace, "name").Error())
	}
}

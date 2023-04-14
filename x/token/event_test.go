package token_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/line/lbm-sdk/x/token"
)

func TestAttributeKeyStringer(t *testing.T) {
	for _, name := range token.AttributeKey_name {
		t.Run(name, func(t *testing.T) {
			value := token.AttributeKey(token.AttributeKey_value[name])
			customName := value.String()
			require.EqualValues(t, value, token.AttributeKeyFromString(customName))
		})
	}
}

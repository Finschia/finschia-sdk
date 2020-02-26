package types

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTokenType(t *testing.T) {
	tokenType := NewBaseTokenType(defaultSymbol, defaultTokenType, defaultName)

	require.Equal(t, `{"symbol":"token001","token_type":"10000001","name":"name"}`, tokenType.String())

	var tokenType2 TokenType
	bz, err := ModuleCdc.MarshalJSON(tokenType)
	require.NoError(t, err)
	err = ModuleCdc.UnmarshalJSON(bz, &tokenType2)
	require.NoError(t, err)

	require.Equal(t, defaultName, tokenType2.GetName())
	require.Equal(t, defaultSymbol, tokenType2.GetSymbol())
	require.Equal(t, defaultTokenType, tokenType2.GetTokenType())

	require.Equal(t, tokenType.GetName(), tokenType2.GetName())
	require.Equal(t, tokenType.GetSymbol(), tokenType2.GetSymbol())
	require.Equal(t, tokenType.GetTokenType(), tokenType2.GetTokenType())

	tokenType3 := tokenType.SetName("testname")
	require.Equal(t, defaultName, tokenType.GetName())
	require.Equal(t, "testname", tokenType3.GetName())
}

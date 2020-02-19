package types

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestToken(t *testing.T) {
	// NewXX with arguments.
	// Serialize it and Deserialize to other variable.
	// Compare both are the same.
	{
		token := NewToken(defaultName, defaultSymbol, defaultTokenURI, sdk.NewInt(defaultDecimals), true)

		var token2 Token
		bz, err := ModuleCdc.MarshalJSON(token)
		require.NoError(t, err)
		err = ModuleCdc.UnmarshalJSON(bz, &token2)
		require.NoError(t, err)

		require.Equal(t, defaultName, token2.GetName())
		require.Equal(t, defaultSymbol, token2.GetSymbol())
		require.Equal(t, defaultSymbol, token2.GetSymbol())
		require.Equal(t, defaultTokenURI, token2.GetTokenURI())
		require.Equal(t, int64(defaultDecimals), token2.GetDecimals().Int64())
		require.Equal(t, true, token2.GetMintable())

		require.Equal(t, token.GetName(), token2.GetName())
		require.Equal(t, token.GetSymbol(), token2.GetSymbol())
		require.Equal(t, token.GetSymbol(), token2.GetSymbol())
		require.Equal(t, token.GetTokenURI(), token2.GetTokenURI())
		require.Equal(t, token.GetDecimals().Int64(), token2.GetDecimals().Int64())
		require.Equal(t, token.GetMintable(), token2.GetMintable())
	}
	{
		tokens := Tokens{
			NewToken(defaultName, defaultSymbol+"1", defaultTokenURI, sdk.NewInt(defaultDecimals), true),
			NewToken(defaultName, defaultSymbol+"2", defaultTokenURI, sdk.NewInt(defaultDecimals), true),
		}
		require.Equal(t, `[{"name":"name","symbol":"linktkn1","token_uri":"token-uri","decimals":"6","mintable":true},{"name":"name","symbol":"linktkn2","token_uri":"token-uri","decimals":"6","mintable":true}]`, tokens.String())
	}
}

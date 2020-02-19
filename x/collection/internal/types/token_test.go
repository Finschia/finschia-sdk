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
	collection := NewCollection(defaultSymbol, defaultName)
	{
		token := NewFT(collection, defaultName, defaultTokenURI, sdk.NewInt(defaultDecimals), true)

		var token2 BaseFT
		bz, err := ModuleCdc.MarshalJSON(token)
		require.NoError(t, err)
		err = ModuleCdc.UnmarshalJSON(bz, &token2)
		require.NoError(t, err)

		require.Equal(t, defaultName, token2.GetName())
		require.Equal(t, defaultSymbol, token2.GetSymbol())
		require.Equal(t, defaultSymbol+token.GetTokenID(), token2.GetDenom())
		require.Equal(t, defaultTokenURI, token2.GetTokenURI())
		require.Equal(t, int64(defaultDecimals), token2.GetDecimals().Int64())
		require.Equal(t, true, token2.GetMintable())

		require.Equal(t, token.GetName(), token2.GetName())
		require.Equal(t, token.GetSymbol(), token2.GetSymbol())
		require.Equal(t, token.GetDenom(), token2.GetDenom())
		require.Equal(t, token.GetTokenURI(), token2.GetTokenURI())
		require.Equal(t, token.GetDecimals().Int64(), token2.GetDecimals().Int64())
		require.Equal(t, token.GetMintable(), token2.GetMintable())
	}
	{
		token := NewNFT(collection, defaultName, defaultTokenType, defaultTokenURI, defaultAddr)

		var token2 BaseNFT
		bz, err := ModuleCdc.MarshalJSON(token)
		require.NoError(t, err)
		err = ModuleCdc.UnmarshalJSON(bz, &token2)
		require.NoError(t, err)

		require.Equal(t, defaultName, token2.GetName())
		require.Equal(t, defaultSymbol, token2.GetSymbol())
		require.Equal(t, defaultSymbol+defaultTokenType+"0001", token2.GetDenom())
		require.Equal(t, defaultTokenURI, token2.GetTokenURI())
		require.Equal(t, defaultAddr, token2.GetOwner())

		require.Equal(t, token.GetName(), token2.GetName())
		require.Equal(t, token.GetSymbol(), token2.GetSymbol())
		require.Equal(t, token.GetDenom(), token2.GetDenom())
		require.Equal(t, token.GetTokenURI(), token2.GetTokenURI())
		require.Equal(t, token.GetOwner(), token2.GetOwner())
	}
}

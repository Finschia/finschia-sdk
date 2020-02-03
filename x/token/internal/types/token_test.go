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
		token := NewFT(defaultName, defaultSymbol, defaultTokenURI, sdk.NewInt(defaultDecimals), true)

		var token2 FT
		bz, err := ModuleCdc.MarshalJSON(token)
		require.NoError(t, err)
		err = ModuleCdc.UnmarshalJSON(bz, &token2)
		require.NoError(t, err)

		require.Equal(t, defaultName, token2.GetName())
		require.Equal(t, defaultSymbol, token2.GetSymbol())
		require.Equal(t, defaultSymbol, token2.GetDenom())
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
		token := NewNFT(defaultName, defaultSymbol, defaultTokenURI, defaultAddr)

		var token2 NFT
		bz, err := ModuleCdc.MarshalJSON(token)
		require.NoError(t, err)
		err = ModuleCdc.UnmarshalJSON(bz, &token2)
		require.NoError(t, err)

		require.Equal(t, defaultName, token2.GetName())
		require.Equal(t, defaultSymbol, token2.GetSymbol())
		require.Equal(t, defaultSymbol, token2.GetDenom())
		require.Equal(t, defaultTokenURI, token2.GetTokenURI())
		require.Equal(t, defaultAddr, token2.GetOwner())

		require.Equal(t, token.GetName(), token2.GetName())
		require.Equal(t, token.GetSymbol(), token2.GetSymbol())
		require.Equal(t, token.GetDenom(), token2.GetDenom())
		require.Equal(t, token.GetTokenURI(), token2.GetTokenURI())
		require.Equal(t, token.GetOwner(), token2.GetOwner())
	}
	collection := NewCollection(defaultSymbol, defaultName)
	{
		token := NewCollectiveFT(collection, defaultName, defaultTokenID, defaultTokenURI, sdk.NewInt(defaultDecimals), true)

		var token2 BaseCollectiveFT
		bz, err := ModuleCdc.MarshalJSON(token)
		require.NoError(t, err)
		err = ModuleCdc.UnmarshalJSON(bz, &token2)
		require.NoError(t, err)

		require.Equal(t, defaultName, token2.GetName())
		require.Equal(t, defaultSymbol, token2.GetSymbol())
		require.Equal(t, defaultSymbol+defaultTokenID, token2.GetDenom())
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
		token := NewCollectiveNFT(collection, defaultName, defaultTokenID, defaultTokenURI, defaultAddr)

		var token2 BaseCollectiveNFT
		bz, err := ModuleCdc.MarshalJSON(token)
		require.NoError(t, err)
		err = ModuleCdc.UnmarshalJSON(bz, &token2)
		require.NoError(t, err)

		require.Equal(t, defaultName, token2.GetName())
		require.Equal(t, defaultSymbol, token2.GetSymbol())
		require.Equal(t, defaultSymbol+defaultTokenID, token2.GetDenom())
		require.Equal(t, defaultTokenURI, token2.GetTokenURI())
		require.Equal(t, defaultAddr, token2.GetOwner())

		require.Equal(t, token.GetName(), token2.GetName())
		require.Equal(t, token.GetSymbol(), token2.GetSymbol())
		require.Equal(t, token.GetDenom(), token2.GetDenom())
		require.Equal(t, token.GetTokenURI(), token2.GetTokenURI())
		require.Equal(t, token.GetOwner(), token2.GetOwner())
	}
}

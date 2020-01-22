package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	"testing"
)

const (
	defaultName     = "name"
	defaultSymbol   = "linktkn"
	defaultTokenURI = "token-uri"
	defaultDecimals = 6
	defaultTokenID  = "00000001"
)

var (
	defaultAddr = sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
)

func TestToken(t *testing.T) {
	// NewXX with arguments.
	// Serialize it and Deserialize to other variable.
	// Compare both are the same.
	{
		token := NewFT(defaultName, defaultSymbol, defaultTokenURI, sdk.NewInt(defaultDecimals), true)
		require.Equal(t, defaultName, token.GetName())
		require.Equal(t, defaultSymbol, token.GetSymbol())
		require.Equal(t, defaultSymbol, token.GetDenom())
		require.Equal(t, defaultTokenURI, token.GetTokenURI())
		require.Equal(t, int64(defaultDecimals), token.GetDecimals().Int64())
		require.Equal(t, true, token.GetMintable())

		var token2 FT
		bz, err := ModuleCdc.MarshalJSON(token)
		require.NoError(t, err)
		err = ModuleCdc.UnmarshalJSON(bz, &token2)
		require.NoError(t, err)

		require.Equal(t, token.GetName(), token2.GetName())
		require.Equal(t, token.GetSymbol(), token2.GetSymbol())
		require.Equal(t, token.GetDenom(), token2.GetDenom())
		require.Equal(t, token.GetTokenURI(), token2.GetTokenURI())
		require.Equal(t, token.GetDecimals().Int64(), token2.GetDecimals().Int64())
		require.Equal(t, token.GetMintable(), token2.GetMintable())

	}
	{
		token := NewNFT(defaultName, defaultSymbol, defaultTokenURI, defaultAddr)
		require.Equal(t, defaultName, token.GetName())
		require.Equal(t, defaultSymbol, token.GetSymbol())
		require.Equal(t, defaultSymbol, token.GetDenom())
		require.Equal(t, defaultTokenURI, token.GetTokenURI())
		require.Equal(t, defaultAddr, token.GetOwner())

		var token2 NFT
		bz, err := ModuleCdc.MarshalJSON(token)
		require.NoError(t, err)
		err = ModuleCdc.UnmarshalJSON(bz, &token2)
		require.NoError(t, err)

		require.Equal(t, token.GetName(), token2.GetName())
		require.Equal(t, token.GetSymbol(), token2.GetSymbol())
		require.Equal(t, token.GetDenom(), token2.GetDenom())
		require.Equal(t, token.GetTokenURI(), token2.GetTokenURI())
		require.Equal(t, token.GetOwner(), token2.GetOwner())
	}
	{
		token := NewIDFT(defaultName, defaultSymbol, defaultTokenURI, sdk.NewInt(defaultDecimals), true, defaultTokenID)
		require.Equal(t, defaultName, token.GetName())
		require.Equal(t, defaultSymbol, token.GetSymbol())
		require.Equal(t, defaultSymbol+defaultTokenID, token.GetDenom())
		require.Equal(t, int64(defaultDecimals), token.GetDecimals().Int64())
		require.Equal(t, true, token.GetMintable())
		require.Equal(t, defaultTokenID, token.GetTokenID())
		require.Equal(t, defaultTokenURI, token.GetTokenURI())

		var token2 FT
		bz, err := ModuleCdc.MarshalJSON(token)
		require.NoError(t, err)
		err = ModuleCdc.UnmarshalJSON(bz, &token2)
		require.NoError(t, err)

		require.Equal(t, token.GetName(), token2.GetName())
		require.Equal(t, token.GetSymbol(), token2.GetSymbol())
		require.Equal(t, token.GetDenom(), token2.GetDenom())
		require.Equal(t, token.GetTokenURI(), token2.GetTokenURI())
		require.Equal(t, token.GetDecimals().Int64(), token2.GetDecimals().Int64())
		require.Equal(t, token.GetMintable(), token2.GetMintable())
		require.Equal(t, token.GetTokenID(), token2.GetTokenID())
	}
	{

		token := NewIDNFT(defaultName, defaultSymbol, defaultTokenURI, defaultAddr, defaultTokenID)
		require.Equal(t, defaultName, token.GetName())
		require.Equal(t, defaultSymbol, token.GetSymbol())
		require.Equal(t, defaultSymbol+defaultTokenID, token.GetDenom())
		require.Equal(t, defaultTokenURI, token.GetTokenURI())
		require.Equal(t, defaultAddr, token.GetOwner())
		require.Equal(t, defaultTokenID, token.GetTokenID())

		var token2 NFT
		bz, err := ModuleCdc.MarshalJSON(token)
		require.NoError(t, err)
		err = ModuleCdc.UnmarshalJSON(bz, &token2)
		require.NoError(t, err)

		require.Equal(t, token.GetName(), token2.GetName())
		require.Equal(t, token.GetSymbol(), token2.GetSymbol())
		require.Equal(t, token.GetDenom(), token2.GetDenom())
		require.Equal(t, token.GetTokenURI(), token2.GetTokenURI())
		require.Equal(t, token.GetOwner(), token2.GetOwner())
		require.Equal(t, token.GetTokenID(), token2.GetTokenID())
	}
}

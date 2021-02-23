package keeper

import (
	"testing"

	"github.com/line/lbm-sdk/x/collection/internal/types"
	"github.com/stretchr/testify/require"
)

func TestKeeper_GetTokenType(t *testing.T) {
	ctx := cacheKeeper()
	t.Log("Prepare Token Type")
	expected := types.NewBaseTokenType(defaultContractID, defaultTokenType, defaultName, defaultMeta)
	{
		store := ctx.KVStore(keeper.storeKey)
		store.Set(types.TokenTypeKey(defaultContractID, defaultTokenType), keeper.cdc.MustMarshalBinaryBare(expected))
	}
	t.Log("Get Token Type")
	{
		actual, err := keeper.GetTokenType(ctx, defaultTokenType)
		require.NoError(t, err)
		verifyTokenTypeFunc(t, expected, actual)
	}
}

func TestKeeper_SetTokenType(t *testing.T) {
	ctx := cacheKeeper()
	t.Log("Prepare collection")
	require.NoError(t, keeper.CreateCollection(ctx, types.NewCollection(defaultContractID, defaultName, defaultMeta, defaultImgURI), addr1))
	t.Log("Set Token Type")
	expected := types.NewBaseTokenType(defaultContractID, defaultTokenType, defaultName, defaultMeta)
	{
		require.NoError(t, keeper.SetTokenType(ctx, expected))
	}
	t.Log("Compare Token Type")
	{
		store := ctx.KVStore(keeper.storeKey)
		bz := store.Get(types.TokenTypeKey(defaultContractID, defaultTokenType))
		actual := keeper.mustDecodeTokenType(bz)
		verifyTokenTypeFunc(t, expected, actual)
	}
}

func TestKeeper_HasTokenType(t *testing.T) {
	ctx := cacheKeeper()
	t.Log("Prepare Token Type")
	expected := types.NewBaseTokenType(defaultContractID, defaultTokenType, defaultName, defaultMeta)
	{
		store := ctx.KVStore(keeper.storeKey)
		store.Set(types.TokenTypeKey(defaultContractID, defaultTokenType), keeper.cdc.MustMarshalBinaryBare(expected))
	}
	t.Log("Get Token Type")
	{
		require.True(t, keeper.HasTokenType(ctx, defaultTokenType))
	}
}

func TestKeeper_UpdateTokenType(t *testing.T) {
	ctx := cacheKeeper()
	t.Log("Prepare collection")
	require.NoError(t, keeper.CreateCollection(ctx, types.NewCollection(defaultContractID, defaultName, defaultMeta, defaultImgURI), addr1))
	t.Log("Set Token Type")
	expected := types.NewBaseTokenType(defaultContractID, defaultTokenType, defaultName, defaultMeta)
	{
		require.NoError(t, keeper.SetTokenType(ctx, expected))
	}
	t.Log("Update Token Type")
	{
		expected.SetName("modifiedname")
		require.NoError(t, keeper.UpdateTokenType(ctx, expected))
	}

	t.Log("Get Token Type")
	{
		actual, err := keeper.GetTokenType(ctx, defaultTokenType)
		require.NoError(t, err)
		verifyTokenTypeFunc(t, expected, actual)
	}
}

func TestKeeper_GetNextTokenType(t *testing.T) {
	ctx := cacheKeeper()
	t.Log("Prepare collection")
	require.NoError(t, keeper.CreateCollection(ctx, types.NewCollection(defaultContractID, defaultName, defaultMeta, defaultImgURI), addr1))
	t.Log("Get Next Token Type")
	{
		tokenType, err := keeper.GetNextTokenType(ctx)
		require.NoError(t, err)
		require.Equal(t, defaultTokenType, tokenType)
	}
	t.Log("Set Token Type")
	{
		require.NoError(t, keeper.SetTokenType(ctx, types.NewBaseTokenType(defaultContractID, defaultTokenType, defaultName, defaultMeta)))
		require.NoError(t, keeper.SetTokenType(ctx, types.NewBaseTokenType(defaultContractID, defaultTokenType2, defaultName, defaultMeta)))
		require.NoError(t, keeper.SetTokenType(ctx, types.NewBaseTokenType(defaultContractID, defaultTokenType3, defaultName, defaultMeta)))
	}
	t.Log("Get TokenTypes")
	{
		tokenTypes, err := keeper.GetTokenTypes(ctx)
		require.NoError(t, err)
		require.Equal(t, tokenTypes[0].GetTokenType(), defaultTokenType)
		require.Equal(t, tokenTypes[1].GetTokenType(), defaultTokenType2)
		require.Equal(t, tokenTypes[2].GetTokenType(), defaultTokenType3)
	}
	t.Log("Get Next Token Type")
	{
		tokenType, err := keeper.GetNextTokenType(ctx)
		require.NoError(t, err)
		require.Equal(t, defaultTokenType4, tokenType)
	}
	t.Log("Set Full")
	{
		keeper.setNextTokenTypeNFT(ctx, "ffffffff")
		_, err := keeper.getNextTokenTypeNFT(ctx)
		require.Error(t, err)
	}
}

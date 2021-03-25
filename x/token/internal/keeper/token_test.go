package keeper

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/line/link-modules/x/token/internal/types"
	"github.com/stretchr/testify/require"
)

func TestKeeper_GetToken(t *testing.T) {
	ctx := cacheKeeper()
	t.Log("Prepare Token")
	expected := types.NewToken(defaultContractID, defaultName, defaultSymbol, defaultMeta, defaultImageURI, sdk.NewInt(defaultDecimals), true)
	{
		store := ctx.KVStore(keeper.storeKey)
		store.Set(types.TokenKey(expected.GetContractID()), keeper.cdc.MustMarshalBinaryBare(expected))
	}
	t.Log("Get Token")
	{
		actual, err := keeper.GetToken(ctx)
		require.NoError(t, err)
		verifyTokenFunc(t, expected, actual)
	}
}

func TestKeeper_SetToken(t *testing.T) {
	ctx := cacheKeeper()
	t.Log("Set Token")
	expected := types.NewToken(defaultContractID, defaultName, defaultSymbol, defaultMeta, defaultImageURI, sdk.NewInt(defaultDecimals), true)
	{
		require.NoError(t, keeper.SetToken(ctx, expected))
	}
	t.Log("Compare Token")
	{
		store := ctx.KVStore(keeper.storeKey)
		bz := store.Get(types.TokenKey(expected.GetContractID()))
		actual := keeper.mustDecodeToken(bz)
		verifyTokenFunc(t, expected, actual)
	}
}

func TestKeeper_UpdateToken(t *testing.T) {
	ctx := cacheKeeper()
	t.Log("Set Token")
	token := types.NewToken(defaultContractID, defaultName, defaultSymbol, defaultMeta, defaultImageURI, sdk.NewInt(defaultDecimals), true)
	{
		require.NoError(t, keeper.SetToken(ctx, token))
	}
	t.Log("Update Token")
	expected := types.NewToken(token.GetContractID(), "modifiedname", "BTC", "{}", "modifiedtokenuri", sdk.NewInt(defaultDecimals), true)
	{
		require.NoError(t, keeper.UpdateToken(ctx, expected))
	}
	t.Log("Compare Token")
	{
		store := ctx.KVStore(keeper.storeKey)
		bz := store.Get(types.TokenKey(token.GetContractID()))
		actual := keeper.mustDecodeToken(bz)
		verifyTokenFunc(t, expected, actual)
	}
}

func TestKeeper_GetAllTokens(t *testing.T) {
	ctx := cacheKeeper()
	t.Log("Prepare Tokens")
	expected := types.Tokens{
		types.NewToken(defaultContractID+"1", defaultName, defaultSymbol, defaultMeta, defaultImageURI, sdk.NewInt(defaultDecimals), true),
		types.NewToken(defaultContractID+"2", defaultName, defaultSymbol, defaultMeta, defaultImageURI, sdk.NewInt(defaultDecimals), true),
		types.NewToken(defaultContractID+"3", defaultName, defaultSymbol, defaultMeta, defaultImageURI, sdk.NewInt(defaultDecimals), true),
		types.NewToken(defaultContractID+"4", defaultName, defaultSymbol, defaultMeta, defaultImageURI, sdk.NewInt(defaultDecimals), true),
	}
	{
		store := ctx.KVStore(keeper.storeKey)
		for _, t := range expected {
			store.Set(types.TokenKey(t.GetContractID()), keeper.cdc.MustMarshalBinaryBare(t))
		}
	}
	t.Log("Compare Tokens")
	{
		actual := keeper.GetAllTokens(ctx)
		for index := range expected {
			verifyTokenFunc(t, expected[index], actual[index])
		}
	}
}

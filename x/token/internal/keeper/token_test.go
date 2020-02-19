package keeper

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/line/link/x/token/internal/types"
	"github.com/stretchr/testify/require"
)

func TestKeeper_GetToken(t *testing.T) {
	ctx := cacheKeeper()
	t.Log("Prepare Token")
	expected := types.NewToken(defaultName, defaultSymbol, defaultTokenURI, sdk.NewInt(defaultDecimals), true)
	{
		store := ctx.KVStore(keeper.storeKey)
		store.Set(types.TokenSymbolKey(expected.GetSymbol()), keeper.cdc.MustMarshalBinaryBare(expected))
	}
	t.Log("Get Token")
	{
		actual, err := keeper.GetToken(ctx, defaultSymbol)
		require.NoError(t, err)
		verifyTokenFunc(t, expected, actual)
	}
}

func TestKeeper_SetToken(t *testing.T) {
	ctx := cacheKeeper()
	t.Log("Set Token")
	expected := types.NewToken(defaultName, defaultSymbol, defaultTokenURI, sdk.NewInt(defaultDecimals), true)
	{
		require.NoError(t, keeper.SetToken(ctx, expected))
	}
	t.Log("Compare Token")
	{
		store := ctx.KVStore(keeper.storeKey)
		bz := store.Get(types.TokenSymbolKey(expected.GetSymbol()))
		actual := keeper.mustDecodeToken(bz)
		verifyTokenFunc(t, expected, actual)
	}
}

func TestKeeper_UpdateToken(t *testing.T) {
	ctx := cacheKeeper()
	t.Log("Set Token")
	token := types.NewToken(defaultName, defaultSymbol, defaultTokenURI, sdk.NewInt(defaultDecimals), true)
	{
		require.NoError(t, keeper.SetToken(ctx, token))
	}
	t.Log("Update Token")
	expected := types.NewToken("modifiedname", token.GetSymbol(), "modifiedtokenuri", sdk.NewInt(defaultDecimals), true)
	{
		require.NoError(t, keeper.UpdateToken(ctx, expected))
	}
	t.Log("Compare Token")
	{
		store := ctx.KVStore(keeper.storeKey)
		bz := store.Get(types.TokenSymbolKey(token.GetSymbol()))
		actual := keeper.mustDecodeToken(bz)
		verifyTokenFunc(t, expected, actual)
	}
}

func TestKeeper_GetSupply(t *testing.T) {
	ctx := cacheKeeper()
	t.Log("No Token. Get Supply")
	token := types.NewToken(defaultName, defaultSymbol, defaultTokenURI, sdk.NewInt(defaultDecimals), true)
	{
		supply, err := keeper.GetSupply(ctx, token.GetSymbol())
		require.Error(t, err)
		require.Equal(t, int64(0), supply.Int64())
	}
	t.Log("Set Token")
	{
		require.NoError(t, keeper.SetToken(ctx, token))
		require.NoError(t, keeper.mintTokens(ctx, sdk.NewCoins(sdk.NewCoin(token.GetSymbol(), sdk.NewInt(defaultAmount))), addr1))
	}
	t.Log("Token Exist. Get Supply")
	{
		supply, err := keeper.GetSupply(ctx, token.GetSymbol())
		require.NoError(t, err)
		require.Equal(t, int64(defaultAmount), supply.Int64())
	}
}

func TestKeeper_GetAllTokens(t *testing.T) {
	ctx := cacheKeeper()
	t.Log("Prepare Tokens")
	expected := types.Tokens{
		types.NewToken(defaultName, defaultSymbol+"1", defaultTokenURI, sdk.NewInt(defaultDecimals), true),
		types.NewToken(defaultName, defaultSymbol+"2", defaultTokenURI, sdk.NewInt(defaultDecimals), true),
		types.NewToken(defaultName, defaultSymbol+"3", defaultTokenURI, sdk.NewInt(defaultDecimals), true),
		types.NewToken(defaultName, defaultSymbol+"4", defaultTokenURI, sdk.NewInt(defaultDecimals), true),
	}
	{
		store := ctx.KVStore(keeper.storeKey)
		for _, t := range expected {
			store.Set(types.TokenSymbolKey(t.GetSymbol()), keeper.cdc.MustMarshalBinaryBare(t))
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

func TestKeeper_ModifyTokenURI(t *testing.T) {
	ctx := cacheKeeper()

	const modifiedTokenURI = "modifiedtokenuri"

	t.Log("Set Token")
	token := types.NewToken(defaultName, defaultSymbol, defaultTokenURI, sdk.NewInt(defaultDecimals), true)
	{
		require.NoError(t, keeper.SetToken(ctx, token))
		modifyTokenURIPermission := types.NewModifyTokenURIPermission(token.GetSymbol())
		keeper.AddPermission(ctx, addr1, modifyTokenURIPermission)
	}
	{
		require.NoError(t, keeper.ModifyTokenURI(ctx, addr1, token.GetSymbol(), modifiedTokenURI))
	}
	t.Log("Compare Token")
	{
		store := ctx.KVStore(keeper.storeKey)
		bz := store.Get(types.TokenSymbolKey(token.GetSymbol()))
		actual := keeper.mustDecodeToken(bz)
		require.Equal(t, modifiedTokenURI, actual.GetTokenURI())
	}
}

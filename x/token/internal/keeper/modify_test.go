package keeper

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	linktype "github.com/line/link/types"
	"github.com/line/link/x/token/internal/types"
	"github.com/stretchr/testify/require"
)

func TestModifyTokenName(t *testing.T) {
	const modifiedTokenName = "modifiedTokenName"
	const modifiedTokenURI = "modifiedTokenURI"
	nameChange := linktype.NewChange("name", modifiedTokenName)
	tokenURIChange := linktype.NewChange("token_uri", modifiedTokenURI)

	ctx := cacheKeeper()
	token := aToken(defaultSymbol)
	tokenWithoutPerm := aToken(defaultSymbol + "2")
	modifyPermission := types.NewModifyPermission(token.GetSymbol())

	// Given Token And Permission
	require.NoError(t, keeper.SetToken(ctx, token))
	keeper.AddPermission(ctx, addr1, modifyPermission)

	t.Logf("Test to modify name for token to %s", modifiedTokenName)
	{
		// When modify token name
		require.NoError(t, keeper.ModifyToken(ctx, addr1, token.GetSymbol(), nameChange))

		// Then token name is modified
		store := ctx.KVStore(keeper.storeKey)
		bz := store.Get(types.TokenSymbolKey(token.GetSymbol()))
		actual := keeper.mustDecodeToken(bz)
		require.Equal(t, modifiedTokenName, actual.GetName())
	}
	t.Logf("Test to modify token uri for token to %s", modifiedTokenURI)
	{
		// When modify token uri
		require.NoError(t, keeper.ModifyToken(ctx, addr1, token.GetSymbol(), tokenURIChange))

		// Then token uri is modified
		store := ctx.KVStore(keeper.storeKey)
		bz := store.Get(types.TokenSymbolKey(token.GetSymbol()))
		actual := keeper.mustDecodeToken(bz)
		require.Equal(t, modifiedTokenURI, actual.GetTokenURI())
	}
	t.Log("Test with nonexistent symbol")
	{
		// Given nonexistent symbol
		nonExistentSymbols := "symbol2"

		// When modify token name with invalid symbol, Then error is occurred
		require.EqualError(t, keeper.ModifyToken(ctx, addr1, nonExistentSymbols, nameChange),
			types.ErrTokenNotExist(types.DefaultCodespace, nonExistentSymbols).Error())
	}
	t.Log("Test without permission")
	{
		// Given Token without Permission
		require.NoError(t, keeper.SetToken(ctx, tokenWithoutPerm))
		invalidPerm := types.NewModifyPermission(tokenWithoutPerm.GetSymbol())

		// When modify token name with invalid permission, Then error is occurred
		require.EqualError(t, keeper.ModifyToken(ctx, addr1, tokenWithoutPerm.GetSymbol(), nameChange),
			types.ErrTokenNoPermission(types.DefaultCodespace, addr1, invalidPerm).Error())
	}
}

func aToken(symbol string) types.Token {
	return types.NewToken(defaultName, symbol, defaultTokenURI, sdk.NewInt(defaultDecimals), true)
}

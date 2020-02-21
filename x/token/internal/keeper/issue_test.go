package keeper

import (
	"strings"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/line/link/x/token/internal/types"
	"github.com/stretchr/testify/require"
)

func TestKeeper_IssueToken(t *testing.T) {
	ctx := cacheKeeper()
	t.Log("Issue Token")
	expected := types.NewToken(defaultName, defaultSymbol, defaultTokenURI, sdk.NewInt(defaultDecimals), true)
	{
		require.NoError(t, keeper.IssueToken(ctx, expected, sdk.NewInt(defaultAmount), addr1))
	}
	t.Log("Get Token")
	{
		actual, err := keeper.GetToken(ctx, defaultSymbol)
		require.NoError(t, err)
		verifyTokenFunc(t, expected, actual)
	}
	t.Log("Permission")
	{
		require.True(t, keeper.HasPermission(ctx, addr1, types.NewModifyTokenURIPermission(defaultSymbol)))
		require.True(t, keeper.HasPermission(ctx, addr1, types.NewMintPermission(defaultSymbol)))
		require.True(t, keeper.HasPermission(ctx, addr1, types.NewBurnPermission(defaultSymbol)))
	}
	t.Log("Permission only addr1 has the permissions")
	{
		require.False(t, keeper.HasPermission(ctx, addr2, types.NewModifyTokenURIPermission(defaultSymbol)))
		require.False(t, keeper.HasPermission(ctx, addr2, types.NewMintPermission(defaultSymbol)))
		require.False(t, keeper.HasPermission(ctx, addr2, types.NewBurnPermission(defaultSymbol)))
	}
	t.Log("Total supply")
	{
		supply, err := keeper.GetSupply(ctx, defaultSymbol)
		require.NoError(t, err)
		require.Equal(t, int64(defaultAmount), supply.Int64())
	}
	t.Log("Balance of Account")
	{
		supply := keeper.GetAccountBalance(ctx, defaultSymbol, addr1)
		require.Equal(t, int64(defaultAmount), supply.Int64())
	}
}

func TestKeeper_IssueTokenNotMintable(t *testing.T) {
	ctx := cacheKeeper()
	t.Log("Issue a Token Not Mintable")
	expected := types.NewToken(defaultName, defaultSymbol, defaultTokenURI, sdk.NewInt(defaultDecimals), false)
	{
		require.NoError(t, keeper.IssueToken(ctx, expected, sdk.NewInt(defaultAmount), addr1))
	}
	{
		actual, err := keeper.GetToken(ctx, defaultSymbol)
		require.NoError(t, err)
		verifyTokenFunc(t, expected, actual)
	}
	t.Log("Permission only addr1 has no mint/burn permisssions")
	{
		require.True(t, keeper.HasPermission(ctx, addr1, types.NewModifyTokenURIPermission(defaultSymbol)))
		require.False(t, keeper.HasPermission(ctx, addr1, types.NewMintPermission(defaultSymbol)))
		require.False(t, keeper.HasPermission(ctx, addr1, types.NewBurnPermission(defaultSymbol)))
	}
}

func TestKeeper_IssueTokenTooLongTokenURI(t *testing.T) {
	ctx := cacheKeeper()

	length1001String := strings.Repeat("Eng글자日本語はスゲ", 91) // 11 * 91 = 1001

	t.Log("issue a token with too long token uri")
	{
		token := types.NewToken(defaultName, defaultSymbol, length1001String, sdk.NewInt(defaultDecimals), true)
		require.EqualError(t, keeper.IssueToken(ctx, token, sdk.NewInt(defaultAmount), addr1), types.ErrInvalidTokenURILength(types.DefaultCodespace, length1001String).Error())
	}

	t.Log("issue a token and update with too long token uri")
	{
		token := types.NewToken(defaultName, defaultSymbol, defaultTokenURI, sdk.NewInt(defaultDecimals), true)
		require.NoError(t, keeper.IssueToken(ctx, token, sdk.NewInt(defaultAmount), addr1))
		require.EqualError(t, keeper.ModifyTokenURI(ctx, addr1, defaultSymbol, length1001String), types.ErrInvalidTokenURILength(types.DefaultCodespace, length1001String).Error())
	}
}

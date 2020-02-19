package keeper

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/line/link/x/token/internal/types"
	"github.com/stretchr/testify/require"
)

func TestKeeper_BurnTokens(t *testing.T) {
	ctx := cacheKeeper()
	t.Log("Issue Token")
	{
		token := types.NewToken(defaultName, defaultSymbol, defaultTokenURI, sdk.NewInt(defaultDecimals), true)
		require.NoError(t, keeper.IssueToken(ctx, token, sdk.NewInt(defaultAmount+defaultAmount), addr1))
	}
	t.Log("Total supply")
	{
		supply, err := keeper.GetSupply(ctx, defaultSymbol)
		require.NoError(t, err)
		require.Equal(t, int64(defaultAmount+defaultAmount), supply.Int64())
	}
	t.Log("Balance of Account")
	{
		supply := keeper.GetAccountBalance(ctx, defaultSymbol, addr1)
		require.Equal(t, int64(defaultAmount+defaultAmount), supply.Int64())
	}

	t.Log("Burn Tokens by addr1")
	{
		err := keeper.BurnTokens(ctx, sdk.NewCoins(sdk.NewCoin(defaultSymbol, sdk.NewInt(defaultAmount))), addr1)
		require.NoError(t, err)
	}
	t.Log("Total supply")
	{
		supply, err := keeper.GetSupply(ctx, defaultSymbol)
		require.Equal(t, int64(defaultAmount), supply.Int64())
		require.NoError(t, err)
	}
	t.Log("Balance of Account 1")
	{
		supply := keeper.GetAccountBalance(ctx, defaultSymbol, addr1)
		require.Equal(t, int64(defaultAmount), supply.Int64())
	}
}

func TestKeeper_BurnTokensWithoutPermissions(t *testing.T) {
	ctx := cacheKeeper()
	t.Log("Issue Token")
	{
		token := types.NewToken(defaultName, defaultSymbol, defaultTokenURI, sdk.NewInt(defaultDecimals), true)
		require.NoError(t, keeper.IssueToken(ctx, token, sdk.NewInt(defaultAmount), addr1))
	}

	t.Log("Transfer Enough Token")
	{
		err := keeper.Transfer(ctx, addr1, addr2, defaultSymbol, sdk.NewInt(defaultAmount))
		require.NoError(t, err)
	}

	t.Log("Burn Tokens by addr2. Expect Fail")
	{
		err := keeper.BurnTokens(ctx, sdk.NewCoins(sdk.NewCoin(defaultSymbol, sdk.NewInt(defaultAmount))), addr2)
		require.Error(t, err)
		require.EqualError(t, err, types.ErrTokenNoPermission(types.DefaultCodespace, addr2, types.NewBurnPermission(defaultSymbol)).Error())
	}
}

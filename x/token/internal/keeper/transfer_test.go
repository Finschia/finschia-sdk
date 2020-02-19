package keeper

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/line/link/x/token/internal/types"
	"github.com/stretchr/testify/require"
)

func TestKeeper_Transfer(t *testing.T) {
	ctx := cacheKeeper()
	t.Log("Issue Token")
	{
		token := types.NewToken(defaultName, defaultSymbol, defaultTokenURI, sdk.NewInt(defaultDecimals), true)
		require.NoError(t, keeper.IssueToken(ctx, token, sdk.NewInt(defaultAmount), addr1))
	}
	t.Log("Total supply")
	{
		supply, err := keeper.GetSupply(ctx, defaultSymbol)
		require.NoError(t, err)
		require.Equal(t, int64(defaultAmount), supply.Int64())
	}
	t.Log("Balance of Account 1")
	{
		supply := keeper.GetAccountBalance(ctx, defaultSymbol, addr1)
		require.Equal(t, int64(defaultAmount), supply.Int64())
	}
	t.Log("Balance of Account 2")
	{
		supply := keeper.GetAccountBalance(ctx, defaultSymbol, addr2)
		require.Equal(t, int64(0), supply.Int64())
	}
	t.Log("Transfer Token")
	{
		err := keeper.Transfer(ctx, addr1, addr2, defaultSymbol, sdk.NewInt(defaultAmount))
		require.NoError(t, err)
	}
	t.Log("Total supply")
	{
		supply, err := keeper.GetSupply(ctx, defaultSymbol)
		require.NoError(t, err)
		require.Equal(t, int64(defaultAmount), supply.Int64())
	}
	t.Log("Balance of Account 1")
	{
		supply := keeper.GetAccountBalance(ctx, defaultSymbol, addr1)
		require.Equal(t, int64(0), supply.Int64())
	}
	t.Log("Balance of Account 2")
	{
		supply := keeper.GetAccountBalance(ctx, defaultSymbol, addr2)
		require.Equal(t, int64(defaultAmount), supply.Int64())
	}
}

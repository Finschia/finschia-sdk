package keeper

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestKeeper_GetAccountSupply(t *testing.T) {
	ctx := cacheKeeper()
	t.Log("Balance of addr1. Expect 0")
	{
		require.Equal(t, int64(0), keeper.GetAccountBalance(ctx, defaultSymbol, defaultTokenIDFT, addr1).Int64())
	}
	t.Log("Set tokens to addr1")
	{
		err := keeper.mintTokens(ctx, sdk.NewCoins(sdk.NewCoin(defaultSymbol+defaultTokenIDFT, sdk.NewInt(defaultAmount))), addr1)
		require.NoError(t, err)
	}
	t.Log("Balance of addr1.")
	{
		require.Equal(t, int64(defaultAmount), keeper.GetAccountBalance(ctx, defaultSymbol, defaultTokenIDFT, addr1).Int64())
	}
}

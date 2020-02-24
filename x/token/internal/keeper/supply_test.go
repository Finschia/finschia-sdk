package keeper

import (
	"testing"

	"github.com/line/link/x/token/internal/types"
	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func verifySupplyFunc(t *testing.T, expected types.Supply, actual types.Supply) {
	require.Equal(t, expected.GetSymbol(), actual.GetSymbol())
	require.Equal(t, expected.GetTotal().Int64(), actual.GetTotal().Int64())
}

func TestKeeper_GetSupplyInt(t *testing.T) {
	ctx := cacheKeeper()
	t.Log("Prepare Supply and Token")
	expected := types.DefaultSupply(defaultSymbol)
	{
		store := ctx.KVStore(keeper.storeKey)
		b := keeper.cdc.MustMarshalBinaryLengthPrefixed(expected)
		store.Set(types.SupplyKey(expected.GetSymbol()), b)
		token := types.NewToken(defaultName, defaultSymbol, defaultTokenURI, sdk.NewInt(defaultDecimals), true)
		store.Set(types.TokenSymbolKey(expected.GetSymbol()), keeper.cdc.MustMarshalBinaryBare(token))
	}
	t.Log("Get Supply")
	{
		actual, err := keeper.getSupply(ctx, defaultSymbol)
		require.NoError(t, err)
		verifySupplyFunc(t, expected, actual)
	}
	t.Log("Get Supply Int")
	{
		actual, err := keeper.GetSupplyInt(ctx, defaultSymbol)
		require.NoError(t, err)
		require.Equal(t, expected.GetTotal().Int64(), actual.Int64())
	}
}

func TestKeeper_MintSupply(t *testing.T) {
	ctx := cacheKeeper()
	t.Log("Prepare Token")
	{
		token := types.NewToken(defaultName, defaultSymbol, defaultTokenURI, sdk.NewInt(defaultDecimals), true)
		require.NoError(t, keeper.SetToken(ctx, token))
	}
	t.Log("Set Account")
	acc := types.NewBaseAccountWithAddress(defaultSymbol, addr1)
	{
		require.NoError(t, keeper.SetAccount(ctx, acc))
	}
	t.Log("Mint Supply")
	{
		require.NoError(t, keeper.MintSupply(ctx, defaultSymbol, addr1, sdk.NewInt(defaultAmount)))
	}
	t.Log("Get Balance")
	{
		balance := keeper.GetBalance(ctx, defaultSymbol, addr1)
		require.Equal(t, sdk.NewInt(defaultAmount).Int64(), balance.Int64())
	}
}

func TestKeeper_BurnSupply(t *testing.T) {
	ctx := cacheKeeper()
	t.Log("Prepare Token")
	{
		token := types.NewToken(defaultName, defaultSymbol, defaultTokenURI, sdk.NewInt(defaultDecimals), true)
		require.NoError(t, keeper.SetToken(ctx, token))
	}
	t.Log("Set Account")
	acc := types.NewBaseAccountWithAddress(defaultSymbol, addr1)
	{
		require.NoError(t, keeper.SetAccount(ctx, acc))
	}
	t.Log("Set Balance And Supply")
	{
		require.NoError(t, keeper.SetBalance(ctx, defaultSymbol, addr1, sdk.NewInt(defaultAmount)))
		keeper.setSupply(ctx, types.DefaultSupply(defaultSymbol).SetTotal(sdk.NewInt(defaultAmount)))
	}
	t.Log("Burn Supply")
	{
		require.NoError(t, keeper.BurnSupply(ctx, defaultSymbol, addr1, sdk.NewInt(defaultAmount)))
	}
	t.Log("Get Balance")
	{
		balance := keeper.GetBalance(ctx, defaultSymbol, addr1)
		require.Equal(t, sdk.ZeroInt().Int64(), balance.Int64())
	}
}

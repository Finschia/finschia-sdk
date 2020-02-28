package keeper

import (
	"testing"

	"github.com/line/link/x/token/internal/types"
	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func verifySupplyFunc(t *testing.T, expected types.Supply, actual types.Supply) {
	require.Equal(t, expected.GetContractID(), actual.GetContractID())
	require.Equal(t, expected.GetTotalSupply().Int64(), actual.GetTotalSupply().Int64())
}

func TestKeeper_GetTotalInt(t *testing.T) {
	ctx := cacheKeeper()
	t.Log("Prepare Supply and Token")
	expected := types.DefaultSupply(defaultContractID)
	{
		store := ctx.KVStore(keeper.storeKey)
		b := keeper.cdc.MustMarshalBinaryLengthPrefixed(expected)
		store.Set(types.SupplyKey(expected.GetContractID()), b)
		token := types.NewToken(defaultContractID, defaultName, defaultSymbol, defaultImageURI, sdk.NewInt(defaultDecimals), true)
		store.Set(types.TokenKey(expected.GetContractID()), keeper.cdc.MustMarshalBinaryBare(token))
	}
	t.Log("Get Supply")
	{
		actual, err := keeper.getSupply(ctx, defaultContractID)
		require.NoError(t, err)
		verifySupplyFunc(t, expected, actual)
	}
	t.Log("Get Total Supply Int")
	{
		actual, err := keeper.GetTotalInt(ctx, defaultContractID, types.QuerySupply)
		require.NoError(t, err)
		require.Equal(t, expected.GetTotalSupply().Int64(), actual.Int64())
	}
	t.Log("Get Total Mint Int")
	{
		actual, err := keeper.GetTotalInt(ctx, defaultContractID, types.QueryMint)
		require.NoError(t, err)
		require.Equal(t, expected.GetTotalMint().Int64(), actual.Int64())
	}
	t.Log("Get Total Burn Int")
	{
		actual, err := keeper.GetTotalInt(ctx, defaultContractID, types.QueryBurn)
		require.NoError(t, err)
		require.Equal(t, expected.GetTotalBurn().Int64(), actual.Int64())
	}
}

func TestKeeper_MintSupply(t *testing.T) {
	ctx := cacheKeeper()
	t.Log("Prepare Token")
	{
		token := types.NewToken(defaultContractID, defaultName, defaultSymbol, defaultImageURI, sdk.NewInt(defaultDecimals), true)
		require.NoError(t, keeper.SetToken(ctx, token))
	}
	t.Log("Set Account")
	acc := types.NewBaseAccountWithAddress(defaultContractID, addr1)
	{
		require.NoError(t, keeper.SetAccount(ctx, acc))
	}
	t.Log("Mint Supply")
	{
		require.NoError(t, keeper.MintSupply(ctx, defaultContractID, addr1, sdk.NewInt(defaultAmount)))
	}
	t.Log("Get Balance")
	{
		balance := keeper.GetBalance(ctx, defaultContractID, addr1)
		require.Equal(t, sdk.NewInt(defaultAmount).Int64(), balance.Int64())
	}
	t.Log("Get Total Supply Int")
	{
		actual, err := keeper.GetTotalInt(ctx, defaultContractID, types.QuerySupply)
		require.NoError(t, err)
		require.Equal(t, sdk.NewInt(defaultAmount), actual)
	}
	t.Log("Get Total Mint Int")
	{
		actual, err := keeper.GetTotalInt(ctx, defaultContractID, types.QueryMint)
		require.NoError(t, err)
		require.Equal(t, sdk.NewInt(defaultAmount), actual)
	}
	t.Log("Get Total Burn Int")
	{
		actual, err := keeper.GetTotalInt(ctx, defaultContractID, types.QueryBurn)
		require.NoError(t, err)
		require.Equal(t, sdk.ZeroInt(), actual)
	}
}

func TestKeeper_BurnSupply(t *testing.T) {
	ctx := cacheKeeper()
	t.Log("Prepare Token")
	{
		token := types.NewToken(defaultContractID, defaultName, defaultSymbol, defaultImageURI, sdk.NewInt(defaultDecimals), true)
		require.NoError(t, keeper.SetToken(ctx, token))
	}
	t.Log("Set Account")
	acc := types.NewBaseAccountWithAddress(defaultContractID, addr1)
	{
		require.NoError(t, keeper.SetAccount(ctx, acc))
	}
	t.Log("Set Balance And Supply")
	{
		require.NoError(t, keeper.SetBalance(ctx, defaultContractID, addr1, sdk.NewInt(defaultAmount)))
		keeper.setSupply(ctx, types.DefaultSupply(defaultContractID).SetTotalSupply(sdk.NewInt(defaultAmount)))
	}
	t.Log("Burn Supply")
	{
		require.NoError(t, keeper.BurnSupply(ctx, defaultContractID, addr1, sdk.NewInt(defaultAmount)))
	}
	t.Log("Get Balance")
	{
		balance := keeper.GetBalance(ctx, defaultContractID, addr1)
		require.Equal(t, sdk.ZeroInt().Int64(), balance.Int64())
	}
	t.Log("Get Total Supply Int")
	{
		actual, err := keeper.GetTotalInt(ctx, defaultContractID, types.QuerySupply)
		require.NoError(t, err)
		require.Equal(t, sdk.ZeroInt(), actual)
	}
	t.Log("Get Total Mint Int")
	{
		actual, err := keeper.GetTotalInt(ctx, defaultContractID, types.QueryMint)
		require.NoError(t, err)
		require.Equal(t, sdk.NewInt(defaultAmount), actual)
	}
	t.Log("Get Total Burn Int")
	{
		actual, err := keeper.GetTotalInt(ctx, defaultContractID, types.QueryBurn)
		require.NoError(t, err)
		require.Equal(t, sdk.NewInt(defaultAmount), actual)
	}
}

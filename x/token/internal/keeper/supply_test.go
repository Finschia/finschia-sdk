package keeper

import (
	"testing"

	"github.com/line/link-modules/x/token/internal/types"
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
		token := types.NewToken(defaultContractID, defaultName, defaultSymbol, defaultMeta, defaultImageURI, sdk.NewInt(defaultDecimals), true)
		store.Set(types.TokenKey(expected.GetContractID()), keeper.cdc.MustMarshalBinaryBare(token))
	}
	t.Log("Get Supply")
	{
		actual, err := keeper.getSupply(ctx)
		require.NoError(t, err)
		verifySupplyFunc(t, expected, actual)
	}
	t.Log("Get Total Supply Int")
	{
		actual, err := keeper.GetTotalInt(ctx, types.QuerySupply)
		require.NoError(t, err)
		require.Equal(t, expected.GetTotalSupply().Int64(), actual.Int64())
	}
	t.Log("Get Total Mint Int")
	{
		actual, err := keeper.GetTotalInt(ctx, types.QueryMint)
		require.NoError(t, err)
		require.Equal(t, expected.GetTotalMint().Int64(), actual.Int64())
	}
	t.Log("Get Total Burn Int")
	{
		actual, err := keeper.GetTotalInt(ctx, types.QueryBurn)
		require.NoError(t, err)
		require.Equal(t, expected.GetTotalBurn().Int64(), actual.Int64())
	}
}

func TestKeeper_MintSupply(t *testing.T) {
	ctx := cacheKeeper()
	t.Log("Prepare Token")
	{
		token := types.NewToken(defaultContractID, defaultName, defaultSymbol, defaultMeta, defaultImageURI, sdk.NewInt(defaultDecimals), true)
		require.NoError(t, keeper.SetToken(ctx, token))
	}
	t.Log("Set Account")
	acc := types.NewBaseAccountWithAddress(defaultContractID, addr1)
	{
		require.NoError(t, keeper.SetAccount(ctx, acc))
	}
	t.Log("Mint Supply")
	{
		require.NoError(t, keeper.MintSupply(ctx, addr1, sdk.NewInt(defaultAmount)))
	}
	t.Log("Get Balance")
	{
		balance := keeper.GetBalance(ctx, addr1)
		require.Equal(t, sdk.NewInt(defaultAmount).Int64(), balance.Int64())
	}
	t.Log("Get Total Supply Int")
	{
		actual, err := keeper.GetTotalInt(ctx, types.QuerySupply)
		require.NoError(t, err)
		require.Equal(t, sdk.NewInt(defaultAmount), actual)
	}
	t.Log("Get Total Mint Int")
	{
		actual, err := keeper.GetTotalInt(ctx, types.QueryMint)
		require.NoError(t, err)
		require.Equal(t, sdk.NewInt(defaultAmount), actual)
	}
	t.Log("Get Total Burn Int")
	{
		actual, err := keeper.GetTotalInt(ctx, types.QueryBurn)
		require.NoError(t, err)
		require.Equal(t, sdk.ZeroInt(), actual)
	}
}

func TestKeeper_BurnSupply(t *testing.T) {
	ctx := cacheKeeper()
	t.Log("Prepare Token")
	{
		token := types.NewToken(defaultContractID, defaultName, defaultSymbol, defaultMeta, defaultImageURI, sdk.NewInt(defaultDecimals), true)
		require.NoError(t, keeper.SetToken(ctx, token))
	}
	t.Log("Set Account")
	acc := types.NewBaseAccountWithAddress(defaultContractID, addr1)
	{
		require.NoError(t, keeper.SetAccount(ctx, acc))
	}
	t.Log("Set Balance And Supply")
	{
		require.NoError(t, keeper.SetBalance(ctx, addr1, sdk.NewInt(defaultAmount)))
		keeper.setSupply(ctx, types.DefaultSupply(defaultContractID).SetTotalSupply(sdk.NewInt(defaultAmount)))
	}
	t.Log("Burn Supply")
	{
		require.NoError(t, keeper.BurnSupply(ctx, addr1, sdk.NewInt(defaultAmount)))
	}
	t.Log("Get Balance")
	{
		balance := keeper.GetBalance(ctx, addr1)
		require.Equal(t, sdk.ZeroInt().Int64(), balance.Int64())
	}
	t.Log("Get Total Supply Int")
	{
		actual, err := keeper.GetTotalInt(ctx, types.QuerySupply)
		require.NoError(t, err)
		require.Equal(t, sdk.ZeroInt(), actual)
	}
	t.Log("Get Total Mint Int")
	{
		actual, err := keeper.GetTotalInt(ctx, types.QueryMint)
		require.NoError(t, err)
		require.Equal(t, sdk.NewInt(defaultAmount), actual)
	}
	t.Log("Get Total Burn Int")
	{
		actual, err := keeper.GetTotalInt(ctx, types.QueryBurn)
		require.NoError(t, err)
		require.Equal(t, sdk.NewInt(defaultAmount), actual)
	}
}

func TestKeeper_Handle_Overflows(t *testing.T) {
	ctx := cacheKeeper()

	t.Log("Prepare Supply and Token")
	expected := types.DefaultSupply(defaultContractID)
	{
		store := ctx.KVStore(keeper.storeKey)
		b := keeper.cdc.MustMarshalBinaryLengthPrefixed(expected)
		store.Set(types.SupplyKey(expected.GetContractID()), b)
		token := types.NewToken(defaultContractID, defaultName, defaultSymbol, defaultImageURI, defaultMeta, sdk.NewInt(defaultDecimals), true)
		store.Set(types.TokenKey(expected.GetContractID()), keeper.cdc.MustMarshalBinaryBare(token))
	}

	// int64 is the set of all signed 64-bit integers.
	// Range: -9223372036854775808 through 9223372036854775807.

	// Int wraps integer with 256 bit range bound
	// Checks overflow, underflow and division by zero
	// Exists in range from -(2^maxBitLen-1) to 2^maxBitLen-1

	t.Log("Set supply less than the overflow limit")
	maxInt64Supply := sdk.NewInt(9223372036854775807)

	initialSupply := maxInt64Supply.Mul(maxInt64Supply).Mul(maxInt64Supply).Mul(maxInt64Supply)
	newSupply := types.NewSupply(defaultContractID, initialSupply)
	keeper.setSupply(ctx, newSupply)

	ts, err := keeper.GetTotalInt(ctx, types.QuerySupply)
	require.NoError(t, err)
	require.Equal(t, newSupply.GetTotalSupply(), ts)

	tm, err := keeper.GetTotalInt(ctx, types.QueryMint)
	require.NoError(t, err)
	require.Equal(t, newSupply.GetTotalMint(), tm)

	tb, err := keeper.GetTotalInt(ctx, types.QueryBurn)
	require.NoError(t, err)
	require.Equal(t, newSupply.GetTotalBurn(), tb)

	// inflate over the overflow limit
	t.Log("Inflate the supply over the overflow limit")
	addToOverflow := initialSupply.Mul(sdk.NewInt(8))
	err = keeper.MintSupply(ctx, addr1, addToOverflow)
	require.Equal(t, types.ErrSupplyOverflow, err)

	// should have not changed
	t.Log("Totals have not changed")
	ts, err = keeper.GetTotalInt(ctx, types.QuerySupply)
	require.NoError(t, err)
	require.Equal(t, newSupply.GetTotalSupply(), ts)

	tm, err = keeper.GetTotalInt(ctx, types.QueryMint)
	require.NoError(t, err)
	require.Equal(t, newSupply.GetTotalMint(), tm)

	tb, err = keeper.GetTotalInt(ctx, types.QueryBurn)
	require.NoError(t, err)
	require.Equal(t, newSupply.GetTotalBurn(), tb)

	// deflate below the overflow limit - it will return insufficient fund instead of panicking
	t.Log("Deflate the supply below the overflow limit - will return insufficient fund instead of panicking")
	subToOverflow := initialSupply.Mul(sdk.NewInt(8))
	err = keeper.BurnSupply(ctx, addr1, subToOverflow)
	require.True(t, types.ErrInsufficientSupply.Is(err))

	// should have not changed
	t.Log("Totals have not changed")
	ts, err = keeper.GetTotalInt(ctx, types.QuerySupply)
	require.NoError(t, err)
	require.Equal(t, newSupply.GetTotalSupply(), ts)

	tm, err = keeper.GetTotalInt(ctx, types.QueryMint)
	require.NoError(t, err)
	require.Equal(t, newSupply.GetTotalMint(), tm)

	tb, err = keeper.GetTotalInt(ctx, types.QueryBurn)
	require.NoError(t, err)
	require.Equal(t, newSupply.GetTotalBurn(), tb)
}

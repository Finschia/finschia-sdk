package keeper

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/line/lbm-sdk/v2/x/collection/internal/types"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/crypto/secp256k1"
)

func TestKeeper_GetTotalInt(t *testing.T) {
	ctx := cacheKeeper()
	addr1 := sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())

	t.Log("Prepare Supply and FT")
	expected := types.DefaultSupply(defaultContractID)
	ft := types.NewFT(defaultContractID, defaultTokenIDFT, defaultName, defaultMeta, sdk.ZeroInt(), true)
	{
		keeper.SetSupply(ctx, expected)
		err := keeper.SetCollection(ctx, types.NewCollection(defaultContractID, defaultName, defaultMeta, defaultImgURI))
		require.NoError(t, err)
		err = keeper.IssueFT(ctx, addr1, addr1, ft, sdk.NewInt(defaultAmount))
		require.NoError(t, err)
	}
	t.Log("Get Total Supply Int")
	{
		actual, err := keeper.GetTotalInt(ctx, defaultTokenIDFT, types.QuerySupply)
		require.NoError(t, err)
		require.Equal(t, int64(defaultAmount), actual.Int64())
	}
	t.Log("Get Total Mint Int")
	{
		actual, err := keeper.GetTotalInt(ctx, defaultTokenIDFT, types.QueryMint)
		require.NoError(t, err)
		require.Equal(t, int64(defaultAmount), actual.Int64())
	}
	t.Log("Get Total Burn Int")
	{
		actual, err := keeper.GetTotalInt(ctx, defaultTokenIDFT, types.QueryBurn)
		require.NoError(t, err)
		require.Equal(t, int64(0), actual.Int64())
	}
}

func TestKeeper_MintSupply(t *testing.T) {
	ctx := cacheKeeper()
	addr1 := sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())

	t.Log("Prepare Supply and FT")
	expected := types.DefaultSupply(defaultContractID)
	ft := types.NewFT(defaultContractID, defaultTokenIDFT, defaultName, defaultMeta, sdk.ZeroInt(), true)
	{
		keeper.SetSupply(ctx, expected)
		err := keeper.SetCollection(ctx, types.NewCollection(defaultContractID, defaultName, defaultMeta, defaultImgURI))
		require.NoError(t, err)
		err = keeper.IssueFT(ctx, addr1, addr1, ft, sdk.NewInt(defaultAmount))
		require.NoError(t, err)
	}
	t.Log("Get Balance")
	{
		balance, err := keeper.GetBalance(ctx, defaultTokenIDFT, addr1)
		require.NoError(t, err)
		require.Equal(t, sdk.NewInt(defaultAmount).Int64(), balance.Int64())
	}
	t.Log("Mint Supply")
	{
		require.NoError(t, keeper.MintSupply(ctx, addr1, types.OneCoins(defaultTokenIDFT)))
	}
	t.Log("Get Balance")
	{
		balance, err := keeper.GetBalance(ctx, defaultTokenIDFT, addr1)
		require.NoError(t, err)
		require.Equal(t, sdk.NewInt(defaultAmount+1).Int64(), balance.Int64())
	}
	t.Log("Get Total Supply Int")
	{
		actual, err := keeper.GetTotalInt(ctx, defaultTokenIDFT, types.QuerySupply)
		require.NoError(t, err)
		require.Equal(t, sdk.NewInt(defaultAmount+1), actual)
	}
	t.Log("Get Total Mint Int")
	{
		actual, err := keeper.GetTotalInt(ctx, defaultTokenIDFT, types.QueryMint)
		require.NoError(t, err)
		require.Equal(t, sdk.NewInt(defaultAmount+1), actual)
	}
	t.Log("Get Total Burn Int")
	{
		actual, err := keeper.GetTotalInt(ctx, defaultTokenIDFT, types.QueryBurn)
		require.NoError(t, err)
		require.Equal(t, sdk.ZeroInt(), actual)
	}
}

func TestKeeper_BurnSupply(t *testing.T) {
	ctx := cacheKeeper()
	addr1 := sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())

	t.Log("Prepare Supply and FT")
	expected := types.DefaultSupply(defaultContractID)
	ft := types.NewFT(defaultContractID, defaultTokenIDFT, defaultName, defaultMeta, sdk.ZeroInt(), true)
	{
		keeper.SetSupply(ctx, expected)
		err := keeper.SetCollection(ctx, types.NewCollection(defaultContractID, defaultName, defaultMeta, defaultImgURI))
		require.NoError(t, err)
		err = keeper.IssueFT(ctx, addr1, addr1, ft, sdk.NewInt(defaultAmount))
		require.NoError(t, err)
	}
	t.Log("Get Balance")
	{
		balance, err := keeper.GetBalance(ctx, defaultTokenIDFT, addr1)
		require.NoError(t, err)
		require.Equal(t, sdk.NewInt(defaultAmount).Int64(), balance.Int64())
	}
	t.Log("Burn Supply")
	{
		require.NoError(t, keeper.BurnSupply(ctx, addr1, types.OneCoins(defaultTokenIDFT)))
	}
	t.Log("Get Balance")
	{
		balance, err := keeper.GetBalance(ctx, defaultTokenIDFT, addr1)
		require.NoError(t, err)
		require.Equal(t, sdk.NewInt(defaultAmount-1).Int64(), balance.Int64())
	}
	t.Log("Get Total Supply Int")
	{
		actual, err := keeper.GetTotalInt(ctx, defaultTokenIDFT, types.QuerySupply)
		require.NoError(t, err)
		require.Equal(t, sdk.NewInt(defaultAmount-1), actual)
	}
	t.Log("Get Total Mint Int")
	{
		actual, err := keeper.GetTotalInt(ctx, defaultTokenIDFT, types.QueryMint)
		require.NoError(t, err)
		require.Equal(t, sdk.NewInt(defaultAmount), actual)
	}
	t.Log("Get Total Burn Int")
	{
		actual, err := keeper.GetTotalInt(ctx, defaultTokenIDFT, types.QueryBurn)
		require.NoError(t, err)
		require.Equal(t, sdk.NewInt(1), actual)
	}
}

func TestKeeper_Handle_Overflows(t *testing.T) {
	ctx := cacheKeeper()
	addr1 := sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())

	// int64 is the set of all signed 64-bit integers.
	// Range: -9223372036854775808 through 9223372036854775807.

	// Int wraps integer with 256 bit range bound
	// Checks overflow, underflow and division by zero
	// Exists in range from -(2^maxBitLen-1) to 2^maxBitLen-1

	t.Log("Prepare Supply and FT less than the overflow limit")
	maxInt64Supply := sdk.NewInt(9223372036854775807)
	initialSupply := maxInt64Supply.Mul(maxInt64Supply).Mul(maxInt64Supply).Mul(maxInt64Supply)
	ft := types.NewFT(defaultContractID, defaultTokenIDFT, defaultName, defaultMeta, initialSupply, true)
	newSupply := types.NewSupply(defaultContractID, types.NewCoins(types.NewCoin(defaultTokenIDFT, initialSupply)))
	{
		keeper.SetSupply(ctx, newSupply)
		err := keeper.SetCollection(ctx, types.NewCollection(defaultContractID, defaultName, defaultMeta, defaultImgURI))
		require.NoError(t, err)
		err = keeper.IssueFT(ctx, addr1, addr1, ft, sdk.ZeroInt())
		require.NoError(t, err)
	}

	ts, err := keeper.GetTotalInt(ctx, defaultTokenIDFT, types.QuerySupply)
	require.NoError(t, err)
	require.Equal(t, newSupply.GetTotalSupply().AmountOf(defaultTokenIDFT), ts)

	tm, err := keeper.GetTotalInt(ctx, defaultTokenIDFT, types.QueryMint)
	require.NoError(t, err)
	require.Equal(t, newSupply.GetTotalMint().AmountOf(defaultTokenIDFT), tm)

	tb, err := keeper.GetTotalInt(ctx, defaultTokenIDFT, types.QueryBurn)
	require.NoError(t, err)
	require.Equal(t, newSupply.GetTotalBurn().AmountOf(defaultTokenIDFT), tb)

	// inflate over the overflow limit
	t.Log("Inflate the supply over the overflow limit")
	addToOverflow := types.NewCoins(types.NewCoin(defaultTokenIDFT, initialSupply.Mul(sdk.NewInt(8))))
	err = keeper.MintSupply(ctx, addr1, addToOverflow)
	require.Equal(t, types.ErrSupplyOverflow, err)

	// should have not changed
	t.Log("Totals have not changed")
	ts, err = keeper.GetTotalInt(ctx, defaultTokenIDFT, types.QuerySupply)
	require.NoError(t, err)
	require.Equal(t, newSupply.GetTotalSupply().AmountOf(defaultTokenIDFT), ts)

	tm, err = keeper.GetTotalInt(ctx, defaultTokenIDFT, types.QueryMint)
	require.NoError(t, err)
	require.Equal(t, newSupply.GetTotalMint().AmountOf(defaultTokenIDFT), tm)

	tb, err = keeper.GetTotalInt(ctx, defaultTokenIDFT, types.QueryBurn)
	require.NoError(t, err)
	require.Equal(t, newSupply.GetTotalBurn().AmountOf(defaultTokenIDFT), tb)

	// deflate below the overflow limit
	t.Log("Deflate the supply below the overflow limit")
	subToOverflow := types.NewCoins(types.NewCoin(defaultTokenIDFT, initialSupply.Mul(sdk.NewInt(8))))
	err = keeper.BurnSupply(ctx, addr1, subToOverflow)
	require.Error(t, err)
	require.Equal(t, types.ErrSupplyOverflow, err)

	// should have not changed
	t.Log("Totals have not changed")
	ts, err = keeper.GetTotalInt(ctx, defaultTokenIDFT, types.QuerySupply)
	require.NoError(t, err)
	require.Equal(t, newSupply.GetTotalSupply().AmountOf(defaultTokenIDFT), ts)

	tm, err = keeper.GetTotalInt(ctx, defaultTokenIDFT, types.QueryMint)
	require.NoError(t, err)
	require.Equal(t, newSupply.GetTotalMint().AmountOf(defaultTokenIDFT), tm)

	tb, err = keeper.GetTotalInt(ctx, defaultTokenIDFT, types.QueryBurn)
	require.NoError(t, err)
	require.Equal(t, newSupply.GetTotalBurn().AmountOf(defaultTokenIDFT), tb)
}

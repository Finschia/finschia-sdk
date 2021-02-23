package keeper

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/line/link-modules/x/collection/internal/types"
	"github.com/stretchr/testify/require"
)

func TestKeeper_GetAccountSupply(t *testing.T) {
	ctx := cacheKeeper()
	t.Log("Balance of addr1. Expect 0")
	{
		_, err := keeper.GetBalance(ctx, defaultTokenIDFT, addr1)
		require.Error(t, err)
	}
	t.Log("Set tokens to addr1")
	{
		_, err := keeper.AddCoins(ctx, addr1, types.NewCoins(types.NewCoin(defaultTokenIDFT, sdk.NewInt(defaultAmount))))
		require.NoError(t, err)
	}
	t.Log("Balance of addr1.")
	{
		balance, err := keeper.GetBalance(ctx, defaultTokenIDFT, addr1)
		require.NoError(t, err)
		require.Equal(t, int64(defaultAmount), balance.Int64())
	}
}

func verifyAccountFunc(t *testing.T, expected types.Account, actual types.Account) {
	require.Equal(t, expected.GetContractID(), actual.GetContractID())
	require.Equal(t, expected.GetAddress(), actual.GetAddress())
	require.Equal(t, expected.GetCoins().String(), actual.GetCoins().String())
}

func TestKeeper_SetAccount(t *testing.T) {
	ctx := cacheKeeper()
	t.Log("Set Account")
	expected := types.NewBaseAccountWithAddress(defaultContractID, addr1)
	{
		require.NoError(t, keeper.SetAccount(ctx, expected))
	}
	t.Log("Compare Account")
	{
		store := ctx.KVStore(keeper.storeKey)
		bz := store.Get(types.AccountKey(expected.GetContractID(), addr1))
		actual := keeper.mustDecodeAccount(bz)
		verifyAccountFunc(t, expected, actual)
	}
}

func TestKeeper_GetAccount(t *testing.T) {
	ctx := cacheKeeper()
	t.Log("Prepare Account")
	expected := types.NewBaseAccountWithAddress(defaultContractID, addr1)
	{
		store := ctx.KVStore(keeper.storeKey)
		store.Set(types.AccountKey(expected.GetContractID(), addr1), keeper.cdc.MustMarshalBinaryBare(expected))
	}
	t.Log("Get Account")
	{
		actual, err := keeper.GetAccount(ctx, addr1)
		require.NoError(t, err)
		verifyAccountFunc(t, expected, actual)
	}
}

func TestKeeper_UpdateAccount(t *testing.T) {
	ctx := cacheKeeper()
	t.Log("Set Account")
	acc := types.NewBaseAccountWithAddress(defaultContractID, addr1)
	{
		require.NoError(t, keeper.SetAccount(ctx, acc))
	}
	t.Log("Update Account")
	var expected types.Account
	expected = types.NewBaseAccountWithAddress(defaultContractID, addr1)
	expected = expected.SetCoins(types.NewCoins(types.NewCoin(defaultTokenIDFT, sdk.OneInt())))
	{
		require.NoError(t, keeper.UpdateAccount(ctx, expected))
	}
	t.Log("Compare Account")
	{
		store := ctx.KVStore(keeper.storeKey)
		bz := store.Get(types.AccountKey(acc.GetContractID(), addr1))
		actual := keeper.mustDecodeAccount(bz)
		verifyAccountFunc(t, expected, actual)
	}
}

func TestKeeper_GetOrNewAccount(t *testing.T) {
	ctx := cacheKeeper()
	t.Log("Prepare Account")
	expected := types.NewBaseAccountWithAddress(defaultContractID, addr1)
	{
		store := ctx.KVStore(keeper.storeKey)
		store.Set(types.AccountKey(expected.GetContractID(), addr1), keeper.cdc.MustMarshalBinaryBare(expected))
	}
	t.Log("Get Account addr1")
	{
		actual, err := keeper.GetOrNewAccount(ctx, addr1)
		require.NoError(t, err)
		verifyAccountFunc(t, expected, actual)
	}

	expected = types.NewBaseAccountWithAddress(defaultContractID, addr2)
	t.Log("Get Account addr2")
	{
		actual, err := keeper.GetOrNewAccount(ctx, addr2)
		require.NoError(t, err)
		verifyAccountFunc(t, expected, actual)
	}
}

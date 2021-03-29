package keeper

import (
	"testing"

	"github.com/line/lbm-sdk/v2/x/token/internal/types"
	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func TestKeeper_GetBalance(t *testing.T) {
	ctx := cacheKeeper()
	t.Log("Set Account")
	acc := types.NewBaseAccountWithAddress(defaultContractID, addr1)
	{
		require.NoError(t, keeper.SetAccount(ctx, acc))
	}
	t.Log("Get Balance")
	{
		balance := keeper.GetBalance(ctx, addr1)
		require.Equal(t, balance.Int64(), acc.GetBalance().Int64())
	}
}

func TestKeeper_HasBalance(t *testing.T) {
	ctx := cacheKeeper()
	t.Log("Set Account")
	var acc types.Account
	acc = types.NewBaseAccountWithAddress(defaultContractID, addr1)
	acc = acc.SetBalance(sdk.NewInt(defaultAmount))
	{
		require.NoError(t, keeper.SetAccount(ctx, acc))
	}
	t.Log("Has Balance")
	{
		require.True(t, keeper.HasBalance(ctx, addr1, sdk.NewInt(defaultAmount)))
	}
}

func TestKeeper_SetBalance(t *testing.T) {
	ctx := cacheKeeper()
	t.Log("Set Account")
	acc := types.NewBaseAccountWithAddress(defaultContractID, addr1)
	{
		require.NoError(t, keeper.SetAccount(ctx, acc))
	}
	t.Log("Set Balance")
	{
		require.NoError(t, keeper.SetBalance(ctx, addr1, sdk.NewInt(defaultAmount)))
	}
	t.Log("Get Balance")
	{
		balance := keeper.GetBalance(ctx, addr1)
		require.Equal(t, sdk.NewInt(defaultAmount).Int64(), balance.Int64())
	}
}

func TestKeeper_AddBalance(t *testing.T) {
	ctx := cacheKeeper()
	t.Log("Set Account")
	acc := types.NewBaseAccountWithAddress(defaultContractID, addr1)
	{
		require.NoError(t, keeper.SetAccount(ctx, acc))
	}
	t.Log("Add Balance")
	{
		added, err := keeper.AddBalance(ctx, addr1, sdk.NewInt(defaultAmount))
		require.NoError(t, err)
		require.Equal(t, sdk.NewInt(defaultAmount).Int64(), added.Int64())
	}
	t.Log("Get Balance")
	{
		balance := keeper.GetBalance(ctx, addr1)
		require.Equal(t, sdk.NewInt(defaultAmount).Int64(), balance.Int64())
	}
}

func TestKeeper_SubtractBalance(t *testing.T) {
	ctx := cacheKeeper()
	t.Log("Set Account")
	acc := types.NewBaseAccountWithAddress(defaultContractID, addr1)
	{
		require.NoError(t, keeper.SetAccount(ctx, acc))
	}
	t.Log("Set Balance")
	{
		require.NoError(t, keeper.SetBalance(ctx, addr1, sdk.NewInt(defaultAmount)))
	}
	t.Log("Subtract Balance")
	{
		sub, err := keeper.SubtractBalance(ctx, addr1, sdk.NewInt(defaultAmount))
		require.NoError(t, err)
		require.Equal(t, sdk.ZeroInt().Int64(), sub.Int64())
	}
	t.Log("Get Balance")
	{
		balance := keeper.GetBalance(ctx, addr1)
		require.Equal(t, sdk.ZeroInt().Int64(), balance.Int64())
	}
}

func TestKeeper_SendToken(t *testing.T) {
	ctx := cacheKeeper()
	t.Log("Set Account")
	acc := types.NewBaseAccountWithAddress(defaultContractID, addr1)
	{
		require.NoError(t, keeper.SetAccount(ctx, acc))
	}
	t.Log("Set Balance")
	{
		require.NoError(t, keeper.SetBalance(ctx, addr1, sdk.NewInt(defaultAmount)))
	}
	t.Log("Send Balance")
	{
		require.NoError(t, keeper.Send(ctx, addr1, addr2, sdk.NewInt(defaultAmount)))
	}
	{
		require.EqualError(t, keeper.Send(ctx, addr3, addr2, sdk.NewInt(1)), sdkerrors.Wrapf(types.ErrInsufficientBalance, "insufficient account funds for token [9be17165]; 0 < 1").Error())
	}
	t.Log("Get Balance")
	{
		require.Equal(t, sdk.ZeroInt().Int64(), keeper.GetBalance(ctx, addr1).Int64())
		require.Equal(t, sdk.NewInt(defaultAmount).Int64(), keeper.GetBalance(ctx, addr2).Int64())
	}
}

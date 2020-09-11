package keeper

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/line/link-modules/x/coin/internal/types"
	"github.com/stretchr/testify/require"
)

func TestKeeper(t *testing.T) {
	input := SetupTestInput()
	ctx := input.Ctx

	addr1 := sdk.AccAddress([]byte("addr1"))
	addr2 := sdk.AccAddress([]byte("addr2"))
	addr3 := sdk.AccAddress([]byte("addr3"))
	acc := input.Ak.NewAccountWithAddress(ctx, addr1)

	// Test GetCoins/SetCoins
	input.Ak.SetAccount(ctx, acc)
	require.True(t, input.K.GetCoins(ctx, addr1).IsEqual(sdk.NewCoins()))

	acc = input.Ak.GetAccount(ctx, acc.GetAddress())
	err := acc.SetCoins(sdk.NewCoins(sdk.NewInt64Coin("fooc", 15)))
	require.NoError(t, err)
	input.Ak.SetAccount(ctx, acc)
	require.True(t, input.K.GetCoins(ctx, addr1).IsEqual(sdk.NewCoins(sdk.NewInt64Coin("fooc", 15))))

	// Test HasCoins
	require.True(t, input.K.HasCoins(ctx, addr1, sdk.NewCoins(sdk.NewInt64Coin("fooc", 15))))
	require.True(t, input.K.HasCoins(ctx, addr1, sdk.NewCoins(sdk.NewInt64Coin("fooc", 5))))
	require.False(t, input.K.HasCoins(ctx, addr1, sdk.NewCoins(sdk.NewInt64Coin("fooc", 20))))
	require.False(t, input.K.HasCoins(ctx, addr1, sdk.NewCoins(sdk.NewInt64Coin("barc", 5))))

	// Test SendCoins
	err = input.K.SendCoins(ctx, addr1, addr2, sdk.NewCoins(sdk.NewInt64Coin("fooc", 5)))
	require.NoError(t, err)
	require.True(t, input.K.GetCoins(ctx, addr1).IsEqual(sdk.NewCoins(sdk.NewInt64Coin("fooc", 10))))
	require.True(t, input.K.GetCoins(ctx, addr2).IsEqual(sdk.NewCoins(sdk.NewInt64Coin("fooc", 5))))

	err = input.K.SendCoins(ctx, addr1, addr2, sdk.NewCoins(sdk.NewInt64Coin("fooc", 50)))
	require.Error(t, err)
	require.True(t, input.K.GetCoins(ctx, addr1).IsEqual(sdk.NewCoins(sdk.NewInt64Coin("fooc", 10))))
	require.True(t, input.K.GetCoins(ctx, addr2).IsEqual(sdk.NewCoins(sdk.NewInt64Coin("fooc", 5))))

	// Test InputOutputCoins
	input1 := types.NewInput(addr2, sdk.NewCoins(sdk.NewInt64Coin("fooc", 2)))
	output1 := types.NewOutput(addr1, sdk.NewCoins(sdk.NewInt64Coin("fooc", 2)))
	err = input.K.InputOutputCoins(ctx, []types.Input{input1}, []types.Output{output1})
	require.NoError(t, err)
	require.True(t, input.K.GetCoins(ctx, addr1).IsEqual(sdk.NewCoins(sdk.NewInt64Coin("fooc", 12))))
	require.True(t, input.K.GetCoins(ctx, addr2).IsEqual(sdk.NewCoins(sdk.NewInt64Coin("fooc", 3))))

	acc = input.Ak.GetAccount(ctx, acc.GetAddress())
	coins := acc.GetCoins().Add(sdk.NewCoins(sdk.NewInt64Coin("barc", 15))...)
	err = acc.SetCoins(coins)
	require.NoError(t, err)
	input.Ak.SetAccount(ctx, acc)
	require.True(t, input.K.GetCoins(ctx, addr1).IsEqual(sdk.NewCoins(sdk.NewInt64Coin("fooc", 12), sdk.NewInt64Coin("barc", 15))))
	require.True(t, input.K.GetCoins(ctx, addr2).IsEqual(sdk.NewCoins(sdk.NewInt64Coin("fooc", 3))))

	inputs := []types.Input{
		types.NewInput(addr1, sdk.NewCoins(sdk.NewInt64Coin("barc", 4), sdk.NewInt64Coin("fooc", 2))),
		types.NewInput(addr2, sdk.NewCoins(sdk.NewInt64Coin("fooc", 3))),
	}

	outputs := []types.Output{
		types.NewOutput(addr1, sdk.NewCoins(sdk.NewInt64Coin("barc", 1))),
		types.NewOutput(addr2, sdk.NewCoins(sdk.NewInt64Coin("barc", 1))),
		types.NewOutput(addr3, sdk.NewCoins(sdk.NewInt64Coin("barc", 2), sdk.NewInt64Coin("fooc", 5))),
	}
	err = input.K.InputOutputCoins(ctx, inputs, outputs)
	require.NoError(t, err)
	require.True(t, input.K.GetCoins(ctx, addr1).IsEqual(sdk.NewCoins(sdk.NewInt64Coin("barc", 12), sdk.NewInt64Coin("fooc", 10))))
	require.True(t, input.K.GetCoins(ctx, addr2).IsEqual(sdk.NewCoins(sdk.NewInt64Coin("barc", 1))))
	require.True(t, input.K.GetCoins(ctx, addr3).IsEqual(sdk.NewCoins(sdk.NewInt64Coin("barc", 2), sdk.NewInt64Coin("fooc", 5))))
}

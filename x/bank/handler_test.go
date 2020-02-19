package bank

import (
	"strings"
	"testing"

	"github.com/line/link/x/bank/internal/keeper"
	"github.com/line/link/x/bank/internal/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/stretchr/testify/require"
)

func TestInvalidMsg(t *testing.T) {
	h := NewHandler(keeper.Keeper{})

	res := h(sdk.NewContext(nil, abci.Header{}, false, nil), sdk.NewTestMsg())
	require.False(t, res.IsOK())
	require.True(t, strings.Contains(res.Log, "unrecognized bank message type"))
}

func TestHandlerSend(t *testing.T) {
	input := keeper.SetupTestInput()
	ctx, _, ak := input.Ctx, input.K, input.Ak

	h := NewHandler(input.K)

	const (
		length3Denom  = "foo"
		length5Denom  = "f2345"
		length6Denom  = "f23456"
		length5Denom2 = "f2346"
	)

	addr1 := sdk.AccAddress([]byte("addr1"))
	addr2 := sdk.AccAddress([]byte("addr1"))

	acc := ak.NewAccountWithAddress(ctx, addr1)

	err := acc.SetCoins(sdk.NewCoins(sdk.NewInt64Coin(length3Denom, 100), sdk.NewInt64Coin(length5Denom, 100)))
	require.NoError(t, err)
	ak.SetAccount(ctx, acc)

	{
		coins := sdk.NewCoins(sdk.NewInt64Coin(length3Denom, 10))
		msg := types.NewMsgSend(addr1, addr2, coins)
		res := h(ctx, msg)
		require.True(t, res.IsOK())
	}

	{
		coins := sdk.NewCoins(sdk.NewInt64Coin(length6Denom, 10))
		msg := types.NewMsgSend(addr1, addr2, coins)
		res := h(ctx, msg)
		require.False(t, res.IsOK())
	}

	{
		inputs := []Input{
			types.NewInput(addr1, sdk.NewCoins(sdk.NewInt64Coin(length3Denom, 4), sdk.NewInt64Coin(length5Denom, 2))),
			types.NewInput(addr2, sdk.NewCoins(sdk.NewInt64Coin(length3Denom, 3))),
		}

		outputs := []Output{
			types.NewOutput(addr1, sdk.NewCoins(sdk.NewInt64Coin(length3Denom, 7))),
			types.NewOutput(addr2, sdk.NewCoins(sdk.NewInt64Coin(length5Denom, 2))),
		}
		msg := types.NewMsgMultiSend(inputs, outputs)
		res := h(ctx, msg)
		require.True(t, res.IsOK())
	}

	{
		inputs := []Input{
			types.NewInput(addr1, sdk.NewCoins(sdk.NewInt64Coin(length3Denom, 4), sdk.NewInt64Coin(length5Denom, 2))),
			types.NewInput(addr2, sdk.NewCoins(sdk.NewInt64Coin(length3Denom, 3))),
		}

		outputs := []Output{
			types.NewOutput(addr1, sdk.NewCoins(sdk.NewInt64Coin(length3Denom, 7))),
			types.NewOutput(addr2, sdk.NewCoins(sdk.NewInt64Coin(length5Denom2, 2))),
		}
		msg := types.NewMsgMultiSend(inputs, outputs)
		require.Panics(t, func() { h(ctx, msg) })
	}
}

func TestHandlerSendRestricted(t *testing.T) {
	input := keeper.SetupTestInput()
	ctx, _, ak := input.Ctx, input.K, input.Ak

	h := NewHandler(input.K)

	const (
		length3Denom = "foo"
		length8Denom = "f2345678"
	)

	addr1 := sdk.AccAddress([]byte("addr1"))
	addr2 := sdk.AccAddress([]byte("addr1"))

	acc := ak.NewAccountWithAddress(ctx, addr1)

	err := acc.SetCoins(sdk.NewCoins(sdk.NewInt64Coin(length3Denom, 100), sdk.NewInt64Coin(length8Denom, 100)))
	require.NoError(t, err)
	ak.SetAccount(ctx, acc)

	{
		coins := sdk.NewCoins(sdk.NewInt64Coin(length3Denom, 10))
		msg := types.NewMsgSend(addr1, addr2, coins)
		require.NoError(t, msg.ValidateBasic())
		res := h(ctx, msg)
		require.True(t, res.IsOK())
	}

	{
		coins := sdk.NewCoins(sdk.NewInt64Coin(length8Denom, 10))
		msg := types.NewMsgSend(addr1, addr2, coins)
		require.Error(t, msg.ValidateBasic())
	}

	{
		inputs := []Input{
			types.NewInput(addr1, sdk.NewCoins(sdk.NewInt64Coin(length3Denom, 4), sdk.NewInt64Coin(length8Denom, 2))),
			types.NewInput(addr2, sdk.NewCoins(sdk.NewInt64Coin(length3Denom, 3))),
		}

		outputs := []Output{
			types.NewOutput(addr1, sdk.NewCoins(sdk.NewInt64Coin(length3Denom, 7))),
			types.NewOutput(addr2, sdk.NewCoins(sdk.NewInt64Coin(length8Denom, 2))),
		}
		msg := types.NewMsgMultiSend(inputs, outputs)
		res := h(ctx, msg)
		require.False(t, res.IsOK())
	}
}

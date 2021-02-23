package types

import (
	"testing"

	sdk "github.com/line/lbm-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestMsgs(t *testing.T) {
	const (
		length3Denom = "foo"
		length5Denom = "f2345"
		length6Denom = "f23456"
		length8Denom = "f2345678"
	)

	addr1 := sdk.AccAddress([]byte("addr1"))
	addr2 := sdk.AccAddress([]byte("addr1"))

	{
		coins := sdk.NewCoins(sdk.NewInt64Coin(length3Denom, 10))
		msg := NewMsgSend(addr1, addr2, coins)
		require.NoError(t, msg.ValidateBasic())
		require.Equal(t, addr1.String(), msg.GetSigners()[0].String())
	}
	{
		coins := sdk.NewCoins(sdk.NewInt64Coin(length3Denom, 10))
		msg := NewMsgSend(addr1, nil, coins)
		require.Error(t, msg.ValidateBasic())
	}
	{
		coins := sdk.NewCoins(sdk.NewInt64Coin(length6Denom, 10))
		msg := NewMsgSend(addr1, addr2, coins)
		require.Error(t, msg.ValidateBasic())
	}

	{
		inputs := []Input{
			NewInput(addr1, sdk.NewCoins(sdk.NewInt64Coin(length3Denom, 4), sdk.NewInt64Coin(length5Denom, 2))),
			NewInput(addr2, sdk.NewCoins(sdk.NewInt64Coin(length3Denom, 3))),
		}

		outputs := []Output{
			NewOutput(addr1, sdk.NewCoins(sdk.NewInt64Coin(length3Denom, 7))),
			NewOutput(addr2, sdk.NewCoins(sdk.NewInt64Coin(length5Denom, 2))),
		}
		msg := NewMsgMultiSend(inputs, outputs)
		require.NoError(t, msg.ValidateBasic())
		require.Equal(t, addr1.String(), msg.GetSigners()[0].String())
		require.Equal(t, addr2.String(), msg.GetSigners()[1].String())
	}
	// InputOutputMismatch
	{
		inputs := []Input{
			NewInput(addr1, sdk.NewCoins(sdk.NewInt64Coin(length3Denom, 4), sdk.NewInt64Coin(length5Denom, 2))),
			NewInput(addr2, sdk.NewCoins(sdk.NewInt64Coin(length3Denom, 3))),
		}

		outputs := []Output{
			NewOutput(addr1, sdk.NewCoins(sdk.NewInt64Coin(length3Denom, 7))),
			NewOutput(addr2, sdk.NewCoins(sdk.NewInt64Coin(length5Denom, 1))),
		}
		msg := NewMsgMultiSend(inputs, outputs)
		require.Error(t, msg.ValidateBasic())
	}
	// Validate Denom
	{
		inputs := []Input{
			NewInput(addr1, sdk.NewCoins(sdk.NewInt64Coin(length6Denom, 4), sdk.NewInt64Coin(length8Denom, 2))),
			NewInput(addr2, sdk.NewCoins(sdk.NewInt64Coin(length6Denom, 3))),
		}

		outputs := []Output{
			NewOutput(addr1, sdk.NewCoins(sdk.NewInt64Coin(length6Denom, 7))),
			NewOutput(addr2, sdk.NewCoins(sdk.NewInt64Coin(length8Denom, 2))),
		}
		msg := NewMsgMultiSend(inputs, outputs)
		require.Error(t, msg.ValidateBasic())
	}
	// NoInput or NoOutput
	{
		inputs := []Input{
			NewInput(addr1, sdk.NewCoins(sdk.NewInt64Coin(length6Denom, 4), sdk.NewInt64Coin(length8Denom, 2))),
			NewInput(addr2, sdk.NewCoins(sdk.NewInt64Coin(length6Denom, 3))),
		}

		outputs := []Output{
			NewOutput(addr1, sdk.NewCoins(sdk.NewInt64Coin(length6Denom, 7))),
			NewOutput(addr2, sdk.NewCoins(sdk.NewInt64Coin(length8Denom, 2))),
		}
		msg := NewMsgMultiSend([]Input{}, outputs)
		require.Error(t, msg.ValidateBasic())
		msg = NewMsgMultiSend(inputs, []Output{})
		require.Error(t, msg.ValidateBasic())
	}
}

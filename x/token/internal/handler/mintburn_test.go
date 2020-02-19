package handler

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/line/link/x/token/internal/types"
	"github.com/stretchr/testify/require"
)

func TestHandleMsgMint(t *testing.T) {
	ctx, h := cacheKeeper()

	t.Log("Prepare Token Issued")
	{
		token := types.NewToken(name, symbol, tokenuri, decimals, true)
		err := k.IssueToken(ctx, token, amount, addr1)
		require.NoError(t, err)
	}

	t.Log("Burn Tokens")
	{
		msg := types.NewMsgMint(addr1, addr1, sdk.NewCoins(sdk.NewInt64Coin(symbol, amount.Int64())))
		require.NoError(t, msg.ValidateBasic())
		res := h(ctx, msg)
		require.True(t, res.Code.IsOK())
		e := sdk.Events{
			sdk.NewEvent("message", sdk.NewAttribute("module", "token")),
			sdk.NewEvent("message", sdk.NewAttribute("sender", addr1.String())),
			sdk.NewEvent("mint", sdk.NewAttribute("amount", amount.String()+symbol)),
			sdk.NewEvent("mint", sdk.NewAttribute("from", addr1.String())),
			sdk.NewEvent("mint", sdk.NewAttribute("to", addr1.String())),
		}
		verifyEventFunc(t, e, res.Events)
	}
}

func TestHandleMsgBurn(t *testing.T) {
	ctx, h := cacheKeeper()
	t.Log("Prepare Token Issued")
	{
		token := types.NewToken(name, symbol, tokenuri, decimals, true)
		err := k.IssueToken(ctx, token, amount, addr1)
		require.NoError(t, err)
	}

	t.Log("Mint Tokens")
	{
		msg := types.NewMsgBurn(addr1, sdk.NewCoins(sdk.NewInt64Coin(symbol, amount.Int64())))
		require.NoError(t, msg.ValidateBasic())
		res := h(ctx, msg)
		require.True(t, res.Code.IsOK())
		e := sdk.Events{
			sdk.NewEvent("message", sdk.NewAttribute("module", "token")),
			sdk.NewEvent("message", sdk.NewAttribute("sender", addr1.String())),
			sdk.NewEvent("burn", sdk.NewAttribute("amount", amount.String()+symbol)),
			sdk.NewEvent("burn", sdk.NewAttribute("from", addr1.String())),
		}
		verifyEventFunc(t, e, res.Events)
	}
}

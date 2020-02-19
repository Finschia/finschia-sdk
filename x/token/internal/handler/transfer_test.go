package handler

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/line/link/x/token/internal/types"
	"github.com/stretchr/testify/require"
)

func TestHandleMsgTransfer(t *testing.T) {
	ctx, h := cacheKeeper()

	t.Log("Prepare Token Issued")
	{
		token := types.NewToken(name, symbol, tokenuri, decimals, true)
		err := k.IssueToken(ctx, token, amount, addr1)
		require.NoError(t, err)
	}

	t.Log("Transfer Token")
	{
		msg := types.NewMsgTransfer(addr1, addr2, symbol, sdk.NewInt(10))
		res := h(ctx, msg)
		require.True(t, res.Code.IsOK())
		e := sdk.Events{
			sdk.NewEvent("message", sdk.NewAttribute("module", "token")),
			sdk.NewEvent("message", sdk.NewAttribute("sender", addr1.String())),
			sdk.NewEvent("transfer", sdk.NewAttribute("from", addr1.String())),
			sdk.NewEvent("transfer", sdk.NewAttribute("to", addr2.String())),
			sdk.NewEvent("transfer", sdk.NewAttribute("symbol", symbol)),
			sdk.NewEvent("transfer", sdk.NewAttribute("amount", sdk.NewInt(10).String())),
		}
		verifyEventFunc(t, e, res.Events)
	}
	t.Log("Transfer Coin. Expect Fail")
	{
		msg := types.NewMsgTransfer(addr1, addr2, coinSymbol, sdk.NewInt(10))
		require.EqualError(t, msg.ValidateBasic(), sdk.ErrInvalidCoins("Only user defined token is possible: link").Error())
		res := h(ctx, msg)
		require.False(t, res.Code.IsOK())
	}
}

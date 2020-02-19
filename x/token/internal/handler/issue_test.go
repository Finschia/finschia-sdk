package handler

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/line/link/x/token/internal/types"
	"github.com/stretchr/testify/require"
)

func TestHandleMsgIssue(t *testing.T) {
	ctx, h := cacheKeeper()

	t.Log("Issue Token")
	{
		msg := types.NewMsgIssue(addr1, name, symbol, tokenuri, amount, decimals, true)
		require.NoError(t, msg.ValidateBasic())
		res := h(ctx, msg)
		require.True(t, res.Code.IsOK())

		e := sdk.Events{
			sdk.NewEvent("message", sdk.NewAttribute("module", "token")),
			sdk.NewEvent("message", sdk.NewAttribute("sender", addr1.String())),
			sdk.NewEvent("grant_perm", sdk.NewAttribute("to", addr1.String())),
			sdk.NewEvent("grant_perm", sdk.NewAttribute("perm_resource", symbol)),
			sdk.NewEvent("grant_perm", sdk.NewAttribute("perm_action", types.ModifyAction)),
			sdk.NewEvent("issue", sdk.NewAttribute("name", name)),
			sdk.NewEvent("issue", sdk.NewAttribute("symbol", symbol)),
			sdk.NewEvent("issue", sdk.NewAttribute("owner", addr1.String())),
			sdk.NewEvent("issue", sdk.NewAttribute("amount", amount.String())),
			sdk.NewEvent("issue", sdk.NewAttribute("mintable", "true")),
			sdk.NewEvent("issue", sdk.NewAttribute("decimals", decimals.String())),
			sdk.NewEvent("grant_perm", sdk.NewAttribute("to", addr1.String())),
			sdk.NewEvent("grant_perm", sdk.NewAttribute("perm_resource", symbol)),
			sdk.NewEvent("grant_perm", sdk.NewAttribute("perm_action", "mint")),
			sdk.NewEvent("grant_perm", sdk.NewAttribute("to", addr1.String())),
			sdk.NewEvent("grant_perm", sdk.NewAttribute("perm_resource", symbol)),
			sdk.NewEvent("grant_perm", sdk.NewAttribute("perm_action", "burn")),
		}
		verifyEventFunc(t, e, res.Events)
	}

	t.Log("Issue Token Again Expect Fail")
	{
		msg := types.NewMsgIssue(addr1, name, symbol, tokenuri, amount, decimals, true)
		res := h(ctx, msg)
		require.False(t, res.Code.IsOK())
		require.Equal(t, types.DefaultCodespace, res.Codespace)
		require.Equal(t, types.CodeTokenExist, res.Code)
	}
}

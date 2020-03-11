package handler

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/line/link/x/token/internal/types"
	"github.com/stretchr/testify/require"
)

func TestHandleMsgGrant(t *testing.T) {
	ctx, h := cacheKeeper()

	t.Log("Prepare Token Issued")
	{
		token := types.NewToken(defaultContractID, defaultName, defaultSymbol, defaultMeta, defaultImageURI, sdk.NewInt(defaultDecimals), true)
		err := k.IssueToken(ctx, token, sdk.NewInt(defaultAmount), addr1, addr1)
		require.NoError(t, err)
	}

	permission := types.NewMintPermission(defaultContractID)
	t.Log("Grant Permission")
	{
		msg := types.NewMsgGrantPermission(addr1, addr2, permission)
		require.NoError(t, msg.ValidateBasic())
		res, err := h(ctx, msg)
		require.NoError(t, err)

		e := sdk.Events{
			sdk.NewEvent("message", sdk.NewAttribute("module", "token")),
			sdk.NewEvent("message", sdk.NewAttribute("sender", addr1.String())),
			sdk.NewEvent("grant_perm", sdk.NewAttribute("from", addr1.String())),
			sdk.NewEvent("grant_perm", sdk.NewAttribute("to", addr2.String())),
			sdk.NewEvent("grant_perm", sdk.NewAttribute("perm_resource", permission.GetResource())),
			sdk.NewEvent("grant_perm", sdk.NewAttribute("perm_action", permission.GetAction())),
		}
		verifyEventFunc(t, e, res.Events)
	}
}

func TestHandleMsgRevoke(t *testing.T) {
	ctx, h := cacheKeeper()

	t.Log("Prepare Token Issued")
	{
		token := types.NewToken(defaultContractID, defaultName, defaultSymbol, defaultMeta, defaultImageURI, sdk.NewInt(defaultDecimals), true)
		err := k.IssueToken(ctx, token, sdk.NewInt(defaultAmount), addr1, addr1)
		require.NoError(t, err)
	}

	permission := types.NewMintPermission(defaultContractID)
	t.Log("Revoke Permission")
	{
		msg := types.NewMsgRevokePermission(addr1, permission)
		require.NoError(t, msg.ValidateBasic())
		res, err := h(ctx, msg)
		require.NoError(t, err)
		e := sdk.Events{
			sdk.NewEvent("message", sdk.NewAttribute("module", "token")),
			sdk.NewEvent("message", sdk.NewAttribute("sender", addr1.String())),
			sdk.NewEvent("revoke_perm", sdk.NewAttribute("from", addr1.String())),
			sdk.NewEvent("revoke_perm", sdk.NewAttribute("perm_resource", permission.GetResource())),
			sdk.NewEvent("revoke_perm", sdk.NewAttribute("perm_action", permission.GetAction())),
		}
		verifyEventFunc(t, e, res.Events)
	}
}

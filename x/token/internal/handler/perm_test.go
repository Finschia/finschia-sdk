package handler

import (
	"testing"

	sdk "github.com/line/lbm-sdk/types"
	"github.com/line/lbm-sdk/x/token/internal/types"
	"github.com/stretchr/testify/require"
)

func TestHandleMsgGrant(t *testing.T) {
	ctx, h := cacheKeeper()

	t.Log("Prepare Token Issued")
	{
		k.NewContractID(ctx)
		token := types.NewToken(defaultContractID, defaultName, defaultSymbol, defaultMeta, defaultImageURI, sdk.NewInt(defaultDecimals), true)
		err := k.IssueToken(ctx, token, sdk.NewInt(defaultAmount), addr1, addr1)
		require.NoError(t, err)
	}

	permission := types.NewMintPermission()
	t.Log("Invalid contract id")
	{
		msg := types.NewMsgGrantPermission(addr1, "1234567890", addr2, permission)
		require.Error(t, msg.ValidateBasic())
	}
	t.Log("Grant Permission")
	{
		msg := types.NewMsgGrantPermission(addr1, defaultContractID, addr2, permission)
		require.NoError(t, msg.ValidateBasic())
		res, err := h(ctx, msg)
		require.NoError(t, err)

		e := sdk.Events{
			sdk.NewEvent("message", sdk.NewAttribute("module", "token")),
			sdk.NewEvent("message", sdk.NewAttribute("sender", addr1.String())),
			sdk.NewEvent("grant_perm", sdk.NewAttribute("from", addr1.String())),
			sdk.NewEvent("grant_perm", sdk.NewAttribute("to", addr2.String())),
			sdk.NewEvent("grant_perm", sdk.NewAttribute("contract_id", defaultContractID)),
			sdk.NewEvent("grant_perm", sdk.NewAttribute("perm", permission.String())),
		}
		verifyEventFunc(t, e, res.Events)
	}
}

func TestHandleMsgRevoke(t *testing.T) {
	ctx, h := cacheKeeper()

	t.Log("Prepare Token Issued")
	{
		k.NewContractID(ctx)
		token := types.NewToken(defaultContractID, defaultName, defaultSymbol, defaultMeta, defaultImageURI, sdk.NewInt(defaultDecimals), true)
		err := k.IssueToken(ctx, token, sdk.NewInt(defaultAmount), addr1, addr1)
		require.NoError(t, err)
	}

	permission := types.NewMintPermission()
	t.Log("Invalid contract id")
	{
		msg := types.NewMsgRevokePermission(addr1, "1234567890", permission)
		require.Error(t, msg.ValidateBasic())
	}
	t.Log("Revoke Permission")
	{
		msg := types.NewMsgRevokePermission(addr1, defaultContractID, permission)
		require.NoError(t, msg.ValidateBasic())
		res, err := h(ctx, msg)
		require.NoError(t, err)
		e := sdk.Events{
			sdk.NewEvent("message", sdk.NewAttribute("module", "token")),
			sdk.NewEvent("message", sdk.NewAttribute("sender", addr1.String())),
			sdk.NewEvent("revoke_perm", sdk.NewAttribute("from", addr1.String())),
			sdk.NewEvent("revoke_perm", sdk.NewAttribute("contract_id", defaultContractID)),
			sdk.NewEvent("revoke_perm", sdk.NewAttribute("perm", permission.String())),
		}
		verifyEventFunc(t, e, res.Events)
	}
}

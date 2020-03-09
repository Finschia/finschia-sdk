package handler

import (
	"testing"

	"github.com/line/link/x/collection/internal/types"
	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestHandleMsgGrant(t *testing.T) {
	ctx, h, contractID := prepareCreateCollection()
	{
		msg := types.NewMsgGrantPermission(addr1, addr2, types.NewIssuePermission(contractID))
		res := h(ctx, msg)
		require.True(t, res.Code.IsOK())

		e := sdk.Events{
			sdk.NewEvent("message", sdk.NewAttribute("module", "collection")),
			sdk.NewEvent("message", sdk.NewAttribute("sender", addr1.String())),
			sdk.NewEvent("grant_perm", sdk.NewAttribute("from", addr1.String())),
			sdk.NewEvent("grant_perm", sdk.NewAttribute("to", addr2.String())),
			sdk.NewEvent("grant_perm", sdk.NewAttribute("perm_resource", contractID)),
			sdk.NewEvent("grant_perm", sdk.NewAttribute("perm_action", "issue")),
		}
		verifyEventFunc(t, e, res.Events)
	}
	{
		msg := types.NewMsgGrantPermission(addr1, addr2, types.NewMintPermission(contractID))
		res := h(ctx, msg)
		require.True(t, res.Code.IsOK())

		e := sdk.Events{
			sdk.NewEvent("message", sdk.NewAttribute("module", "collection")),
			sdk.NewEvent("message", sdk.NewAttribute("sender", addr1.String())),
			sdk.NewEvent("grant_perm", sdk.NewAttribute("from", addr1.String())),
			sdk.NewEvent("grant_perm", sdk.NewAttribute("to", addr2.String())),
			sdk.NewEvent("grant_perm", sdk.NewAttribute("perm_resource", contractID)),
			sdk.NewEvent("grant_perm", sdk.NewAttribute("perm_action", "mint")),
		}
		verifyEventFunc(t, e, res.Events)
	}
	{
		msg := types.NewMsgGrantPermission(addr1, addr2, types.NewBurnPermission(contractID))
		res := h(ctx, msg)
		require.True(t, res.Code.IsOK())

		e := sdk.Events{
			sdk.NewEvent("message", sdk.NewAttribute("module", "collection")),
			sdk.NewEvent("message", sdk.NewAttribute("sender", addr1.String())),
			sdk.NewEvent("grant_perm", sdk.NewAttribute("from", addr1.String())),
			sdk.NewEvent("grant_perm", sdk.NewAttribute("to", addr2.String())),
			sdk.NewEvent("grant_perm", sdk.NewAttribute("perm_resource", contractID)),
			sdk.NewEvent("grant_perm", sdk.NewAttribute("perm_action", "burn")),
		}
		verifyEventFunc(t, e, res.Events)
	}
	{
		msg := types.NewMsgGrantPermission(addr1, addr2, types.NewModifyPermission(contractID))
		res := h(ctx, msg)
		require.True(t, res.Code.IsOK())

		e := sdk.Events{
			sdk.NewEvent("message", sdk.NewAttribute("module", "collection")),
			sdk.NewEvent("message", sdk.NewAttribute("sender", addr1.String())),
			sdk.NewEvent("grant_perm", sdk.NewAttribute("from", addr1.String())),
			sdk.NewEvent("grant_perm", sdk.NewAttribute("to", addr2.String())),
			sdk.NewEvent("grant_perm", sdk.NewAttribute("perm_resource", contractID)),
			sdk.NewEvent("grant_perm", sdk.NewAttribute("perm_action", "modify")),
		}
		verifyEventFunc(t, e, res.Events)
	}
}

func TestHandleMsgRevoke(t *testing.T) {
	ctx, h, contractID := prepareCreateCollection()
	msg := types.NewMsgGrantPermission(addr1, addr2, types.NewIssuePermission(contractID))
	_ = h(ctx, msg)
	msg = types.NewMsgGrantPermission(addr1, addr2, types.NewMintPermission(contractID))
	_ = h(ctx, msg)
	msg = types.NewMsgGrantPermission(addr1, addr2, types.NewBurnPermission(contractID))
	_ = h(ctx, msg)
	msg = types.NewMsgGrantPermission(addr1, addr2, types.NewModifyPermission(contractID))
	_ = h(ctx, msg)

	{
		msg := types.NewMsgRevokePermission(addr2, types.NewIssuePermission(contractID))
		res := h(ctx, msg)
		require.True(t, res.Code.IsOK())

		e := sdk.Events{
			sdk.NewEvent("message", sdk.NewAttribute("module", "collection")),
			sdk.NewEvent("message", sdk.NewAttribute("sender", addr2.String())),
			sdk.NewEvent("revoke_perm", sdk.NewAttribute("from", addr2.String())),
			sdk.NewEvent("revoke_perm", sdk.NewAttribute("perm_resource", contractID)),
			sdk.NewEvent("revoke_perm", sdk.NewAttribute("perm_action", "issue")),
		}
		verifyEventFunc(t, e, res.Events)
	}
	{
		msg := types.NewMsgRevokePermission(addr2, types.NewMintPermission(contractID))
		res := h(ctx, msg)
		require.True(t, res.Code.IsOK())

		e := sdk.Events{
			sdk.NewEvent("message", sdk.NewAttribute("module", "collection")),
			sdk.NewEvent("message", sdk.NewAttribute("sender", addr2.String())),
			sdk.NewEvent("revoke_perm", sdk.NewAttribute("from", addr2.String())),
			sdk.NewEvent("revoke_perm", sdk.NewAttribute("perm_resource", contractID)),
			sdk.NewEvent("revoke_perm", sdk.NewAttribute("perm_action", "mint")),
		}
		verifyEventFunc(t, e, res.Events)
	}
	{
		msg := types.NewMsgRevokePermission(addr2, types.NewBurnPermission(contractID))
		res := h(ctx, msg)
		require.True(t, res.Code.IsOK())

		e := sdk.Events{
			sdk.NewEvent("message", sdk.NewAttribute("module", "collection")),
			sdk.NewEvent("message", sdk.NewAttribute("sender", addr2.String())),
			sdk.NewEvent("revoke_perm", sdk.NewAttribute("from", addr2.String())),
			sdk.NewEvent("revoke_perm", sdk.NewAttribute("perm_resource", contractID)),
			sdk.NewEvent("revoke_perm", sdk.NewAttribute("perm_action", "burn")),
		}
		verifyEventFunc(t, e, res.Events)
	}
	{
		msg := types.NewMsgRevokePermission(addr2, types.NewModifyPermission(contractID))
		res := h(ctx, msg)
		require.True(t, res.Code.IsOK())

		e := sdk.Events{
			sdk.NewEvent("message", sdk.NewAttribute("module", "collection")),
			sdk.NewEvent("message", sdk.NewAttribute("sender", addr2.String())),
			sdk.NewEvent("revoke_perm", sdk.NewAttribute("from", addr2.String())),
			sdk.NewEvent("revoke_perm", sdk.NewAttribute("perm_resource", contractID)),
			sdk.NewEvent("revoke_perm", sdk.NewAttribute("perm_action", "modify")),
		}
		verifyEventFunc(t, e, res.Events)
	}
}

package handler

import (
	"testing"

	"github.com/line/link/x/collection/internal/types"
	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestHandleMsgGrant(t *testing.T) {
	ctx, h, contractID := prepareCreateCollection(t)
	{
		msg := types.NewMsgGrantPermission(addr1, addr2, types.NewIssuePermission(contractID))
		res, err := h(ctx, msg)
		require.NoError(t, err)

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
		res, err := h(ctx, msg)
		require.NoError(t, err)

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
		res, err := h(ctx, msg)
		require.NoError(t, err)

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
		res, err := h(ctx, msg)
		require.NoError(t, err)

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
	ctx, h, contractID := prepareCreateCollection(t)
	msg := types.NewMsgGrantPermission(addr1, addr2, types.NewIssuePermission(contractID))
	_, err := h(ctx, msg)
	require.NoError(t, err)

	msg = types.NewMsgGrantPermission(addr1, addr2, types.NewMintPermission(contractID))
	_, err = h(ctx, msg)
	require.NoError(t, err)

	msg = types.NewMsgGrantPermission(addr1, addr2, types.NewBurnPermission(contractID))
	_, err = h(ctx, msg)
	require.NoError(t, err)

	msg = types.NewMsgGrantPermission(addr1, addr2, types.NewModifyPermission(contractID))
	_, err = h(ctx, msg)
	require.NoError(t, err)

	{
		msg := types.NewMsgRevokePermission(addr2, types.NewIssuePermission(contractID))
		res, err := h(ctx, msg)
		require.NoError(t, err)

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
		res, err := h(ctx, msg)
		require.NoError(t, err)

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
		res, err := h(ctx, msg)
		require.NoError(t, err)

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
		res, err := h(ctx, msg)
		require.NoError(t, err)

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

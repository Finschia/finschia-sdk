package handler

import (
	"testing"

	"github.com/line/link-modules/x/collection/internal/types"
	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestHandleMsgGrant(t *testing.T) {
	ctx, h, contractID := prepareCreateCollection(t)
	{
		msg := types.NewMsgGrantPermission(addr1, contractID, addr2, types.NewIssuePermission())
		res, err := h(ctx, msg)
		require.NoError(t, err)

		e := sdk.Events{
			sdk.NewEvent("message", sdk.NewAttribute("module", "collection")),
			sdk.NewEvent("message", sdk.NewAttribute("sender", addr1.String())),
			sdk.NewEvent("grant_perm", sdk.NewAttribute("from", addr1.String())),
			sdk.NewEvent("grant_perm", sdk.NewAttribute("to", addr2.String())),
			sdk.NewEvent("grant_perm", sdk.NewAttribute("contract_id", contractID)),
			sdk.NewEvent("grant_perm", sdk.NewAttribute("perm", "issue")),
		}
		verifyEventFunc(t, e, res.Events)
	}
	{
		msg := types.NewMsgGrantPermission(addr1, contractID, addr2, types.NewMintPermission())
		res, err := h(ctx, msg)
		require.NoError(t, err)

		e := sdk.Events{
			sdk.NewEvent("message", sdk.NewAttribute("module", "collection")),
			sdk.NewEvent("message", sdk.NewAttribute("sender", addr1.String())),
			sdk.NewEvent("grant_perm", sdk.NewAttribute("from", addr1.String())),
			sdk.NewEvent("grant_perm", sdk.NewAttribute("to", addr2.String())),
			sdk.NewEvent("grant_perm", sdk.NewAttribute("contract_id", contractID)),
			sdk.NewEvent("grant_perm", sdk.NewAttribute("perm", "mint")),
		}
		verifyEventFunc(t, e, res.Events)
	}
	{
		msg := types.NewMsgGrantPermission(addr1, contractID, addr2, types.NewBurnPermission())
		res, err := h(ctx, msg)
		require.NoError(t, err)

		e := sdk.Events{
			sdk.NewEvent("message", sdk.NewAttribute("module", "collection")),
			sdk.NewEvent("message", sdk.NewAttribute("sender", addr1.String())),
			sdk.NewEvent("grant_perm", sdk.NewAttribute("from", addr1.String())),
			sdk.NewEvent("grant_perm", sdk.NewAttribute("to", addr2.String())),
			sdk.NewEvent("grant_perm", sdk.NewAttribute("contract_id", contractID)),
			sdk.NewEvent("grant_perm", sdk.NewAttribute("perm", "burn")),
		}
		verifyEventFunc(t, e, res.Events)
	}
	{
		msg := types.NewMsgGrantPermission(addr1, contractID, addr2, types.NewModifyPermission())
		res, err := h(ctx, msg)
		require.NoError(t, err)

		e := sdk.Events{
			sdk.NewEvent("message", sdk.NewAttribute("module", "collection")),
			sdk.NewEvent("message", sdk.NewAttribute("sender", addr1.String())),
			sdk.NewEvent("grant_perm", sdk.NewAttribute("from", addr1.String())),
			sdk.NewEvent("grant_perm", sdk.NewAttribute("to", addr2.String())),
			sdk.NewEvent("grant_perm", sdk.NewAttribute("contract_id", contractID)),
			sdk.NewEvent("grant_perm", sdk.NewAttribute("perm", "modify")),
		}
		verifyEventFunc(t, e, res.Events)
	}
	t.Log("Invalid contract id")
	{
		msg := types.NewMsgGrantPermission(addr1, "1234567890", addr2, types.NewModifyPermission())
		require.Error(t, msg.ValidateBasic())
	}
}

func TestHandleMsgRevoke(t *testing.T) {
	ctx, h, contractID := prepareCreateCollection(t)
	msg := types.NewMsgGrantPermission(addr1, contractID, addr2, types.NewIssuePermission())
	_, err := h(ctx, msg)
	require.NoError(t, err)

	msg = types.NewMsgGrantPermission(addr1, contractID, addr2, types.NewMintPermission())
	_, err = h(ctx, msg)
	require.NoError(t, err)

	msg = types.NewMsgGrantPermission(addr1, contractID, addr2, types.NewBurnPermission())
	_, err = h(ctx, msg)
	require.NoError(t, err)

	msg = types.NewMsgGrantPermission(addr1, contractID, addr2, types.NewModifyPermission())
	_, err = h(ctx, msg)
	require.NoError(t, err)

	{
		msg := types.NewMsgRevokePermission(addr2, contractID, types.NewIssuePermission())
		res, err := h(ctx, msg)
		require.NoError(t, err)

		e := sdk.Events{
			sdk.NewEvent("message", sdk.NewAttribute("module", "collection")),
			sdk.NewEvent("message", sdk.NewAttribute("sender", addr2.String())),
			sdk.NewEvent("revoke_perm", sdk.NewAttribute("from", addr2.String())),
			sdk.NewEvent("revoke_perm", sdk.NewAttribute("contract_id", contractID)),
			sdk.NewEvent("revoke_perm", sdk.NewAttribute("perm", "issue")),
		}
		verifyEventFunc(t, e, res.Events)
	}
	{
		msg := types.NewMsgRevokePermission(addr2, contractID, types.NewMintPermission())
		res, err := h(ctx, msg)
		require.NoError(t, err)

		e := sdk.Events{
			sdk.NewEvent("message", sdk.NewAttribute("module", "collection")),
			sdk.NewEvent("message", sdk.NewAttribute("sender", addr2.String())),
			sdk.NewEvent("revoke_perm", sdk.NewAttribute("from", addr2.String())),
			sdk.NewEvent("revoke_perm", sdk.NewAttribute("contract_id", contractID)),
			sdk.NewEvent("revoke_perm", sdk.NewAttribute("perm", "mint")),
		}
		verifyEventFunc(t, e, res.Events)
	}
	{
		msg := types.NewMsgRevokePermission(addr2, contractID, types.NewBurnPermission())
		res, err := h(ctx, msg)
		require.NoError(t, err)

		e := sdk.Events{
			sdk.NewEvent("message", sdk.NewAttribute("module", "collection")),
			sdk.NewEvent("message", sdk.NewAttribute("sender", addr2.String())),
			sdk.NewEvent("revoke_perm", sdk.NewAttribute("from", addr2.String())),
			sdk.NewEvent("revoke_perm", sdk.NewAttribute("contract_id", contractID)),
			sdk.NewEvent("revoke_perm", sdk.NewAttribute("perm", "burn")),
		}
		verifyEventFunc(t, e, res.Events)
	}
	{
		msg := types.NewMsgRevokePermission(addr2, contractID, types.NewModifyPermission())
		res, err := h(ctx, msg)
		require.NoError(t, err)

		e := sdk.Events{
			sdk.NewEvent("message", sdk.NewAttribute("module", "collection")),
			sdk.NewEvent("message", sdk.NewAttribute("sender", addr2.String())),
			sdk.NewEvent("revoke_perm", sdk.NewAttribute("from", addr2.String())),
			sdk.NewEvent("revoke_perm", sdk.NewAttribute("contract_id", contractID)),
			sdk.NewEvent("revoke_perm", sdk.NewAttribute("perm", "modify")),
		}
		verifyEventFunc(t, e, res.Events)
	}
	t.Log("Invalid contract id")
	{
		msg := types.NewMsgRevokePermission(addr1, "1234567890", types.NewModifyPermission())
		require.Error(t, msg.ValidateBasic())
	}
}

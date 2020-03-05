package handler

import (
	"testing"

	"github.com/line/link/x/collection/internal/types"
	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestHandleAttach(t *testing.T) {
	t.Log("implement me - ", t.Name())
}

func TestHandleAttachFrom(t *testing.T) {
	t.Log("implement me - ", t.Name())
}

func TestHandleDetach(t *testing.T) {
	t.Log("implement me - ", t.Name())
}

func TestHandleDetachFrom(t *testing.T) {
	t.Log("implement me - ", t.Name())
}

func attach(t *testing.T, ctx sdk.Context, h sdk.Handler, contractID string) {
	msg := types.NewMsgAttach(addr1, contractID, defaultTokenID1, defaultTokenID2)
	res := h(ctx, msg)
	require.True(t, res.Code.IsOK())
	e := sdk.Events{
		sdk.NewEvent("message", sdk.NewAttribute("module", "collection")),
		sdk.NewEvent("message", sdk.NewAttribute("sender", addr1.String())),
		sdk.NewEvent("attach", sdk.NewAttribute("contract_id", contractID)),
		sdk.NewEvent("attach", sdk.NewAttribute("from", addr1.String())),
		sdk.NewEvent("attach", sdk.NewAttribute("to_token_id", defaultTokenID1)),
		sdk.NewEvent("attach", sdk.NewAttribute("token_id", defaultTokenID2)),
		sdk.NewEvent("attach", sdk.NewAttribute("old_root_token_id", defaultTokenID2)),
		sdk.NewEvent("attach", sdk.NewAttribute("new_root_token_id", defaultTokenID1)),
		sdk.NewEvent("operation_root_changed", sdk.NewAttribute("token_id", defaultTokenID2)),
	}
	verifyEventFunc(t, e, res.Events)
}

func TestHandleAttachDetach(t *testing.T) {
	ctx, h := cacheKeeper()

	var contractID string
	{
		createMsg := types.NewMsgCreateCollection(addr1, defaultName, defaultImgURI)
		res := h(ctx, createMsg)
		require.True(t, res.Code.IsOK())
		contractID = GetMadeContractID(res.Events)

		msg := types.NewMsgIssueNFT(addr1, contractID, defaultName)
		res = h(ctx, msg)
		require.True(t, res.Code.IsOK())
		msg2 := types.NewMsgMintNFT(addr1, contractID, addr1, defaultName, defaultTokenType)
		res = h(ctx, msg2)
		require.True(t, res.Code.IsOK())
		msg2 = types.NewMsgMintNFT(addr1, contractID, addr1, defaultName, defaultTokenType)
		res = h(ctx, msg2)
		require.True(t, res.Code.IsOK())
	}

	attach(t, ctx, h, contractID)

	{
		msg2 := types.NewMsgDetach(addr1, contractID, defaultTokenID2)
		res2 := h(ctx, msg2)
		require.True(t, res2.Code.IsOK())
		e := sdk.Events{
			sdk.NewEvent("message", sdk.NewAttribute("module", "collection")),
			sdk.NewEvent("message", sdk.NewAttribute("sender", addr1.String())),
			sdk.NewEvent("detach", sdk.NewAttribute("contract_id", contractID)),
			sdk.NewEvent("detach", sdk.NewAttribute("from", addr1.String())),
			sdk.NewEvent("detach", sdk.NewAttribute("from_token_id", defaultTokenID1)),
			sdk.NewEvent("detach", sdk.NewAttribute("token_id", defaultTokenID2)),
			sdk.NewEvent("detach", sdk.NewAttribute("old_root_token_id", defaultTokenID1)),
			sdk.NewEvent("detach", sdk.NewAttribute("new_root_token_id", defaultTokenID2)),
			sdk.NewEvent("operation_root_changed", sdk.NewAttribute("token_id", defaultTokenID2)),
		}
		verifyEventFunc(t, e, res2.Events)
	}

	//Attach again
	attach(t, ctx, h, contractID)

	//Burn token
	{
		msg := types.NewMsgBurnNFT(addr1, contractID, defaultTokenID1)
		res := h(ctx, msg)
		require.True(t, res.Code.IsOK())
		e := sdk.Events{
			sdk.NewEvent("message", sdk.NewAttribute("module", "collection")),
			sdk.NewEvent("message", sdk.NewAttribute("sender", addr1.String())),
			sdk.NewEvent("burn_nft", sdk.NewAttribute("contract_id", contractID)),
			sdk.NewEvent("burn_nft", sdk.NewAttribute("from", addr1.String())),
			sdk.NewEvent("burn_nft", sdk.NewAttribute("token_id", defaultTokenID1)),
			sdk.NewEvent("operation_burn_nft", sdk.NewAttribute("token_id", defaultTokenID1)),
			sdk.NewEvent("operation_burn_nft", sdk.NewAttribute("token_id", defaultTokenID2)),
		}
		verifyEventFunc(t, e, res.Events)
	}
}

func TestHandleAttachFromDetachFrom(t *testing.T) {
	ctx, h := cacheKeeper()

	var contractID string
	{
		createMsg := types.NewMsgCreateCollection(addr1, defaultName, defaultImgURI)
		res := h(ctx, createMsg)
		require.True(t, res.Code.IsOK())
		contractID = GetMadeContractID(res.Events)

		msg := types.NewMsgIssueNFT(addr1, contractID, defaultName)
		res = h(ctx, msg)
		require.True(t, res.Code.IsOK())
		msg2 := types.NewMsgMintNFT(addr1, contractID, addr1, defaultName, defaultTokenType)
		res = h(ctx, msg2)
		require.True(t, res.Code.IsOK())
		msg2 = types.NewMsgMintNFT(addr1, contractID, addr1, defaultName, defaultTokenType)
		res = h(ctx, msg2)
		require.True(t, res.Code.IsOK())
		msg3 := types.NewMsgApprove(addr1, contractID, addr2)
		res = h(ctx, msg3)
		require.True(t, res.Code.IsOK())
	}

	msg := types.NewMsgAttachFrom(addr2, contractID, addr1, defaultTokenID1, defaultTokenID2)
	res := h(ctx, msg)
	require.True(t, res.Code.IsOK())
	e := sdk.Events{
		sdk.NewEvent("message", sdk.NewAttribute("module", "collection")),
		sdk.NewEvent("message", sdk.NewAttribute("sender", addr2.String())),
		sdk.NewEvent("attach_from", sdk.NewAttribute("contract_id", contractID)),
		sdk.NewEvent("attach_from", sdk.NewAttribute("proxy", addr2.String())),
		sdk.NewEvent("attach_from", sdk.NewAttribute("from", addr1.String())),
		sdk.NewEvent("attach_from", sdk.NewAttribute("to_token_id", defaultTokenID1)),
		sdk.NewEvent("attach_from", sdk.NewAttribute("token_id", defaultTokenID2)),
		sdk.NewEvent("attach_from", sdk.NewAttribute("old_root_token_id", defaultTokenID2)),
		sdk.NewEvent("attach_from", sdk.NewAttribute("new_root_token_id", defaultTokenID1)),
		sdk.NewEvent("operation_root_changed", sdk.NewAttribute("token_id", defaultTokenID2)),
	}
	verifyEventFunc(t, e, res.Events)

	msg2 := types.NewMsgDetachFrom(addr2, contractID, addr1, defaultTokenID2)
	res2 := h(ctx, msg2)
	require.True(t, res2.Code.IsOK())
	e = sdk.Events{
		sdk.NewEvent("message", sdk.NewAttribute("module", "collection")),
		sdk.NewEvent("message", sdk.NewAttribute("sender", addr2.String())),
		sdk.NewEvent("detach_from", sdk.NewAttribute("contract_id", contractID)),
		sdk.NewEvent("detach_from", sdk.NewAttribute("proxy", addr2.String())),
		sdk.NewEvent("detach_from", sdk.NewAttribute("from", addr1.String())),
		sdk.NewEvent("detach_from", sdk.NewAttribute("from_token_id", defaultTokenID1)),
		sdk.NewEvent("detach_from", sdk.NewAttribute("token_id", defaultTokenID2)),
		sdk.NewEvent("detach_from", sdk.NewAttribute("old_root_token_id", defaultTokenID1)),
		sdk.NewEvent("detach_from", sdk.NewAttribute("new_root_token_id", defaultTokenID2)),
		sdk.NewEvent("operation_root_changed", sdk.NewAttribute("token_id", defaultTokenID2)),
	}
	verifyEventFunc(t, e, res2.Events)
}

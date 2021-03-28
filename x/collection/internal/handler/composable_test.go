package handler

import (
	"testing"

	"github.com/line/lbm-sdk/v2/x/collection/internal/types"
	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestHandleAttach(t *testing.T) {
	ctx, h, contractID := prepareNFT(t, addr1)

	{
		msg := types.NewMsgAttach(addr1, contractID, defaultTokenID1, defaultTokenID2)
		res, err := h(ctx, msg)
		require.NoError(t, err)
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
}

func TestHandleAttachFrom(t *testing.T) {
	ctx, h, contractID := prepareNFT(t, addr2)
	approve(t, addr2, addr1, contractID, ctx, h)
	{
		msg := types.NewMsgAttachFrom(addr1, contractID, addr2, defaultTokenID1, defaultTokenID2)
		res, err := h(ctx, msg)
		require.NoError(t, err)
		e := sdk.Events{
			sdk.NewEvent("message", sdk.NewAttribute("module", "collection")),
			sdk.NewEvent("message", sdk.NewAttribute("sender", addr1.String())),
			sdk.NewEvent("attach_from", sdk.NewAttribute("contract_id", contractID)),
			sdk.NewEvent("attach_from", sdk.NewAttribute("proxy", addr1.String())),
			sdk.NewEvent("attach_from", sdk.NewAttribute("from", addr2.String())),
			sdk.NewEvent("attach_from", sdk.NewAttribute("to_token_id", defaultTokenID1)),
			sdk.NewEvent("attach_from", sdk.NewAttribute("token_id", defaultTokenID2)),
			sdk.NewEvent("attach_from", sdk.NewAttribute("old_root_token_id", defaultTokenID2)),
			sdk.NewEvent("attach_from", sdk.NewAttribute("new_root_token_id", defaultTokenID1)),
			sdk.NewEvent("operation_root_changed", sdk.NewAttribute("token_id", defaultTokenID2)),
		}
		verifyEventFunc(t, e, res.Events)
	}
}

func prepareForDetaching(t *testing.T, mintTo sdk.AccAddress) (sdk.Context, sdk.Handler, string) {
	ctx, h, contractID := prepareNFT(t, mintTo)

	msg := types.NewMsgAttach(mintTo, contractID, defaultTokenID1, defaultTokenID2)
	_, err := h(ctx, msg)
	require.NoError(t, err)
	return ctx, h, contractID
}

func TestHandleDetach(t *testing.T) {
	ctx, h, contractID := prepareForDetaching(t, addr1)

	{
		msg := types.NewMsgDetach(addr1, contractID, defaultTokenID2)
		res, err := h(ctx, msg)
		require.NoError(t, err)
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
		verifyEventFunc(t, e, res.Events)
	}
}

func TestHandleDetachFrom(t *testing.T) {
	ctx, h, contractID := prepareForDetaching(t, addr2)
	approve(t, addr2, addr1, contractID, ctx, h)
	{
		msg := types.NewMsgDetachFrom(addr1, contractID, addr2, defaultTokenID2)
		res, err := h(ctx, msg)
		require.NoError(t, err)
		e := sdk.Events{
			sdk.NewEvent("message", sdk.NewAttribute("module", "collection")),
			sdk.NewEvent("message", sdk.NewAttribute("sender", addr1.String())),
			sdk.NewEvent("detach_from", sdk.NewAttribute("contract_id", contractID)),
			sdk.NewEvent("detach_from", sdk.NewAttribute("proxy", addr1.String())),
			sdk.NewEvent("detach_from", sdk.NewAttribute("from", addr2.String())),
			sdk.NewEvent("detach_from", sdk.NewAttribute("from_token_id", defaultTokenID1)),
			sdk.NewEvent("detach_from", sdk.NewAttribute("token_id", defaultTokenID2)),
			sdk.NewEvent("detach_from", sdk.NewAttribute("old_root_token_id", defaultTokenID1)),
			sdk.NewEvent("detach_from", sdk.NewAttribute("new_root_token_id", defaultTokenID2)),
			sdk.NewEvent("operation_root_changed", sdk.NewAttribute("token_id", defaultTokenID2)),
		}
		verifyEventFunc(t, e, res.Events)
	}
}

func attach(t *testing.T, ctx sdk.Context, h sdk.Handler, contractID string) {
	msg := types.NewMsgAttach(addr1, contractID, defaultTokenID1, defaultTokenID2)
	res, err := h(ctx, msg)
	require.NoError(t, err)

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
		createMsg := types.NewMsgCreateCollection(addr1, defaultName, defaultMeta, defaultImgURI)
		res, err := h(ctx, createMsg)
		require.NoError(t, err)
		contractID = GetMadeContractID(res.Events)

		msg := types.NewMsgIssueNFT(addr1, contractID, defaultName, defaultMeta)
		_, err = h(ctx, msg)
		require.NoError(t, err)
		param := types.NewMintNFTParam(defaultName, defaultMeta, defaultTokenType)
		msg2 := types.NewMsgMintNFT(addr1, contractID, addr1, param)
		_, err = h(ctx, msg2)
		require.NoError(t, err)
		msg2 = types.NewMsgMintNFT(addr1, contractID, addr1, param)
		_, err = h(ctx, msg2)
		require.NoError(t, err)
	}

	attach(t, ctx, h, contractID)

	{
		msg2 := types.NewMsgDetach(addr1, contractID, defaultTokenID2)
		res2, err2 := h(ctx, msg2)
		require.NoError(t, err2)
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

	// Attach again
	attach(t, ctx, h, contractID)

	// Burn token
	{
		msg := types.NewMsgBurnNFT(addr1, contractID, defaultTokenID1)
		res, err := h(ctx, msg)
		require.NoError(t, err)
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

func TestHandleAttachFromDetachFromScenario(t *testing.T) {
	ctx, h := cacheKeeper()

	var contractID string
	{
		createMsg := types.NewMsgCreateCollection(addr1, defaultName, defaultMeta, defaultImgURI)
		res, err := h(ctx, createMsg)
		require.NoError(t, err)
		contractID = GetMadeContractID(res.Events)

		msg := types.NewMsgIssueNFT(addr1, contractID, defaultName, defaultMeta)
		_, err = h(ctx, msg)
		require.NoError(t, err)
		param := types.NewMintNFTParam(defaultName, defaultMeta, defaultTokenType)
		msg2 := types.NewMsgMintNFT(addr1, contractID, addr1, param)
		_, err = h(ctx, msg2)
		require.NoError(t, err)
		msg2 = types.NewMsgMintNFT(addr1, contractID, addr1, param)
		_, err = h(ctx, msg2)
		require.NoError(t, err)
		msg3 := types.NewMsgApprove(addr1, contractID, addr2)
		_, err = h(ctx, msg3)
		require.NoError(t, err)
	}

	msg := types.NewMsgAttachFrom(addr2, contractID, addr1, defaultTokenID1, defaultTokenID2)
	res, err := h(ctx, msg)
	require.NoError(t, err)
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
	res2, err2 := h(ctx, msg2)
	require.NoError(t, err2)
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

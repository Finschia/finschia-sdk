package handler

import (
	"testing"

	"github.com/line/link/x/collection/internal/types"
	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestHandleTransferFT(t *testing.T) {
	ctx, h := cacheKeeper()

	var contractID string
	{
		createMsg := types.NewMsgCreateCollection(addr1, defaultName, defaultImgURI)
		res := h(ctx, createMsg)
		require.True(t, res.Code.IsOK())
		contractID = GetMadeContractID(res.Events)

		msg := types.NewMsgIssueFT(addr1, addr1, contractID, defaultName, sdk.NewInt(defaultAmount), sdk.NewInt(defaultDecimals), true)
		res = h(ctx, msg)
		require.True(t, res.Code.IsOK())
	}

	msg := types.NewMsgTransferFT(addr1, contractID, addr2, types.NewCoin(defaultTokenIDFT, sdk.NewInt(defaultAmount)))
	res := h(ctx, msg)
	require.True(t, res.Code.IsOK())
	e := sdk.Events{
		sdk.NewEvent("message", sdk.NewAttribute("module", "collection")),
		sdk.NewEvent("message", sdk.NewAttribute("sender", addr1.String())),
		sdk.NewEvent("transfer_ft", sdk.NewAttribute("contract_id", contractID)),
		sdk.NewEvent("transfer_ft", sdk.NewAttribute("from", addr1.String())),
		sdk.NewEvent("transfer_ft", sdk.NewAttribute("to", addr2.String())),
		sdk.NewEvent("transfer_ft", sdk.NewAttribute("amount", types.NewCoin(defaultTokenIDFT, sdk.NewInt(defaultAmount)).String())),
	}
	verifyEventFunc(t, e, res.Events)
}

func TestHandleTransferFTFrom(t *testing.T) {
	ctx, h := cacheKeeper()

	var contractID string
	{
		createMsg := types.NewMsgCreateCollection(addr1, defaultName, defaultImgURI)
		res := h(ctx, createMsg)
		require.True(t, res.Code.IsOK())
		contractID = GetMadeContractID(res.Events)

		msg := types.NewMsgIssueFT(addr1, addr1, contractID, defaultName, sdk.NewInt(defaultAmount), sdk.NewInt(defaultDecimals), true)
		res = h(ctx, msg)
		require.True(t, res.Code.IsOK())
		msg2 := types.NewMsgApprove(addr1, contractID, addr2)
		res = h(ctx, msg2)
		require.True(t, res.Code.IsOK())
	}

	msg := types.NewMsgTransferFTFrom(addr2, contractID, addr1, addr2, types.NewCoin(defaultTokenIDFT, sdk.NewInt(defaultAmount)))
	res := h(ctx, msg)
	require.True(t, res.Code.IsOK())
	e := sdk.Events{
		sdk.NewEvent("message", sdk.NewAttribute("module", "collection")),
		sdk.NewEvent("message", sdk.NewAttribute("sender", addr2.String())),
		sdk.NewEvent("transfer_ft_from", sdk.NewAttribute("contract_id", contractID)),
		sdk.NewEvent("transfer_ft_from", sdk.NewAttribute("proxy", addr2.String())),
		sdk.NewEvent("transfer_ft_from", sdk.NewAttribute("from", addr1.String())),
		sdk.NewEvent("transfer_ft_from", sdk.NewAttribute("to", addr2.String())),
		sdk.NewEvent("transfer_ft_from", sdk.NewAttribute("amount", types.NewCoin(defaultTokenIDFT, sdk.NewInt(defaultAmount)).String())),
	}
	verifyEventFunc(t, e, res.Events)
}

func TestHandleTransferNFT(t *testing.T) {
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
	}

	msg := types.NewMsgTransferNFT(addr1, contractID, addr2, defaultTokenID1)
	res := h(ctx, msg)
	require.True(t, res.Code.IsOK())
	e := sdk.Events{
		sdk.NewEvent("message", sdk.NewAttribute("module", "collection")),
		sdk.NewEvent("message", sdk.NewAttribute("sender", addr1.String())),
		sdk.NewEvent("transfer_nft", sdk.NewAttribute("contract_id", contractID)),
		sdk.NewEvent("transfer_nft", sdk.NewAttribute("from", addr1.String())),
		sdk.NewEvent("transfer_nft", sdk.NewAttribute("to", addr2.String())),
		sdk.NewEvent("transfer_nft", sdk.NewAttribute("token_id", defaultTokenID1)),
		sdk.NewEvent("operation_transfer_nft", sdk.NewAttribute("token_id", defaultTokenID1)),
	}
	verifyEventFunc(t, e, res.Events)
}

func TestHandleTransferNFTFrom(t *testing.T) {
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
		msg3 := types.NewMsgApprove(addr1, contractID, addr2)
		res = h(ctx, msg3)
		require.True(t, res.Code.IsOK())
	}

	msg := types.NewMsgTransferNFTFrom(addr2, contractID, addr1, addr2, defaultTokenID1)
	res := h(ctx, msg)
	require.True(t, res.Code.IsOK())
	e := sdk.Events{
		sdk.NewEvent("message", sdk.NewAttribute("module", "collection")),
		sdk.NewEvent("message", sdk.NewAttribute("sender", addr2.String())),
		sdk.NewEvent("transfer_nft_from", sdk.NewAttribute("contract_id", contractID)),
		sdk.NewEvent("transfer_nft_from", sdk.NewAttribute("proxy", addr2.String())),
		sdk.NewEvent("transfer_nft_from", sdk.NewAttribute("from", addr1.String())),
		sdk.NewEvent("transfer_nft_from", sdk.NewAttribute("to", addr2.String())),
		sdk.NewEvent("transfer_nft_from", sdk.NewAttribute("token_id", defaultTokenID1)),
		sdk.NewEvent("operation_transfer_nft", sdk.NewAttribute("token_id", defaultTokenID1)),
	}
	verifyEventFunc(t, e, res.Events)
}

func TestHandleTransferNFTChild(t *testing.T) {
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
		msg2 = types.NewMsgMintNFT(addr1, contractID, addr1, defaultName, defaultTokenType)
		res = h(ctx, msg2)
		require.True(t, res.Code.IsOK())
		msg3 := types.NewMsgAttach(addr1, contractID, defaultTokenID1, defaultTokenID2)
		res = h(ctx, msg3)
		require.True(t, res.Code.IsOK())
		msg3 = types.NewMsgAttach(addr1, contractID, defaultTokenID2, defaultTokenID3)
		res = h(ctx, msg3)
		require.True(t, res.Code.IsOK())
	}

	msg := types.NewMsgTransferNFT(addr1, contractID, addr2, defaultTokenID1)
	res := h(ctx, msg)
	require.True(t, res.Code.IsOK())
	e := sdk.Events{
		sdk.NewEvent("message", sdk.NewAttribute("module", "collection")),
		sdk.NewEvent("message", sdk.NewAttribute("sender", addr1.String())),
		sdk.NewEvent("transfer_nft", sdk.NewAttribute("contract_id", contractID)),
		sdk.NewEvent("transfer_nft", sdk.NewAttribute("from", addr1.String())),
		sdk.NewEvent("transfer_nft", sdk.NewAttribute("to", addr2.String())),
		sdk.NewEvent("transfer_nft", sdk.NewAttribute("token_id", defaultTokenID1)),
		sdk.NewEvent("operation_transfer_nft", sdk.NewAttribute("token_id", defaultTokenID1)),
		sdk.NewEvent("operation_transfer_nft", sdk.NewAttribute("token_id", defaultTokenID2)),
		sdk.NewEvent("operation_transfer_nft", sdk.NewAttribute("token_id", defaultTokenID3)),
	}
	verifyEventFunc(t, e, res.Events)
}

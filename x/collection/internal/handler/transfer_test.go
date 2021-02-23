package handler

import (
	"testing"

	"github.com/line/lbm-sdk/x/collection/internal/types"
	"github.com/stretchr/testify/require"

	sdk "github.com/line/lbm-sdk/types"
)

func TestHandleTransferFT(t *testing.T) {
	ctx, h := cacheKeeper()

	var contractID string
	{
		createMsg := types.NewMsgCreateCollection(addr1, defaultName, defaultMeta, defaultImgURI)
		res, err := h(ctx, createMsg)
		require.NoError(t, err)
		contractID = GetMadeContractID(res.Events)

		msg := types.NewMsgIssueFT(addr1, addr1, contractID, defaultName, defaultMeta, sdk.NewInt(defaultAmount), sdk.NewInt(defaultDecimals), true)
		_, err = h(ctx, msg)
		require.NoError(t, err)
	}

	msg := types.NewMsgTransferFT(addr1, contractID, addr2, types.NewCoin(defaultTokenIDFT, sdk.NewInt(defaultAmount)))
	res, err := h(ctx, msg)
	require.NoError(t, err)
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
		createMsg := types.NewMsgCreateCollection(addr1, defaultName, defaultMeta, defaultImgURI)
		res, err := h(ctx, createMsg)
		require.NoError(t, err)
		contractID = GetMadeContractID(res.Events)

		msg := types.NewMsgIssueFT(addr1, addr1, contractID, defaultName, defaultMeta, sdk.NewInt(defaultAmount), sdk.NewInt(defaultDecimals), true)
		_, err = h(ctx, msg)
		require.NoError(t, err)
		msg2 := types.NewMsgApprove(addr1, contractID, addr2)
		_, err = h(ctx, msg2)
		require.NoError(t, err)
	}

	msg := types.NewMsgTransferFTFrom(addr2, contractID, addr1, addr2, types.NewCoin(defaultTokenIDFT, sdk.NewInt(defaultAmount)))
	res, err := h(ctx, msg)
	require.NoError(t, err)
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
	}

	msg := types.NewMsgTransferNFT(addr1, contractID, addr2, defaultTokenID1)
	res, err := h(ctx, msg)
	require.NoError(t, err)
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
		msg3 := types.NewMsgApprove(addr1, contractID, addr2)
		_, err = h(ctx, msg3)
		require.NoError(t, err)
	}

	msg := types.NewMsgTransferNFTFrom(addr2, contractID, addr1, addr2, defaultTokenID1)
	res, err := h(ctx, msg)
	require.NoError(t, err)
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
		msg2 = types.NewMsgMintNFT(addr1, contractID, addr1, param)
		_, err = h(ctx, msg2)
		require.NoError(t, err)
		msg3 := types.NewMsgAttach(addr1, contractID, defaultTokenID1, defaultTokenID2)
		_, err = h(ctx, msg3)
		require.NoError(t, err)
		msg3 = types.NewMsgAttach(addr1, contractID, defaultTokenID2, defaultTokenID3)
		_, err = h(ctx, msg3)
		require.NoError(t, err)
	}

	msg := types.NewMsgTransferNFT(addr1, contractID, addr2, defaultTokenID1)
	res, err := h(ctx, msg)
	require.NoError(t, err)
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

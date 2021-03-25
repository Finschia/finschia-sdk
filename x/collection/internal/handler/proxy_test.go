package handler

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/line/link-modules/x/collection/internal/types"
	"github.com/stretchr/testify/require"
)

func approve(t *testing.T, approver, proxy sdk.AccAddress, contractID string, ctx sdk.Context, h sdk.Handler) {
	approveMsg := types.NewMsgApprove(approver, contractID, proxy)
	_, err := h(ctx, approveMsg)
	require.NoError(t, err)
}

func TestHandleMsgApprove(t *testing.T) {
	t.Log("implement me - ", t.Name())
}
func TestHandleMsgDisapprove(t *testing.T) {
	t.Log("implement me - ", t.Name())
}

func TestHandleApproveDisapprove(t *testing.T) {
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
		msg3 := types.NewMsgIssueFT(addr1, addr1, contractID, defaultName, defaultMeta, sdk.NewInt(defaultAmount), sdk.NewInt(defaultDecimals), true)
		_, err = h(ctx, msg3)
		require.NoError(t, err)
		msg4 := types.NewMsgMintFT(addr1, contractID, addr1, types.NewCoin(defaultTokenIDFT, sdk.NewInt(defaultAmount)))
		_, err = h(ctx, msg4)
		require.NoError(t, err)
	}

	msg := types.NewMsgTransferNFTFrom(addr2, contractID, addr1, addr2, defaultTokenID1)
	_, err := h(ctx, msg)
	require.Error(t, err)

	{
		msg3 := types.NewMsgApprove(addr1, contractID, addr2)
		_, err = h(ctx, msg3)
		require.NoError(t, err)
	}

	msg = types.NewMsgTransferNFTFrom(addr2, contractID, addr1, addr2, defaultTokenID1)
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

	msg2 := types.NewMsgBurnNFTFrom(addr2, contractID, addr1, defaultTokenID2)
	_, err = h(ctx, msg2)
	require.Error(t, err) // addr2 does not have the burn permission
	msg3 := types.NewMsgBurnFTFrom(addr2, contractID, addr1, types.NewCoin(defaultTokenIDFT, sdk.NewInt(defaultAmount)))
	_, err = h(ctx, msg3)
	require.Error(t, err) // addr2 does not have the burn permission

	{
		permission := types.NewBurnPermission()
		msg := types.NewMsgGrantPermission(addr1, contractID, addr2, permission)
		_, err := h(ctx, msg)
		require.NoError(t, err)
	}

	msg2 = types.NewMsgBurnNFTFrom(addr2, contractID, addr1, defaultTokenID2)
	res, err = h(ctx, msg2)
	require.NoError(t, err) // addr2 does not have the burn permission

	e = sdk.Events{
		sdk.NewEvent("message", sdk.NewAttribute("module", "collection")),
		sdk.NewEvent("message", sdk.NewAttribute("sender", addr2.String())),
		sdk.NewEvent("burn_nft_from", sdk.NewAttribute("contract_id", contractID)),
		sdk.NewEvent("burn_nft_from", sdk.NewAttribute("proxy", addr2.String())),
		sdk.NewEvent("burn_nft_from", sdk.NewAttribute("from", addr1.String())),
		sdk.NewEvent("burn_nft_from", sdk.NewAttribute("token_id", defaultTokenID2)),
		sdk.NewEvent("operation_burn_nft", sdk.NewAttribute("token_id", defaultTokenID2)),
	}
	verifyEventFunc(t, e, res.Events)

	msg3 = types.NewMsgBurnFTFrom(addr2, contractID, addr1, types.NewCoin(defaultTokenIDFT, sdk.NewInt(defaultAmount)))
	_, err = h(ctx, msg3)
	require.NoError(t, err)

	{
		permission := types.NewBurnPermission()
		msg := types.NewMsgGrantPermission(addr1, contractID, addr2, permission)
		_, err := h(ctx, msg)
		require.NoError(t, err)
	}

	msg3 = types.NewMsgBurnFTFrom(addr2, contractID, addr1, types.NewCoin(defaultTokenIDFT, sdk.NewInt(defaultAmount)))
	res, err = h(ctx, msg3)
	require.NoError(t, err)
	e = sdk.Events{
		sdk.NewEvent("message", sdk.NewAttribute("module", "collection")),
		sdk.NewEvent("message", sdk.NewAttribute("sender", addr2.String())),
		sdk.NewEvent("burn_ft_from", sdk.NewAttribute("contract_id", contractID)),
		sdk.NewEvent("burn_ft_from", sdk.NewAttribute("proxy", addr2.String())),
		sdk.NewEvent("burn_ft_from", sdk.NewAttribute("from", addr1.String())),
		sdk.NewEvent("burn_ft_from", sdk.NewAttribute("amount", types.NewCoins(types.NewCoin(defaultTokenIDFT, sdk.NewInt(defaultAmount))).String())),
	}
	verifyEventFunc(t, e, res.Events)

	{
		msg3 := types.NewMsgDisapprove(addr1, contractID, addr2)
		_, err = h(ctx, msg3)
		require.NoError(t, err)
	}

	msg = types.NewMsgTransferNFTFrom(addr2, contractID, addr1, addr2, defaultTokenID1)
	_, err = h(ctx, msg)
	require.Error(t, err)

	msg2 = types.NewMsgBurnNFTFrom(addr2, contractID, addr1, defaultTokenID1)
	_, err = h(ctx, msg2)
	require.Error(t, err)

	msg3 = types.NewMsgBurnFTFrom(addr2, contractID, addr1, types.NewCoin(defaultTokenIDFT, sdk.NewInt(defaultAmount)))
	_, err = h(ctx, msg3)
	require.Error(t, err)
}

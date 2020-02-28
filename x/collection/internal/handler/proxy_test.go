package handler

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/line/link/x/collection/internal/types"
	"github.com/stretchr/testify/require"
)

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
		msg3 := types.NewMsgIssueFT(addr1, contractID, defaultName, sdk.NewInt(defaultAmount), sdk.NewInt(defaultDecimals), true)
		res = h(ctx, msg3)
		require.True(t, res.Code.IsOK())
		msg4 := types.NewMsgMintFT(addr1, contractID, addr1, types.NewCoin(defaultTokenIDFT, sdk.NewInt(defaultAmount)))
		res = h(ctx, msg4)
		require.True(t, res.Code.IsOK())
	}

	msg := types.NewMsgTransferNFTFrom(addr2, contractID, addr1, addr2, defaultTokenID1)
	res := h(ctx, msg)
	require.False(t, res.Code.IsOK())

	{
		msg3 := types.NewMsgApprove(addr1, contractID, addr2)
		res = h(ctx, msg3)
		require.True(t, res.Code.IsOK())
	}

	msg = types.NewMsgTransferNFTFrom(addr2, contractID, addr1, addr2, defaultTokenID1)
	res = h(ctx, msg)
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

	msg2 := types.NewMsgBurnNFTFrom(addr2, contractID, addr1, defaultTokenID2)
	res = h(ctx, msg2)
	require.False(t, res.Code.IsOK()) // addr2 does not have the burn permission
	msg3 := types.NewMsgBurnFTFrom(addr2, contractID, addr1, types.NewCoin(defaultTokenIDFT, sdk.NewInt(defaultAmount)))
	res = h(ctx, msg3)
	require.False(t, res.Code.IsOK()) // addr2 does not have the burn permission

	{
		permission := types.Permission{
			Action:   "burn",
			Resource: contractID,
		}
		msg := types.NewMsgGrantPermission(addr1, addr2, permission)
		res := h(ctx, msg)
		require.True(t, res.Code.IsOK())
	}

	msg2 = types.NewMsgBurnNFTFrom(addr2, contractID, addr1, defaultTokenID2)
	res = h(ctx, msg2)
	require.True(t, res.Code.IsOK()) // addr2 does not have the burn permission

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
	res = h(ctx, msg3)
	require.True(t, res.Code.IsOK())

	{
		permission := types.Permission{
			Action:   "burn",
			Resource: contractID,
		}
		msg := types.NewMsgGrantPermission(addr1, addr2, permission)
		res := h(ctx, msg)
		require.True(t, res.Code.IsOK())
	}

	msg3 = types.NewMsgBurnFTFrom(addr2, contractID, addr1, types.NewCoin(defaultTokenIDFT, sdk.NewInt(defaultAmount)))
	res = h(ctx, msg3)
	require.True(t, res.Code.IsOK())
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
		res = h(ctx, msg3)
		require.True(t, res.Code.IsOK())
	}

	msg = types.NewMsgTransferNFTFrom(addr2, contractID, addr1, addr2, defaultTokenID1)
	res = h(ctx, msg)
	require.False(t, res.Code.IsOK())

	msg2 = types.NewMsgBurnNFTFrom(addr2, contractID, addr1, defaultTokenID1)
	res = h(ctx, msg2)
	require.False(t, res.Code.IsOK())

	msg3 = types.NewMsgBurnFTFrom(addr2, contractID, addr1, types.NewCoin(defaultTokenIDFT, sdk.NewInt(defaultAmount)))
	res = h(ctx, msg3)
	require.False(t, res.Code.IsOK())
}

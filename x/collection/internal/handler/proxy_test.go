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

	{
		createMsg := types.NewMsgCreateCollection(addr1, defaultName, defaultSymbol)
		res := h(ctx, createMsg)
		require.True(t, res.Code.IsOK())
		msg := types.NewMsgIssueNFT(addr1, defaultSymbol, defaultName)
		res = h(ctx, msg)
		require.True(t, res.Code.IsOK())
		msg2 := types.NewMsgMintNFT(addr1, addr1, defaultName, defaultSymbol, defaultTokenURI, defaultTokenType)
		res = h(ctx, msg2)
		require.True(t, res.Code.IsOK())
		msg2 = types.NewMsgMintNFT(addr1, addr1, defaultName, defaultSymbol, defaultTokenURI, defaultTokenType)
		res = h(ctx, msg2)
		require.True(t, res.Code.IsOK())
		msg3 := types.NewMsgIssueFT(addr1, defaultName, defaultSymbol, defaultTokenURI, sdk.NewInt(defaultAmount), sdk.NewInt(defaultDecimals), true)
		res = h(ctx, msg3)
		require.True(t, res.Code.IsOK())
		msg4 := types.NewMsgMintFT(defaultSymbol, addr1, addr1, types.NewCoin(defaultTokenIDFT, sdk.NewInt(defaultAmount)))
		res = h(ctx, msg4)
		require.True(t, res.Code.IsOK())
	}

	msg := types.NewMsgTransferNFTFrom(addr2, addr1, addr2, defaultSymbol, defaultTokenID1)
	res := h(ctx, msg)
	require.False(t, res.Code.IsOK())

	{
		msg3 := types.NewMsgApprove(addr1, addr2, defaultSymbol)
		res = h(ctx, msg3)
		require.True(t, res.Code.IsOK())
	}

	msg = types.NewMsgTransferNFTFrom(addr2, addr1, addr2, defaultSymbol, defaultTokenID1)
	res = h(ctx, msg)
	require.True(t, res.Code.IsOK())

	e := sdk.Events{
		sdk.NewEvent("message", sdk.NewAttribute("module", "collection")),
		sdk.NewEvent("message", sdk.NewAttribute("sender", addr2.String())),
		sdk.NewEvent("transfer_nft_from", sdk.NewAttribute("proxy", addr2.String())),
		sdk.NewEvent("transfer_nft_from", sdk.NewAttribute("from", addr1.String())),
		sdk.NewEvent("transfer_nft_from", sdk.NewAttribute("to", addr2.String())),
		sdk.NewEvent("transfer_nft_from", sdk.NewAttribute("symbol", defaultSymbol)),
		sdk.NewEvent("transfer_nft_from", sdk.NewAttribute("token_id", defaultTokenID1)),
		sdk.NewEvent("operation_transfer_nft", sdk.NewAttribute("token_id", defaultTokenID1)),
	}
	verifyEventFunc(t, e, res.Events)

	msg2 := types.NewMsgBurnNFTFrom(addr2, addr1, defaultSymbol, defaultTokenID2)
	res = h(ctx, msg2)
	require.False(t, res.Code.IsOK()) // addr2 does not have the burn permission
	msg3 := types.NewMsgBurnFTFrom(defaultSymbol, addr2, addr1, types.NewCoin(defaultTokenIDFT, sdk.NewInt(defaultAmount)))
	res = h(ctx, msg3)
	require.False(t, res.Code.IsOK()) // addr2 does not have the burn permission

	{
		permission := types.Permission{
			Action:   "burn",
			Resource: defaultSymbol,
		}
		msg := types.NewMsgGrantPermission(addr1, addr2, permission)
		res := h(ctx, msg)
		require.True(t, res.Code.IsOK())
	}

	msg2 = types.NewMsgBurnNFTFrom(addr2, addr1, defaultSymbol, defaultTokenID2)
	res = h(ctx, msg2)
	require.True(t, res.Code.IsOK()) // addr2 does not have the burn permission

	e = sdk.Events{
		sdk.NewEvent("message", sdk.NewAttribute("module", "collection")),
		sdk.NewEvent("message", sdk.NewAttribute("sender", addr2.String())),
		sdk.NewEvent("burn_nft_from", sdk.NewAttribute("proxy", addr2.String())),
		sdk.NewEvent("burn_nft_from", sdk.NewAttribute("from", addr1.String())),
		sdk.NewEvent("burn_nft_from", sdk.NewAttribute("symbol", defaultSymbol)),
		sdk.NewEvent("burn_nft_from", sdk.NewAttribute("token_id", defaultTokenID2)),
		sdk.NewEvent("operation_burn_nft", sdk.NewAttribute("token_id", defaultTokenID2)),
	}
	verifyEventFunc(t, e, res.Events)

	msg3 = types.NewMsgBurnFTFrom(defaultSymbol, addr2, addr1, types.NewCoin(defaultTokenIDFT, sdk.NewInt(defaultAmount)))
	res = h(ctx, msg3)
	require.True(t, res.Code.IsOK())
	e = sdk.Events{
		sdk.NewEvent("message", sdk.NewAttribute("module", "collection")),
		sdk.NewEvent("message", sdk.NewAttribute("sender", addr2.String())),
		sdk.NewEvent("burn_ft_from", sdk.NewAttribute("proxy", addr2.String())),
		sdk.NewEvent("burn_ft_from", sdk.NewAttribute("from", addr1.String())),
		sdk.NewEvent("burn_ft_from", sdk.NewAttribute("amount", types.NewCoins(types.NewCoin(defaultTokenIDFT, sdk.NewInt(defaultAmount))).String())),
	}
	verifyEventFunc(t, e, res.Events)

	{
		msg3 := types.NewMsgDisapprove(addr1, addr2, defaultSymbol)
		res = h(ctx, msg3)
		require.True(t, res.Code.IsOK())
	}

	msg = types.NewMsgTransferNFTFrom(addr2, addr1, addr2, defaultSymbol, defaultTokenID1)
	res = h(ctx, msg)
	require.False(t, res.Code.IsOK())

	msg2 = types.NewMsgBurnNFTFrom(addr2, addr1, defaultSymbol, defaultTokenID1)
	res = h(ctx, msg2)
	require.False(t, res.Code.IsOK())

	msg3 = types.NewMsgBurnFTFrom(defaultSymbol, addr2, addr1, types.NewCoin(defaultTokenIDFT, sdk.NewInt(defaultAmount)))
	res = h(ctx, msg3)
	require.False(t, res.Code.IsOK())
}

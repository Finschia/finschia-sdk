package handler

import (
	"testing"

	"github.com/line/link/x/collection/internal/types"
	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestHandleTransferCFT(t *testing.T) {
	ctx, h := cacheKeeper()

	{
		createMsg := types.NewMsgCreateCollection(addr1, defaultName, defaultSymbol)
		res := h(ctx, createMsg)
		require.True(t, res.Code.IsOK())
		msg := types.NewMsgIssueCFT(addr1, defaultName, defaultSymbol, defaultTokenURI, sdk.NewInt(defaultAmount), sdk.NewInt(defaultDecimals), true)
		res = h(ctx, msg)
		require.True(t, res.Code.IsOK())
	}

	msg := types.NewMsgTransferCFT(addr1, addr2, defaultSymbol, defaultTokenIDFT, sdk.NewInt(defaultAmount))
	res := h(ctx, msg)
	require.True(t, res.Code.IsOK())
	e := sdk.Events{
		sdk.NewEvent("message", sdk.NewAttribute("module", "collection")),
		sdk.NewEvent("message", sdk.NewAttribute("sender", addr1.String())),
		sdk.NewEvent("transfer_cft", sdk.NewAttribute("from", addr1.String())),
		sdk.NewEvent("transfer_cft", sdk.NewAttribute("to", addr2.String())),
		sdk.NewEvent("transfer_cft", sdk.NewAttribute("symbol", defaultSymbol)),
		sdk.NewEvent("transfer_cft", sdk.NewAttribute("token_id", defaultTokenIDFT)),
		sdk.NewEvent("transfer_cft", sdk.NewAttribute("amount", sdk.NewInt(defaultAmount).String())),
	}
	verifyEventFunc(t, e, res.Events)
}

func TestHandleTransferCFTFrom(t *testing.T) {
	ctx, h := cacheKeeper()

	{
		createMsg := types.NewMsgCreateCollection(addr1, defaultName, defaultSymbol)
		res := h(ctx, createMsg)
		require.True(t, res.Code.IsOK())
		msg := types.NewMsgIssueCFT(addr1, defaultName, defaultSymbol, defaultTokenURI, sdk.NewInt(defaultAmount), sdk.NewInt(defaultDecimals), true)
		res = h(ctx, msg)
		require.True(t, res.Code.IsOK())
		msg2 := types.NewMsgApprove(addr1, addr2, defaultSymbol)
		res = h(ctx, msg2)
		require.True(t, res.Code.IsOK())
	}

	msg := types.NewMsgTransferCFTFrom(addr2, addr1, addr2, defaultSymbol, defaultTokenIDFT, sdk.NewInt(defaultAmount))
	res := h(ctx, msg)
	require.True(t, res.Code.IsOK())
	e := sdk.Events{
		sdk.NewEvent("message", sdk.NewAttribute("module", "collection")),
		sdk.NewEvent("message", sdk.NewAttribute("sender", addr2.String())),
		sdk.NewEvent("transfer_cft_from", sdk.NewAttribute("proxy", addr2.String())),
		sdk.NewEvent("transfer_cft_from", sdk.NewAttribute("from", addr1.String())),
		sdk.NewEvent("transfer_cft_from", sdk.NewAttribute("to", addr2.String())),
		sdk.NewEvent("transfer_cft_from", sdk.NewAttribute("symbol", defaultSymbol)),
		sdk.NewEvent("transfer_cft_from", sdk.NewAttribute("token_id", defaultTokenIDFT)),
		sdk.NewEvent("transfer_cft_from", sdk.NewAttribute("amount", sdk.NewInt(defaultAmount).String())),
	}
	verifyEventFunc(t, e, res.Events)
}

func TestHandleTransferCNFT(t *testing.T) {
	ctx, h := cacheKeeper()

	{
		createMsg := types.NewMsgCreateCollection(addr1, defaultName, defaultSymbol)
		res := h(ctx, createMsg)
		require.True(t, res.Code.IsOK())
		msg := types.NewMsgIssueCNFT(addr1, defaultSymbol)
		res = h(ctx, msg)
		require.True(t, res.Code.IsOK())
		msg2 := types.NewMsgMintCNFT(addr1, addr1, defaultName, defaultSymbol, defaultTokenURI, defaultTokenType)
		res = h(ctx, msg2)
		require.True(t, res.Code.IsOK())
	}

	msg := types.NewMsgTransferCNFT(addr1, addr2, defaultSymbol, defaultTokenID1)
	res := h(ctx, msg)
	require.True(t, res.Code.IsOK())
	e := sdk.Events{
		sdk.NewEvent("message", sdk.NewAttribute("module", "collection")),
		sdk.NewEvent("message", sdk.NewAttribute("sender", addr1.String())),
		sdk.NewEvent("transfer_cnft", sdk.NewAttribute("from", addr1.String())),
		sdk.NewEvent("transfer_cnft", sdk.NewAttribute("to", addr2.String())),
		sdk.NewEvent("transfer_cnft", sdk.NewAttribute("symbol", defaultSymbol)),
		sdk.NewEvent("transfer_cnft", sdk.NewAttribute("token_id", defaultTokenID1)),
	}
	verifyEventFunc(t, e, res.Events)
}

func TestHandleTransferCNFTFrom(t *testing.T) {
	ctx, h := cacheKeeper()

	{
		createMsg := types.NewMsgCreateCollection(addr1, defaultName, defaultSymbol)
		res := h(ctx, createMsg)
		require.True(t, res.Code.IsOK())
		msg := types.NewMsgIssueCNFT(addr1, defaultSymbol)
		res = h(ctx, msg)
		require.True(t, res.Code.IsOK())
		msg2 := types.NewMsgMintCNFT(addr1, addr1, defaultName, defaultSymbol, defaultTokenURI, defaultTokenType)
		res = h(ctx, msg2)
		require.True(t, res.Code.IsOK())
		msg3 := types.NewMsgApprove(addr1, addr2, defaultSymbol)
		res = h(ctx, msg3)
		require.True(t, res.Code.IsOK())
	}

	msg := types.NewMsgTransferCNFTFrom(addr2, addr1, addr2, defaultSymbol, defaultTokenID1)
	res := h(ctx, msg)
	require.True(t, res.Code.IsOK())
	e := sdk.Events{
		sdk.NewEvent("message", sdk.NewAttribute("module", "collection")),
		sdk.NewEvent("message", sdk.NewAttribute("sender", addr2.String())),
		sdk.NewEvent("transfer_cnft_from", sdk.NewAttribute("proxy", addr2.String())),
		sdk.NewEvent("transfer_cnft_from", sdk.NewAttribute("from", addr1.String())),
		sdk.NewEvent("transfer_cnft_from", sdk.NewAttribute("to", addr2.String())),
		sdk.NewEvent("transfer_cnft_from", sdk.NewAttribute("symbol", defaultSymbol)),
		sdk.NewEvent("transfer_cnft_from", sdk.NewAttribute("token_id", defaultTokenID1)),
	}
	verifyEventFunc(t, e, res.Events)
}

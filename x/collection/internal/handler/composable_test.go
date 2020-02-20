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

func TestHandleAttachDetach(t *testing.T) {
	ctx, h := cacheKeeper()

	{
		createMsg := types.NewMsgCreateCollection(addr1, defaultName, defaultSymbol)
		res := h(ctx, createMsg)
		require.True(t, res.Code.IsOK())
		msg := types.NewMsgIssueCNFT(addr1, defaultSymbol)
		res = h(ctx, msg)
		require.True(t, res.Code.IsOK())
		msg2 := types.NewMsgMintCNFT(addr1, addr1, defaultName, defaultSymbol, defaultTokenURI, "1001")
		res = h(ctx, msg2)
		require.True(t, res.Code.IsOK())
		msg2 = types.NewMsgMintCNFT(addr1, addr1, defaultName, defaultSymbol, defaultTokenURI, "1001")
		res = h(ctx, msg2)
		require.True(t, res.Code.IsOK())
	}

	{
		msg := types.NewMsgAttach(addr1, defaultSymbol, defaultTokenID1, defaultTokenID2)
		res := h(ctx, msg)
		require.True(t, res.Code.IsOK())
		e := sdk.Events{
			sdk.NewEvent("message", sdk.NewAttribute("module", "collection")),
			sdk.NewEvent("message", sdk.NewAttribute("sender", addr1.String())),
			sdk.NewEvent("attach", sdk.NewAttribute("from", addr1.String())),
			sdk.NewEvent("attach", sdk.NewAttribute("symbol", defaultSymbol)),
			sdk.NewEvent("attach", sdk.NewAttribute("to_token_id", defaultTokenID1)),
			sdk.NewEvent("attach", sdk.NewAttribute("token_id", defaultTokenID2)),
		}
		verifyEventFunc(t, e, res.Events)
	}

	{
		msg2 := types.NewMsgDetach(addr1, defaultSymbol, defaultTokenID2)
		res2 := h(ctx, msg2)
		require.True(t, res2.Code.IsOK())
		e := sdk.Events{
			sdk.NewEvent("message", sdk.NewAttribute("module", "collection")),
			sdk.NewEvent("message", sdk.NewAttribute("sender", addr1.String())),
			sdk.NewEvent("detach", sdk.NewAttribute("from", addr1.String())),
			sdk.NewEvent("detach", sdk.NewAttribute("symbol", defaultSymbol)),
			sdk.NewEvent("detach", sdk.NewAttribute("token_id", defaultTokenID2)),
		}
		verifyEventFunc(t, e, res2.Events)
	}

	//Attach again
	{
		msg := types.NewMsgAttach(addr1, defaultSymbol, defaultTokenID1, defaultTokenID2)
		res := h(ctx, msg)
		require.True(t, res.Code.IsOK())
		e := sdk.Events{
			sdk.NewEvent("message", sdk.NewAttribute("module", "collection")),
			sdk.NewEvent("message", sdk.NewAttribute("sender", addr1.String())),
			sdk.NewEvent("attach", sdk.NewAttribute("from", addr1.String())),
			sdk.NewEvent("attach", sdk.NewAttribute("symbol", defaultSymbol)),
			sdk.NewEvent("attach", sdk.NewAttribute("to_token_id", defaultTokenID1)),
			sdk.NewEvent("attach", sdk.NewAttribute("token_id", defaultTokenID2)),
		}
		verifyEventFunc(t, e, res.Events)
	}
	//Burn token
	{
		msg := types.NewMsgBurnCNFT(addr1, defaultSymbol, defaultTokenID1)
		res := h(ctx, msg)
		require.True(t, res.Code.IsOK())
		e := sdk.Events{
			sdk.NewEvent("message", sdk.NewAttribute("module", "collection")),
			sdk.NewEvent("message", sdk.NewAttribute("sender", addr1.String())),
			sdk.NewEvent("burn_cnft", sdk.NewAttribute("from", addr1.String())),
			sdk.NewEvent("burn_cnft", sdk.NewAttribute("symbol", defaultSymbol)),
			sdk.NewEvent("burn_cnft", sdk.NewAttribute("token_id", defaultTokenID1)),
		}
		verifyEventFunc(t, e, res.Events)
	}
}

func TestHandleAttachFromDetachFrom(t *testing.T) {
	ctx, h := cacheKeeper()

	{
		createMsg := types.NewMsgCreateCollection(addr1, defaultName, defaultSymbol)
		res := h(ctx, createMsg)
		require.True(t, res.Code.IsOK())
		msg := types.NewMsgIssueCNFT(addr1, defaultSymbol)
		res = h(ctx, msg)
		require.True(t, res.Code.IsOK())
		msg2 := types.NewMsgMintCNFT(addr1, addr1, defaultName, defaultSymbol, defaultTokenURI, "1001")
		res = h(ctx, msg2)
		require.True(t, res.Code.IsOK())
		msg2 = types.NewMsgMintCNFT(addr1, addr1, defaultName, defaultSymbol, defaultTokenURI, "1001")
		res = h(ctx, msg2)
		require.True(t, res.Code.IsOK())
		msg3 := types.NewMsgApprove(addr1, addr2, defaultSymbol)
		res = h(ctx, msg3)
		require.True(t, res.Code.IsOK())
	}

	msg := types.NewMsgAttachFrom(addr2, addr1, defaultSymbol, defaultTokenID1, defaultTokenID2)
	res := h(ctx, msg)
	require.True(t, res.Code.IsOK())
	e := sdk.Events{
		sdk.NewEvent("message", sdk.NewAttribute("module", "collection")),
		sdk.NewEvent("message", sdk.NewAttribute("sender", addr2.String())),
		sdk.NewEvent("attach_from", sdk.NewAttribute("proxy", addr2.String())),
		sdk.NewEvent("attach_from", sdk.NewAttribute("from", addr1.String())),
		sdk.NewEvent("attach_from", sdk.NewAttribute("symbol", defaultSymbol)),
		sdk.NewEvent("attach_from", sdk.NewAttribute("to_token_id", defaultTokenID1)),
		sdk.NewEvent("attach_from", sdk.NewAttribute("token_id", defaultTokenID2)),
	}
	verifyEventFunc(t, e, res.Events)

	msg2 := types.NewMsgDetachFrom(addr2, addr1, defaultSymbol, defaultTokenID2)
	res2 := h(ctx, msg2)
	require.True(t, res2.Code.IsOK())
	e = sdk.Events{
		sdk.NewEvent("message", sdk.NewAttribute("module", "collection")),
		sdk.NewEvent("message", sdk.NewAttribute("sender", addr2.String())),
		sdk.NewEvent("detach_from", sdk.NewAttribute("proxy", addr2.String())),
		sdk.NewEvent("detach_from", sdk.NewAttribute("from", addr1.String())),
		sdk.NewEvent("detach_from", sdk.NewAttribute("symbol", defaultSymbol)),
		sdk.NewEvent("detach_from", sdk.NewAttribute("token_id", defaultTokenID2)),
	}
	verifyEventFunc(t, e, res2.Events)
}

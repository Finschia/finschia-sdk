package handler

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/line/link-modules/x/collection/internal/types"
	"github.com/stretchr/testify/require"
)

func TestHandleMsgBurnFT(t *testing.T) {
	ctx, h, contractID := prepareFT(t)

	{
		// invalid user
		burnMsg := types.NewMsgBurnFT(addr2, contractID, types.NewCoin("0000000100000000", sdk.NewInt(100)))
		_, err := h(ctx, burnMsg)
		require.Error(t, err)
	}
	{
		// burn non-exist token
		burnMsg := types.NewMsgBurnFT(addr1, contractID, types.NewCoin("0000000200000000", sdk.NewInt(100)))
		_, err := h(ctx, burnMsg)
		require.Error(t, err)
	}
	{
		// burn tokens over the being supplied
		burnMsg := types.NewMsgBurnFT(addr1, contractID, types.NewCoin("0000000100000000", sdk.NewInt(1001)))
		_, err := h(ctx, burnMsg)
		require.Error(t, err)
	}
	{
		// burn tokens with invalid contractID
		burnMsg := types.NewMsgBurnFT(addr1, "abcd11234", types.NewCoin("0000000100000000", sdk.NewInt(1000)))
		_, err := h(ctx, burnMsg)
		require.Error(t, err)
	}
	{
		// success case
		burnMsg := types.NewMsgBurnFT(addr1, contractID, types.NewCoin("0000000100000000", sdk.NewInt(100)))
		res, err := h(ctx, burnMsg)
		require.NoError(t, err)

		e := sdk.Events{
			sdk.NewEvent("message", sdk.NewAttribute("module", "collection")),
			sdk.NewEvent("message", sdk.NewAttribute("sender", addr1.String())),
			sdk.NewEvent("burn_ft", sdk.NewAttribute("contract_id", contractID)),
			sdk.NewEvent("burn_ft", sdk.NewAttribute("from", addr1.String())),
			sdk.NewEvent("burn_ft", sdk.NewAttribute("amount", types.NewCoins(types.NewCoin("0000000100000000", sdk.NewInt(100))).String())),
		}
		verifyEventFunc(t, e, res.Events)
	}
}

func TestHandleMsgBurnFTFrom(t *testing.T) {
	ctx, h, contractID := prepareFT(t)

	sendMsg := types.NewMsgTransferFT(addr1, contractID, addr2, types.NewCoin("0000000100000000", sdk.NewInt(1000)))
	_, err := h(ctx, sendMsg)
	require.NoError(t, err)

	{
		// not approved
		burnMsg := types.NewMsgBurnFTFrom(addr1, contractID, addr2, types.NewCoin("0000000100000000", sdk.NewInt(100)))
		_, err := h(ctx, burnMsg)
		require.Error(t, err)
	}

	approve(t, addr2, addr1, contractID, ctx, h)
	{
		// invalid user
		burnMsg := types.NewMsgBurnFTFrom(addr2, contractID, addr2, types.NewCoin("0000000100000000", sdk.NewInt(100)))
		_, err := h(ctx, burnMsg)
		require.Error(t, err)
	}
	{
		// burn non-exist token
		burnMsg := types.NewMsgBurnFTFrom(addr1, contractID, addr2, types.NewCoin("0000000200000000", sdk.NewInt(100)))
		_, err := h(ctx, burnMsg)
		require.Error(t, err)
	}
	{
		// burn tokens over the being supplied
		burnMsg := types.NewMsgBurnFTFrom(addr1, contractID, addr2, types.NewCoin("0000000100000000", sdk.NewInt(1001)))
		_, err := h(ctx, burnMsg)
		require.Error(t, err)
	}
	{
		// burn tokens with invalid contractID
		burnMsg := types.NewMsgBurnFTFrom(addr1, "abcd11234", addr2, types.NewCoin("0000000100000000", sdk.NewInt(1000)))
		_, err := h(ctx, burnMsg)
		require.Error(t, err)
	}
	{
		// success case
		burnMsg := types.NewMsgBurnFTFrom(addr1, contractID, addr2, types.NewCoin("0000000100000000", sdk.NewInt(100)))
		res, err := h(ctx, burnMsg)
		require.NoError(t, err)

		e := sdk.Events{
			sdk.NewEvent("message", sdk.NewAttribute("module", "collection")),
			sdk.NewEvent("message", sdk.NewAttribute("sender", addr1.String())),
			sdk.NewEvent("burn_ft_from", sdk.NewAttribute("contract_id", contractID)),
			sdk.NewEvent("burn_ft_from", sdk.NewAttribute("proxy", addr1.String())),
			sdk.NewEvent("burn_ft_from", sdk.NewAttribute("from", addr2.String())),
			sdk.NewEvent("burn_ft_from", sdk.NewAttribute("amount", types.NewCoins(types.NewCoin("0000000100000000", sdk.NewInt(100))).String())),
		}
		verifyEventFunc(t, e, res.Events)
	}
}

func TestHandleMsgBurnNFT(t *testing.T) {
	ctx, h, contractID := prepareNFT(t, addr1)

	{
		// invalid user
		burnMsg := types.NewMsgBurnNFT(addr2, contractID, "1000000100000001")
		_, err := h(ctx, burnMsg)
		require.Error(t, err)
	}
	{
		// burn non-exist token
		burnMsg := types.NewMsgBurnNFT(addr1, contractID, "1000000200000001")
		_, err := h(ctx, burnMsg)
		require.Error(t, err)
	}
	{
		// burn tokens with invalid contractID
		burnMsg := types.NewMsgBurnNFT(addr1, "abcd11234", "1000000100000001")
		_, err := h(ctx, burnMsg)
		require.Error(t, err)
	}
	{
		// success case
		burnMsg := types.NewMsgBurnNFT(addr1, contractID, "1000000100000001")
		res, err := h(ctx, burnMsg)
		require.NoError(t, err)

		e := sdk.Events{
			sdk.NewEvent("message", sdk.NewAttribute("module", "collection")),
			sdk.NewEvent("message", sdk.NewAttribute("sender", addr1.String())),
			sdk.NewEvent("burn_nft", sdk.NewAttribute("contract_id", contractID)),
			sdk.NewEvent("burn_nft", sdk.NewAttribute("from", addr1.String())),
			sdk.NewEvent("burn_nft", sdk.NewAttribute("token_id", "1000000100000001")),
			sdk.NewEvent("operation_burn_nft", sdk.NewAttribute("token_id", "1000000100000001")),
		}
		verifyEventFunc(t, e, res.Events)
	}
	{
		// burn already burned
		burnMsg := types.NewMsgBurnNFT(addr1, contractID, "1000000100000001")
		_, err := h(ctx, burnMsg)
		require.Error(t, err)
	}
}

func TestHandleMsgBurnNFTFrom(t *testing.T) {
	ctx, h, contractID := prepareNFT(t, addr2)

	{
		// not approved
		burnMsg := types.NewMsgBurnNFTFrom(addr1, contractID, addr2, "1000000100000001")
		_, err := h(ctx, burnMsg)
		require.Error(t, err)
	}

	approve(t, addr2, addr1, contractID, ctx, h)
	{
		// invalid user
		burnMsg := types.NewMsgBurnNFTFrom(addr2, contractID, addr2, "1000000100000001")
		_, err := h(ctx, burnMsg)
		require.Error(t, err)
	}
	{
		// burn non-exist token
		burnMsg := types.NewMsgBurnNFTFrom(addr1, contractID, addr2, "1000000200000001")
		_, err := h(ctx, burnMsg)
		require.Error(t, err)
	}
	{
		// burn tokens with invalid contractID
		burnMsg := types.NewMsgBurnNFTFrom(addr1, "abcd11234", addr2, "1000000100000001")
		_, err := h(ctx, burnMsg)
		require.Error(t, err)
	}
	{
		// success case
		burnMsg := types.NewMsgBurnNFTFrom(addr1, contractID, addr2, "1000000100000001")
		res, err := h(ctx, burnMsg)
		require.NoError(t, err)

		e := sdk.Events{
			sdk.NewEvent("message", sdk.NewAttribute("module", "collection")),
			sdk.NewEvent("message", sdk.NewAttribute("sender", addr1.String())),
			sdk.NewEvent("burn_nft_from", sdk.NewAttribute("contract_id", contractID)),
			sdk.NewEvent("burn_nft_from", sdk.NewAttribute("proxy", addr1.String())),
			sdk.NewEvent("burn_nft_from", sdk.NewAttribute("from", addr2.String())),
			sdk.NewEvent("burn_nft_from", sdk.NewAttribute("token_id", "1000000100000001")),
			sdk.NewEvent("operation_burn_nft", sdk.NewAttribute("token_id", "1000000100000001")),
		}
		verifyEventFunc(t, e, res.Events)
	}
	{
		// burn already burned
		burnMsg := types.NewMsgBurnNFTFrom(addr1, contractID, addr2, "1000000100000001")
		_, err := h(ctx, burnMsg)
		require.Error(t, err)
	}
}

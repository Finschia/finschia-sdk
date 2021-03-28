package handler

import (
	"fmt"
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/line/lbm-sdk/v2/x/collection/internal/types"
	"github.com/stretchr/testify/require"
)

func TestHandleMsgMintFT(t *testing.T) {
	ctx, h, contractID := prepareFT(t)

	{
		burnMsg := types.NewMsgMintFT(addr1, contractID, addr1, types.NewCoin(defaultTokenIDFT, sdk.NewInt(100)))
		res, err := h(ctx, burnMsg)
		require.NoError(t, err)

		e := sdk.Events{
			sdk.NewEvent("message", sdk.NewAttribute("module", "collection")),
			sdk.NewEvent("message", sdk.NewAttribute("sender", addr1.String())),
			sdk.NewEvent("mint_ft", sdk.NewAttribute("contract_id", contractID)),
			sdk.NewEvent("mint_ft", sdk.NewAttribute("from", addr1.String())),
			sdk.NewEvent("mint_ft", sdk.NewAttribute("to", addr1.String())),
			sdk.NewEvent("mint_ft", sdk.NewAttribute("amount", types.NewCoins(types.NewCoin("0000000100000000", sdk.NewInt(100))).String())),
		}
		verifyEventFunc(t, e, res.Events)
	}
}

func TestHandleMsgMintNFT(t *testing.T) {
	ctx, h, contractID := prepareNFT(t, addr1)

	{
		param := types.NewMintNFTParam("shield", "", "10000001")
		msg := types.NewMsgMintNFT(addr1, contractID, addr1, param)
		res, err := h(ctx, msg)
		require.NoError(t, err)

		e := sdk.Events{
			sdk.NewEvent("message", sdk.NewAttribute("module", "collection")),
			sdk.NewEvent("message", sdk.NewAttribute("sender", addr1.String())),
			sdk.NewEvent("mint_nft", sdk.NewAttribute("contract_id", contractID)),
			sdk.NewEvent("mint_nft", sdk.NewAttribute("name", "shield")),
			sdk.NewEvent("mint_nft", sdk.NewAttribute("token_id", defaultTokenID3)),
			sdk.NewEvent("mint_nft", sdk.NewAttribute("from", addr1.String())),
			sdk.NewEvent("mint_nft", sdk.NewAttribute("to", addr1.String())),
		}
		verifyEventFunc(t, e, res.Events)
	}
}

func TestHandleMsgMintNFTPerformance(t *testing.T) {
	ctx, h, contractID := prepareNFT(t, addr1)
	var mean int64
	var sum int64 = 0
	{
		param := types.NewMintNFTParam("shield", "", "10000001")
		msg := types.NewMsgMintNFT(addr1, contractID, addr1, param)
		for jdx := 0; jdx < 10; jdx++ {
			startTime := time.Now()
			for idx := 0; idx < 1000; idx++ {
				_, err := h(ctx, msg)
				require.NoError(t, err)
			}
			duration := time.Since(startTime)
			sum += duration.Nanoseconds()
		}
		mean = sum / 10
	}

	ctx, h, contractID = prepareNFT(t, addr1)
	{
		param := types.NewMintNFTParam("shield", "", "10000001")
		msg := types.NewMsgMintNFT(addr1, contractID, addr1, param)
		for jdx := 0; jdx < 10; jdx++ {
			startTime := time.Now()
			for idx := 0; idx < 1000; idx++ {
				_, err := h(ctx, msg)
				require.NoError(t, err)
			}
			duration := time.Since(startTime)
			t.Log(fmt.Sprintf("MintNFT %s", duration.String()))
			require.Less(t, duration.Nanoseconds(), mean*2)
		}
	}
}

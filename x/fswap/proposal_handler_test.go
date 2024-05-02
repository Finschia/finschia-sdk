package fswap_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	"github.com/Finschia/finschia-sdk/simapp"
	sdk "github.com/Finschia/finschia-sdk/types"
	"github.com/Finschia/finschia-sdk/x/fswap"
	"github.com/Finschia/finschia-sdk/x/fswap/types"
)

func testProposal(swap types.Swap) *types.MakeSwapProposal {
	return types.NewMakeSwapProposal("Test", "description", swap)
}

func TestProposalHandlerPassed(t *testing.T) {
	app := simapp.Setup(false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})

	swap := types.Swap{
		FromDenom:           "aaa",
		ToDenom:             "bbb",
		AmountCapForToDenom: sdk.NewInt(100),
		SwapMultiple:        sdk.NewInt(10),
	}
	tp := testProposal(swap)
	hdlr := fswap.NewSwapHandler(app.FswapKeeper)
	require.NoError(t, hdlr(ctx, tp))

	// todo check contents
}

// todo check failed
// func TestProposalHandlerFailed(t *testing.T) {}

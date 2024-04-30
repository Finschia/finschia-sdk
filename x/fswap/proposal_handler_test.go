package fswap_test

import (
	"testing"

	"github.com/Finschia/finschia-sdk/x/fswap/types"
)

func testProposal(swapInit types.SwapInit) *types.SwapInitProposal {
	return types.NewSwapInitProposal("Test", "description", swapInit)
}

func TestProposalHandlerPassed(t *testing.T) {
	//app := simapp.Setup(false)
	//ctx := app.BaseApp.NewContext(false, tmproto.Header{})
	//
	//fswapInit := types.NewSwapInit("aaa", "bbb", sdk.NewInt(1000), sdk.NewInt(100000))
	//tp := testProposal(fswapInit)
	//hdlr := fswap.NewSwapInitHandler(app.FswapKeeper)
	//require.NoError(t, hdlr(ctx, tp))

	// todo check contents
}

// todo check failed
// func TestProposalHandlerFailed(t *testing.T) {}

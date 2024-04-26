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

func testProposal(fswapInit types.FswapInit) *types.FswapInitProposal {
	return types.NewFswapInitProposal("Test", "description", fswapInit)
}

func TestProposalHandlerPassed(t *testing.T) {
	app := simapp.Setup(false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})

	fswapInit := types.NewFswapInit("aaa", "bbb", sdk.NewInt(1000), sdk.NewInt(100000))
	tp := testProposal(fswapInit)
	hdlr := fswap.NewFswapInitHandler(app.FswapKeeper)
	require.NoError(t, hdlr(ctx, tp))

	// todo check contents
}

// todo check failed
// func TestProposalHandlerFailed(t *testing.T) {}

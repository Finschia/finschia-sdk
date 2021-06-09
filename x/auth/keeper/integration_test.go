package keeper_test

import (
	ostproto "github.com/line/ostracon/proto/ostracon/types"

	"github.com/line/lfb-sdk/simapp"
	sdk "github.com/line/lfb-sdk/types"
	authtypes "github.com/line/lfb-sdk/x/auth/types"
)

// returns context and app with params set on account keeper
func createTestApp(isCheckTx bool) (*simapp.SimApp, sdk.Context) {
	app := simapp.Setup(isCheckTx)
	ctx := app.BaseApp.NewContext(isCheckTx, ostproto.Header{})
	app.AccountKeeper.SetParams(ctx, authtypes.DefaultParams())

	return app, ctx
}

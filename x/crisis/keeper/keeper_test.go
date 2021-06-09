package keeper_test

import (
	"testing"

	abci "github.com/line/ostracon/abci/types"
	ostproto "github.com/line/ostracon/proto/ostracon/types"
	"github.com/stretchr/testify/require"

	"github.com/line/lfb-sdk/simapp"
	sdk "github.com/line/lfb-sdk/types"
)

func TestLogger(t *testing.T) {
	app := simapp.Setup(false)

	ctx := app.NewContext(true, ostproto.Header{})
	require.Equal(t, ctx.Logger(), app.CrisisKeeper.Logger(ctx))
}

func TestInvariants(t *testing.T) {
	app := simapp.Setup(false)
	app.Commit()
	app.BeginBlock(abci.RequestBeginBlock{Header: ostproto.Header{Height: app.LastBlockHeight() + 1}})

	require.Equal(t, app.CrisisKeeper.InvCheckPeriod(), uint(5))

	// SimApp has 11 registered invariants
	orgInvRoutes := app.CrisisKeeper.Routes()
	app.CrisisKeeper.RegisterRoute("testModule", "testRoute", func(sdk.Context) (string, bool) { return "", false })
	require.Equal(t, len(app.CrisisKeeper.Routes()), len(orgInvRoutes)+1)
}

func TestAssertInvariants(t *testing.T) {
	app := simapp.Setup(false)
	app.Commit()
	app.BeginBlock(abci.RequestBeginBlock{Header: ostproto.Header{Height: app.LastBlockHeight() + 1}})

	ctx := app.NewContext(true, ostproto.Header{})

	app.CrisisKeeper.RegisterRoute("testModule", "testRoute1", func(sdk.Context) (string, bool) { return "", false })
	require.NotPanics(t, func() { app.CrisisKeeper.AssertInvariants(ctx) })

	app.CrisisKeeper.RegisterRoute("testModule", "testRoute2", func(sdk.Context) (string, bool) { return "", true })
	require.Panics(t, func() { app.CrisisKeeper.AssertInvariants(ctx) })
}

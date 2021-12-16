package keeper_test

import (
	"testing"

	ocproto "github.com/line/ostracon/proto/ostracon/types"
	"github.com/stretchr/testify/require"

	"github.com/line/lbm-sdk/simapp"
	"github.com/line/lbm-sdk/x/consortium/types"
)

func TestGetSetParams(t *testing.T) {
	app := simapp.Setup(false)
	ctx := app.BaseApp.NewContext(false, ocproto.Header{})

	k := app.ConsortiumKeeper

	params := &types.Params{
		Enabled: true,
	}
	k.SetParams(ctx, params)
	require.Equal(t, params, k.GetParams(ctx))
	require.Equal(t, params.Enabled, k.GetEnabled(ctx))
}

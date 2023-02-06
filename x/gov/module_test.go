package gov_test

import (
	"testing"

	ocproto "github.com/line/ostracon/proto/ostracon/types"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/line/lbm-sdk/simapp"
	authtypes "github.com/line/lbm-sdk/x/auth/types"
	"github.com/line/lbm-sdk/x/gov/types"
)

func TestItCreatesModuleAccountOnInitBlock(t *testing.T) {
	app := simapp.Setup(false)
	ctx := app.BaseApp.NewContext(false, ocproto.Header{})

	app.InitChain(
		abci.RequestInitChain{
			AppStateBytes: []byte("{}"),
			ChainId:       "test-chain-id",
		},
	)

	acc := app.AccountKeeper.GetAccount(ctx, authtypes.NewModuleAddress(types.ModuleName))
	require.NotNil(t, acc)
}

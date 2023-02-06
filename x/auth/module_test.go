package auth_test

import (
	"testing"

	ocproto "github.com/line/ostracon/proto/ostracon/types"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/line/lbm-sdk/simapp"
	"github.com/line/lbm-sdk/x/auth/types"
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

	acc := app.AccountKeeper.GetAccount(ctx, types.NewModuleAddress(types.FeeCollectorName))
	require.NotNil(t, acc)
}

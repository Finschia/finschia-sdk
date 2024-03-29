package simapp

import (
	"testing"

	"github.com/Finschia/finschia-sdk/baseapp"
	"github.com/Finschia/finschia-sdk/client"
	"github.com/Finschia/finschia-sdk/codec"
	"github.com/Finschia/finschia-sdk/simapp/params"
	bankkeeper "github.com/Finschia/finschia-sdk/x/bank/keeper"
	stakingkeeper "github.com/Finschia/finschia-sdk/x/staking/keeper"
)

type TestSupport struct {
	t   testing.TB
	app *SimApp
}

func NewTestSupport(tb testing.TB, app *SimApp) *TestSupport {
	tb.Helper()
	return &TestSupport{t: tb, app: app}
}

func (s TestSupport) AppCodec() codec.Codec {
	return s.app.appCodec
}

func (s TestSupport) StakingKeeper() stakingkeeper.Keeper {
	return s.app.StakingKeeper
}

func (s TestSupport) BankKeeper() bankkeeper.Keeper {
	return s.app.BankKeeper
}

func (s TestSupport) GetBaseApp() *baseapp.BaseApp {
	return s.app.BaseApp
}

func (s TestSupport) GetTxConfig() client.TxConfig {
	return params.MakeTestEncodingConfig().TxConfig
}

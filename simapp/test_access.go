package simapp

import (
	"testing"

	"github.com/line/lbm-sdk/baseapp"
	"github.com/line/lbm-sdk/client"

	"github.com/line/lbm-sdk/simapp/params"

	"github.com/line/lbm-sdk/codec"
	bankkeeper "github.com/line/lbm-sdk/x/bank/keeper"
	stakingkeeper "github.com/line/lbm-sdk/x/staking/keeper"
)

type TestSupport struct {
	t   testing.TB
	app *SimApp
}

func NewTestSupport(t testing.TB, app *SimApp) *TestSupport {
	return &TestSupport{t: t, app: app}
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

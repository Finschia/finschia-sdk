package app

import (
	"testing"

	"github.com/line/lfb-sdk/codec"
	bankkeeper "github.com/line/lfb-sdk/x/bank/keeper"
	capabilitykeeper "github.com/line/lfb-sdk/x/capability/keeper"
	ibctransferkeeper "github.com/line/lfb-sdk/x/ibc/applications/transfer/keeper"
	ibckeeper "github.com/line/lfb-sdk/x/ibc/core/keeper"
	stakingkeeper "github.com/line/lfb-sdk/x/staking/keeper"
	"github.com/line/lfb-sdk/x/wasm"
)

type TestSupport struct {
	t   *testing.T
	app *LinkApp
}

func NewTestSupport(t *testing.T, app *LinkApp) *TestSupport {
	return &TestSupport{t: t, app: app}
}

func (s TestSupport) IBCKeeper() ibckeeper.Keeper {
	return *s.app.IBCKeeper
}

func (s TestSupport) WasmKeeper() wasm.Keeper {
	return s.app.wasmKeeper
}

func (s TestSupport) AppCodec() codec.Marshaler {
	return s.app.appCodec
}
func (s TestSupport) ScopedWasmIBCKeeper() capabilitykeeper.ScopedKeeper {
	return s.app.ScopedWasmKeeper
}

func (s TestSupport) ScopeIBCKeeper() capabilitykeeper.ScopedKeeper {
	return s.app.ScopedIBCKeeper
}

func (s TestSupport) ScopedTransferKeeper() capabilitykeeper.ScopedKeeper {
	return s.app.ScopedTransferKeeper
}

func (s TestSupport) StakingKeeper() stakingkeeper.Keeper {
	return s.app.StakingKeeper
}

func (s TestSupport) BankKeeper() bankkeeper.Keeper {
	return s.app.BankKeeper
}

func (s TestSupport) TransferKeeper() ibctransferkeeper.Keeper {
	return s.app.TransferKeeper
}

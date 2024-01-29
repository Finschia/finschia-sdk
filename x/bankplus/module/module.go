package bankplus

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/types/module"
	accountkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/bank/exported"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"

	"github.com/Finschia/finschia-sdk/x/bankplus/keeper"
)

var (
	_ module.AppModule           = AppModule{}
	_ module.AppModuleSimulation = AppModule{}
)

type AppModule struct {
	bank.AppModule

	bankKeeper     bankkeeper.Keeper
	legacySubspace exported.Subspace
}

func NewAppModule(cdc codec.Codec, keeper bankkeeper.Keeper, accountKeeper accountkeeper.AccountKeeper, ss exported.Subspace) AppModule {
	return AppModule{
		AppModule:      bank.NewAppModule(cdc, keeper, accountKeeper, ss),
		bankKeeper:     keeper,
		legacySubspace: ss,
	}
}

func (am AppModule) RegisterServices(cfg module.Configurator) {
	banktypes.RegisterMsgServer(cfg.MsgServer(), bankkeeper.NewMsgServerImpl(am.bankKeeper))
	banktypes.RegisterQueryServer(cfg.QueryServer(), am.bankKeeper)

	m := bankkeeper.NewMigrator(am.bankKeeper.(keeper.BaseKeeper).BaseKeeper, am.legacySubspace)
	if err := cfg.RegisterMigration(banktypes.ModuleName, 1, m.Migrate1to2); err != nil {
		panic(fmt.Sprintf("failed to migrate x/bank from version 1 to 2: %v", err))
	}
}

package bankplus

import (
	"fmt"

	"github.com/Finschia/finschia-sdk/codec"
	"github.com/Finschia/finschia-sdk/types/module"
	accountkeeper "github.com/Finschia/finschia-sdk/x/auth/keeper"
	"github.com/Finschia/finschia-sdk/x/bank"
	bankkeeper "github.com/Finschia/finschia-sdk/x/bank/keeper"
	banktypes "github.com/Finschia/finschia-sdk/x/bank/types"
	"github.com/Finschia/finschia-sdk/x/bankplus/keeper"
)

var (
	_ module.AppModule           = AppModule{}
	_ module.AppModuleSimulation = AppModule{}
)

type AppModule struct {
	bank.AppModule

	bankKeeper bankkeeper.Keeper
}

func NewAppModule(cdc codec.Codec, keeper bankkeeper.Keeper, accountKeeper accountkeeper.AccountKeeper) AppModule {
	return AppModule{
		AppModule:  bank.NewAppModule(cdc, keeper, accountKeeper),
		bankKeeper: keeper,
	}
}

func (am AppModule) RegisterServices(cfg module.Configurator) {
	banktypes.RegisterMsgServer(cfg.MsgServer(), bankkeeper.NewMsgServerImpl(am.bankKeeper))
	banktypes.RegisterQueryServer(cfg.QueryServer(), am.bankKeeper)

	m := bankkeeper.NewMigrator(am.bankKeeper.(keeper.BaseKeeper).BaseKeeper)
	if err := cfg.RegisterMigration(banktypes.ModuleName, 1, m.Migrate1to2); err != nil {
		panic(fmt.Sprintf("failed to migrate x/bank from version 1 to 2: %v", err))
	}
}

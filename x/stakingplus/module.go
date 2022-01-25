package stakingplus

import (
	"github.com/line/lbm-sdk/codec"
	"github.com/line/lbm-sdk/types/module"

	sdk "github.com/line/lbm-sdk/types"

	"github.com/line/lbm-sdk/x/stakingplus/keeper"
	"github.com/line/lbm-sdk/x/stakingplus/types"

	"github.com/line/lbm-sdk/x/staking"
	stakingkeeper "github.com/line/lbm-sdk/x/staking/keeper"
	stakingtypes "github.com/line/lbm-sdk/x/staking/types"
)

var (
	_ module.AppModule           = AppModule{}
	_ module.AppModuleSimulation = AppModule{}
)

// AppModule implements an application module for the stakingplus module.
type AppModule struct {
	staking.AppModule

	keeper stakingkeeper.Keeper
	ck     types.ConsortiumKeeper
}

// NewAppModule creates a new AppModule object
func NewAppModule(cdc codec.Marshaler, keeper stakingkeeper.Keeper, ak stakingtypes.AccountKeeper, bk stakingtypes.BankKeeper, ck types.ConsortiumKeeper) AppModule {
	return AppModule{
		AppModule: staking.NewAppModule(cdc, keeper, ak, bk),
		keeper:    keeper,
		ck:        ck,
	}
}

// RegisterServices registers module services.
func (am AppModule) RegisterServices(cfg module.Configurator) {
	stakingtypes.RegisterMsgServer(cfg.MsgServer(), keeper.NewMsgServerImpl(am.keeper, am.ck))
	querier := stakingkeeper.Querier{Keeper: am.keeper}
	stakingtypes.RegisterQueryServer(cfg.QueryServer(), querier)

	// m := keeper.NewMigrator(am.keeper)
	// if err := cfg.RegisterMigration(types.ModuleName, 1, m.Migrate1to2); err != nil {
	// 	panic(fmt.Sprintf("failed to migrate x/staking from version 1 to 2: %v", err))
	// }
}

// Route returns the message routing key for the stakingplus module.
func (am AppModule) Route() sdk.Route {
	return sdk.NewRoute(stakingtypes.RouterKey, NewHandler(am.keeper, am.ck))
}

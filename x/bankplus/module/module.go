package bankplus

import (
	"fmt"

	modulev1 "cosmossdk.io/api/cosmos/bank/module/v1"
	"cosmossdk.io/core/appmodule"
	"cosmossdk.io/core/store"
	"cosmossdk.io/depinject"
	"cosmossdk.io/log"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/types/module"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/bank/exported"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"

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

func NewAppModule(cdc codec.Codec, keeper bankkeeper.Keeper, accountKeeper banktypes.AccountKeeper, ss exported.Subspace) AppModule {
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

func init() {
	appmodule.Register(
		&modulev1.Module{},
		appmodule.Provide(ProvideModule),
	)
}

type ModuleInputs struct {
	depinject.In

	Config       *modulev1.Module
	Cdc          codec.Codec
	StoreService store.KVStoreService
	Logger       log.Logger

	AccountKeeper banktypes.AccountKeeper
	paramtypes.Subspace
	DeactMultiSend bool

	// LegacySubspace is used solely for migration of x/params managed parameters
	LegacySubspace exported.Subspace `optional:"true"`
}

type ModuleOutputs struct {
	depinject.Out
	Keeper keeper.Keeper
	Module appmodule.AppModule
}

func ProvideModule(in ModuleInputs) ModuleOutputs {
	// Configure blocked module accounts.
	//
	// Default behavior for blockedAddresses is to regard any module mentioned in
	// AccountKeeper's module account permissions as blocked.
	blockedAddresses := make(map[string]bool)
	if len(in.Config.BlockedModuleAccountsOverride) > 0 {
		for _, moduleName := range in.Config.BlockedModuleAccountsOverride {
			blockedAddresses[authtypes.NewModuleAddress(moduleName).String()] = true
		}
	} else {
		for _, permission := range in.AccountKeeper.GetModulePermissions() {
			blockedAddresses[permission.GetAddress().String()] = true
		}
	}

	// default to governance authority if not provided
	authority := authtypes.NewModuleAddress(govtypes.ModuleName)
	if in.Config.Authority != "" {
		authority = authtypes.NewModuleAddressOrBech32Address(in.Config.Authority)
	}

	k := keeper.NewBaseKeeper(
		in.Cdc,
		in.StoreService,
		in.AccountKeeper,
		blockedAddresses,
		in.DeactMultiSend,
		authority.String(),
		in.Logger,
	)
	m := NewAppModule(in.Cdc, k, in.AccountKeeper, in.LegacySubspace)
	return ModuleOutputs{Keeper: k, Module: m}
}

package bankplus

import (
	"encoding/json"
	"fmt"

	modulev1 "cosmossdk.io/api/cosmos/bank/module/v1"
	"cosmossdk.io/core/appmodule"
	"cosmossdk.io/core/store"
	"cosmossdk.io/depinject"
	"cosmossdk.io/log"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/bank/exported"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	"github.com/Finschia/finschia-sdk/x/bankplus/keeper"
)

var (
	_ module.AppModuleBasic = AppModule{}
	_ module.HasGenesis     = AppModule{}
	_ module.HasServices    = AppModule{}
	_ module.HasInvariants  = AppModule{}

	_ appmodule.AppModule = AppModule{}
)

// AppModuleBasic defines the basic application module used by the bankplus module.
type AppModuleBasic struct {
	bank.AppModuleBasic
}

// AppModule implements an application module for the bankplus module.
type AppModule struct {
	AppModuleBasic
	bankAppModule bank.AppModule

	bankKeeper     bankkeeper.Keeper
	accountKeeper  banktypes.AccountKeeper
	legacySubspace exported.Subspace
}

func NewAppModule(cdc codec.Codec, keeper bankkeeper.Keeper, accKeeper banktypes.AccountKeeper, ss exported.Subspace) AppModule {
	appModule := bank.NewAppModule(cdc, keeper, accKeeper, ss)
	return AppModule{
		AppModuleBasic: AppModuleBasic{
			AppModuleBasic: appModule.AppModuleBasic,
		},
		bankAppModule:  appModule,
		bankKeeper:     keeper,
		accountKeeper:  accKeeper,
		legacySubspace: ss,
	}
}

func (a AppModule) IsOnePerModuleType() {
}

func (a AppModule) IsAppModule() {
}

func (a AppModule) RegisterInvariants(ir sdk.InvariantRegistry) {
	a.bankAppModule.RegisterInvariants(ir)
}

func (a AppModule) RegisterServices(cfg module.Configurator) {
	banktypes.RegisterMsgServer(cfg.MsgServer(), bankkeeper.NewMsgServerImpl(a.bankKeeper))
	banktypes.RegisterQueryServer(cfg.QueryServer(), a.bankKeeper)
	m := bankkeeper.NewMigrator(a.bankKeeper.(keeper.BaseKeeper).BaseKeeper, a.legacySubspace)
	if err := cfg.RegisterMigration(banktypes.ModuleName, 1, m.Migrate1to2); err != nil {
		panic(fmt.Sprintf("failed to migrate x/bank from version 1 to 2: %v", err))
	}
}

func (a AppModule) InitGenesis(ctx sdk.Context, cdc codec.JSONCodec, data json.RawMessage) {
	a.bankAppModule.InitGenesis(ctx, cdc, data)
}

func (a AppModule) ExportGenesis(ctx sdk.Context, codec codec.JSONCodec) json.RawMessage {
	return a.bankAppModule.ExportGenesis(ctx, codec)
}

func (a AppModule) GenerateGenesisState(simState *module.SimulationState) {
	a.bankAppModule.GenerateGenesisState(simState)
}

func (a AppModule) RegisterStoreDecoder(sdr simtypes.StoreDecoderRegistry) {
	sdr[banktypes.StoreKey] = simtypes.NewStoreDecoderFuncFromCollectionsSchema(a.bankKeeper.(keeper.BaseKeeper).Schema)
}

func (a AppModule) WeightedOperations(simState module.SimulationState) []simtypes.WeightedOperation {
	return a.bankAppModule.WeightedOperations(simState)
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
	//DeactMultiSend bool // FIXME: inject properly

	// LegacySubspace is used solely for migration of x/params managed parameters
	LegacySubspace exported.Subspace `optional:"true"`
}

func ProvideModule(in ModuleInputs) bank.ModuleOutputs {
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

	bankKeeper := keeper.NewBaseKeeper(
		in.Cdc,
		in.StoreService,
		in.AccountKeeper,
		blockedAddresses,
		true, //in.DeactMultiSend, // FIXME: inject properly
		authority.String(),
		in.Logger,
	)
	m := NewAppModule(in.Cdc, bankKeeper, in.AccountKeeper, in.LegacySubspace)
	return bank.ModuleOutputs{
		BankKeeper: bankKeeper.BaseKeeper,
		Module:     m,
	}
}

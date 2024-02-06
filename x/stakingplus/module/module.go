package module

import (
	"context"
	"encoding/json"
	"fmt"

	abci "github.com/cometbft/cometbft/abci/types"

	modulev1 "cosmossdk.io/api/cosmos/staking/module/v1"
	"cosmossdk.io/core/appmodule"
	"cosmossdk.io/core/store"
	"cosmossdk.io/depinject"

	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/cosmos/cosmos-sdk/x/staking/exported"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/Finschia/finschia-sdk/x/stakingplus"
	"github.com/Finschia/finschia-sdk/x/stakingplus/keeper"
)

var (
	_ module.AppModuleBasic  = AppModuleBasic{}
	_ module.HasServices     = AppModule{}
	_ module.HasInvariants   = AppModule{}
	_ module.HasABCIGenesis  = AppModule{}
	_ module.HasABCIEndBlock = AppModule{}

	_ appmodule.AppModule       = AppModule{}
	_ appmodule.HasBeginBlocker = AppModule{}
)

// AppModuleBasic defines the basic application module used by the stakingplus module.
type AppModuleBasic struct {
	staking.AppModuleBasic
}

func (b AppModuleBasic) RegisterInterfaces(registry codectypes.InterfaceRegistry) {
	b.AppModuleBasic.RegisterInterfaces(registry)
	stakingplus.RegisterInterfaces(registry)
}

// ____________________________________________________________________________

// AppModule implements an application module for the stakingplus module.
type AppModule struct {
	AppModuleBasic
	impl staking.AppModule

	keeper *stakingkeeper.Keeper
	ak     stakingtypes.AccountKeeper
	bk     stakingtypes.BankKeeper
	fk     stakingplus.FoundationKeeper
	ls     exported.Subspace
}

// NewAppModule creates a new AppModule object
func NewAppModule(cdc codec.Codec, keeper *stakingkeeper.Keeper, ak stakingtypes.AccountKeeper, bk stakingtypes.BankKeeper, fk stakingplus.FoundationKeeper, ls exported.Subspace) AppModule {
	impl := staking.NewAppModule(cdc, keeper, ak, bk, ls)
	return AppModule{
		AppModuleBasic: AppModuleBasic{
			impl.AppModuleBasic,
		},
		impl:   impl,
		keeper: keeper,
		ak:     ak,
		bk:     bk,
		fk:     fk,
	}
}

// IsOnePerModuleType implements the depinject.OnePerModuleType interface.
func (am AppModule) IsOnePerModuleType() {
}

// IsAppModule implements the appmodule.AppModule interface.
func (am AppModule) IsAppModule() {}

// RegisterInvariants registers the staking module invariants.
func (am AppModule) RegisterInvariants(ir sdk.InvariantRegistry) {
	am.impl.RegisterInvariants(ir)
}

// RegisterServices registers module services.
func (am AppModule) RegisterServices(cfg module.Configurator) {
	stakingtypes.RegisterMsgServer(cfg.MsgServer(), keeper.NewMsgServerImpl(stakingkeeper.NewMsgServerImpl(am.keeper), am.fk, am.keeper.ValidatorAddressCodec()))
	querier := stakingkeeper.Querier{Keeper: am.keeper}
	stakingtypes.RegisterQueryServer(cfg.QueryServer(), querier)

	m := stakingkeeper.NewMigrator(am.keeper, am.ls)
	if err := cfg.RegisterMigration(stakingtypes.ModuleName, 1, m.Migrate1to2); err != nil {
		panic(fmt.Sprintf("failed to migrate x/%s from version 1 to 2: %v", stakingtypes.ModuleName, err))
	}
	if err := cfg.RegisterMigration(stakingtypes.ModuleName, 2, m.Migrate2to3); err != nil {
		panic(fmt.Sprintf("failed to migrate x/%s from version 2 to 3: %v", stakingtypes.ModuleName, err))
	}
	if err := cfg.RegisterMigration(stakingtypes.ModuleName, 3, m.Migrate3to4); err != nil {
		panic(fmt.Sprintf("failed to migrate x/%s from version 3 to 4: %v", stakingtypes.ModuleName, err))
	}
	if err := cfg.RegisterMigration(stakingtypes.ModuleName, 4, m.Migrate4to5); err != nil {
		panic(fmt.Sprintf("failed to migrate x/%s from version 4 to 5: %v", stakingtypes.ModuleName, err))
	}
}

// InitGenesis performs genesis initialization for the stakingplus module. It returns
// no validator updates.
func (am AppModule) InitGenesis(ctx sdk.Context, cdc codec.JSONCodec, data json.RawMessage) []abci.ValidatorUpdate {
	return am.impl.InitGenesis(ctx, cdc, data)
}

// ExportGenesis returns the exported genesis state as raw bytes for the stakingplus
// module.
func (am AppModule) ExportGenesis(ctx sdk.Context, cdc codec.JSONCodec) json.RawMessage {
	return am.impl.ExportGenesis(ctx, cdc)
}

// ConsensusVersion implements AppModule/ConsensusVersion.
func (am AppModule) ConsensusVersion() uint64 {
	return am.impl.ConsensusVersion()
}

// BeginBlock returns the begin blocker for the stakingplus module.
func (am AppModule) BeginBlock(ctx context.Context) error {
	return am.impl.BeginBlock(ctx)
}

// EndBlock returns the end blocker for the stakingplus module. It returns no validator
// updates.
func (am AppModule) EndBlock(ctx context.Context) ([]abci.ValidatorUpdate, error) {
	return am.impl.EndBlock(ctx)
}

func init() {
	appmodule.Register(
		&modulev1.Module{},
		appmodule.Provide(ProvideModule),
		appmodule.Invoke(staking.InvokeSetStakingHooks),
	)
}

type StakingplusInputs struct {
	depinject.In

	Config                *modulev1.Module
	ValidatorAddressCodec runtime.ValidatorAddressCodec
	ConsensusAddressCodec runtime.ConsensusAddressCodec
	AccountKeeper         stakingtypes.AccountKeeper
	BankKeeper            stakingtypes.BankKeeper
	FoundationKeeper      stakingplus.FoundationKeeper
	Cdc                   codec.Codec
	StoreService          store.KVStoreService

	// LegacySubspace is used solely for migration of x/params managed parameters
	LegacySubspace exported.Subspace `optional:"true"`
}

func ProvideModule(in StakingplusInputs) staking.ModuleOutputs {
	// default to governance authority if not provided
	authority := authtypes.NewModuleAddress(govtypes.ModuleName)

	if in.Config.Authority != "" {
		authority = authtypes.NewModuleAddressOrBech32Address(in.Config.Authority)
	}

	k := stakingkeeper.NewKeeper(
		in.Cdc,
		in.StoreService,
		in.AccountKeeper,
		in.BankKeeper,
		authority.String(),
		in.ValidatorAddressCodec,
		in.ConsensusAddressCodec,
	)
	m := NewAppModule(in.Cdc, k, in.AccountKeeper, in.BankKeeper, in.FoundationKeeper, in.LegacySubspace)
	return staking.ModuleOutputs{StakingKeeper: k, Module: m}
}

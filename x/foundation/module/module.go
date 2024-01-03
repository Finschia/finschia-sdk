package module

import (
	"context"
	"encoding/json"
	"fmt"

	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/spf13/cobra"

	"cosmossdk.io/core/appmodule"
	"cosmossdk.io/core/store"
	"cosmossdk.io/depinject"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"

	modulev1 "github.com/Finschia/finschia-sdk/api/lbm/foundation/module/v1"
	"github.com/Finschia/finschia-sdk/x/foundation"
	"github.com/Finschia/finschia-sdk/x/foundation/client/cli"
	"github.com/Finschia/finschia-sdk/x/foundation/keeper"
)

const (
	consensusVersion uint64 = 2
)

var (
	_ module.AppModule          = AppModule{}
	_ module.AppModuleBasic     = AppModuleBasic{}
	_ appmodule.HasBeginBlocker = AppModule{}
	_ appmodule.HasEndBlocker   = AppModule{}
)

// AppModuleBasic defines the basic application module used by the foundation module.
type AppModuleBasic struct {
	cdc codec.Codec
}

// Name returns the ModuleName
func (AppModuleBasic) Name() string {
	return foundation.ModuleName
}

// RegisterLegacyAminoCodec registers the foundation types on the LegacyAmino codec
func (AppModuleBasic) RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {}

// DefaultGenesis returns default genesis state as raw bytes for the foundation
// module.
func (AppModuleBasic) DefaultGenesis(cdc codec.JSONCodec) json.RawMessage {
	return cdc.MustMarshalJSON(foundation.DefaultGenesisState())
}

// ValidateGenesis performs genesis state validation for the foundation module.
func (am AppModuleBasic) ValidateGenesis(cdc codec.JSONCodec, config client.TxEncodingConfig, bz json.RawMessage) error {
	var data foundation.GenesisState
	if err := cdc.UnmarshalJSON(bz, &data); err != nil {
		return fmt.Errorf("failed to unmarshal %s genesis state: %w", foundation.ModuleName, err)
	}

	return foundation.ValidateGenesis(data, am.cdc.InterfaceRegistry().SigningContext().AddressCodec())
}

// RegisterGRPCGatewayRoutes registers the gRPC Gateway routes for the foundation module.
func (AppModuleBasic) RegisterGRPCGatewayRoutes(clientCtx client.Context, mux *runtime.ServeMux) {
	if err := foundation.RegisterQueryHandlerClient(context.Background(), mux, foundation.NewQueryClient(clientCtx)); err != nil {
		panic(err)
	}
}

// GetQueryCmd returns the cli query commands for this module
func (AppModuleBasic) GetQueryCmd() *cobra.Command {
	return cli.NewQueryCmd()
}

// GetTxCmd returns the transaction commands for this module
func (AppModuleBasic) GetTxCmd() *cobra.Command {
	return cli.NewTxCmd()
}

func (am AppModuleBasic) RegisterInterfaces(registry codectypes.InterfaceRegistry) {
	foundation.RegisterInterfaces(registry)
}

// ____________________________________________________________________________

// AppModule implements an application module for the foundation module.
type AppModule struct {
	AppModuleBasic

	keeper keeper.Keeper
}

// NewAppModule creates a new AppModule object
func NewAppModule(cdc codec.Codec, keeper keeper.Keeper) AppModule {
	return AppModule{
		AppModuleBasic: AppModuleBasic{
			cdc: cdc,
		},
		keeper: keeper,
	}
}

// RegisterInvariants does nothing, there are no invariants to enforce
func (am AppModule) RegisterInvariants(ir sdk.InvariantRegistry) {
	keeper.RegisterInvariants(ir, am.keeper)
}

// RegisterServices registers a GRPC query service to respond to the
// module-specific GRPC queries.
func (am AppModule) RegisterServices(cfg module.Configurator) {
	foundation.RegisterMsgServer(cfg.MsgServer(), keeper.NewMsgServer(am.keeper))
	foundation.RegisterQueryServer(cfg.QueryServer(), keeper.NewQueryServer(am.keeper))
	err := keeper.NewMigrator(am.keeper).Register(cfg.RegisterMigration)
	if err != nil {
		panic(err)
	}
}

// InitGenesis performs genesis initialization for the foundation module. It returns
// no validator updates.
func (am AppModule) InitGenesis(ctx sdk.Context, cdc codec.JSONCodec, data json.RawMessage) []abci.ValidatorUpdate {
	var genesisState foundation.GenesisState
	cdc.MustUnmarshalJSON(data, &genesisState)
	if err := am.keeper.InitGenesis(ctx, &genesisState); err != nil {
		panic(err)
	}
	return []abci.ValidatorUpdate{}
}

// ExportGenesis returns the exported genesis state as raw bytes for the foundation
// module.
func (am AppModule) ExportGenesis(ctx sdk.Context, cdc codec.JSONCodec) json.RawMessage {
	gs := am.keeper.ExportGenesis(ctx)
	return cdc.MustMarshalJSON(gs)
}

// ConsensusVersion implements AppModule/ConsensusVersion.
func (AppModule) ConsensusVersion() uint64 { return consensusVersion }

// BeginBlock performs a no-op.
func (am AppModule) BeginBlock(ctx context.Context) error {
	return keeper.BeginBlocker(sdk.UnwrapSDKContext(ctx), am.keeper)
}

// EndBlock performs a no-op.
func (am AppModule) EndBlock(ctx context.Context) error {
	return keeper.EndBlocker(sdk.UnwrapSDKContext(ctx), am.keeper)
}

// ____________________________________________________________________________

// AppModuleSimulation functions

// // GenerateGenesisState creates a randomized GenState of the foundation module.
// func (AppModule) GenerateGenesisState(simState *module.SimulationState) {
// 	simulation.RandomizedGenState(simState)
// }

// // ProposalContents returns all the foundation content functions used to
// // simulate foundation proposals.
// func (AppModule) ProposalContents(simState module.SimulationState) []simtypes.WeightedProposalContent {
// 	return simulation.ProposalContents()
// }

// // RandomizedParams creates randomized foundation param changes for the simulator.
// func (AppModule) RandomizedParams(r *rand.Rand) []simtypes.ParamChange {
// 	return simulation.ParamChanges(r)
// }

// // RegisterStoreDecoder registers a decoder for foundation module's types
// func (am AppModule) RegisterStoreDecoder(sdr sdk.StoreDecoderRegistry) {
// 	sdr[types.StoreKey] = simulation.NewDecodeStore(am.cdc)
// }

// // WeightedOperations returns the all the foundation module operations with their respective weights.
// func (am AppModule) WeightedOperations(simState module.SimulationState) []simtypes.WeightedOperation {
// 	return simulation.WeightedOperations(
// 		simState.AppParams, simState.Cdc,
// 		am.stakingKeeper, am.slashingKeeper, am.keeper, simState.Contents,
// 	)
// }

// IsOnePerModuleType implements the depinject.OnePerModuleType interface.
func (am AppModule) IsOnePerModuleType() {}

// IsAppModule implements the appmodule.AppModule interface.
func (am AppModule) IsAppModule() {}

func init() {
	appmodule.Register(&modulev1.Module{},
		appmodule.Provide(ProvideModule, ProvideKeyTable),
	)
}

type FoundationInputs struct {
	depinject.In

	StoreService store.KVStoreService
	Cdc          codec.Codec
	Config       *modulev1.Module

	AuthKeeper       foundation.AuthKeeper
	BankKeeper       foundation.BankKeeper
	Subspace         paramstypes.Subspace
	MsgServiceRouter baseapp.MessageRouter
}

type FoundationOutputs struct {
	depinject.Out

	Keeper keeper.Keeper
	Module appmodule.AppModule
}

func ProvideModule(in FoundationInputs) FoundationOutputs {
	// default to governance authority if not provided
	authority := authtypes.NewModuleAddress(govtypes.ModuleName)
	if in.Config.Authority != "" {
		authority = authtypes.NewModuleAddressOrBech32Address(in.Config.Authority)
	}

	feeCollectorName := in.Config.FeeCollectorName
	if feeCollectorName == "" {
		feeCollectorName = authtypes.FeeCollectorName
	}

	config := foundation.Config{
		MaxExecutionPeriod: in.Config.MaxExecutionPeriod.AsDuration(),
		MaxMetadataLen:     in.Config.MaxMetadataLen,
	}

	addressCodec := in.Cdc.InterfaceRegistry().SigningContext().AddressCodec()
	authorityStr, err := addressCodec.BytesToString(authority)
	if err != nil {
		panic(err)
	}

	k := keeper.NewKeeper(in.Cdc, in.StoreService, in.MsgServiceRouter, in.AuthKeeper, in.BankKeeper, feeCollectorName, config, authorityStr, in.Subspace)
	m := NewAppModule(in.Cdc, k)

	return FoundationOutputs{
		Keeper: k,
		Module: m,
	}
}

func ProvideKeyTable() paramstypes.KeyTable {
	return foundation.ParamKeyTable()
}

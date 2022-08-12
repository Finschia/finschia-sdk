package lbm

import (
	"math/rand"

	"encoding/json"
	"github.com/line/lbm-sdk/codec"
	sdk "github.com/line/lbm-sdk/types"
	"github.com/line/lbm-sdk/types/module"
	simtypes "github.com/line/lbm-sdk/types/simulation"
	simKeeper "github.com/line/lbm-sdk/x/simulation"
	"github.com/line/lbm-sdk/x/wasm"
	wasmkeeper "github.com/line/lbm-sdk/x/wasm/keeper"
	lbmwasmkeeper "github.com/line/lbm-sdk/x/wasm/lbm/keeper"
	lbmwasmtypes "github.com/line/lbm-sdk/x/wasm/lbm/types"
	wasmsim "github.com/line/lbm-sdk/x/wasm/simulation"
	wasmsimulation "github.com/line/lbm-sdk/x/wasm/simulation"
	wasmtypes "github.com/line/lbm-sdk/x/wasm/types"
	abci "github.com/line/ostracon/abci/types"
)

// AppModule implements an application module for the wasm module.
type AppModule struct {
	wasm.AppModuleBasic
	cdc                codec.Codec
	keeper             *lbmwasmkeeper.Keeper
	validatorSetSource wasmkeeper.ValidatorSetSource
	accountKeeper      wasmtypes.AccountKeeper
	bankKeeper         simKeeper.BankKeeper
}

// ConsensusVersion is a sequence number for state-breaking change of the
// module. It should be incremented on each consensus-breaking change
// introduced by the module. To avoid wrong/empty versions, the initial version
// should be set to 1.
func (AppModule) ConsensusVersion() uint64 { return 1 }

// NewAppModule creates a new AppModule object
func NewAppModule(
	cdc codec.Codec,
	keeper *lbmwasmkeeper.Keeper,
	validatorSetSource wasmkeeper.ValidatorSetSource,
	ak wasmtypes.AccountKeeper,
	bk simKeeper.BankKeeper,
) AppModule {
	return AppModule{
		AppModuleBasic:     wasm.AppModuleBasic{},
		cdc:                cdc,
		keeper:             keeper,
		validatorSetSource: validatorSetSource,
		accountKeeper:      ak,
		bankKeeper:         bk,
	}
}

// RegisterServices registers all services.
func (am AppModule) RegisterServices(cfg module.Configurator) {
	wasmtypes.RegisterMsgServer(cfg.MsgServer(), wasmkeeper.NewMsgServerImpl(wasmkeeper.NewDefaultPermissionKeeper(am.keeper)))
	wasmtypes.RegisterQueryServer(cfg.QueryServer(), wasm.NewQuerier(&am.keeper.Keeper))
	lbmwasmtypes.RegisterQueryServer(cfg.QueryServer(), lbmwasmkeeper.Querier(am.keeper))
}

// LegacyQuerierHandler returns the auth module sdk.Querier.
func (am AppModule) LegacyQuerierHandler(*codec.LegacyAmino) sdk.Querier {
	return wasmkeeper.NewLegacyQuerier(am.keeper.Keeper, am.keeper.QueryGasLimit())
}

// RegisterInvariants registers the wasm module invariants.
func (am AppModule) RegisterInvariants(ir sdk.InvariantRegistry) {}

// Route returns the message routing key for the wasm module.
func (am AppModule) Route() sdk.Route {
	return sdk.NewRoute(wasm.RouterKey, wasm.NewHandler(wasmkeeper.NewDefaultPermissionKeeper(am.keeper)))
}

// QuerierRoute returns the wasm module's querier route name.
func (AppModule) QuerierRoute() string {
	return wasm.QuerierRoute
}

// InitGenesis performs genesis initialization for the wasm module. It returns
// no validator updates.
func (am AppModule) InitGenesis(ctx sdk.Context, cdc codec.JSONCodec, data json.RawMessage) []abci.ValidatorUpdate {
	var genesisState lbmwasmtypes.GenesisState
	cdc.MustUnmarshalJSON(data, &genesisState)
	validators, err := wasmkeeper.InitGenesis(ctx, &am.keeper.Keeper, genesisState, am.validatorSetSource, am.Route().Handler())
	if err != nil {
		panic(err)
	}
	return validators
}

// ExportGenesis returns the exported genesis state as raw bytes for the wasm module.
func (am AppModule) ExportGenesis(ctx sdk.Context, cdc codec.JSONCodec) json.RawMessage {
	gs := wasmkeeper.ExportGenesis(ctx, &am.keeper.Keeper)
	return cdc.MustMarshalJSON(gs)
}

// BeginBlock returns the begin blocker for the wasm module.
func (am AppModule) BeginBlock(_ sdk.Context, _ abci.RequestBeginBlock) {}

// EndBlock returns the end blocker for the wasm module. It returns no validator
// updates.
func (AppModule) EndBlock(_ sdk.Context, _ abci.RequestEndBlock) []abci.ValidatorUpdate {
	return []abci.ValidatorUpdate{}
}

// ____________________________________________________________________________

// AppModuleSimulation functions

// GenerateGenesisState creates a randomized GenState of the bank module.
func (AppModule) GenerateGenesisState(simState *module.SimulationState) {
	wasmsim.RandomizedGenState(simState)
}

// ProposalContents doesn't return any content functions for governance proposals.
func (AppModule) ProposalContents(simState module.SimulationState) []simtypes.WeightedProposalContent {
	return nil
}

// RandomizedParams creates randomized bank param changes for the simulator.
func (am AppModule) RandomizedParams(r *rand.Rand) []simtypes.ParamChange {
	return wasmsimulation.ParamChanges(r, am.cdc)
}

// RegisterStoreDecoder registers a decoder for supply module's types
func (am AppModule) RegisterStoreDecoder(sdr sdk.StoreDecoderRegistry) {
}

// WeightedOperations returns the all the gov module operations with their respective weights.
func (am AppModule) WeightedOperations(simState module.SimulationState) []simtypes.WeightedOperation {
	return wasmsimulation.WeightedOperations(&simState, am.accountKeeper, am.bankKeeper, am.keeper)
}

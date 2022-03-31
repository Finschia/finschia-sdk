package module

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/gorilla/mux"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	abci "github.com/line/ostracon/abci/types"
	"github.com/spf13/cobra"

	"github.com/line/lbm-sdk/client"
	"github.com/line/lbm-sdk/codec"
	codectypes "github.com/line/lbm-sdk/codec/types"
	sdk "github.com/line/lbm-sdk/types"
	"github.com/line/lbm-sdk/types/module"
	"github.com/line/lbm-sdk/x/consortium"
	"github.com/line/lbm-sdk/x/consortium/client/cli"
	"github.com/line/lbm-sdk/x/consortium/keeper"
)

var (
	_ module.AppModule      = AppModule{}
	_ module.AppModuleBasic = AppModuleBasic{}
)

// AppModuleBasic defines the basic application module used by the consortium module.
type AppModuleBasic struct {
	cdc codec.Marshaler
}

// NewAppModuleBasic creates a new AppModuleBasic object
func NewAppModuleBasic() AppModuleBasic {
	return AppModuleBasic{}
}

// Name returns the ModuleName
func (AppModuleBasic) Name() string {
	return consortium.ModuleName
}

// RegisterLegacyAminoCodec registers the consortium types on the LegacyAmino codec
func (AppModuleBasic) RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {}

// DefaultGenesis returns default genesis state as raw bytes for the consortium
// module.
func (AppModuleBasic) DefaultGenesis(cdc codec.JSONMarshaler) json.RawMessage {
	return cdc.MustMarshalJSON(consortium.DefaultGenesisState())
}

// ValidateGenesis performs genesis state validation for the consortium module.
func (AppModuleBasic) ValidateGenesis(cdc codec.JSONMarshaler, config client.TxEncodingConfig, bz json.RawMessage) error {
	var data consortium.GenesisState
	if err := cdc.UnmarshalJSON(bz, &data); err != nil {
		return fmt.Errorf("failed to unmarshal %s genesis state: %w", consortium.ModuleName, err)
	}

	return consortium.ValidateGenesis(data)
}

// RegisterRESTRoutes registers all REST query handlers
func (AppModuleBasic) RegisterRESTRoutes(clientCtx client.Context, r *mux.Router) {}

// RegisterGRPCGatewayRoutes registers the gRPC Gateway routes for the consortium module.
func (AppModuleBasic) RegisterGRPCGatewayRoutes(clientCtx client.Context, mux *runtime.ServeMux) {
	if err := consortium.RegisterQueryHandlerClient(context.Background(), mux, consortium.NewQueryClient(clientCtx)); err != nil {
		panic(err)
	}
}

// GetQueryCmd returns the cli query commands for this module
func (AppModuleBasic) GetQueryCmd() *cobra.Command {
	return cli.NewQueryCmd()
}

// GetTxCmd returns the transaction commands for this module
func (AppModuleBasic) GetTxCmd() *cobra.Command {
	return nil
}

func (b AppModuleBasic) RegisterInterfaces(registry codectypes.InterfaceRegistry) {
	consortium.RegisterInterfaces(registry)
}

//____________________________________________________________________________

// AppModule implements an application module for the consortium module.
type AppModule struct {
	AppModuleBasic

	keeper        keeper.Keeper
	stakingKeeper consortium.StakingKeeper
}

// NewAppModule creates a new AppModule object
func NewAppModule(cdc codec.Marshaler, keeper keeper.Keeper, stk consortium.StakingKeeper) AppModule {
	return AppModule{
		AppModuleBasic: AppModuleBasic{cdc: cdc},
		keeper:         keeper,
		stakingKeeper:  stk,
	}
}

// RegisterInvariants does nothing, there are no invariants to enforce
func (AppModule) RegisterInvariants(_ sdk.InvariantRegistry) {}

// Route is empty, as we do not handle Messages (just proposals)
func (AppModule) Route() sdk.Route { return sdk.Route{} }

// QuerierRoute returns the route we respond to for abci queries
func (AppModule) QuerierRoute() string { return "" }

// LegacyQuerierHandler registers a query handler to respond to the module-specific queries
func (am AppModule) LegacyQuerierHandler(legacyQuerierCdc *codec.LegacyAmino) sdk.Querier {
	return nil
}

// RegisterServices registers a GRPC query service to respond to the
// module-specific GRPC queries.
func (am AppModule) RegisterServices(cfg module.Configurator) {
	consortium.RegisterQueryServer(cfg.QueryServer(), keeper.NewQueryServer(am.keeper))

	/* m := keeper.NewMigrator(am.keeper)
	if err := cfg.RegisterMigration(consortium.ModuleName, 1, m.Migrate1to2); err != nil {
		panic(fmt.Sprintf("failed to migrate x/consortium from version 1 to 2: %v", err))
	} */
}

// InitGenesis performs genesis initialization for the consortium module. It returns
// no validator updates.
func (am AppModule) InitGenesis(ctx sdk.Context, cdc codec.JSONMarshaler, data json.RawMessage) []abci.ValidatorUpdate {
	var genesisState consortium.GenesisState
	cdc.MustUnmarshalJSON(data, &genesisState)
	am.keeper.InitGenesis(ctx, am.stakingKeeper, &genesisState)
	return []abci.ValidatorUpdate{}
}

// ExportGenesis returns the exported genesis state as raw bytes for the consortium
// module.
func (am AppModule) ExportGenesis(ctx sdk.Context, cdc codec.JSONMarshaler) json.RawMessage {
	gs := am.keeper.ExportGenesis(ctx)
	return cdc.MustMarshalJSON(gs)
}

// ConsensusVersion implements AppModule/ConsensusVersion.
func (AppModule) ConsensusVersion() uint64 { return 1 }

// BeginBlock performs a no-op.
func (am AppModule) BeginBlock(_ sdk.Context, _ abci.RequestBeginBlock) {}

// EndBlock performs a no-op.
func (am AppModule) EndBlock(ctx sdk.Context, _ abci.RequestEndBlock) []abci.ValidatorUpdate {
	return []abci.ValidatorUpdate{}
}

//____________________________________________________________________________

// AppModuleSimulation functions

// // GenerateGenesisState creates a randomized GenState of the consortium module.
// func (AppModule) GenerateGenesisState(simState *module.SimulationState) {
// 	simulation.RandomizedGenState(simState)
// }

// // ProposalContents returns all the consortium content functions used to
// // simulate consortium proposals.
// func (AppModule) ProposalContents(simState module.SimulationState) []simtypes.WeightedProposalContent {
// 	return simulation.ProposalContents()
// }

// // RandomizedParams creates randomized consortium param changes for the simulator.
// func (AppModule) RandomizedParams(r *rand.Rand) []simtypes.ParamChange {
// 	return simulation.ParamChanges(r)
// }

// // RegisterStoreDecoder registers a decoder for consortium module's types
// func (am AppModule) RegisterStoreDecoder(sdr sdk.StoreDecoderRegistry) {
// 	sdr[types.StoreKey] = simulation.NewDecodeStore(am.cdc)
// }

// // WeightedOperations returns the all the consortium module operations with their respective weights.
// func (am AppModule) WeightedOperations(simState module.SimulationState) []simtypes.WeightedOperation {
// 	return simulation.WeightedOperations(
// 		simState.AppParams, simState.Cdc,
// 		am.stakingKeeper, am.slashingKeeper, am.keeper, simState.Contents,
// 	)
// }

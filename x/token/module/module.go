package module

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/spf13/cobra"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/line/lbm-sdk/client"
	"github.com/line/lbm-sdk/codec"
	codectypes "github.com/line/lbm-sdk/codec/types"
	sdk "github.com/line/lbm-sdk/types"
	"github.com/line/lbm-sdk/types/module"
	"github.com/line/lbm-sdk/x/token"
	"github.com/line/lbm-sdk/x/token/client/cli"
	"github.com/line/lbm-sdk/x/token/keeper"
)

var (
	_ module.AppModule      = AppModule{}
	_ module.AppModuleBasic = AppModuleBasic{}
)

// AppModuleBasic defines the basic application module used by the token module.
type AppModuleBasic struct{}

// Name returns the ModuleName
func (AppModuleBasic) Name() string {
	return token.ModuleName
}

// RegisterLegacyAminoCodec registers the token types on the LegacyAmino codec
func (AppModuleBasic) RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {}

// DefaultGenesis returns default genesis state as raw bytes for the token
// module.
func (AppModuleBasic) DefaultGenesis(cdc codec.JSONCodec) json.RawMessage {
	return cdc.MustMarshalJSON(token.DefaultGenesisState())
}

// ValidateGenesis performs genesis state validation for the token module.
func (AppModuleBasic) ValidateGenesis(cdc codec.JSONCodec, config client.TxEncodingConfig, bz json.RawMessage) error {
	var data token.GenesisState
	if err := cdc.UnmarshalJSON(bz, &data); err != nil {
		return fmt.Errorf("failed to unmarshal %s genesis state: %w", token.ModuleName, err)
	}

	return token.ValidateGenesis(data)
}

// RegisterGRPCGatewayRoutes registers the gRPC Gateway routes for the token module.
func (AppModuleBasic) RegisterGRPCGatewayRoutes(clientCtx client.Context, mux *runtime.ServeMux) {
	if err := token.RegisterQueryHandlerClient(context.Background(), mux, token.NewQueryClient(clientCtx)); err != nil {
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

func (b AppModuleBasic) RegisterInterfaces(registry codectypes.InterfaceRegistry) {
	token.RegisterInterfaces(registry)
}

//____________________________________________________________________________

// AppModule implements an application module for the token module.
type AppModule struct {
	AppModuleBasic

	keeper keeper.Keeper
}

// NewAppModule creates a new AppModule object
func NewAppModule(cdc codec.Codec, keeper keeper.Keeper) AppModule {
	return AppModule{
		keeper: keeper,
	}
}

// RegisterInvariants does nothing, there are no invariants to enforce
func (AppModule) RegisterInvariants(_ sdk.InvariantRegistry) {}

// Route returns the message routing key for the token module.
func (am AppModule) Route() sdk.Route { return sdk.Route{} }

// QuerierRoute returns the route we respond to for abci queries
func (AppModule) QuerierRoute() string { return "" }

// LegacyQuerierHandler registers a query handler to respond to the module-specific queries
func (am AppModule) LegacyQuerierHandler(legacyQuerierCdc *codec.LegacyAmino) sdk.Querier {
	return nil
}

// RegisterServices registers a GRPC query service to respond to the
// module-specific GRPC queries.
func (am AppModule) RegisterServices(cfg module.Configurator) {
	token.RegisterMsgServer(cfg.MsgServer(), keeper.NewMsgServer(am.keeper))
	token.RegisterQueryServer(cfg.QueryServer(), keeper.NewQueryServer(am.keeper))

	// m := keeper.NewMigrator(am.keeper)
	// migrations := map[uint64]func(sdk.Context) error{}
	// for ver, handler := range migrations {
	// 	if err := cfg.RegisterMigration(token.ModuleName, ver, handler); err != nil {
	// 		panic(fmt.Sprintf("failed to migrate x/%s from version %d to %d: %v", token.ModuleName, ver, ver+1, err))
	// 	}
	// }
}

// InitGenesis performs genesis initialization for the token module. It returns
// no validator updates.
func (am AppModule) InitGenesis(ctx sdk.Context, cdc codec.JSONCodec, data json.RawMessage) []abci.ValidatorUpdate {
	var genesisState token.GenesisState
	cdc.MustUnmarshalJSON(data, &genesisState)
	am.keeper.InitGenesis(ctx, &genesisState)
	return []abci.ValidatorUpdate{}
}

// ExportGenesis returns the exported genesis state as raw bytes for the token
// module.
func (am AppModule) ExportGenesis(ctx sdk.Context, cdc codec.JSONCodec) json.RawMessage {
	gs := am.keeper.ExportGenesis(ctx)
	return cdc.MustMarshalJSON(gs)
}

// ConsensusVersion implements AppModule/ConsensusVersion.
func (AppModule) ConsensusVersion() uint64 { return 1 }

//____________________________________________________________________________

// AppModuleSimulation functions

// // GenerateGenesisState creates a randomized GenState of the token module.
// func (AppModule) GenerateGenesisState(simState *module.SimulationState) {
// 	simulation.RandomizedGenState(simState)
// }

// // ProposalContents returns all the token content functions used to
// // simulate token proposals.
// func (AppModule) ProposalContents(simState module.SimulationState) []simtypes.WeightedProposalContent {
// 	return simulation.ProposalContents()
// }

// // RandomizedParams creates randomized token param changes for the simulator.
// func (AppModule) RandomizedParams(r *rand.Rand) []simtypes.ParamChange {
// 	return simulation.ParamChanges(r)
// }

// // RegisterStoreDecoder registers a decoder for token module's types
// func (am AppModule) RegisterStoreDecoder(sdr sdk.StoreDecoderRegistry) {
// 	sdr[types.StoreKey] = simulation.NewDecodeStore(am.cdc)
// }

// // WeightedOperations returns the all the token module operations with their respective weights.
// func (am AppModule) WeightedOperations(simState module.SimulationState) []simtypes.WeightedOperation {
// 	return simulation.WeightedOperations(
// 		simState.AppParams, simState.Cdc,
// 		am.stakingKeeper, am.slashingKeeper, am.keeper, simState.Contents,
// 	)
// }

package module

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/spf13/cobra"

	"cosmossdk.io/core/appmodule"
	"cosmossdk.io/core/store"
	"cosmossdk.io/depinject"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"

	modulev1 "github.com/Finschia/finschia-sdk/api/lbm/collection/module/v1"
	"github.com/Finschia/finschia-sdk/x/collection"
	"github.com/Finschia/finschia-sdk/x/collection/client/cli"
	"github.com/Finschia/finschia-sdk/x/collection/keeper"
)

// ConsensusVersion defines the current x/collection module consensus version.
const ConsensusVersion = 2

var (
	_ module.AppModuleBasic = AppModule{}
	_ module.HasGenesis     = AppModule{}
	_ module.HasServices    = AppModule{}

	_ appmodule.AppModule = AppModule{}
)

// AppModuleBasic defines the basic application module used by the collection module.
type AppModuleBasic struct {
	cdc codec.Codec
}

// Name returns the ModuleName
func (AppModuleBasic) Name() string {
	return collection.ModuleName
}

// RegisterLegacyAminoCodec registers the collection types on the LegacyAmino codec
func (AppModuleBasic) RegisterLegacyAminoCodec(_ *codec.LegacyAmino) {}

// RegisterInterfaces registers the collection module's interface types
func (AppModuleBasic) RegisterInterfaces(registry codectypes.InterfaceRegistry) {
	collection.RegisterInterfaces(registry)
}

// DefaultGenesis returns default genesis state as raw bytes for the collection
// module.
func (AppModuleBasic) DefaultGenesis(cdc codec.JSONCodec) json.RawMessage {
	return cdc.MustMarshalJSON(collection.DefaultGenesisState())
}

// ValidateGenesis performs genesis state validation for the collection module.
func (ab AppModuleBasic) ValidateGenesis(cdc codec.JSONCodec, _ client.TxEncodingConfig, bz json.RawMessage) error {
	var data collection.GenesisState
	if err := cdc.UnmarshalJSON(bz, &data); err != nil {
		return fmt.Errorf("failed to unmarshal %s genesis state: %w", collection.ModuleName, err)
	}

	return collection.ValidateGenesis(data, ab.cdc.InterfaceRegistry().SigningContext().AddressCodec())
}

// RegisterGRPCGatewayRoutes registers the gRPC Gateway routes for the collection module.
func (AppModuleBasic) RegisterGRPCGatewayRoutes(clientCtx client.Context, mux *runtime.ServeMux) {
	if err := collection.RegisterQueryHandlerClient(context.Background(), mux, collection.NewQueryClient(clientCtx)); err != nil {
		panic(err)
	}
}

// GetQueryCmd returns the cli query commands for this module
func (ab AppModuleBasic) GetQueryCmd() *cobra.Command {
	return cli.NewQueryCmd()
}

// GetTxCmd returns the transaction commands for this module
func (ab AppModuleBasic) GetTxCmd() *cobra.Command {
	return cli.NewTxCmd()
}

// ____________________________________________________________________________

// AppModule implements an application module for the collection module.
type AppModule struct {
	AppModuleBasic

	keeper keeper.Keeper
}

// NewAppModule creates a new AppModule object
func NewAppModule(cdc codec.Codec, keeper keeper.Keeper) AppModule {
	return AppModule{
		AppModuleBasic: AppModuleBasic{cdc},
		keeper:         keeper,
	}
}

// IsOnePerModuleType implements the depinject.OnePerModuleType interface.
func (am AppModule) IsOnePerModuleType() {}

// IsAppModule implements the appmodule.AppModule interface.
func (am AppModule) IsAppModule() {}

// RegisterServices registers a GRPC query service to respond to the
// module-specific GRPC queries.
func (am AppModule) RegisterServices(cfg module.Configurator) {
	collection.RegisterMsgServer(cfg.MsgServer(), keeper.NewMsgServer(am.keeper))
	collection.RegisterQueryServer(cfg.QueryServer(), keeper.NewQueryServer(am.keeper))

	m := keeper.NewMigrator(am.keeper)
	if err := cfg.RegisterMigration(collection.ModuleName, 1, m.Migrate1to2); err != nil {
		panic(fmt.Sprintf("failed to migrate x/%s from version 1 to 2: %v", collection.ModuleName, err))
	}
	if err := cfg.RegisterMigration(collection.ModuleName, 2, m.Migrate2to3); err != nil {
		panic(fmt.Sprintf("failed to migrate x/%s from version 2 to 3: %v", collection.ModuleName, err))
	}
}

// InitGenesis performs genesis initialization for the collection module. It returns
// no validator updates.
func (am AppModule) InitGenesis(ctx sdk.Context, cdc codec.JSONCodec, data json.RawMessage) {
	var genesisState collection.GenesisState
	cdc.MustUnmarshalJSON(data, &genesisState)
	am.keeper.InitGenesis(ctx, &genesisState)
}

// ExportGenesis returns the exported genesis state as raw bytes for the collection
// module.
func (am AppModule) ExportGenesis(ctx sdk.Context, cdc codec.JSONCodec) json.RawMessage {
	gs := am.keeper.ExportGenesis(ctx)
	return cdc.MustMarshalJSON(gs)
}

// ConsensusVersion implements AppModule/ConsensusVersion.
func (AppModule) ConsensusVersion() uint64 { return ConsensusVersion }

//
// App Wiring Setup
//

func init() {
	appmodule.Register(&modulev1.Module{},
		appmodule.Provide(ProvideModule),
	)
}

type CollectionInputs struct {
	depinject.In

	Cdc          codec.Codec
	StoreService store.KVStoreService
}

type CollectionOutputs struct {
	depinject.Out

	Keeper keeper.Keeper
	Module appmodule.AppModule
}

func ProvideModule(in CollectionInputs) CollectionOutputs {
	k := keeper.NewKeeper(in.Cdc, in.StoreService)
	m := NewAppModule(in.Cdc, k)

	return CollectionOutputs{
		Keeper: k,
		Module: m,
	}
}

package bank

import (
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/gorilla/mux"
	"github.com/link-chain/link/x/bank/client/rest"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"

	"github.com/link-chain/link/x/bank/internal/keeper"
)

var (
	_ module.AppModule      = AppModule{}
	_ module.AppModuleBasic = AppModuleBasic{}
)

// app module basics object
type AppModuleBasic struct {
	CosmosAppModuleBasic
}

// register module codec
func (amb AppModuleBasic) RegisterCodec(cdc *codec.Codec) {
	RegisterCodec(cdc)
}

// register rest routes
func (amb AppModuleBasic) RegisterRESTRoutes(ctx context.CLIContext, rtr *mux.Router) {
	rest.RegisterRoutes(ctx, rtr)
	amb.CosmosAppModuleBasic.RegisterRESTRoutes(ctx, rtr)
}

//___________________________
// app module
type AppModule struct {
	CosmosAppModule
	appModuleBasic AppModuleBasic
	keeper         Keeper
}

// NewAppModule creates a new AppModule object
func NewAppModule(keeper Keeper, accountKeeper AccountKeeper) AppModule {
	cosmosAppModule := NewCosmosAppModule(keeper, accountKeeper)
	return AppModule{
		CosmosAppModule: cosmosAppModule,
		appModuleBasic:  AppModuleBasic{},
		keeper:          keeper,
	}
}

// register module codec
func (am AppModule) RegisterCodec(cdc *codec.Codec) {
	am.appModuleBasic.RegisterCodec(cdc)
}

// register rest routes
func (am AppModule) RegisterRESTRoutes(ctx context.CLIContext, rtr *mux.Router) {
	am.appModuleBasic.RegisterRESTRoutes(ctx, rtr)
}

// module querier
func (am AppModule) NewQuerierHandler() sdk.Querier {
	fallbackQuerier := am.CosmosAppModule.NewQuerierHandler()
	return keeper.NewQuerier(am.keeper, fallbackQuerier)
}

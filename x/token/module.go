package token

import (
	"encoding/json"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/x/upgrade"
	"github.com/line/lbm-sdk/v2/x/token/client/cli"
	"github.com/line/lbm-sdk/v2/x/token/client/rest"
	"github.com/line/lbm-sdk/v2/x/token/internal/handler"
	"github.com/line/lbm-sdk/v2/x/token/internal/keeper"
	"github.com/line/lbm-sdk/v2/x/token/internal/legacy"
	"github.com/line/lbm-sdk/v2/x/token/internal/querier"

	"github.com/gorilla/mux"
	"github.com/spf13/cobra"

	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
)

var (
	_ module.AppModule      = AppModule{}
	_ module.AppModuleBasic = AppModuleBasic{}
)

// app module basics object
type AppModuleBasic struct{}

// module name
func (AppModuleBasic) Name() string { return ModuleName }

// register module codec
func (AppModuleBasic) RegisterCodec(cdc *codec.Codec) { RegisterCodec(cdc) }

// default genesis state
func (AppModuleBasic) DefaultGenesis() json.RawMessage {
	return ModuleCdc.MustMarshalJSON(DefaultGenesisState())
}

// module validate genesis
func (AppModuleBasic) ValidateGenesis(bz json.RawMessage) error {
	var data GenesisState
	err := ModuleCdc.UnmarshalJSON(bz, &data)
	if err != nil {
		return err
	}
	return ValidateGenesis(data)
}

// register rest routes
func (AppModuleBasic) RegisterRESTRoutes(ctx context.CLIContext, rtr *mux.Router) {
	rest.RegisterRoutes(ctx, rtr)
}

// get the root tx command of this module
func (AppModuleBasic) GetTxCmd(cdc *codec.Codec) *cobra.Command {
	return cli.GetTxCmd(cdc)
}

// get the root query command of this module
func (AppModuleBasic) GetQueryCmd(cdc *codec.Codec) *cobra.Command {
	return cli.GetQueryCmd(cdc)
}

func (AppModuleBasic) GetUpgradeHandler(version string) upgrade.UpgradeHandler {
	return legacy.UpgradeHandler(version)
}

// ___________________________
// app module
type AppModule struct {
	AppModuleBasic
	keeper keeper.Keeper
}

// NewAppModule creates a new AppModule object
func NewAppModule(keeper keeper.Keeper) AppModule {
	return AppModule{
		AppModuleBasic: AppModuleBasic{},
		keeper:         keeper,
	}
}

// module name
func (AppModule) Name() string { return ModuleName }

// register invariants
// TODO: should this module need invariants?
func (am AppModule) RegisterInvariants(_ sdk.InvariantRegistry) {}

// module message route name
func (AppModule) Route() string { return RouterKey }

// module handler
func (am AppModule) NewHandler() sdk.Handler { return handler.NewHandler(am.keeper) }

// module querier route name
func (AppModule) QuerierRoute() string { return RouterKey }

// module querier
func (am AppModule) NewQuerierHandler() sdk.Querier {
	return querier.NewQuerier(am.keeper)
}

// module init-genesis
func (am AppModule) InitGenesis(ctx sdk.Context, data json.RawMessage) []abci.ValidatorUpdate {
	var genesisState GenesisState
	ModuleCdc.MustUnmarshalJSON(data, &genesisState)
	InitGenesis(ctx, am.keeper, genesisState)
	return []abci.ValidatorUpdate{}
}

// module export genesis
func (am AppModule) ExportGenesis(ctx sdk.Context) json.RawMessage {
	gs := ExportGenesis(ctx, am.keeper)
	return ModuleCdc.MustMarshalJSON(gs)
}

// module begin-block
func (AppModule) BeginBlock(_ sdk.Context, _ abci.RequestBeginBlock) {}

// module end-block
func (AppModule) EndBlock(_ sdk.Context, _ abci.RequestEndBlock) []abci.ValidatorUpdate {
	return []abci.ValidatorUpdate{}
}

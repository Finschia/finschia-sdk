package app

import (
	"io"
	"os"
	"path/filepath"

	"github.com/cosmos/cosmos-sdk/x/gov"
	"github.com/line/lbm-sdk/x/account"
	"github.com/line/lbm-sdk/x/coin"
	"github.com/line/lbm-sdk/x/contract"
	"github.com/line/lbm-sdk/x/token"
	"github.com/line/lbm-sdk/x/wasm"
	wasmclient "github.com/line/lbm-sdk/x/wasm/client"
	"github.com/spf13/viper"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/cli"
	"github.com/tendermint/tendermint/libs/log"
	tmos "github.com/tendermint/tendermint/libs/os"
	dbm "github.com/tendermint/tm-db"

	bam "github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/simapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/vesting"
	distr "github.com/cosmos/cosmos-sdk/x/distribution"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	"github.com/cosmos/cosmos-sdk/x/params"
	paramsclient "github.com/cosmos/cosmos-sdk/x/params/client"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/cosmos/cosmos-sdk/x/supply"

	"github.com/cosmos/cosmos-sdk/version"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/line/lbm-sdk/x/collection"
)

const appName = "LinkWasmApp"

var (
	// DefaultCLIHome for linkcli
	DefaultCLIHome = os.ExpandEnv("$HOME/.linkwasmcli")

	// DefaultNodeHome for linkd
	DefaultNodeHome = os.ExpandEnv("$HOME/.linkwasmd")

	// ModuleBasics is in charge of setting up basic,
	// non-dependant module elements, such as codec registration
	// and genesis verification.
	ModuleBasics = module.NewBasicManager(
		genutil.AppModuleBasic{},
		auth.AppModuleBasic{},
		coin.AppModuleBasic{},
		staking.AppModuleBasic{},
		params.AppModuleBasic{},
		supply.AppModuleBasic{},
		gov.NewAppModuleBasic(
			append(
				wasmclient.ProposalHandlers,
				paramsclient.ProposalHandler,
			)...,
		),
		token.AppModuleBasic{},
		collection.AppModuleBasic{},
		account.AppModuleBasic{},
		wasm.AppModuleBasic{},
	)

	// module account permissions
	maccPerms = map[string][]string{
		auth.FeeCollectorName:     nil,
		staking.BondedPoolName:    {supply.Burner, supply.Staking},
		staking.NotBondedPoolName: {supply.Burner, supply.Staking},
		token.ModuleName:          {supply.Minter, supply.Burner},
		collection.ModuleName:     {supply.Minter, supply.Burner},
		gov.ModuleName:            {supply.Burner},
	}
)

// custom tx codec
func MakeCodec() *codec.Codec {
	var cdc = codec.New()

	ModuleBasics.RegisterCodec(cdc)
	vesting.RegisterCodec(cdc)
	sdk.RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)
	codec.RegisterEvidences(cdc)

	return cdc
}

// Extended ABCI application
type LinkApp struct {
	*bam.BaseApp
	cdc *codec.Codec

	// keys to access the substores
	keys  map[string]*sdk.KVStoreKey
	tkeys map[string]*sdk.TransientStoreKey

	// subspaces
	subspaces map[string]params.Subspace

	// keepers
	accountKeeper    auth.AccountKeeper
	bankKeeper       bank.Keeper
	coinKeeper       coin.Keeper
	supplyKeeper     supply.Keeper
	stakingKeeper    staking.Keeper
	paramsKeeper     params.Keeper
	govKeeper        gov.Keeper
	tokenKeeper      token.Keeper
	collectionKeeper collection.Keeper
	wasmKeeper       wasm.Keeper

	// the module manager
	mm *module.Manager

	// simulation manager
	sm *module.SimulationManager
}

type WasmWrapper struct {
	Wasm wasm.Config `mapstructure:"wasm"`
}

// NewLinkApp returns a reference to an initialized LinkApp.
func NewLinkApp(logger log.Logger, db dbm.DB, traceStore io.Writer, loadLatest bool, skipUpgradeHeights map[int64]bool,
	invCheckPeriod uint, baseAppOptions ...func(*bam.BaseApp)) *LinkApp {
	cdc := MakeCodec()

	bApp := bam.NewBaseApp(appName, logger, db, auth.DefaultTxDecoder(cdc), baseAppOptions...)
	bApp.SetCommitMultiStoreTracer(traceStore)
	bApp.SetAppVersion(version.Version)

	keys := sdk.NewKVStoreKeys(
		bam.MainStoreKey,
		auth.StoreKey,
		staking.StoreKey,
		supply.StoreKey,
		params.StoreKey,
		gov.StoreKey,
		token.StoreKey,
		collection.StoreKey,
		coin.StoreKey,
		contract.StoreKey,
		wasm.StoreKey,
	)
	tkeys := sdk.NewTransientStoreKeys(staking.TStoreKey, params.TStoreKey)

	app := &LinkApp{
		BaseApp:   bApp,
		cdc:       cdc,
		keys:      keys,
		tkeys:     tkeys,
		subspaces: make(map[string]params.Subspace),
	}

	// init params keeper and subspaces
	app.paramsKeeper = params.NewKeeper(app.cdc, keys[params.StoreKey], tkeys[params.TStoreKey])
	app.subspaces[auth.ModuleName] = app.paramsKeeper.Subspace(auth.DefaultParamspace)
	app.subspaces[bank.ModuleName] = app.paramsKeeper.Subspace(bank.DefaultParamspace)
	app.subspaces[staking.ModuleName] = app.paramsKeeper.Subspace(staking.DefaultParamspace)
	app.subspaces[collection.ModuleName] = app.paramsKeeper.Subspace(collection.DefaultParamspace)
	app.subspaces[gov.ModuleName] = app.paramsKeeper.Subspace(gov.DefaultParamspace).WithKeyTable(gov.ParamKeyTable())
	app.subspaces[wasm.ModuleName] = app.paramsKeeper.Subspace(wasm.DefaultParamspace)

	// add keepers
	app.accountKeeper = auth.NewAccountKeeper(app.cdc, keys[auth.StoreKey], app.subspaces[auth.ModuleName], auth.ProtoBaseAccount)
	app.bankKeeper = bank.NewBaseKeeper(app.accountKeeper, app.subspaces[bank.ModuleName], app.ModuleAccountAddrs())
	app.coinKeeper = coin.NewKeeper(app.bankKeeper, keys[coin.StoreKey])
	app.supplyKeeper = supply.NewKeeper(app.cdc, keys[supply.StoreKey], app.accountKeeper, app.bankKeeper, maccPerms)
	app.stakingKeeper = staking.NewKeeper(app.cdc, keys[staking.StoreKey], app.supplyKeeper, app.subspaces[staking.ModuleName])

	contractKeeper := contract.NewContractKeeper(cdc, keys[contract.StoreKey])
	app.tokenKeeper = token.NewKeeper(app.cdc, app.accountKeeper, contractKeeper, keys[token.StoreKey])
	app.collectionKeeper = collection.NewKeeper(
		app.cdc,
		app.accountKeeper,
		contractKeeper,
		app.subspaces[collection.ModuleName],
		keys[collection.StoreKey],
	)

	// just re-use the full router - do we want to limit this more?
	var wasmRouter = bApp.Router()

	// encodeRouter
	tokenEncodeHandler := token.NewMsgEncodeHandler(app.tokenKeeper)
	collectionEncoder := collection.NewMsgEncodeHandler(app.collectionKeeper)
	var encodeRouter = wasm.NewRouter()
	encodeRouter.AddRoute(token.EncodeRouterKey, tokenEncodeHandler)
	encodeRouter.AddRoute(collection.EncodeRouterKey, collectionEncoder)

	// queryRouter
	tokenQuerier := token.NewQuerier(app.tokenKeeper)
	tokenQueryEncoder := token.NewQueryEncoder(tokenQuerier)
	collectionQuerier := collection.NewQuerier(app.collectionKeeper)
	collectionQueryEncoder := collection.NewQueryEncoder(collectionQuerier)
	var querierRouter = wasm.NewQuerierRouter()
	querierRouter.AddRoute(token.EncodeRouterKey, tokenQueryEncoder)
	querierRouter.AddRoute(collection.EncodeRouterKey, collectionQueryEncoder)

	// better way to get this dir???
	homeDir := viper.GetString(cli.HomeFlag)
	wasmDir := filepath.Join(homeDir, "wasm")

	wasmWrap := WasmWrapper{Wasm: wasm.DefaultWasmConfig()}
	err := viper.Unmarshal(&wasmWrap)
	if err != nil {
		panic("error while reading wasm config: " + err.Error())
	}
	wasmConfig := wasmWrap.Wasm
	supportedFeatures := "staking,link"

	app.wasmKeeper = wasm.NewKeeper(
		app.cdc,
		keys[wasm.StoreKey],
		app.subspaces[wasm.ModuleName],
		app.accountKeeper,
		app.coinKeeper,
		app.stakingKeeper,
		distr.Keeper{},
		wasmRouter,
		encodeRouter,
		querierRouter,
		wasmDir,
		wasmConfig,
		supportedFeatures,
		nil,
		nil,
	)

	// register the proposal types
	govRouter := gov.NewRouter()
	govRouter.AddRoute(gov.RouterKey, gov.ProposalHandler).
		AddRoute(params.RouterKey, params.NewParamChangeProposalHandler(app.paramsKeeper)).
		AddRoute(wasm.RouterKey, wasm.NewWasmProposalHandler(app.wasmKeeper, wasm.EnableAllProposals))
	app.govKeeper = gov.NewKeeper(
		app.cdc, keys[gov.StoreKey], app.subspaces[gov.ModuleName], app.supplyKeeper,
		&app.stakingKeeper, govRouter,
	)

	// NOTE: Any module instantiated in the module manager that is later modified
	// must be passed by reference here.
	app.mm = module.NewManager(
		genutil.NewAppModule(app.accountKeeper, app.stakingKeeper, app.BaseApp.DeliverTx),
		auth.NewAppModule(app.accountKeeper),
		coin.NewAppModule(app.coinKeeper),
		supply.NewAppModule(app.supplyKeeper, app.accountKeeper),
		staking.NewAppModule(app.stakingKeeper, app.accountKeeper, app.supplyKeeper),
		gov.NewAppModule(app.govKeeper, app.accountKeeper, app.supplyKeeper),
		token.NewAppModule(app.tokenKeeper),
		collection.NewAppModule(app.collectionKeeper),
		account.NewAppModule(app.accountKeeper),
		wasm.NewAppModule(app.wasmKeeper),
	)
	app.mm.SetOrderEndBlockers(gov.ModuleName, staking.ModuleName)

	// NOTE: The genutils module must occur after staking so that pools are
	// properly initialized with tokens from genesis accounts.
	app.mm.SetOrderInitGenesis(
		staking.ModuleName,
		auth.ModuleName,
		gov.ModuleName,
		supply.ModuleName,
		coin.ModuleName,
		genutil.ModuleName,
		token.ModuleName,
		collection.ModuleName,
		account.ModuleName,
		wasm.ModuleName,
	)

	app.mm.RegisterRoutes(app.Router(), app.QueryRouter())

	// create the simulation manager and define the order of the modules for deterministic simulations
	//
	// NOTE: This is not required for apps that don't use the simulator for fuzz testing
	// transactions.
	app.sm = module.NewSimulationManager(
		auth.NewAppModule(app.accountKeeper),
		supply.NewAppModule(app.supplyKeeper, app.accountKeeper),
		staking.NewAppModule(app.stakingKeeper, app.accountKeeper, app.supplyKeeper),
		collection.NewAppModule(app.collectionKeeper),
		gov.NewAppModule(app.govKeeper, app.accountKeeper, app.supplyKeeper),
		// TODO: Implement AppModuleSimulation interface in each module.
		// bank.NewAppModule(app.coinKeeper),
		// token.NewAppModule(app.tokenKeeper),
		// account.NewAppModule(app.accountKeeper),
	)

	app.sm.RegisterStoreDecoders()

	// initialize stores
	app.MountKVStores(keys)
	app.MountTransientStores(tkeys)

	// initialize BaseApp
	app.SetInitChainer(app.InitChainer)
	app.SetBeginBlocker(app.BeginBlocker)
	app.SetAnteHandler(auth.NewAnteHandler(app.accountKeeper, app.supplyKeeper, auth.DefaultSigVerificationGasConsumer))
	app.SetEndBlocker(app.EndBlocker)

	if loadLatest {
		err := app.LoadLatestVersion(app.keys[bam.MainStoreKey])
		if err != nil {
			tmos.Exit(err.Error())
		}
	}

	return app
}

// application updates every begin block
func (app *LinkApp) BeginBlocker(ctx sdk.Context, req abci.RequestBeginBlock) abci.ResponseBeginBlock {
	return app.mm.BeginBlock(ctx, req)
}

// application updates every end block
func (app *LinkApp) EndBlocker(ctx sdk.Context, req abci.RequestEndBlock) abci.ResponseEndBlock {
	return app.mm.EndBlock(ctx, req)
}

// application update at chain initialization
func (app *LinkApp) InitChainer(ctx sdk.Context, req abci.RequestInitChain) abci.ResponseInitChain {
	var genesisState simapp.GenesisState
	app.cdc.MustUnmarshalJSON(req.AppStateBytes, &genesisState)

	return app.mm.InitGenesis(ctx, genesisState)
}

// load a particular height
func (app *LinkApp) LoadHeight(height int64) error {
	return app.LoadVersion(height, app.keys[bam.MainStoreKey])
}

// ModuleAccountAddrs returns all the app's module account addresses.
func (app *LinkApp) ModuleAccountAddrs() map[string]bool {
	modAccAddrs := make(map[string]bool)
	for acc := range maccPerms {
		modAccAddrs[supply.NewModuleAddress(acc).String()] = true
	}

	return modAccAddrs
}

// Codec returns the application's sealed codec.
func (app *LinkApp) Codec() *codec.Codec {
	return app.cdc
}

// SimulationManager implements the SimulationApp interface
func (app *LinkApp) SimulationManager() *module.SimulationManager {
	return app.sm
}

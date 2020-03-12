package app

import (
	"io"
	"os"

	"github.com/line/link/x/account"
	"github.com/line/link/x/bank"
	"github.com/line/link/x/contract"
	"github.com/line/link/x/iam"
	"github.com/line/link/x/token"

	abci "github.com/tendermint/tendermint/abci/types"
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
	"github.com/cosmos/cosmos-sdk/x/genutil"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/cosmos/cosmos-sdk/x/supply"

	cbank "github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/line/link/version"
	"github.com/line/link/x/collection"
)

const appName = "LinkApp"

var (
	// DefaultCLIHome for linkcli
	DefaultCLIHome = os.ExpandEnv("$HOME/.linkcli")

	// DefaultNodeHome for linkd
	DefaultNodeHome = os.ExpandEnv("$HOME/.linkd")

	// ModuleBasics is in charge of setting up basic,
	// non-dependant module elements, such as codec registration
	// and genesis verification.
	ModuleBasics = module.NewBasicManager(
		genutil.AppModuleBasic{},
		auth.AppModuleBasic{},
		bank.AppModuleBasic{},
		staking.AppModuleBasic{},
		params.AppModuleBasic{},
		supply.AppModuleBasic{},
		token.AppModuleBasic{},
		collection.AppModuleBasic{},
		iam.AppModuleBasic{},
		account.AppModuleBasic{},
	)

	// module account permissions
	maccPerms = map[string][]string{
		auth.FeeCollectorName:     nil,
		staking.BondedPoolName:    {supply.Burner, supply.Staking},
		staking.NotBondedPoolName: {supply.Burner, supply.Staking},
		token.ModuleName:          {supply.Minter, supply.Burner},
		collection.ModuleName:     {supply.Minter, supply.Burner},
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

	// keepers
	accountKeeper    auth.AccountKeeper
	cbankKeeper      cbank.Keeper
	bankKeeper       bank.Keeper
	supplyKeeper     supply.Keeper
	stakingKeeper    staking.Keeper
	paramsKeeper     params.Keeper
	tokenKeeper      token.Keeper
	collectionKeeper collection.Keeper
	iamKeeper        iam.Keeper

	// the module manager
	mm *module.Manager

	// simulation manager
	sm *module.SimulationManager
}

// NewLinkApp returns a reference to an initialized LinkApp.
func NewLinkApp(logger log.Logger, db dbm.DB, traceStore io.Writer, loadLatest bool,
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
		token.StoreKey,
		collection.StoreKey,
		iam.StoreKey,
		bank.StoreKey,
		contract.StoreKey,
	)
	tkeys := sdk.NewTransientStoreKeys(staking.TStoreKey, params.TStoreKey)

	app := &LinkApp{
		BaseApp: bApp,
		cdc:     cdc,
		keys:    keys,
		tkeys:   tkeys,
	}

	// init params keeper and subspaces
	app.paramsKeeper = params.NewKeeper(app.cdc, keys[params.StoreKey], tkeys[params.TStoreKey])
	authSubspace := app.paramsKeeper.Subspace(auth.DefaultParamspace)
	cbankSubspace := app.paramsKeeper.Subspace(cbank.DefaultParamspace)
	stakingSubspace := app.paramsKeeper.Subspace(staking.DefaultParamspace)

	app.iamKeeper = iam.NewKeeper(cdc, keys[iam.StoreKey])

	// add keepers
	app.accountKeeper = auth.NewAccountKeeper(app.cdc, keys[auth.StoreKey], authSubspace, auth.ProtoBaseAccount)
	app.cbankKeeper = cbank.NewBaseKeeper(app.accountKeeper, cbankSubspace, app.ModuleAccountAddrs())
	app.bankKeeper = bank.NewKeeper(app.cbankKeeper, keys[bank.StoreKey])
	app.supplyKeeper = supply.NewKeeper(app.cdc, keys[supply.StoreKey], app.accountKeeper, app.cbankKeeper, maccPerms)
	app.stakingKeeper = staking.NewKeeper(app.cdc, keys[staking.StoreKey], app.supplyKeeper, stakingSubspace)

	contractKeeper := contract.NewContractKeeper(cdc, keys[contract.StoreKey])
	app.tokenKeeper = token.NewKeeper(app.cdc, app.accountKeeper, app.iamKeeper.WithPrefix(token.ModuleName), contractKeeper, keys[token.StoreKey])
	app.collectionKeeper = collection.NewKeeper(app.cdc, app.accountKeeper, app.iamKeeper.WithPrefix(collection.ModuleName), contractKeeper, keys[collection.StoreKey])

	// NOTE: Any module instantiated in the module manager that is later modified
	// must be passed by reference here.
	app.mm = module.NewManager(
		genutil.NewAppModule(app.accountKeeper, app.stakingKeeper, app.BaseApp.DeliverTx),
		auth.NewAppModule(app.accountKeeper),
		bank.NewAppModule(app.bankKeeper),
		supply.NewAppModule(app.supplyKeeper, app.accountKeeper),
		staking.NewAppModule(app.stakingKeeper, app.accountKeeper, app.supplyKeeper),
		token.NewAppModule(app.tokenKeeper),
		collection.NewAppModule(app.collectionKeeper),
		account.NewAppModule(app.accountKeeper),
	)
	app.mm.SetOrderEndBlockers(staking.ModuleName)

	// NOTE: The genutils module must occur after staking so that pools are
	// properly initialized with tokens from genesis accounts.
	app.mm.SetOrderInitGenesis(
		staking.ModuleName,
		auth.ModuleName,
		supply.ModuleName,
		bank.ModuleName,
		genutil.ModuleName,
		token.ModuleName,
		collection.ModuleName,
		account.ModuleName,
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
		// TODO: Implement AppModuleSimulation interface in each module.
		//bank.NewAppModule(app.bankKeeper),
		//token.NewAppModule(app.tokenKeeper),
		//collection.NewAppModule(app.collectionKeeper),
		//account.NewAppModule(app.accountKeeper),
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
